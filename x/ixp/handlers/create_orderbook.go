package handlers

import (
	"encoding/json"

	"github.com/KiraCore/sekai/x/ixp/keeper"
	"github.com/KiraCore/sekai/x/ixp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// CreateOrderBookResponse describes create order response struct
type CreateOrderBookResponse struct {
	ID string
}

// HandlerMsgCreateOrderBook is a function to handler create orderbook message
func HandlerMsgCreateOrderBook(context sdk.Context, keeper keeper.Keeper, message *types.MsgCreateOrderBook) (*sdk.Result, error) {
	id, err := keeper.CreateOrderBook(context, message.Quote, message.Base, message.Curator, message.Mnemonic)
	res := CreateOrderBookResponse{
		ID: id,
	}
	jsonData, err := json.Marshal(res)
	if err != nil {
		return &sdk.Result{}, err
	}

	return &sdk.Result{Data: jsonData}, err
}
