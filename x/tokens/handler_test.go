package tokens_test

import (
	"testing"

	"github.com/KiraCore/sekai/simapp"
	"github.com/KiraCore/sekai/x/gov/types"
	customgovtypes "github.com/KiraCore/sekai/x/gov/types"
	tokens "github.com/KiraCore/sekai/x/tokens"
	tokenstypes "github.com/KiraCore/sekai/x/tokens/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

func setPermissionToAddr(t *testing.T, app *simapp.SimApp, ctx sdk.Context, addr sdk.AccAddress, perm types.PermValue) error {
	proposerActor := customgovtypes.NewDefaultActor(addr)
	err := proposerActor.Permissions.AddToWhitelist(perm)
	require.NoError(t, err)

	app.CustomGovKeeper.SaveNetworkActor(ctx, proposerActor)

	return nil
}

func TestNewHandler_MsgUpsertTokenAlias(t *testing.T) {

	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})
	handler := tokens.NewHandler(app.TokensKeeper, app.CustomGovKeeper)

	tests := []struct {
		name        string
		constructor func(sdk.AccAddress) (*tokenstypes.MsgUpsertTokenAlias, error)
		handlerErr  string
	}{
		{
			name: "good permission test",
			constructor: func(addr sdk.AccAddress) (*tokenstypes.MsgUpsertTokenAlias, error) {
				err := setPermissionToAddr(t, app, ctx, addr, types.PermSetPermissions)
				require.NoError(t, err)
				return tokenstypes.NewMsgUpsertTokenAlias(
					addr,
					0, 0,
					[]tokenstypes.VoteType{tokenstypes.VoteType_no, tokenstypes.VoteType_yes},
					"ukex",
					"Kira",
					"icon",
					6,
					[]string{"ukex"},
					tokenstypes.ProposalStatus_active,
				), nil
			},
		},
		{
			name: "lack permission test",
			constructor: func(addr sdk.AccAddress) (*tokenstypes.MsgUpsertTokenAlias, error) {
				return tokenstypes.NewMsgUpsertTokenAlias(
					addr,
					0, 0,
					[]tokenstypes.VoteType{tokenstypes.VoteType_no, tokenstypes.VoteType_yes},
					"ukex",
					"Kira",
					"icon",
					6,
					[]string{"ukex"},
					tokenstypes.ProposalStatus_active,
				), nil
			},
		},
	}
	for i, tt := range tests {
		addr := NewAccountByIndex(i)
		theMsg, err := tt.constructor(addr)
		require.NoError(t, err)

		_, err = handler(ctx, theMsg)
		if len(tt.handlerErr) != 0 {
			require.Error(t, err)
			require.Contains(t, err.Error(), tt.handlerErr)
		} else {
			require.NoError(t, err)
		}
	}
}
