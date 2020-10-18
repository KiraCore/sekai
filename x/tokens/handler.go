package tokens

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"

	customgovtypes "github.com/KiraCore/sekai/x/gov/types"
	"github.com/KiraCore/sekai/x/tokens/keeper"
	"github.com/KiraCore/sekai/x/tokens/types"
)

func NewHandler(ck keeper.Keeper, cgk types.CustomGovKeeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case *types.MsgUpsertTokenAlias:
			return handleUpsertTokenAlias(ctx, ck, cgk, msg)
		case *types.MsgUpsertTokenRate:
			return handleUpsertTokenRate(ctx, ck, cgk, msg)
		default:
			return nil, errors.Wrapf(errors.ErrUnknownRequest, "unrecognized %s message type: %T", types.ModuleName, msg)
		}
	}
}

func handleUpsertTokenAlias(ctx sdk.Context, ck keeper.Keeper, cgk types.CustomGovKeeper, msg *types.MsgUpsertTokenAlias) (*sdk.Result, error) {
	isAllowed := cgk.CheckIfAllowedPermission(ctx, msg.Proposer, customgovtypes.PermUpsertTokenAlias)
	if !isAllowed {
		return nil, errors.Wrap(customgovtypes.ErrNotEnoughPermissions, "PermUpsertTokenAlias")
	}

	err := ck.UpsertTokenAlias(ctx, *types.NewTokenAlias(
		msg.Expiration,
		msg.Enactment,
		msg.AllowedVoteTypes,
		msg.Symbol,
		msg.Name,
		msg.Icon,
		msg.Decimals,
		msg.Denoms,
		msg.Status,
	))
	return &sdk.Result{}, err
}

func handleUpsertTokenRate(ctx sdk.Context, ck keeper.Keeper, cgk types.CustomGovKeeper, msg *types.MsgUpsertTokenRate) (*sdk.Result, error) {
	isAllowed := cgk.CheckIfAllowedPermission(ctx, msg.Proposer, customgovtypes.PermUpsertTokenRate)
	if !isAllowed {
		return nil, errors.Wrap(customgovtypes.ErrNotEnoughPermissions, "PermUpsertTokenRate")
	}

	err := ck.UpsertTokenRate(ctx, *types.NewTokenRate(
		msg.Denom,
		msg.Rate,
		msg.FeePayments,
	))
	return &sdk.Result{}, err
}
