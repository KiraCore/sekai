package tokens

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"

	customgovtypes "github.com/KiraCore/sekai/x/gov/types"
	"github.com/KiraCore/sekai/x/staking/types"
	"github.com/KiraCore/sekai/x/tokens/keeper"
	tokenstypes "github.com/KiraCore/sekai/x/tokens/types"
)

func NewHandler(ck keeper.Keeper, cgk tokenstypes.CustomGovKeeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case *tokenstypes.MsgUpsertTokenAlias:
			return handleUpsertTokenAlias(ctx, ck, cgk, msg)
		default:
			return nil, errors.Wrapf(errors.ErrUnknownRequest, "unrecognized %s message type: %T", types.ModuleName, msg)
		}
	}
}

func handleUpsertTokenAlias(ctx sdk.Context, ck keeper.Keeper, cgk tokenstypes.CustomGovKeeper, msg *tokenstypes.MsgUpsertTokenAlias) (*sdk.Result, error) {
	isAllowed := cgk.CheckIfAllowedPermission(ctx, msg.Proposer, customgovtypes.PermUpsertTokenAlias)
	if !isAllowed {
		return nil, errors.Wrap(customgovtypes.ErrNotEnoughPermissions, "PermUpsertTokenAlias")
	}

	ck.UpsertTokenAlias(ctx, msg.Address)
	return &sdk.Result{}, nil
}
