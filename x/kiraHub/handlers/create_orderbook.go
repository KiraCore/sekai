package handlers

import (
	sdk "github.com/KiraCore/cosmos-sdk/types"
	"github.com/KiraCore/sekai/x/kiraHub/keeper"
	"github.com/KiraCore/sekai/x/kiraHub/types"
)

func HandlerMsgCreateOrderBook(context sdk.Context, keeper keeper.Keeper, message *types.MsgCreateOrderBook) (*sdk.Result, error) {
	keeper.CreateOrderBook(context, message.Quote, message.Base, message.Curator, message.Mnemonic)
	return &sdk.Result{}, nil
}
