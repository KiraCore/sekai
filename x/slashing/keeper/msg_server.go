package keeper

import (
	"context"

	"github.com/KiraCore/sekai/x/slashing/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the slashing MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

// Activate implements MsgServer.Activate method.
// Validators must submit a transaction to activate itself after
// having been inactivated (and thus unbonded) for downtime
func (k msgServer) Activate(goCtx context.Context, msg *types.MsgActivate) (*types.MsgActivateResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	valAddr, valErr := sdk.ValAddressFromBech32(msg.ValidatorAddr)
	if valErr != nil {
		return nil, valErr
	}
	err := k.Keeper.Activate(ctx, valAddr)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.ValidatorAddr),
		),
	)

	return &types.MsgActivateResponse{}, nil
}
