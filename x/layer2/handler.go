package layer2

import (
	"github.com/KiraCore/sekai/x/layer2/keeper"
	"github.com/KiraCore/sekai/x/layer2/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
)

// NewHandler returns new instance of handler
func NewHandler(k keeper.Keeper) sdk.Handler {
	msgServer := keeper.NewMsgServerImpl(k)

	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case *types.MsgCreateDappProposal:
			res, err := msgServer.CreateDappProposal(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgBondDappProposal:
			res, err := msgServer.BondDappProposal(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgReclaimDappBondProposal:
			res, err := msgServer.ReclaimDappBondProposal(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgJoinDappVerifierWithBond:
			res, err := msgServer.JoinDappVerifierWithBond(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgExitDapp:
			res, err := msgServer.ExitDapp(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgRedeemDappPoolTx:
			res, err := msgServer.RedeemDappPoolTx(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgSwapDappPoolTx:
			res, err := msgServer.SwapDappPoolTx(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgConvertDappPoolTx:
			res, err := msgServer.ConvertDappPoolTx(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgPauseDappTx:
			res, err := msgServer.PauseDappTx(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgUnPauseDappTx:
			res, err := msgServer.UnPauseDappTx(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgReactivateDappTx:
			res, err := msgServer.ReactivateDappTx(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgExecuteDappTx:
			res, err := msgServer.ExecuteDappTx(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgDenounceLeaderTx:
			res, err := msgServer.DenounceLeaderTx(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgTransitionDappTx:
			res, err := msgServer.TransitionDappTx(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgApproveDappTransitionTx:
			res, err := msgServer.ApproveDappTransitionTx(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgRejectDappTransitionTx:
			res, err := msgServer.RejectDappTransitionTx(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgTransferDappTx:
			res, err := msgServer.TransferDappTx(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgAckTransferDappTx:
			res, err := msgServer.AckTransferDappTx(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgMintCreateFtTx:
			res, err := msgServer.MintCreateFtTx(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgMintCreateNftTx:
			res, err := msgServer.MintCreateNftTx(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgMintIssueTx:
			res, err := msgServer.MintIssueTx(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgMintBurnTx:
			res, err := msgServer.MintBurnTx(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		default:
			return nil, errors.Wrapf(errors.ErrUnknownRequest, "unrecognized %s message type: %T", types.ModuleName, msg)
		}
	}
}
