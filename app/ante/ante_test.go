package ante_test

import (
	"errors"
	"fmt"

	"github.com/tendermint/tendermint/crypto"

	"github.com/KiraCore/sekai/types"
	customgovtypes "github.com/KiraCore/sekai/x/gov/types"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Test that simulate transaction process execution fee correctly on ante handler step
func (suite *AnteTestSuite) TestCustomAnteHandlerExecutionFee() {
	suite.SetupTest(false) // reset

	// set execution fee for set network properties
	suite.app.CustomGovKeeper.SetExecutionFee(suite.ctx, &customgovtypes.ExecutionFee{
		Name:              types.MsgTypeSetNetworkProperties,
		TransactionType:   types.MsgTypeSetNetworkProperties,
		ExecutionFee:      10000,
		FailureFee:        1000,
		Timeout:           0,
		DefaultParameters: 0,
	})
	suite.app.CustomGovKeeper.SetNetworkProperties(suite.ctx, &customgovtypes.NetworkProperties{
		MinTxFee:                 2,
		MaxTxFee:                 10000,
		EnableForeignFeePayments: true,
	})

	// Same data for every test cases
	accounts := suite.CreateTestAccounts(5)

	suite.app.BankKeeper.SetBalance(suite.ctx, accounts[0].acc.GetAddress(), sdk.NewInt64Coin("ukex", 10000))
	suite.app.BankKeeper.SetBalance(suite.ctx, accounts[1].acc.GetAddress(), sdk.NewInt64Coin("ukex", 10000))
	suite.app.BankKeeper.SetBalance(suite.ctx, accounts[2].acc.GetAddress(), sdk.NewInt64Coin("ukex", 10000))
	suite.app.BankKeeper.SetBalance(suite.ctx, accounts[3].acc.GetAddress(), sdk.NewInt64Coin("ukex", 1))
	suite.app.BankKeeper.SetBalance(suite.ctx, accounts[4].acc.GetAddress(), sdk.NewInt64Coin("ukex", 10000))
	suite.app.BankKeeper.SetBalance(suite.ctx, accounts[4].acc.GetAddress(), sdk.NewInt64Coin("ubtc", 10000))
	defaultFee := sdk.NewCoins(sdk.NewInt64Coin("ukex", 100))
	gasLimit := testdata.NewTestGasLimit()
	privs := []crypto.PrivKey{accounts[0].priv, accounts[1].priv, accounts[2].priv, accounts[3].priv, accounts[4].priv}
	accNums := []uint64{0, 1, 2, 3, 4}

	testCases := []TestCase{
		{
			"insufficient max execution fee set",
			func() ([]sdk.Msg, []crypto.PrivKey, []uint64, []uint64, sdk.Coins) {
				msgs := []sdk.Msg{
					customgovtypes.NewMsgSetNetworkProperties(accounts[0].acc.GetAddress(), &customgovtypes.NetworkProperties{
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
			func() ([]sdk.Msg, []crypto.PrivKey, []uint64, []uint64, sdk.Coins) {
				msgs := []sdk.Msg{
					customgovtypes.NewMsgSetNetworkProperties(accounts[1].acc.GetAddress(), &customgovtypes.NetworkProperties{
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
			func() ([]sdk.Msg, []crypto.PrivKey, []uint64, []uint64, sdk.Coins) {
				msgs := []sdk.Msg{
					customgovtypes.NewMsgSetExecutionFee(
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
			func() ([]sdk.Msg, []crypto.PrivKey, []uint64, []uint64, sdk.Coins) {
				msgs := []sdk.Msg{
					customgovtypes.NewMsgSetExecutionFee(
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
			func() ([]sdk.Msg, []crypto.PrivKey, []uint64, []uint64, sdk.Coins) {
				msgs := []sdk.Msg{
					customgovtypes.NewMsgSetExecutionFee(
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
			func() ([]sdk.Msg, []crypto.PrivKey, []uint64, []uint64, sdk.Coins) {
				suite.app.CustomGovKeeper.SetNetworkProperties(suite.ctx, &customgovtypes.NetworkProperties{
					MinTxFee:                 2,
					MaxTxFee:                 10000,
					EnableForeignFeePayments: true,
				})
				msgs := []sdk.Msg{
					customgovtypes.NewMsgSetExecutionFee(
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
			func() ([]sdk.Msg, []crypto.PrivKey, []uint64, []uint64, sdk.Coins) {

				suite.app.CustomGovKeeper.SetNetworkProperties(suite.ctx, &customgovtypes.NetworkProperties{
					MinTxFee:                 2,
					MaxTxFee:                 10000,
					EnableForeignFeePayments: false,
				})
				msgs := []sdk.Msg{
					customgovtypes.NewMsgSetExecutionFee(
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
