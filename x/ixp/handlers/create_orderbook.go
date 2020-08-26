package handlers

import (
	"github.com/KiraCore/sekai/x/ixp/keeper"
	"github.com/KiraCore/sekai/x/ixp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// HandlerMsgCreateOrderBook is a function to handler create orderbook message
func HandlerMsgCreateOrderBook(context sdk.Context, keeper keeper.Keeper, message *types.MsgCreateOrderBook) (*sdk.Result, error) {
	keeper.CreateOrderBook(context, message.Quote, message.Base, message.Curator, message.Mnemonic)
	return &sdk.Result{}, nil
}
