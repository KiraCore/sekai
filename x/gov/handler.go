package gov

import (
	customgovtypes "github.com/KiraCore/sekai/x/gov/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/KiraCore/sekai/x/gov/keeper"
	"github.com/KiraCore/sekai/x/staking/types"
)

func NewHandler(ck keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case *customgovtypes.MsgWhitelistPermissions:
			return handleWhitelistPermissions(ctx, ck, msg)
		case *customgovtypes.MsgSetNetworkProperties:
			return handleSetNetworkProperties(ctx, ck, msg)
		case *customgovtypes.MsgBlacklistPermissions:
			return handleBlacklistPermissions(ctx, ck, msg)
		default:
			return nil, errors.Wrapf(errors.ErrUnknownRequest, "unrecognized %s message type: %T", types.ModuleName, msg)
		}
	}
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

func handleSetNetworkProperties(ctx sdk.Context, ck keeper.Keeper, msg *customgovtypes.MsgSetNetworkProperties) (*sdk.Result, error) {
	isAllowed := keeper.CheckIfAllowedPermission(ctx, ck, msg.Proposer, customgovtypes.PermChangeTxFee)
	if !isAllowed {
		return nil, errors.Wrap(customgovtypes.ErrNotEnoughPermissions, "PermChangeTxFee")
	}
	ck.SetNetworkProperties(ctx, msg.NetworkProperties)
	return &sdk.Result{}, nil
}
