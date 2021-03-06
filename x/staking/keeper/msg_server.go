package keeper

import (
	"context"
	"strings"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"

	govkeeper "github.com/KiraCore/sekai/x/gov/keeper"
	customgovtypes "github.com/KiraCore/sekai/x/gov/types"
	"github.com/KiraCore/sekai/x/staking/types"
	"github.com/cosmos/cosmos-sdk/types/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type msgServer struct {
	keeper    Keeper
	govKeeper govkeeper.Keeper
}

// NewMsgServerImpl returns an implementation of the bank MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper, govKeeper govkeeper.Keeper) types.MsgServer {
	return &msgServer{
		keeper:    keeper,
		govKeeper: govKeeper,
	}
}

var _ types.MsgServer = msgServer{}

func (k msgServer) ClaimValidator(goCtx context.Context, msg *types.MsgClaimValidator) (*types.MsgClaimValidatorResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	isAllowed := govkeeper.CheckIfAllowedPermission(ctx, k.govKeeper, sdk.AccAddress(msg.ValKey), customgovtypes.PermClaimValidator)
	if !isAllowed {
		return nil, errors.Wrap(customgovtypes.ErrNotEnoughPermissions, "PermClaimValidator")
	}

	pk, ok := msg.PubKey.GetCachedValue().(cryptotypes.PubKey)
	if !ok {
		return nil, errors.Wrapf(errors.ErrInvalidPubKey, "Expecting cryptotypes.PubKey, got %T", pk)
	}

	_, err := k.keeper.GetValidator(ctx, msg.ValKey)
	if err == nil {
		return nil, types.ErrValidatorAlreadyClaimed
	}

	_, err = k.keeper.GetValidatorByMoniker(ctx, msg.Moniker)
	if err == nil {
		return nil, types.ErrValidatorMonikerExists
	}

	validator, err := types.NewValidator(strings.Trim(msg.Moniker, " "), msg.Commission, msg.ValKey, pk)
	if err != nil {
		return nil, err
	}

	k.keeper.AddPendingValidator(ctx, validator)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeClaimValidator,
			sdk.NewAttribute(types.AttributeKeyMoniker, msg.Moniker),
			sdk.NewAttribute(types.AttributeKeyCommission, msg.Commission.String()),
			sdk.NewAttribute(types.AttributeKeyValKey, msg.ValKey.String()),
			sdk.NewAttribute(types.AttributeKeyPubKey, msg.PubKey.String()),
		),
	)
	return &types.MsgClaimValidatorResponse{}, nil
}
