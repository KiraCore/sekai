package staking

import (
	customgovtypes "github.com/KiraCore/sekai/x/gov/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/tendermint/tendermint/crypto"

	govkeeper "github.com/KiraCore/sekai/x/gov/keeper"
	customkeeper "github.com/KiraCore/sekai/x/staking/keeper"
	"github.com/KiraCore/sekai/x/staking/types"
)

func NewHandler(ck customkeeper.Keeper, govkeeper govkeeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case *types.MsgClaimValidator:
			return handleMsgClaimValidator(ctx, ck, govkeeper, msg)
		default:
			return nil, errors.Wrapf(errors.ErrUnknownRequest, "unrecognized %s message type: %T", types.ModuleName, msg)
		}
	}
}

func handleMsgClaimValidator(ctx sdk.Context, k customkeeper.Keeper, gk govkeeper.Keeper, msg *types.MsgClaimValidator) (*sdk.Result, error) {
	isAllowed := govkeeper.CheckIfAllowedPermission(ctx, gk, sdk.AccAddress(msg.ValKey), customgovtypes.PermClaimValidator)
	if !isAllowed {
		return nil, errors.Wrap(customgovtypes.ErrNotEnoughPermissions, "PermClaimValidator")
	}

	pk, ok := msg.PubKey.GetCachedValue().(crypto.PubKey)
	if !ok {
		return nil, errors.Wrap(errors.ErrInvalidPubKey, "Invalid pub key")
	}

	validator, err := types.NewValidator(msg.Moniker, msg.Website, msg.Social, msg.Identity, msg.Commission, msg.ValKey, pk)
	if err != nil {
		return nil, err
	}

	k.AddValidator(ctx, validator)

	return &sdk.Result{}, nil
}
