package spending

import (
	"github.com/KiraCore/sekai/x/spending/keeper"
	"github.com/KiraCore/sekai/x/spending/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
)

// NewHandler returns new instance of handler
func NewHandler(ck keeper.Keeper, cgk types.CustomGovKeeper, bk types.BankKeeper) sdk.Handler {
	msgServer := keeper.NewMsgServerImpl(ck, cgk, bk)

	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case *types.MsgCreateSpendingPool:
			res, err := msgServer.CreateSpendingPool(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgDepositSpendingPool:
			res, err := msgServer.DepositSpendingPool(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgClaimSpendingPool:
			res, err := msgServer.ClaimSpendingPool(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgRegisterSpendingPoolBeneficiary:
			res, err := msgServer.RegisterSpendingPoolBeneficiary(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		default:
			return nil, errors.Wrapf(errors.ErrUnknownRequest, "unrecognized %s message type: %T", types.ModuleName, msg)
		}
	}
}
