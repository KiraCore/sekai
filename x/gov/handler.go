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
		default:
			return nil, errors.Wrapf(errors.ErrUnknownRequest, "unrecognized %s message type: %T", types.ModuleName, msg)
		}
	}
}

func handleWhitelistPermissions(ctx sdk.Context, ck keeper.Keeper, msg *customgovtypes.MsgWhitelistPermissions) (*sdk.Result, error) {
	// Check if proposer have permissions to SetPermissions.
	proposer, err := ck.GetNetworkActorByAddress(ctx, msg.Proposer)
	if err != nil {
		return nil, errors.Wrap(customgovtypes.ErrNotEnoughPermissions, "SetPermissions")
	}

	if !proposer.HasRole(customgovtypes.RoleSudo) {
		if proposer.Permissions.IsBlacklisted(customgovtypes.PermSetPermissions) || !proposer.Permissions.IsWhitelisted(customgovtypes.PermSetPermissions) {
			return nil, errors.Wrap(customgovtypes.ErrNotEnoughPermissions, "SetPermissions")
		}
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
