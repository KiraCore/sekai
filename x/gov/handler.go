package gov

import (
	types2 "github.com/KiraCore/sekai/x/gov/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/KiraCore/sekai/x/gov/keeper"
	"github.com/KiraCore/sekai/x/staking/types"
)

func NewHandler(ck keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case *types2.MsgWhitelistPermissions:
			return handleWhitelistPermissions(ctx, ck, msg)
		default:
			return nil, errors.Wrapf(errors.ErrUnknownRequest, "unrecognized %s message type: %T", types.ModuleName, msg)
		}
	}
}

func handleWhitelistPermissions(ctx sdk.Context, ck keeper.Keeper, msg *types2.MsgWhitelistPermissions) (*sdk.Result, error) {
	actor, err := ck.GetNetworkActorByAddress(ctx, msg.Address)
	if err != nil {
		actor = types2.NewDefaultActor(msg.Address)
	}

	for _, perm := range msg.Permissions {
		err := actor.Permissions.AddToWhitelist(types2.PermValue(perm))
		if err != nil {
			return nil, errors.Wrapf(types2.ErrSetPermissions, "error setting %d to whitelist", perm)
		}
	}

	ck.SaveNetworkActor(ctx, actor)

	return &sdk.Result{}, nil
}
