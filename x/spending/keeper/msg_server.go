package keeper

import (
	"context"
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"

	govtypes "github.com/KiraCore/sekai/x/gov/types"
	"github.com/KiraCore/sekai/x/spending/types"
)

type msgServer struct {
	keeper Keeper
	cgk    types.CustomGovKeeper
}

// NewMsgServerImpl returns an implementation of the bank MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper, cgk types.CustomGovKeeper) types.MsgServer {
	return &msgServer{
		keeper: keeper,
		cgk:    cgk,
	}
}

var _ types.MsgServer = msgServer{}

func (k msgServer) UpsertTokenAlias(
	goCtx context.Context,
	msg *types.MsgUpsertTokenAlias,
) (*types.MsgUpsertTokenAliasResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	isAllowed := k.cgk.CheckIfAllowedPermission(ctx, msg.Proposer, govtypes.PermUpsertTokenAlias)
	if !isAllowed {
		return nil, errors.Wrap(govtypes.ErrNotEnoughPermissions, govtypes.PermUpsertTokenAlias.String())
	}

	err := k.keeper.UpsertTokenAlias(ctx, *types.NewTokenAlias(
		msg.Symbol,
		msg.Name,
		msg.Icon,
		msg.Decimals,
		msg.Denoms,
	))
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeUpsertTokenAlias,
			sdk.NewAttribute(types.AttributeKeyProposer, msg.Proposer.String()),
			sdk.NewAttribute(types.AttributeKeySymbol, msg.Symbol),
			sdk.NewAttribute(types.AttributeKeyName, msg.Name),
			sdk.NewAttribute(types.AttributeKeyIcon, msg.Icon),
			sdk.NewAttribute(types.AttributeKeyDecimals, fmt.Sprintf("%d", msg.Decimals)),
			sdk.NewAttribute(types.AttributeKeyDenoms, strings.Join(msg.Denoms, ",")),
		),
	)
	return &types.MsgUpsertTokenAliasResponse{}, err
}
