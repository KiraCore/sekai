package middleware_test

import (
	"bytes"
	"os"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/KiraCore/sekai/app"
	"github.com/KiraCore/sekai/middleware"
	"github.com/KiraCore/sekai/simapp"
	"github.com/KiraCore/sekai/types"
	"github.com/KiraCore/sekai/x/gov"
	customgovtypes "github.com/KiraCore/sekai/x/gov/types"
)

func TestMain(m *testing.M) {
	app.SetConfig()
	os.Exit(m.Run())
}

func TestNewHandler_SetNetworkProperties(t *testing.T) {
	changeFeeAddr, err := sdk.AccAddressFromBech32("kira15ky9du8a2wlstz6fpx3p4mqpjyrm5cgqzp4f3d")
	require.NoError(t, err)

	sudoAddr, err := sdk.AccAddressFromBech32("kira1alzyfq40zjsveat87jlg8jxetwqmr0a29sgd0f")
	require.NoError(t, err)

	tests := []struct {
		name       string
		msg        sdk.Msg
		desiredErr string
	}{
		{
			name: "Success run with ChangeTxFee permission",
			msg: &customgovtypes.MsgSetNetworkProperties{
				NetworkProperties: &customgovtypes.NetworkProperties{
					MinTxFee: 100,
					MaxTxFee: 1000,
				},
				Proposer: changeFeeAddr,
			},
			desiredErr: "",
		},
		{
			name: "Failure run without ChangeTxFee permission",
			msg: &customgovtypes.MsgSetNetworkProperties{
				NetworkProperties: &customgovtypes.NetworkProperties{
					MinTxFee: 100,
					MaxTxFee: 1000,
				},
				Proposer: sudoAddr,
			},
			desiredErr: "not enough permissions",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			app := simapp.Setup(false)
			ctx := app.NewContext(false, tmproto.Header{})

			app.BankKeeper.SetBalance(ctx, sudoAddr, sdk.NewInt64Coin("ukex", 100000))
			app.BankKeeper.SetBalance(ctx, changeFeeAddr, sdk.NewInt64Coin("ukex", 100000))

			// First we set Role Sudo to proposer Actor
			proposerActor := customgovtypes.NewDefaultActor(sudoAddr)
			proposerActor.SetRole(customgovtypes.RoleSudo)
			require.NoError(t, err)
			app.CustomGovKeeper.SaveNetworkActor(ctx, proposerActor)

			handler := gov.NewHandler(app.CustomGovKeeper)

			// set change fee permission to addr
			_, err = handler(ctx, &customgovtypes.MsgWhitelistPermissions{
				Proposer:   sudoAddr,
				Address:    changeFeeAddr,
				Permission: uint32(customgovtypes.PermChangeTxFee),
			})
			require.NoError(t, err)

			// set execution fee
			_, err = handler(ctx, &customgovtypes.MsgSetExecutionFee{
				Proposer:          changeFeeAddr,
				Name:              types.MsgTypeSetNetworkProperties,
				TransactionType:   types.MsgTypeSetNetworkProperties,
				ExecutionFee:      10000,
				FailureFee:        1000,
				Timeout:           1,
				DefaultParameters: 2,
			})
			require.NoError(t, err)

			app.FeeProcessingKeeper.AddExecutionStart(ctx, tt.msg)

			// test message with new middleware handler
			newHandler := middleware.NewRoute(customgovtypes.ModuleName, gov.NewHandler(app.CustomGovKeeper)).Handler()
			_, err = newHandler(ctx, tt.msg)

			if tt.desiredErr == "" {
				require.NoError(t, err)

				// check success flag change after successful run
				executions := app.FeeProcessingKeeper.GetExecutionsStatus(ctx)
				successExist := false
				for _, exec := range executions {
					if exec.Success == true && exec.MsgType == tt.msg.Type() && bytes.Equal(exec.FeePayer, tt.msg.GetSigners()[0]) {
						successExist = true
						break
					}
				}
				require.True(t, successExist)
			} else {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.desiredErr)
			}
		})
	}
}
