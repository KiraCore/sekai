package handlers

import (
	"github.com/KiraCore/sekai/x/ixp/keeper"
	"github.com/KiraCore/sekai/x/ixp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// HandlerMsgCreateOrder is a function to handler create order message
func HandlerMsgCreateOrder(context sdk.Context, keeper keeper.Keeper, message *types.MsgCreateOrder) (*sdk.Result, error) {
	keeper.CreateOrder(context, message.OrderBookID, message.OrderType, message.Amount, message.LimitPrice, message.ExpiryTime, message.Curator)
	return &sdk.Result{}, nil
}
