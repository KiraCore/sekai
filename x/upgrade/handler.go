package upgrade

import (
	"github.com/KiraCore/sekai/x/upgrade/keeper"
	"github.com/KiraCore/sekai/x/upgrade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
)

// NewHandler returns new instance of handler
func NewHandler(ck keeper.Keeper, cgk types.CustomGovKeeper) sdk.Handler {
	// msgServer := keeper.NewMsgServerImpl(ck, cgk)

	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		// case *types.MsgProposalSoftwareUpgradeRequest:
		// 	res, err := msgServer.ProposalSoftwareUpgrade(sdk.WrapSDKContext(ctx), msg)
		// 	return sdk.WrapServiceResult(ctx, res, err)

		// case *types.MsgProposalCancelSoftwareUpgradeRequest:
		// 	res, err := msgServer.ProposalCancelSoftwareUpgrade(sdk.WrapSDKContext(ctx), msg)
		// 	return sdk.WrapServiceResult(ctx, res, err)

		default:
			return nil, errors.Wrapf(errors.ErrUnknownRequest, "unrecognized %s message type: %T", types.ModuleName, msg)
		}
	}
}
