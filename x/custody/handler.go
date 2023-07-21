package custody

import (
	"github.com/KiraCore/sekai/x/custody/keeper"
	"github.com/KiraCore/sekai/x/custody/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
)

// NewHandler returns new instance of handler
func NewHandler(ck keeper.Keeper, cgk types.CustomGovKeeper, bk types.BankKeeper) sdk.Handler {
	msgServer := keeper.NewMsgServerImpl(ck, cgk, bk)

	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case *types.MsgCreateCustodyRecord:
			res, err := msgServer.CreateCustody(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgDisableCustodyRecord:
			res, err := msgServer.DisableCustody(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgDropCustodyRecord:
			res, err := msgServer.DropCustody(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgAddToCustodyCustodians:
			res, err := msgServer.AddToCustodians(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgRemoveFromCustodyCustodians:
			res, err := msgServer.RemoveFromCustodians(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgDropCustodyCustodians:
			res, err := msgServer.DropCustodians(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgAddToCustodyWhiteList:
			res, err := msgServer.AddToWhiteList(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgRemoveFromCustodyWhiteList:
			res, err := msgServer.RemoveFromWhiteList(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgDropCustodyWhiteList:
			res, err := msgServer.DropWhiteList(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgAddToCustodyLimits:
			res, err := msgServer.AddToLimits(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgRemoveFromCustodyLimits:
			res, err := msgServer.RemoveFromLimits(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgDropCustodyLimits:
			res, err := msgServer.DropLimits(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgSend:
			res, err := msgServer.Send(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgApproveCustodyTransaction:
			res, err := msgServer.ApproveTransaction(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgDeclineCustodyTransaction:
			res, err := msgServer.DeclineTransaction(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgPasswordConfirmTransaction:
			res, err := msgServer.PasswordConfirm(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		default:
			return nil, errors.Wrapf(errors.ErrUnknownRequest, "unrecognized %s message type: %T", types.ModuleName, msg)
		}
	}
}
