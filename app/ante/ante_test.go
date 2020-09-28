package ante_test

import (
	"fmt"

	"github.com/tendermint/tendermint/crypto"

	customgovtypes "github.com/KiraCore/sekai/x/gov/types"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Test that simulate transaction process execution fee correctly on ante handler step
func (suite *AnteTestSuite) TestCustomAnteHandlerExecutionFee() {
	suite.SetupTest(true) // reset

	// set execution fee for set network properties
	suite.app.CustomGovKeeper.SetExecutionFee(suite.ctx, &customgovtypes.ExecutionFee{
		Name:              customgovtypes.SetNetworkProperties,
		TransactionType:   "B",
		ExecutionFee:      10000,
		FailureFee:        1000,
		Timeout:           0,
		DefaultParameters: 0,
	})

	// Same data for every test cases
	accounts := suite.CreateTestAccounts(3)
	feeAmount := testdata.NewTestFeeAmount()
	gasLimit := testdata.NewTestGasLimit()
	accSeqs := []uint64{0, 0, 0}
	privs := []crypto.PrivKey{accounts[0].priv, accounts[1].priv, accounts[2].priv}
	accNums := []uint64{0, 1, 2}

	testCases := []TestCase{
		{
			"failure fee reduction",
			func() []sdk.Msg {
				// TODO should this buildTest logic correctly based on README for execution fee
				suite.txBuilder.SetFeeAmount(feeAmount)
				suite.txBuilder.SetGasLimit(gasLimit)

				networkActor := customgovtypes.NewNetworkActor(
					accounts[0].acc.GetAddress(),
					nil,
					1,
					nil,
					customgovtypes.NewPermissions(nil, nil),
					1,
				)
				suite.app.CustomGovKeeper.SaveNetworkActor(suite.ctx, networkActor)
				msgs := []sdk.Msg{
					customgovtypes.NewMsgSetNetworkProperties(accounts[0].acc.GetAddress(), &customgovtypes.NetworkProperties{
						MinTxFee: 2,
						MaxTxFee: 10000,
					}),
				}
				return msgs
			},
			true,
			true,
			nil,
		},
		{
			"no deduction when does not exist",
			func() []sdk.Msg {
				// TODO should this buildTest logic correctly based on README for execution fee
				simulatedGas := suite.ctx.GasMeter().GasConsumed()

				accSeqs = []uint64{1, 1, 1}
				suite.txBuilder.SetFeeAmount(feeAmount)
				suite.txBuilder.SetGasLimit(simulatedGas)

				networkActor := customgovtypes.NewNetworkActor(
					accounts[0].acc.GetAddress(),
					nil,
					1,
					nil,
					customgovtypes.NewPermissions(nil, nil),
					1,
				)
				suite.app.CustomGovKeeper.SaveNetworkActor(suite.ctx, networkActor)

				msgs := []sdk.Msg{
					customgovtypes.NewMsgSetNetworkProperties(accounts[0].acc.GetAddress(), &customgovtypes.NetworkProperties{
						MinTxFee: 2,
						MaxTxFee: 10000,
					}),
				}
				return msgs
			},
			false,
			true,
			nil,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.desc), func() {
			suite.txBuilder = suite.clientCtx.TxConfig.NewTxBuilder()
			msgs := tc.buildTest()

			// this runs multi signature transaction with the params provided
			suite.RunTestCase(privs, msgs, feeAmount, gasLimit, accNums, accSeqs, suite.ctx.ChainID(), tc)
			// TODO should check balance change after a transaction (two cases describe below)
		})
	}
}
