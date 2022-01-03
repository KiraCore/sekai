package ante_test

import (
	"errors"
	"fmt"

	customante "github.com/KiraCore/sekai/app/ante"
	"github.com/KiraCore/sekai/types"
	govtypes "github.com/KiraCore/sekai/x/gov/types"
	tokenstypes "github.com/KiraCore/sekai/x/tokens/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bank "github.com/cosmos/cosmos-sdk/x/bank/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
)

func (suite *AnteTestSuite) SetBalance(addr sdk.AccAddress, coin sdk.Coin) {
	coins := sdk.Coins{coin}
	suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, coins)
	suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, addr, coins)
}

// Test that simulate transaction process execution fee correctly on ante handler step
func (suite *AnteTestSuite) TestCustomAnteHandlerExecutionFee() {
	suite.SetupTest(false) // reset

	// set execution fee for set network properties
	suite.app.CustomGovKeeper.SetExecutionFee(suite.ctx, &govtypes.ExecutionFee{
		Name:              types.MsgTypeSetNetworkProperties,
		TransactionType:   types.MsgTypeSetNetworkProperties,
		ExecutionFee:      10000,
		FailureFee:        1000,
		Timeout:           0,
		DefaultParameters: 0,
	})
	suite.app.CustomGovKeeper.SetNetworkProperties(suite.ctx, &govtypes.NetworkProperties{
		MinTxFee:                 2,
		MaxTxFee:                 10000,
		EnableForeignFeePayments: true,
	})

	// Same data for every test cases
	accounts := suite.CreateTestAccounts(5)

	suite.SetBalance(accounts[0].acc.GetAddress(), sdk.NewInt64Coin("ukex", 10000))
	suite.SetBalance(accounts[1].acc.GetAddress(), sdk.NewInt64Coin("ukex", 10000))
	suite.SetBalance(accounts[2].acc.GetAddress(), sdk.NewInt64Coin("ukex", 10000))
	suite.SetBalance(accounts[3].acc.GetAddress(), sdk.NewInt64Coin("ukex", 1))
	suite.SetBalance(accounts[4].acc.GetAddress(), sdk.NewInt64Coin("ukex", 10000))
	suite.SetBalance(accounts[4].acc.GetAddress(), sdk.NewInt64Coin("ubtc", 10000))

	defaultFee := sdk.NewCoins(sdk.NewInt64Coin("ukex", 100))
	gasLimit := testdata.NewTestGasLimit()
	privs := []cryptotypes.PrivKey{accounts[0].priv, accounts[1].priv, accounts[2].priv, accounts[3].priv, accounts[4].priv}
	accNums := []uint64{0, 1, 2, 3, 4}

	testCases := []TestCase{
		{
			"insufficient max execution fee set",
			func() ([]sdk.Msg, []cryptotypes.PrivKey, []uint64, []uint64, sdk.Coins) {
				msgs := []sdk.Msg{
					govtypes.NewMsgSetNetworkProperties(accounts[0].acc.GetAddress(), &govtypes.NetworkProperties{
						MinTxFee:                 2,
						MaxTxFee:                 10000,
						EnableForeignFeePayments: true,
					}),
				}
				return msgs, privs[0:1], accNums[0:1], []uint64{0}, defaultFee
			},
			true,
			false,
			errors.New("fee 100ukex(100) is less than max execution fee 10000ukex: invalid request"),
		},
		{
			"execution failure fee deduction",
			func() ([]sdk.Msg, []cryptotypes.PrivKey, []uint64, []uint64, sdk.Coins) {
				msgs := []sdk.Msg{
					govtypes.NewMsgSetNetworkProperties(accounts[1].acc.GetAddress(), &govtypes.NetworkProperties{
						MinTxFee:                 2,
						MaxTxFee:                 10000,
						EnableForeignFeePayments: true,
					}),
				}
				return msgs, privs[1:2], accNums[1:2], []uint64{0}, sdk.NewCoins(sdk.NewInt64Coin("ukex", 10000))
			},
			true,
			true,
			nil,
		},
		{
			"no execution fee deduction when does not exist",
			func() ([]sdk.Msg, []cryptotypes.PrivKey, []uint64, []uint64, sdk.Coins) {
				msgs := []sdk.Msg{
					govtypes.NewMsgSetExecutionFee(
						types.MsgTypeSetNetworkProperties,
						types.MsgTypeSetNetworkProperties,
						10000,
						1000,
						0,
						0,
						accounts[2].acc.GetAddress(),
					),
				}
				return msgs, privs[2:3], accNums[2:3], []uint64{0}, defaultFee
			},
			false,
			true,
			nil,
		},
		{
			"insufficient balance to pay for fee",
			func() ([]sdk.Msg, []cryptotypes.PrivKey, []uint64, []uint64, sdk.Coins) {
				msgs := []sdk.Msg{
					govtypes.NewMsgSetExecutionFee(
						types.MsgTypeSetNetworkProperties,
						types.MsgTypeSetNetworkProperties,
						10000,
						1000,
						0,
						0,
						accounts[3].acc.GetAddress(),
					),
				}
				return msgs, privs[3:4], accNums[3:4], []uint64{0}, sdk.NewCoins(sdk.NewInt64Coin("ukex", 10))
			},
			false,
			false,
			errors.New("1ukex is smaller than 10ukex: insufficient funds"),
		},
		{
			"fee out of range",
			func() ([]sdk.Msg, []cryptotypes.PrivKey, []uint64, []uint64, sdk.Coins) {
				msgs := []sdk.Msg{
					govtypes.NewMsgSetExecutionFee(
						types.MsgTypeSetNetworkProperties,
						types.MsgTypeSetNetworkProperties,
						10000,
						1000,
						0,
						0,
						accounts[4].acc.GetAddress(),
					),
				}
				return msgs, privs[4:5], accNums[4:5], []uint64{0}, sdk.NewCoins(sdk.NewInt64Coin("ukex", 1))
			},
			false,
			false,
			errors.New("fee 1ukex(1) is out of range [2, 10000]ukex: invalid request"),
		},
		{
			"foreign currency as fee payment when EnableForeignFeePayments is enabled by governance",
			func() ([]sdk.Msg, []cryptotypes.PrivKey, []uint64, []uint64, sdk.Coins) {
				suite.app.CustomGovKeeper.SetNetworkProperties(suite.ctx, &govtypes.NetworkProperties{
					MinTxFee:                 2,
					MaxTxFee:                 10000,
					EnableForeignFeePayments: true,
				})
				msgs := []sdk.Msg{
					govtypes.NewMsgSetExecutionFee(
						types.MsgTypeSetNetworkProperties,
						types.MsgTypeSetNetworkProperties,
						10000,
						1000,
						0,
						0,
						accounts[4].acc.GetAddress(),
					),
				}
				return msgs, privs[4:5], accNums[4:5], []uint64{0}, sdk.NewCoins(sdk.NewInt64Coin("ubtc", 10))
			},
			false,
			true,
			nil,
		},
		{
			"foreign currency as fee payment when EnableForeignFeePayments is disabled by governance",
			func() ([]sdk.Msg, []cryptotypes.PrivKey, []uint64, []uint64, sdk.Coins) {
				suite.app.CustomGovKeeper.SetNetworkProperties(suite.ctx, &govtypes.NetworkProperties{
					MinTxFee:                 2,
					MaxTxFee:                 10000,
					EnableForeignFeePayments: false,
				})
				msgs := []sdk.Msg{
					govtypes.NewMsgSetExecutionFee(
						types.MsgTypeSetNetworkProperties,
						types.MsgTypeSetNetworkProperties,
						10000,
						1000,
						0,
						0,
						accounts[4].acc.GetAddress(),
					),
				}
				return msgs, privs[4:5], accNums[4:5], []uint64{0}, sdk.NewCoins(sdk.NewInt64Coin("ubtc", 10))
			},
			false,
			false,
			errors.New("foreign fee payments is disabled by governance: invalid request"),
		},
		{
			"try sending non bond denom coins on poor network",
			func() ([]sdk.Msg, []cryptotypes.PrivKey, []uint64, []uint64, sdk.Coins) {
				suite.app.CustomGovKeeper.SetNetworkProperties(suite.ctx, &govtypes.NetworkProperties{
					MinTxFee:                 2,
					MaxTxFee:                 10000,
					EnableForeignFeePayments: true,
					MinValidators:            100,
				})
				msgs := []sdk.Msg{
					bank.NewMsgSend(
						accounts[4].acc.GetAddress(),
						accounts[3].acc.GetAddress(),
						sdk.NewCoins(sdk.NewInt64Coin("ubtc", 10)),
					),
				}
				return msgs, privs[4:5], accNums[4:5], []uint64{0}, sdk.NewCoins(sdk.NewInt64Coin("ukex", 10))
			},
			false,
			false,
			errors.New("only bond denom is allowed on poor network: invalid request"),
		},
		{
			"try sending more bond denom than restricted amount on poor network",
			func() ([]sdk.Msg, []cryptotypes.PrivKey, []uint64, []uint64, sdk.Coins) {
				suite.app.CustomGovKeeper.SetNetworkProperties(suite.ctx, &govtypes.NetworkProperties{
					MinTxFee:                 2,
					MaxTxFee:                 10000,
					EnableForeignFeePayments: true,
					MinValidators:            100,
					PoorNetworkMaxBankSend:   1000,
				})
				msgs := []sdk.Msg{
					bank.NewMsgSend(
						accounts[4].acc.GetAddress(),
						accounts[3].acc.GetAddress(),
						sdk.NewCoins(sdk.NewInt64Coin("ukex", 10000000)),
					),
				}
				return msgs, privs[4:5], accNums[4:5], []uint64{0}, sdk.NewCoins(sdk.NewInt64Coin("ukex", 10))
			},
			false,
			false,
			errors.New("only restricted amount send is allowed on poor network: invalid request"),
		},
		{
			"try sending lower than restriction amount on poor network",
			func() ([]sdk.Msg, []cryptotypes.PrivKey, []uint64, []uint64, sdk.Coins) {
				suite.app.CustomGovKeeper.SetNetworkProperties(suite.ctx, &govtypes.NetworkProperties{
					MinTxFee:                 2,
					MaxTxFee:                 10000,
					EnableForeignFeePayments: true,
					MinValidators:            100,
					PoorNetworkMaxBankSend:   1000,
				})
				msgs := []sdk.Msg{
					bank.NewMsgSend(
						accounts[4].acc.GetAddress(),
						accounts[3].acc.GetAddress(),
						sdk.NewCoins(sdk.NewInt64Coin("ukex", 1000)),
					),
				}
				return msgs, privs[4:5], accNums[4:5], []uint64{1}, sdk.NewCoins(sdk.NewInt64Coin("ubtc", 10))
			},
			false,
			true,
			nil,
		},
		{
			"try sending enabled message on poor network",
			func() ([]sdk.Msg, []cryptotypes.PrivKey, []uint64, []uint64, sdk.Coins) {
				suite.app.CustomGovKeeper.SetNetworkProperties(suite.ctx, &govtypes.NetworkProperties{
					MinTxFee:                 2,
					MaxTxFee:                 10000,
					EnableForeignFeePayments: true,
					MinValidators:            100,
				})
				msgs := []sdk.Msg{
					govtypes.NewMsgSetNetworkProperties(
						accounts[4].acc.GetAddress(),
						&govtypes.NetworkProperties{
							MinTxFee:                 2,
							MaxTxFee:                 10000,
							EnableForeignFeePayments: true,
							MinValidators:            100,
						},
					),
				}
				return msgs, privs[4:5], accNums[4:5], []uint64{2}, sdk.NewCoins(sdk.NewInt64Coin("ubtc", 1000))
			},
			false,
			true,
			nil,
		},
		{
			"try sending not enabled message on poor network",
			func() ([]sdk.Msg, []cryptotypes.PrivKey, []uint64, []uint64, sdk.Coins) {
				suite.app.CustomGovKeeper.SetNetworkProperties(suite.ctx, &govtypes.NetworkProperties{
					MinTxFee:                 2,
					MaxTxFee:                 10000,
					EnableForeignFeePayments: true,
					MinValidators:            100,
				})
				msgs := []sdk.Msg{
					govtypes.NewMsgSetExecutionFee(
						types.MsgTypeSetNetworkProperties,
						types.MsgTypeSetNetworkProperties,
						10000,
						1000,
						0,
						0,
						accounts[4].acc.GetAddress(),
					),
				}
				return msgs, privs[4:5], accNums[4:5], []uint64{0}, sdk.NewCoins(sdk.NewInt64Coin("ubtc", 10))
			},
			false,
			false,
			errors.New("invalid transaction type on poor network: invalid request"),
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.desc), func() {
			suite.txBuilder = suite.clientCtx.TxConfig.NewTxBuilder()
			msgs, privs, accNums, accSeqs, feeAmount := tc.buildTest()

			// this runs multi signature transaction with the params provided
			suite.RunTestCase(privs, msgs, feeAmount, gasLimit, accNums, accSeqs, suite.ctx.ChainID(), tc)
		})
	}
}

// Test that simulate transaction process fee range decorator correctly on ante handler step
func (suite *AnteTestSuite) TestValidateFeeRangeDecorator() {
	suite.SetupTest(false) // reset

	// set execution fee for set network properties
	suite.app.CustomGovKeeper.SetNetworkProperties(suite.ctx, &govtypes.NetworkProperties{
		MinTxFee:                 2,
		MaxTxFee:                 10000,
		EnableForeignFeePayments: true,
		EnableTokenBlacklist:     true,
		EnableTokenWhitelist:     true,
	})

	suite.app.CustomGovKeeper.SetExecutionFee(suite.ctx, &govtypes.ExecutionFee{
		Name:              types.MsgTypeSetNetworkProperties,
		TransactionType:   types.MsgTypeSetNetworkProperties,
		ExecutionFee:      10000,
		FailureFee:        1000,
		Timeout:           0,
		DefaultParameters: 0,
	})

	// Same data for every test cases
	accounts := suite.CreateTestAccounts(5)

	suite.SetBalance(accounts[0].acc.GetAddress(), sdk.NewInt64Coin("ukex", 10000))
	suite.SetBalance(accounts[0].acc.GetAddress(), sdk.NewInt64Coin("frozen", 10000))
	suite.SetBalance(accounts[0].acc.GetAddress(), sdk.NewInt64Coin("nofeetoken", 10000))
	suite.SetBalance(accounts[1].acc.GetAddress(), sdk.NewInt64Coin("ukex", 10000))
	suite.SetBalance(accounts[2].acc.GetAddress(), sdk.NewInt64Coin("ukex", 10000))
	suite.SetBalance(accounts[3].acc.GetAddress(), sdk.NewInt64Coin("ukex", 1))
	suite.SetBalance(accounts[4].acc.GetAddress(), sdk.NewInt64Coin("ukex", 10000))
	suite.SetBalance(accounts[4].acc.GetAddress(), sdk.NewInt64Coin("ubtc", 10000))
	gasLimit := testdata.NewTestGasLimit()
	privs := []cryptotypes.PrivKey{accounts[0].priv, accounts[1].priv, accounts[2].priv, accounts[3].priv, accounts[4].priv}
	accNums := []uint64{0, 1, 2, 3, 4}

	testCases := []TestCase{
		{
			"frozen fee set test",
			func() ([]sdk.Msg, []cryptotypes.PrivKey, []uint64, []uint64, sdk.Coins) {
				msgs := []sdk.Msg{
					govtypes.NewMsgSetNetworkProperties(accounts[0].acc.GetAddress(), &govtypes.NetworkProperties{
						MinTxFee:                 2,
						MaxTxFee:                 10000,
						EnableForeignFeePayments: true,
					}),
				}
				return msgs, privs[0:1], accNums[0:1], []uint64{0}, sdk.NewCoins(sdk.NewInt64Coin("frozen", 100))
			},
			true,
			false,
			errors.New("currency you are trying to use as fee is frozen: invalid request"),
		},
		{
			"not whitelisted token",
			func() ([]sdk.Msg, []cryptotypes.PrivKey, []uint64, []uint64, sdk.Coins) {
				msgs := []sdk.Msg{
					govtypes.NewMsgSetNetworkProperties(accounts[0].acc.GetAddress(), &govtypes.NetworkProperties{
						MinTxFee:                 2,
						MaxTxFee:                 10000,
						EnableForeignFeePayments: true,
					}),
				}
				return msgs, privs[0:1], accNums[0:1], []uint64{0}, sdk.NewCoins(sdk.NewInt64Coin("nofeetoken", 100))
			},
			true,
			false,
			errors.New("currency you are trying to use was not whitelisted as fee payment: invalid request"),
		},
		{
			"foreign fee payment disable check",
			func() ([]sdk.Msg, []cryptotypes.PrivKey, []uint64, []uint64, sdk.Coins) {
				err := suite.app.CustomGovKeeper.SetNetworkProperty(suite.ctx, govtypes.EnableForeignFeePayments, govtypes.NetworkPropertyValue{Value: 0})
				suite.Require().NoError(err)
				msgs := []sdk.Msg{
					govtypes.NewMsgSetNetworkProperties(accounts[4].acc.GetAddress(), &govtypes.NetworkProperties{
						MinTxFee:                 2,
						MaxTxFee:                 10000,
						EnableForeignFeePayments: true,
					}),
				}
				return msgs, privs[0:1], accNums[0:1], []uint64{0}, sdk.NewCoins(sdk.NewInt64Coin("ubtc", 100))
			},
			true,
			false,
			errors.New("foreign fee payments is disabled by governance: invalid request"),
		},
		{
			"fee out of range for low amount",
			func() ([]sdk.Msg, []cryptotypes.PrivKey, []uint64, []uint64, sdk.Coins) {
				err := suite.app.CustomGovKeeper.SetNetworkProperty(suite.ctx, govtypes.EnableForeignFeePayments, govtypes.NetworkPropertyValue{Value: 0})
				suite.Require().NoError(err)
				msgs := []sdk.Msg{
					govtypes.NewMsgSetNetworkProperties(accounts[4].acc.GetAddress(), &govtypes.NetworkProperties{
						MinTxFee:                 2,
						MaxTxFee:                 10000,
						EnableForeignFeePayments: true,
					}),
				}
				return msgs, privs[0:1], accNums[0:1], []uint64{0}, sdk.NewCoins(sdk.NewInt64Coin("ukex", 1))
			},
			true,
			false,
			errors.New("fee 1ukex(1) is out of range [2, 10000]ukex: invalid request"),
		},
		{
			"fee out of range for big amount",
			func() ([]sdk.Msg, []cryptotypes.PrivKey, []uint64, []uint64, sdk.Coins) {
				err := suite.app.CustomGovKeeper.SetNetworkProperty(suite.ctx, govtypes.EnableForeignFeePayments, govtypes.NetworkPropertyValue{Value: 0})
				suite.Require().NoError(err)
				msgs := []sdk.Msg{
					govtypes.NewMsgSetNetworkProperties(accounts[4].acc.GetAddress(), &govtypes.NetworkProperties{
						MinTxFee:                 2,
						MaxTxFee:                 10000,
						EnableForeignFeePayments: true,
					}),
				}
				return msgs, privs[0:1], accNums[0:1], []uint64{0}, sdk.NewCoins(sdk.NewInt64Coin("ukex", 10001))
			},
			true,
			false,
			errors.New("fee 10001ukex(10001) is out of range [2, 10000]ukex: invalid request"),
		},
		{
			"fee should be bigger than max of execution and failure fee",
			func() ([]sdk.Msg, []cryptotypes.PrivKey, []uint64, []uint64, sdk.Coins) {
				err := suite.app.CustomGovKeeper.SetNetworkProperty(suite.ctx, govtypes.EnableForeignFeePayments, govtypes.NetworkPropertyValue{Value: 0})
				suite.Require().NoError(err)
				msgs := []sdk.Msg{
					govtypes.NewMsgSetNetworkProperties(accounts[4].acc.GetAddress(), &govtypes.NetworkProperties{
						MinTxFee:                 2,
						MaxTxFee:                 10000,
						EnableForeignFeePayments: true,
					}),
				}
				return msgs, privs[0:1], accNums[0:1], []uint64{0}, sdk.NewCoins(sdk.NewInt64Coin("ukex", 3))
			},
			true,
			false,
			errors.New("fee 3ukex(3) is less than max execution fee 10000ukex: invalid request"),
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.desc), func() {
			suite.txBuilder = suite.clientCtx.TxConfig.NewTxBuilder()
			msgs, privs, accNums, accSeqs, feeAmount := tc.buildTest()

			// this runs multi signature transaction with the params provided
			suite.RunTestCase(privs, msgs, feeAmount, gasLimit, accNums, accSeqs, suite.ctx.ChainID(), tc)
		})
	}
}

// Test that simulate transaction process poor network manager correctly on ante handler step
func (suite *AnteTestSuite) TestPoorNetworkManagementDecorator() {
	suite.SetupTest(false) // reset

	// set execution fee for set network properties
	suite.app.CustomGovKeeper.SetNetworkProperties(suite.ctx, &govtypes.NetworkProperties{
		MinTxFee:                 2,
		MaxTxFee:                 10000,
		EnableForeignFeePayments: true,
		EnableTokenBlacklist:     true,
		EnableTokenWhitelist:     false,
		MinValidators:            10,
		PoorNetworkMaxBankSend:   1000,
	})

	suite.app.CustomGovKeeper.SetExecutionFee(suite.ctx, &govtypes.ExecutionFee{
		Name:              types.MsgTypeSetNetworkProperties,
		TransactionType:   types.MsgTypeSetNetworkProperties,
		ExecutionFee:      10000,
		FailureFee:        1000,
		Timeout:           0,
		DefaultParameters: 0,
	})

	// Same data for every test cases
	accounts := suite.CreateTestAccounts(5)

	suite.SetBalance(accounts[0].acc.GetAddress(), sdk.NewInt64Coin("ukex", 10000))
	suite.SetBalance(accounts[0].acc.GetAddress(), sdk.NewInt64Coin("frozen", 10000))
	suite.SetBalance(accounts[0].acc.GetAddress(), sdk.NewInt64Coin("nofeetoken", 10000))
	suite.SetBalance(accounts[1].acc.GetAddress(), sdk.NewInt64Coin("ukex", 10000))
	suite.SetBalance(accounts[2].acc.GetAddress(), sdk.NewInt64Coin("ukex", 10000))
	suite.SetBalance(accounts[3].acc.GetAddress(), sdk.NewInt64Coin("ukex", 1))
	suite.SetBalance(accounts[4].acc.GetAddress(), sdk.NewInt64Coin("ukex", 10000))
	suite.SetBalance(accounts[4].acc.GetAddress(), sdk.NewInt64Coin("ubtc", 10000))
	gasLimit := testdata.NewTestGasLimit()
	privs := []cryptotypes.PrivKey{accounts[0].priv, accounts[1].priv, accounts[2].priv, accounts[3].priv, accounts[4].priv}
	accNums := []uint64{0, 1, 2, 3, 4}

	testCases := []TestCase{
		{
			"only bond denom is allowed on poor network",
			func() ([]sdk.Msg, []cryptotypes.PrivKey, []uint64, []uint64, sdk.Coins) {
				msgs := []sdk.Msg{
					bank.NewMsgSend(accounts[4].acc.GetAddress(), accounts[3].acc.GetAddress(), sdk.Coins{sdk.NewInt64Coin("ubtc", 1)}),
				}
				return msgs, privs[0:1], accNums[0:1], []uint64{0}, sdk.NewCoins(sdk.NewInt64Coin("ukex", 100))
			},
			true,
			false,
			errors.New("only bond denom is allowed on poor network"),
		},
		{
			"only restricted amount send is allowed on poor network",
			func() ([]sdk.Msg, []cryptotypes.PrivKey, []uint64, []uint64, sdk.Coins) {
				msgs := []sdk.Msg{
					bank.NewMsgSend(accounts[4].acc.GetAddress(), accounts[3].acc.GetAddress(), sdk.Coins{sdk.NewInt64Coin("ukex", 2000)}),
				}
				return msgs, privs[0:1], accNums[0:1], []uint64{0}, sdk.NewCoins(sdk.NewInt64Coin("ukex", 100))
			},
			true,
			false,
			errors.New("only restricted amount send is allowed on poor network"),
		},
		{
			"invalid transaction type on poor network",
			func() ([]sdk.Msg, []cryptotypes.PrivKey, []uint64, []uint64, sdk.Coins) {
				msgs := []sdk.Msg{
					tokenstypes.NewMsgUpsertTokenRate(accounts[4].acc.GetAddress(), "foo", sdk.NewDec(1), true),
				}
				return msgs, privs[0:1], accNums[0:1], []uint64{0}, sdk.NewCoins(sdk.NewInt64Coin("ukex", 100))
			},
			true,
			false,
			errors.New("invalid transaction type on poor network"),
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.desc), func() {
			suite.txBuilder = suite.clientCtx.TxConfig.NewTxBuilder()
			msgs, privs, accNums, accSeqs, feeAmount := tc.buildTest()

			// this runs multi signature transaction with the params provided
			suite.RunTestCase(privs, msgs, feeAmount, gasLimit, accNums, accSeqs, suite.ctx.ChainID(), tc)
		})
	}
}

// Test that simulate transaction process black/white tokens on transfer
func (suite *AnteTestSuite) TestBlackWhiteTokensCheckDecorator() {
	suite.SetupTest(false) // reset

	// set execution fee for set network properties
	suite.app.CustomGovKeeper.SetNetworkProperties(suite.ctx, &govtypes.NetworkProperties{
		MinTxFee:                 2,
		MaxTxFee:                 10000,
		EnableForeignFeePayments: true,
		EnableTokenBlacklist:     true,
		EnableTokenWhitelist:     false,
		MinValidators:            0,
		PoorNetworkMaxBankSend:   1000,
	})

	suite.app.CustomGovKeeper.SetExecutionFee(suite.ctx, &govtypes.ExecutionFee{
		Name:              types.MsgTypeSetNetworkProperties,
		TransactionType:   types.MsgTypeSetNetworkProperties,
		ExecutionFee:      10000,
		FailureFee:        1000,
		Timeout:           0,
		DefaultParameters: 0,
	})

	// Same data for every test cases
	accounts := suite.CreateTestAccounts(5)

	suite.SetBalance(accounts[0].acc.GetAddress(), sdk.NewInt64Coin("ukex", 10000))
	suite.SetBalance(accounts[0].acc.GetAddress(), sdk.NewInt64Coin("frozen", 10000))
	suite.SetBalance(accounts[0].acc.GetAddress(), sdk.NewInt64Coin("nofeetoken", 10000))
	suite.SetBalance(accounts[1].acc.GetAddress(), sdk.NewInt64Coin("ukex", 10000))
	suite.SetBalance(accounts[2].acc.GetAddress(), sdk.NewInt64Coin("ukex", 10000))
	suite.SetBalance(accounts[3].acc.GetAddress(), sdk.NewInt64Coin("ukex", 1))
	suite.SetBalance(accounts[4].acc.GetAddress(), sdk.NewInt64Coin("ukex", 10000))
	suite.SetBalance(accounts[4].acc.GetAddress(), sdk.NewInt64Coin("ubtc", 10000))
	suite.SetBalance(accounts[4].acc.GetAddress(), sdk.NewInt64Coin("frozen", 10000))
	gasLimit := testdata.NewTestGasLimit()
	privs := []cryptotypes.PrivKey{accounts[0].priv, accounts[1].priv, accounts[2].priv, accounts[3].priv, accounts[4].priv}
	accNums := []uint64{0, 1, 2, 3, 4}

	testCases := []TestCase{
		{
			"token frozen check",
			func() ([]sdk.Msg, []cryptotypes.PrivKey, []uint64, []uint64, sdk.Coins) {
				msgs := []sdk.Msg{
					bank.NewMsgSend(accounts[4].acc.GetAddress(), accounts[3].acc.GetAddress(), sdk.Coins{sdk.NewInt64Coin("frozen", 1)}),
				}
				return msgs, privs[0:1], accNums[0:1], []uint64{0}, sdk.NewCoins(sdk.NewInt64Coin("ukex", 100))
			},
			true,
			false,
			errors.New("token is frozen"),
		},
		{
			"no frozen",
			func() ([]sdk.Msg, []cryptotypes.PrivKey, []uint64, []uint64, sdk.Coins) {
				msgs := []sdk.Msg{
					bank.NewMsgSend(accounts[4].acc.GetAddress(), accounts[3].acc.GetAddress(), sdk.Coins{sdk.NewInt64Coin("ukex", 1)}),
				}
				return msgs, privs[0:1], accNums[0:1], []uint64{0}, sdk.NewCoins(sdk.NewInt64Coin("ukex", 100))
			},
			true,
			true,
			nil,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.desc), func() {
			suite.txBuilder = suite.clientCtx.TxConfig.NewTxBuilder()
			msgs, privs, accNums, accSeqs, feeAmount := tc.buildTest()

			// this runs multi signature transaction with the params provided
			suite.RunTestCase(privs, msgs, feeAmount, gasLimit, accNums, accSeqs, suite.ctx.ChainID(), tc)
		})
	}
}

// Test that simulate transaction process execution fee registration process
func (suite *AnteTestSuite) TestExecutionFeeRegistrationDecorator() {
	suite.SetupTest(false) // reset

	// set execution fee for set network properties
	suite.app.CustomGovKeeper.SetNetworkProperties(suite.ctx, &govtypes.NetworkProperties{
		MinTxFee:                 2,
		MaxTxFee:                 10000,
		EnableForeignFeePayments: true,
		EnableTokenBlacklist:     true,
		EnableTokenWhitelist:     false,
		MinValidators:            0,
		PoorNetworkMaxBankSend:   1000,
	})

	suite.app.CustomGovKeeper.SetExecutionFee(suite.ctx, &govtypes.ExecutionFee{
		Name:              types.MsgTypeSetNetworkProperties,
		TransactionType:   types.MsgTypeSetNetworkProperties,
		ExecutionFee:      10000,
		FailureFee:        1000,
		Timeout:           0,
		DefaultParameters: 0,
	})

	// Same data for every test cases
	accounts := suite.CreateTestAccounts(5)

	suite.SetBalance(accounts[0].acc.GetAddress(), sdk.NewInt64Coin("ukex", 10000))
	suite.SetBalance(accounts[0].acc.GetAddress(), sdk.NewInt64Coin("frozen", 10000))
	suite.SetBalance(accounts[0].acc.GetAddress(), sdk.NewInt64Coin("nofeetoken", 10000))
	suite.SetBalance(accounts[1].acc.GetAddress(), sdk.NewInt64Coin("ukex", 10000))
	suite.SetBalance(accounts[2].acc.GetAddress(), sdk.NewInt64Coin("ukex", 10000))
	suite.SetBalance(accounts[3].acc.GetAddress(), sdk.NewInt64Coin("ukex", 1))
	suite.SetBalance(accounts[4].acc.GetAddress(), sdk.NewInt64Coin("ukex", 10000))
	suite.SetBalance(accounts[4].acc.GetAddress(), sdk.NewInt64Coin("ubtc", 10000))
	suite.SetBalance(accounts[4].acc.GetAddress(), sdk.NewInt64Coin("frozen", 10000))
	gasLimit := testdata.NewTestGasLimit()
	privs := []cryptotypes.PrivKey{accounts[0].priv, accounts[1].priv, accounts[2].priv, accounts[3].priv, accounts[4].priv}
	accNums := []uint64{0, 1, 2, 3, 4}

	testCases := []TestCase{
		{
			"check correctly add executions",
			func() ([]sdk.Msg, []cryptotypes.PrivKey, []uint64, []uint64, sdk.Coins) {
				msgs := []sdk.Msg{
					govtypes.NewMsgSetNetworkProperties(accounts[0].acc.GetAddress(), &govtypes.NetworkProperties{
						MinTxFee:                 2,
						MaxTxFee:                 10000,
						EnableForeignFeePayments: true,
					}),
				}
				return msgs, privs[0:1], accNums[0:1], []uint64{0}, sdk.NewCoins(sdk.NewInt64Coin("ukex", 10000))
			},
			true,
			true,
			nil,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.desc), func() {
			suite.txBuilder = suite.clientCtx.TxConfig.NewTxBuilder()
			msgs, privs, accNums, accSeqs, feeAmount := tc.buildTest()

			// this runs multi signature transaction with the params provided
			suite.RunTestCase(privs, msgs, feeAmount, gasLimit, accNums, accSeqs, suite.ctx.ChainID(), tc)
			execs := suite.app.FeeProcessingKeeper.GetExecutionsStatus(suite.ctx)
			suite.Require().Len(execs, 1)
		})
	}
}

// Test that simulate transaction set gas limit correctly on ante handler step
func (suite *AnteTestSuite) TestInfiniteGasMeterDecorator() {
	suite.SetupTest(true) // setup
	suite.txBuilder = suite.clientCtx.TxConfig.NewTxBuilder()

	// keys and addresses
	priv1, _, addr1 := testdata.KeyTestPubAddr()

	// msg and signatures
	msg := testdata.NewTestMsg(addr1)
	feeAmount := testdata.NewTestFeeAmount()
	gasLimit := testdata.NewTestGasLimit()
	suite.Require().NoError(suite.txBuilder.SetMsgs(msg))
	suite.txBuilder.SetFeeAmount(feeAmount)
	suite.txBuilder.SetGasLimit(gasLimit)

	privs, accNums, accSeqs := []cryptotypes.PrivKey{priv1}, []uint64{0}, []uint64{0}
	tx, err := suite.CreateTestTx(privs, accNums, accSeqs, suite.ctx.ChainID())
	suite.Require().NoError(err)

	sud := customante.NewZeroGasMeterDecorator()
	antehandler := sdk.ChainAnteDecorators(sud)

	// Set height to non-zero value for GasMeter to be set
	suite.ctx = suite.ctx.WithBlockHeight(1)

	// Context GasMeter Limit not set
	suite.Require().Equal(uint64(0), suite.ctx.GasMeter().Limit(), "GasMeter set with limit before setup")

	newCtx, err := antehandler(suite.ctx, tx, false)
	suite.Require().Nil(err, "InfiniteGasMeterDecorator returned error")

	// Context GasMeter Limit should be set after InfiniteGasMeterDecorator runs
	suite.Require().Equal(uint64(0x0), newCtx.GasMeter().Limit(), "GasMeter not set correctly")

	sud = customante.NewZeroGasMeterDecorator()
	antehandler = sdk.ChainAnteDecorators(sud, OutOfGasDecorator{})

	// Set height to non-zero value for GasMeter to be set
	suite.ctx = suite.ctx.WithBlockHeight(1)
	newCtx, err = antehandler(suite.ctx, tx, false)
	suite.Require().Nil(err, "no error for gas overflow")

	antehandler = sdk.ChainAnteDecorators(sud, PanicDecorator{})
	suite.Require().Panics(func() { antehandler(suite.ctx, tx, false) }, "Recovered from non-Out-of-Gas panic") // nolint:errcheck
}

type OutOfGasDecorator struct{}

// AnteDecorator that will throw OutOfGas panic
func (ogd OutOfGasDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	overLimit := ctx.GasMeter().Limit() + 1

	// Should panic with outofgas error
	ctx.GasMeter().ConsumeGas(overLimit, "test panic")

	// not reached
	return next(ctx, tx, simulate)
}

type PanicDecorator struct{}

func (pd PanicDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	panic("random error")
}
