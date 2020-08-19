package kiraHub

import (
	"github.com/KiraCore/sekai/x/kiraHub/handlers"
	"github.com/KiraCore/sekai/x/kiraHub/keeper"
	"github.com/KiraCore/sekai/x/kiraHub/types"
	"github.com/pkg/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewHandler(k keeper.Keeper) sdk.Handler {
	return func(context sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		context = context.WithEventManager(sdk.NewEventManager())

		switch message := msg.(type) {
		case *types.MsgCreateOrderBook:
			return handlers.HandlerMsgCreateOrderBook(context, k, message)

		case *types.MsgCreateOrder:
			return handlers.HandlerMsgCreateOrder(context, k, message)

		case *types.MsgUpsertSignerKey:
			return handlers.HandlerMsgUpsertSignerKey(context, k, message)

		default:
			return nil, errors.Wrapf(types.UnknownMessageCode, "%T", msg)
		}
	}
}
