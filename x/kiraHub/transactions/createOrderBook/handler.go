package createOrderBook

import sdk "github.com/KiraCore/cosmos-sdk/types"

func HandleMessage(context sdk.Context, keeper Keeper, message Message) (*sdk.Result, error) {
	keeper.CreateOrderBook(context, message.Quote, message.Base, message.Curator, message.Mnemonic)
	return &sdk.Result{}, nil
}