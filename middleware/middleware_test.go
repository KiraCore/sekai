package middleware_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/KiraCore/sekai/app"
	simapp "github.com/KiraCore/sekai/app"
	"github.com/KiraCore/sekai/middleware"
	"github.com/KiraCore/sekai/types"
	kiratypes "github.com/KiraCore/sekai/types"
	"github.com/KiraCore/sekai/x/gov"
	govtypes "github.com/KiraCore/sekai/x/gov/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

func TestMain(m *testing.M) {
	app.SetConfig()
	os.Exit(m.Run())
}

func Test_Middleware_SetNetworkProperties(t *testing.T) {
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
			msg: &govtypes.MsgSetNetworkProperties{
				NetworkProperties: &govtypes.NetworkProperties{
					MinTxFee: 100,
					MaxTxFee: 1000,
				},
				Proposer: changeFeeAddr,
			},
			desiredErr: "",
		},
		{
			name: "Failure run without ChangeTxFee permission",
			msg: &govtypes.MsgSetNetworkProperties{
				NetworkProperties: &govtypes.NetworkProperties{
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

			coins := sdk.Coins{sdk.NewInt64Coin("ukex", 100000)}
			app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, coins)
			app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, sudoAddr, coins)
			app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, coins)
			app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, changeFeeAddr, coins)

			// First we set Role Sudo to proposer Actor
			proposerActor := govtypes.NewDefaultActor(sudoAddr)
			proposerActor.SetRole(govtypes.RoleSudo)
			require.NoError(t, err)
			app.CustomGovKeeper.SaveNetworkActor(ctx, proposerActor)

			handler := gov.NewHandler(app.CustomGovKeeper)

			// set change fee permission to addr
			_, err = handler(ctx, &govtypes.MsgWhitelistPermissions{
				Proposer:   sudoAddr,
				Address:    changeFeeAddr,
				Permission: uint32(govtypes.PermChangeTxFee),
			})
			require.NoError(t, err)

			// set execution fee
			_, err = handler(ctx, &govtypes.MsgSetExecutionFee{
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
			newHandler := middleware.NewRoute(govtypes.ModuleName, gov.NewHandler(app.CustomGovKeeper)).Handler()
			_, err = newHandler(ctx, tt.msg)

			if tt.desiredErr == "" {
				require.NoError(t, err)

				// check success flag change after successful run
				executions := app.FeeProcessingKeeper.GetExecutionsStatus(ctx)
				successExist := false
				for _, exec := range executions {
					if exec.Success == true && exec.MsgType == kiratypes.MsgType(tt.msg) && bytes.Equal(exec.FeePayer, tt.msg.GetSigners()[0]) {
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
