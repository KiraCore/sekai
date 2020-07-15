package kiraHub

import (
	constants "github.com/KiraCore/sekai/x/kiraHub/constants"
	"github.com/KiraCore/sekai/x/kiraHub/transactions/createOrder"
	"github.com/KiraCore/sekai/x/kiraHub/transactions/createOrderBook"
	"github.com/pkg/errors"

	sdk "github.com/KiraCore/cosmos-sdk/types"
)


func NewHandler(keeper Keeper) sdk.Handler {
	return func(context sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		context = context.WithEventManager(sdk.NewEventManager())

		switch message := msg.(type) {
		case createOrderBook.Message:
			return createOrderBook.HandleMessage(context, keeper.getCreateOrderBookKeeper(), message)

		case createOrder.Message:
			return createOrder.HandleMessage(context, keeper.getCreateOrderKeeper(), message)

		default:
			return nil, errors.Wrapf(constants.UnknownMessageCode, "%T", msg)
		}
	}
}

