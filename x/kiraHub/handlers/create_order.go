package handlers

import (
	sdk "github.com/KiraCore/cosmos-sdk/types"
	"github.com/KiraCore/sekai/x/kiraHub/keeper"
	"github.com/KiraCore/sekai/x/kiraHub/types"
)

func HandlerMsgCreateOrder(context sdk.Context, keeper keeper.Keeper, message *types.MsgCreateOrder) (*sdk.Result, error) {
	keeper.CreateOrder(context, message.OrderBookID, message.OrderType, message.Amount, message.LimitPrice, message.ExpiryTime, message.Curator)
	return &sdk.Result{}, nil
}
