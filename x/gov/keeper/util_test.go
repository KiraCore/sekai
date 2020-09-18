package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/KiraCore/sekai/app"
	"github.com/KiraCore/sekai/simapp"
	"github.com/KiraCore/sekai/x/gov/keeper"
	"github.com/KiraCore/sekai/x/gov/types"
)

func TestMain(m *testing.M) {
	app.SetConfig()
	m.Run()
}

func TestCheckIfAllowedPermission(t *testing.T) {
	addr, err := sdk.AccAddressFromBech32("kira15ky9du8a2wlstz6fpx3p4mqpjyrm5cgqzp4f3d")
	require.NoError(t, err)

	tests := []struct {
		name            string
		prepareScenario func(ctx sdk.Context, keeper keeper.Keeper)
		isAllowed       bool
	}{
		{
			name: "actor not present",
			prepareScenario: func(ctx sdk.Context, keeper keeper.Keeper) {
				return
			},
			isAllowed: false,
		},
		{
			name: "actor has permission individually blacklisted",
			prepareScenario: func(ctx sdk.Context, keeper keeper.Keeper) {
				actor := types.NewDefaultActor(addr)
				require.NoError(t, actor.Permissions.AddToBlacklist(types.PermClaimValidator))
				keeper.SaveNetworkActor(ctx, actor)
			},
			isAllowed: false,
		},
		{
			name: "actor has permission blacklisted in role",
			prepareScenario: func(ctx sdk.Context, keeper keeper.Keeper) {
				roleWithBlacklistedValue := types.Role(123)
				keeper.SetPermissionsForRole(ctx, roleWithBlacklistedValue, types.NewPermissions(nil, []types.PermValue{
					types.PermClaimValidator,
				}))

				actor := types.NewDefaultActor(addr)
				actor.SetRole(roleWithBlacklistedValue)

				keeper.SaveNetworkActor(ctx, actor)
			},
			isAllowed: false,
		},
		{
			name: "actor has permission whitelisted in role",
			prepareScenario: func(ctx sdk.Context, keeper keeper.Keeper) {
				roleWithWhitelistedValue := types.Role(123)
				keeper.SetPermissionsForRole(ctx, roleWithWhitelistedValue, types.NewPermissions([]types.PermValue{
					types.PermClaimValidator,
				}, nil))

				actor := types.NewDefaultActor(addr)
				actor.SetRole(roleWithWhitelistedValue)

				keeper.SaveNetworkActor(ctx, actor)
			},
			isAllowed: true,
		},
		{
			name: "actor has permission whitelisted individually",
			prepareScenario: func(ctx sdk.Context, keeper keeper.Keeper) {
				actor := types.NewDefaultActor(addr)
				require.NoError(t, actor.Permissions.AddToWhitelist(types.PermClaimValidator))

				keeper.SaveNetworkActor(ctx, actor)
			},
			isAllowed: true,
		},
		{
			name: "actor not whitelisted or blacklisted",
			prepareScenario: func(ctx sdk.Context, keeper keeper.Keeper) {
				actor := types.NewDefaultActor(addr)
				keeper.SaveNetworkActor(ctx, actor)
			},
			isAllowed: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			app := simapp.Setup(false)
			ctx := app.NewContext(false, tmproto.Header{})

			tt.prepareScenario(ctx, app.CustomGovKeeper)

			allowed := keeper.CheckIfAllowedPermission(ctx, app.CustomGovKeeper, addr, types.PermClaimValidator)
			require.Equal(t, tt.isAllowed, allowed)
		})
	}
}
