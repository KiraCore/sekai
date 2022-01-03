package slashing_test

import (
	"os"
	"strings"
	"testing"

	"github.com/KiraCore/sekai/app"
	simapp "github.com/KiraCore/sekai/app"
	"github.com/KiraCore/sekai/x/gov"
	govtypes "github.com/KiraCore/sekai/x/gov/types"
	"github.com/KiraCore/sekai/x/slashing"
	"github.com/KiraCore/sekai/x/slashing/keeper"
	"github.com/KiraCore/sekai/x/slashing/types"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

func TestMain(m *testing.M) {
	app.SetConfig()
	os.Exit(m.Run())
}

func TestInvalidMsg(t *testing.T) {
	k := keeper.Keeper{}
	h := slashing.NewHandler(k)

	res, err := h(sdk.NewContext(nil, tmproto.Header{}, false, nil), testdata.NewTestMsg())
	require.Error(t, err)
	require.Nil(t, res)
	require.True(t, strings.Contains(err.Error(), "unrecognized customslashing message type"))
}

func TestHandler_CreateProposalResetWholeValidatorRank(t *testing.T) {
	proposerAddr, err := sdk.AccAddressFromBech32("kira1alzyfq40zjsveat87jlg8jxetwqmr0a29sgd0f")
	require.NoError(t, err)

	tests := []struct {
		name         string
		content      govtypes.Content
		preparePerms func(t *testing.T, app *simapp.SekaiApp, ctx sdk.Context)
		expectedErr  error
	}{
		{
			"Proposer does not have Perm",
			types.NewResetWholeValidatorRankProposal(
				proposerAddr,
			),
			func(t *testing.T, app *simapp.SekaiApp, ctx sdk.Context) {},
			errors.Wrap(govtypes.ErrNotEnoughPermissions, "PERMISSION_CREATE_RESET_WHOLE_VALIDATOR_RANK_PROPOSAL"),
		},
		{
			"Proposer has permission",
			types.NewResetWholeValidatorRankProposal(
				proposerAddr,
			),
			func(t *testing.T, app *simapp.SekaiApp, ctx sdk.Context) {
				proposerActor := govtypes.NewDefaultActor(proposerAddr)
				err2 := app.CustomGovKeeper.AddWhitelistPermission(ctx, proposerActor, govtypes.PermCreateResetWholeValidatorRankProposal)
				require.NoError(t, err2)
			},
			nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			app := simapp.Setup(false)
			ctx := app.NewContext(false, tmproto.Header{})

			tt.preparePerms(t, app, ctx)

			handler := gov.NewHandler(app.CustomGovKeeper)
			msg, err := govtypes.NewMsgSubmitProposal(proposerAddr, "title", "some desc", tt.content)
			require.NoError(t, err)
			_, err = handler(ctx, msg)
			if tt.expectedErr == nil {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tt.expectedErr.Error())
			}
		})
	}
}
