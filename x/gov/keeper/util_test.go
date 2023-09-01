package keeper_test

import (
	"os"
	"testing"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	simapp "github.com/KiraCore/sekai/app"
	appparams "github.com/KiraCore/sekai/app/params"
	"github.com/KiraCore/sekai/x/gov/keeper"
	"github.com/KiraCore/sekai/x/gov/types"
)

func TestMain(m *testing.M) {
	appparams.SetConfig()
	os.Exit(m.Run())
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
				roleWithBlacklistedValue := uint64(123)
				keeper.SetRole(ctx, types.Role{
					Id:          uint32(roleWithBlacklistedValue),
					Sid:         "123",
					Description: "123",
				})
				err2 := keeper.BlacklistRolePermission(ctx, roleWithBlacklistedValue, types.PermClaimValidator)
				require.NoError(t, err2)

				actor := types.NewDefaultActor(addr)
				keeper.AssignRoleToActor(ctx, actor, roleWithBlacklistedValue)
			},
			isAllowed: false,
		},
		{
			name: "actor has permission whitelisted in role",
			prepareScenario: func(ctx sdk.Context, keeper keeper.Keeper) {
				roleWithWhitelistedValue := uint64(123)
				keeper.SetRole(ctx, types.Role{
					Id:          uint32(roleWithWhitelistedValue),
					Sid:         "123",
					Description: "123",
				})

				err2 := keeper.WhitelistRolePermission(ctx, roleWithWhitelistedValue, types.PermClaimValidator)
				require.NoError(t, err2)

				actor := types.NewDefaultActor(addr)
				keeper.AssignRoleToActor(ctx, actor, roleWithWhitelistedValue)
			},
			isAllowed: true,
		},
		{
			name: "actor has permission whitelisted individually",
			prepareScenario: func(ctx sdk.Context, keeper keeper.Keeper) {
				actor := types.NewDefaultActor(addr)
				require.NoError(t, keeper.AddWhitelistPermission(ctx, actor, types.PermClaimValidator))
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
