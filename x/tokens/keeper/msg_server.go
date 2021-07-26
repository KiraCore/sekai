package keeper

import (
	"context"
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"

	govtypes "github.com/KiraCore/sekai/x/gov/types"
	"github.com/KiraCore/sekai/x/tokens/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
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

func (k msgServer) UpsertTokenRate(goCtx context.Context, msg *types.MsgUpsertTokenRate) (*types.MsgUpsertTokenRateResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	err := msg.ValidateBasic()
	if err != nil {
		return nil, errors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	isAllowed := k.cgk.CheckIfAllowedPermission(ctx, msg.Proposer, govtypes.PermUpsertTokenRate)
	if !isAllowed {
		return nil, errors.Wrap(govtypes.ErrNotEnoughPermissions, govtypes.PermUpsertTokenRate.String())
	}

	err = k.keeper.UpsertTokenRate(ctx, *types.NewTokenRate(
		msg.Denom,
		msg.Rate,
		msg.FeePayments,
	))

	if err != nil {
		return nil, errors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeUpsertTokenRate,
			sdk.NewAttribute(types.AttributeKeyProposer, msg.Proposer.String()),
			sdk.NewAttribute(types.AttributeKeyDenom, msg.Denom),
			sdk.NewAttribute(types.AttributeKeyRate, msg.Rate.String()),
			sdk.NewAttribute(types.AttributeKeyFeePayments, fmt.Sprintf("%t", msg.FeePayments)),
		),
	)

	return &types.MsgUpsertTokenRateResponse{}, nil
}
