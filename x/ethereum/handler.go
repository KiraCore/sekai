package ethereum

import (
	"github.com/KiraCore/sekai/x/ethereum/keeper"
	"github.com/KiraCore/sekai/x/ethereum/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
)

// NewHandler returns new instance of handler
func NewHandler(ck keeper.Keeper, cgk types.CustomGovKeeper, bk types.BankKeeper) sdk.Handler {
	msgServer := keeper.NewMsgServerImpl(ck, cgk, bk)

	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case *types.MsgRelay:
			{
				res, err := msgServer.Relay(sdk.WrapSDKContext(ctx), msg)
				return sdk.WrapServiceResult(ctx, res, err)
			}
		default:
			return nil, errors.Wrapf(errors.ErrUnknownRequest, "unrecognized %s message type: %T", types.ModuleName, msg)
		}
	}
}
