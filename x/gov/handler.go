package gov

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/KiraCore/sekai/x/gov/keeper"
	customgovtypes "github.com/KiraCore/sekai/x/gov/types"
	"github.com/KiraCore/sekai/x/staking/types"
)

func NewHandler(ck keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case *customgovtypes.MsgWhitelistPermissions:
			return handleWhitelistPermissions(ctx, ck, msg)
		case *customgovtypes.MsgBlacklistPermissions:
			return handleBlacklistPermissions(ctx, ck, msg)
		case *customgovtypes.MsgClaimCouncilor:
			return handleClaimCouncilor(ctx, ck, msg)
		case *customgovtypes.MsgWhitelistRolePermission:
			return handleWhitelistRolePermission(ctx, ck, msg)
		case *customgovtypes.MsgBlacklistRolePermission:
			return handleBlacklistRolePermission(ctx, ck, msg)
		case *customgovtypes.MsgRemoveWhitelistRolePermission:
			return handleRemoveWhitelistRolePermission(ctx, ck, msg)
		default:
			return nil, errors.Wrapf(errors.ErrUnknownRequest, "unrecognized %s message type: %T", types.ModuleName, msg)
		}
	}
}

func handleRemoveWhitelistRolePermission(ctx sdk.Context, ck keeper.Keeper, msg *customgovtypes.MsgRemoveWhitelistRolePermission) (*sdk.Result, error) {
	perms, err := validateAndGetPermissionsForRole(ctx, ck, msg.Proposer, customgovtypes.Role(msg.Role))
	if err != nil {
		return nil, err
	}

	err = perms.RemoveFromWhitelist(customgovtypes.PermValue(msg.Permission))
	if err != nil {
		return nil, errors.Wrap(customgovtypes.ErrWhitelisting, err.Error())
	}

	ck.SetPermissionsForRole(ctx, customgovtypes.Role(msg.Role), perms)

	return &sdk.Result{}, nil
}

func handleBlacklistRolePermission(ctx sdk.Context, ck keeper.Keeper, msg *customgovtypes.MsgBlacklistRolePermission) (*sdk.Result, error) {
	perms, err := validateAndGetPermissionsForRole(ctx, ck, msg.Proposer, customgovtypes.Role(msg.Role))
	if err != nil {
		return nil, err
	}

	err = perms.AddToBlacklist(customgovtypes.PermValue(msg.Permission))
	if err != nil {
		return nil, errors.Wrap(customgovtypes.ErrBlacklisting, err.Error())
	}

	ck.SetPermissionsForRole(ctx, customgovtypes.Role(msg.Role), perms)

	return &sdk.Result{}, nil
}

func handleWhitelistRolePermission(ctx sdk.Context, ck keeper.Keeper, msg *customgovtypes.MsgWhitelistRolePermission) (*sdk.Result, error) {
	perms, err := validateAndGetPermissionsForRole(ctx, ck, msg.Proposer, customgovtypes.Role(msg.Role))
	if err != nil {
		return nil, err
	}

	err = perms.AddToWhitelist(customgovtypes.PermValue(msg.Permission))
	if err != nil {
		return nil, errors.Wrap(customgovtypes.ErrWhitelisting, err.Error())
	}

	ck.SetPermissionsForRole(ctx, customgovtypes.Role(msg.Role), perms)

	return &sdk.Result{}, nil
}

func handleWhitelistPermissions(ctx sdk.Context, ck keeper.Keeper, msg *customgovtypes.MsgWhitelistPermissions) (*sdk.Result, error) {
	isAllowed := keeper.CheckIfAllowedPermission(ctx, ck, msg.Proposer, customgovtypes.PermSetPermissions)
	if !isAllowed {
		return nil, errors.Wrap(customgovtypes.ErrNotEnoughPermissions, "PermSetPermissions")
	}

	actor, err := ck.GetNetworkActorByAddress(ctx, msg.Address)
	if err != nil {
		actor = customgovtypes.NewDefaultActor(msg.Address)
	}

	err = actor.Permissions.AddToWhitelist(customgovtypes.PermValue(msg.Permission))
	if err != nil {
		return nil, errors.Wrapf(customgovtypes.ErrSetPermissions, "error setting %d to whitelist", msg.Permission)
	}

	ck.SaveNetworkActor(ctx, actor)

	return &sdk.Result{}, nil
}

func handleBlacklistPermissions(ctx sdk.Context, ck keeper.Keeper, msg *customgovtypes.MsgBlacklistPermissions) (*sdk.Result, error) {
	isAllowed := keeper.CheckIfAllowedPermission(ctx, ck, msg.Proposer, customgovtypes.PermSetPermissions)
	if !isAllowed {
		return nil, errors.Wrap(customgovtypes.ErrNotEnoughPermissions, "PermSetPermissions")
	}

	actor, err := ck.GetNetworkActorByAddress(ctx, msg.Address)
	if err != nil {
		actor = customgovtypes.NewDefaultActor(msg.Address)
	}

	err = actor.Permissions.AddToBlacklist(customgovtypes.PermValue(msg.Permission))
	if err != nil {
		return nil, errors.Wrapf(customgovtypes.ErrSetPermissions, "error setting %d to whitelist", msg.Permission)
	}

	ck.SaveNetworkActor(ctx, actor)

	return &sdk.Result{}, nil
}

func handleClaimCouncilor(ctx sdk.Context, ck keeper.Keeper, msg *customgovtypes.MsgClaimCouncilor) (*sdk.Result, error) {
	isAllowed := keeper.CheckIfAllowedPermission(ctx, ck, msg.Address, customgovtypes.PermClaimCouncilor)
	if !isAllowed {
		return nil, errors.Wrap(customgovtypes.ErrNotEnoughPermissions, "PermClaimCouncilor")
	}

	councilor := customgovtypes.NewCouncilor(msg.Moniker, msg.Website, msg.Social, msg.Identity, msg.Address)

	ck.SaveCouncilor(ctx, councilor)

	return &sdk.Result{}, nil
}

// validateAndGetPermissionsForRole checks if:
// - Proposer has permissions to SetPermissions.
// - Role exists.
// And returns the permissions.
func validateAndGetPermissionsForRole(
	ctx sdk.Context,
	ck keeper.Keeper,
	proposer sdk.AccAddress,
	role customgovtypes.Role,
) (*customgovtypes.Permissions, error) {
	isAllowed := keeper.CheckIfAllowedPermission(ctx, ck, proposer, customgovtypes.PermSetPermissions)
	if !isAllowed {
		return nil, errors.Wrap(customgovtypes.ErrNotEnoughPermissions, "PermSetPermissions")
	}

	perms, err := ck.GetPermissionsForRole(ctx, role)
	if err != nil {
		return nil, customgovtypes.ErrRoleDoesNotExist
	}

	return perms, nil
}
