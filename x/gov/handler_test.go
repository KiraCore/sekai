package gov_test

import (
	"testing"

	"github.com/KiraCore/sekai/x/gov"

	"github.com/KiraCore/sekai/app"

	types2 "github.com/cosmos/cosmos-sdk/types"

	"github.com/stretchr/testify/require"

	"github.com/KiraCore/sekai/x/gov/types"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/KiraCore/sekai/simapp"
)

func TestMain(m *testing.M) {
	app.SetConfig()
	m.Run()
}

// When a network actor has not been saved before, it creates one with default params
// and sets the permissions.
func TestNewHandler_SetPermissions_ActorWithoutPerms(t *testing.T) {
	addr, err := types2.AccAddressFromBech32("kira15ky9du8a2wlstz6fpx3p4mqpjyrm5cgqzp4f3d")
	require.NoError(t, err)

	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	handler := gov.NewHandler(app.CustomGovKeeper)

	_, err = handler(ctx, &types.MsgWhitelistPermissions{
		Address: addr,
		Permissions: []uint32{
			uint32(types.PermClaimValidator),
		},
	})
	require.NoError(t, err)

	actor, err := app.CustomGovKeeper.GetNetworkActorByAddress(ctx, addr)
	require.NoError(t, err)

	require.True(t, actor.Permissions.IsWhitelisted(types.PermClaimValidator))
}
