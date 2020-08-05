package staking

import (
	sdk "github.com/KiraCore/cosmos-sdk/types"
	"github.com/KiraCore/cosmos-sdk/x/staking"
	stakingkeeper "github.com/KiraCore/cosmos-sdk/x/staking/keeper"

	customkeeper "github.com/KiraCore/sekai/x/staking/keeper"
	"github.com/KiraCore/sekai/x/staking/types"
)

func NewHandler(k stakingkeeper.Keeper, ck customkeeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case *types.MsgClaimValidator:
			return handleMsgClaimValidator(ctx, ck, msg)
		default:
			return staking.NewHandler(k)(ctx, msg)
		}
	}
}

func handleMsgClaimValidator(ctx sdk.Context, k customkeeper.Keeper, msg *types.MsgClaimValidator) (*sdk.Result, error) {
	validator, err := types.NewValidator(msg.Moniker, msg.Website, msg.Social, msg.Identity, msg.Comission, msg.ValKey, msg.PubKey)
	if err != nil {
		return nil, err
	}

	k.AddValidator(ctx, validator)

	return &sdk.Result{}, nil
}
