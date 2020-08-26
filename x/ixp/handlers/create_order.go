package handlers

import (
	"encoding/json"

	"github.com/KiraCore/sekai/x/ixp/keeper"
	"github.com/KiraCore/sekai/x/ixp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// CreateOrderResponse describes create order response struct
type CreateOrderResponse struct {
	ID string
}

// HandlerMsgCreateOrder is a function to handler create order message
func HandlerMsgCreateOrder(context sdk.Context, keeper keeper.Keeper, message *types.MsgCreateOrder) (*sdk.Result, error) {
	id, err := keeper.CreateOrder(context, message.OrderBookID, message.OrderType, message.Amount, message.LimitPrice, message.ExpiryTime, message.Curator)
	res := CreateOrderResponse{
		ID: id,
	}
	jsonData, err := json.Marshal(res)
	if err != nil {
		return &sdk.Result{}, err
	}

	return &sdk.Result{Data: jsonData}, err
}
