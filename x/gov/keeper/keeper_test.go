package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/KiraCore/sekai/simapp"
	"github.com/KiraCore/sekai/x/gov/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

func TestKeeper_SaveGetPermissionsForRole(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	perm := types.NewPermissions(
		nil, []types.PermValue{types.PermClaimValidator},
	)

	app.CustomGovKeeper.SetPermissionsForRole(ctx, types.RoleCouncilor, perm)

	savedPerms := app.CustomGovKeeper.GetPermissionsForRole(ctx, types.RoleCouncilor)
	require.Equal(t, perm, savedPerms)
}

func TestKeeper_RoleHasPermissionsFor(t *testing.T) {

}
