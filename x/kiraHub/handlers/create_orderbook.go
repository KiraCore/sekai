package handlers

import (
	"github.com/KiraCore/sekai/x/kiraHub/keeper"
	"github.com/KiraCore/sekai/x/kiraHub/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func HandlerMsgCreateOrderBook(context sdk.Context, keeper keeper.Keeper, message *types.MsgCreateOrderBook) (*sdk.Result, error) {
	keeper.CreateOrderBook(context, message.Quote, message.Base, message.Curator, message.Mnemonic)
	return &sdk.Result{}, nil
}
