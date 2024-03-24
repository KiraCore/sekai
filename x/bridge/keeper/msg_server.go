package keeper

import (
	"context"
	"github.com/KiraCore/sekai/x/bridge/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type msgServer struct {
	keeper Keeper
	bk     types.BankKeeper
}

// NewMsgServerImpl returns an implementation of the bank MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper, bk types.BankKeeper) types.MsgServer {
	return &msgServer{
		keeper: keeper,
		bk:     bk,
	}
}

var _ types.MsgServer = msgServer{}

func (s msgServer) ChangeCosmosEthereum(goCtx context.Context, msg *types.MsgChangeCosmosEthereum) (*types.MsgChangeCosmosEthereumResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	record := types.ChangeCosmosEthereumRecord{
		From:      msg.From,
		To:        msg.To,
		InAmount:  msg.InAmount,
		OutAmount: msg.OutAmount,
	}

	if err := s.bk.IsSendEnabledCoins(ctx, msg.InAmount...); err != nil {
		return nil, err
	}

	err := s.bk.SendCoinsFromAccountToModule(ctx, msg.From, types.ModuleName, msg.InAmount)
	if err != nil {
		return nil, err
	}

	s.keeper.SetChangeCosmosEthereumRecord(ctx, record)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.From.String()),
		),
	)

	return &types.MsgChangeCosmosEthereumResponse{}, nil
}

func (s msgServer) ChangeEthereumCosmos(goCtx context.Context, msg *types.MsgChangeEthereumCosmos) (*types.MsgChangeEthereumCosmosResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	record := types.ChangeEthereumCosmosRecord{
		From:      msg.From,
		To:        msg.To,
		InAmount:  msg.InAmount,
		OutAmount: msg.OutAmount,
	}

	if err := s.bk.IsSendEnabledCoins(ctx, msg.OutAmount...); err != nil {
		return nil, err
	}

	err := s.bk.SendCoinsFromModuleToAccount(ctx, types.ModuleName, msg.To, msg.OutAmount)
	if err != nil {
		return nil, err
	}

	s.keeper.SetChangeEthereumCosmosRecord(ctx, record)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.To.String()),
		),
	)

	return &types.MsgChangeEthereumCosmosResponse{}, nil
}
