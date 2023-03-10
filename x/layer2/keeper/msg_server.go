package keeper

import (
	"context"

	govtypes "github.com/KiraCore/sekai/x/gov/types"
	"github.com/KiraCore/sekai/x/layer2/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type msgServer struct {
	keeper Keeper
}

// NewMsgServerImpl returns an implementation of the bank MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{
		keeper: keeper,
	}
}

var _ types.MsgServer = msgServer{}

func (k msgServer) CreateDappProposal(goCtx context.Context, msg *types.MsgCreateDappProposal) (*types.MsgCreateDappProposalResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	properties := k.keeper.gk.GetNetworkProperties(ctx)
	addr := sdk.MustAccAddressFromBech32(msg.Sender)

	// permission check PermCreateDappProposalWithoutBond
	isAllowed := k.keeper.CheckIfAllowedPermission(ctx, addr, govtypes.PermCreateDappProposalWithoutBond)
	if !isAllowed {
		minDappBond := properties.MinDappBond
		if msg.Bond.Denom != k.keeper.BondDenom(ctx) {
			return nil, types.ErrInvalidDappBondDenom
		}
		// check 1% of properties.MinDappBond
		if msg.Bond.Amount.Mul(sdk.NewInt(100)).LT(sdk.NewInt(int64(minDappBond)).Mul(sdk.NewInt(1000_000))) {
			return nil, types.ErrLowAmountToCreateDappProposal
		}
	}

	// send initial bond to module account
	if msg.Bond.IsPositive() {
		err := k.keeper.bk.SendCoinsFromAccountToModule(ctx, addr, types.ModuleName, sdk.Coins{msg.Bond})
		if err != nil {
			return nil, err
		}
	}

	dapp := k.keeper.GetDapp(ctx, msg.Dapp.Name)
	if dapp.Name != "" {
		return nil, types.ErrDappAlreadyExists
	}

	// create dapp object
	msg.Dapp.TotalBond = msg.Bond
	msg.Dapp.CreationTime = uint64(ctx.BlockTime().Unix())
	k.keeper.SetDapp(ctx, msg.Dapp)
	k.keeper.SetUserDappBond(ctx, types.UserDappBond{
		DappName: msg.Dapp.Name,
		User:     msg.Sender,
		Bond:     msg.Bond,
	})

	return &types.MsgCreateDappProposalResponse{}, nil
}

func (k msgServer) BondDappProposal(goCtx context.Context, msg *types.MsgBondDappProposal) (*types.MsgBondDappProposalResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	dapp := k.keeper.GetDapp(ctx, msg.DappName)
	if dapp.Name == "" {
		return nil, types.ErrDappDoesNotExist
	}

	if k.keeper.BondDenom(ctx) != msg.Bond.Denom {
		return nil, types.ErrInvalidDappBondDenom
	}

	// send initial bond to module account
	addr := sdk.MustAccAddressFromBech32(msg.Sender)
	err := k.keeper.bk.SendCoinsFromAccountToModule(ctx, addr, types.ModuleName, sdk.Coins{msg.Bond})
	if err != nil {
		return nil, err
	}

	properties := k.keeper.gk.GetNetworkProperties(ctx)
	if dapp.TotalBond.Amount.GTE(sdk.NewInt(int64(properties.MaxDappBond)).Mul(sdk.NewInt(1000_000))) {
		return nil, types.ErrMaxDappBondReached
	}

	dapp.TotalBond = dapp.TotalBond.Add(msg.Bond)
	k.keeper.SetDapp(ctx, dapp)

	userDappBond := k.keeper.GetUserDappBond(ctx, msg.DappName, msg.Sender)
	if userDappBond.User != "" {
		userDappBond.Bond = userDappBond.Bond.Add(msg.Bond)
	} else {
		userDappBond = types.UserDappBond{
			User:     msg.Sender,
			DappName: msg.DappName,
			Bond:     msg.Bond,
		}
	}
	k.keeper.SetUserDappBond(ctx, userDappBond)

	return &types.MsgBondDappProposalResponse{}, nil
}

func (k msgServer) ReclaimDappBondProposal(goCtx context.Context, msg *types.MsgReclaimDappBondProposal) (*types.MsgReclaimDappBondProposalResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	userDappBond := k.keeper.GetUserDappBond(ctx, msg.DappName, msg.Sender)
	if userDappBond.DappName == "" {
		return nil, types.ErrUserDappBondDoesNotExist
	}
	if userDappBond.Bond.Denom != msg.Bond.Denom {
		return nil, types.ErrInvalidDappBondDenom
	}
	if userDappBond.Bond.Amount.LT(msg.Bond.Amount) {
		return nil, types.ErrNotEnoughUserDappBond
	}

	userDappBond.Bond.Amount = userDappBond.Bond.Amount.Sub(msg.Bond.Amount)
	k.keeper.SetUserDappBond(ctx, userDappBond)

	// send tokens back to user
	addr := sdk.MustAccAddressFromBech32(msg.Sender)
	err := k.keeper.bk.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr, sdk.Coins{msg.Bond})
	if err != nil {
		return nil, err
	}

	return &types.MsgReclaimDappBondProposalResponse{}, nil
}

func (k msgServer) JoinDappTx(goCtx context.Context, msg *types.MsgJoinDappTx) (*types.MsgJoinDappTxResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	_ = ctx

	return &types.MsgJoinDappTxResponse{}, nil
}

func (k msgServer) ExitDapp(goCtx context.Context, msg *types.MsgExitDapp) (*types.MsgExitDappResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	_ = ctx

	return &types.MsgExitDappResponse{}, nil
}

func (k msgServer) VoteDappOperatorTx(goCtx context.Context, msg *types.MsgVoteDappOperatorTx) (*types.MsgVoteDappOperatorTxResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	_ = ctx

	return &types.MsgVoteDappOperatorTxResponse{}, nil
}

func (k msgServer) RedeemDappPoolTx(goCtx context.Context, msg *types.MsgRedeemDappPoolTx) (*types.MsgRedeemDappPoolTxResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	_ = ctx

	return &types.MsgRedeemDappPoolTxResponse{}, nil
}

func (k msgServer) SwapDappPoolTx(goCtx context.Context, msg *types.MsgSwapDappPoolTx) (*types.MsgSwapDappPoolTxResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	_ = ctx

	return &types.MsgSwapDappPoolTxResponse{}, nil
}

func (k msgServer) ConvertDappPoolTx(goCtx context.Context, msg *types.MsgConvertDappPoolTx) (*types.MsgConvertDappPoolTxResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	_ = ctx

	return &types.MsgConvertDappPoolTxResponse{}, nil
}

func (k msgServer) PauseDappTx(goCtx context.Context, msg *types.MsgPauseDappTx) (*types.MsgPauseDappTxResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	_ = ctx

	return &types.MsgPauseDappTxResponse{}, nil
}

func (k msgServer) UnPauseDappTx(goCtx context.Context, msg *types.MsgUnPauseDappTx) (*types.MsgUnPauseDappTxResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	_ = ctx

	return &types.MsgUnPauseDappTxResponse{}, nil
}

func (k msgServer) ReactivateDappTx(goCtx context.Context, msg *types.MsgReactivateDappTx) (*types.MsgReactivateDappTxResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	_ = ctx

	return &types.MsgReactivateDappTxResponse{}, nil
}

func (k msgServer) ExecuteDappTx(goCtx context.Context, msg *types.MsgExecuteDappTx) (*types.MsgExecuteDappTxResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	_ = ctx

	return &types.MsgExecuteDappTxResponse{}, nil
}

func (k msgServer) DenounceLeaderTx(goCtx context.Context, msg *types.MsgDenounceLeaderTx) (*types.MsgDenounceLeaderTxResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	_ = ctx

	return &types.MsgDenounceLeaderTxResponse{}, nil
}

func (k msgServer) TransitionDappTx(goCtx context.Context, msg *types.MsgTransitionDappTx) (*types.MsgTransitionDappTxResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	_ = ctx

	return &types.MsgTransitionDappTxResponse{}, nil
}

func (k msgServer) ApproveDappTransitionTx(goCtx context.Context, msg *types.MsgApproveDappTransitionTx) (*types.MsgApproveDappTransitionTxResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	_ = ctx

	return &types.MsgApproveDappTransitionTxResponse{}, nil
}

func (k msgServer) RejectDappTransitionTx(goCtx context.Context, msg *types.MsgRejectDappTransitionTx) (*types.MsgRejectDappTransitionTxResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	_ = ctx

	return &types.MsgRejectDappTransitionTxResponse{}, nil
}

func (k msgServer) UpsertDappProposalTx(goCtx context.Context, msg *types.MsgUpsertDappProposalTx) (*types.MsgUpsertDappProposalTxResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	_ = ctx

	return &types.MsgUpsertDappProposalTxResponse{}, nil
}

func (k msgServer) VoteUpsertDappProposalTx(goCtx context.Context, msg *types.MsgVoteUpsertDappProposalTx) (*types.MsgVoteUpsertDappProposalTxResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	_ = ctx

	return &types.MsgVoteUpsertDappProposalTxResponse{}, nil
}

func (k msgServer) TransferDappTx(goCtx context.Context, msg *types.MsgTransferDappTx) (*types.MsgTransferDappTxResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	_ = ctx

	return &types.MsgTransferDappTxResponse{}, nil
}

func (k msgServer) MintCreateFtTx(goCtx context.Context, msg *types.MsgMintCreateFtTx) (*types.MsgMintCreateFtTxResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	_ = ctx

	return &types.MsgMintCreateFtTxResponse{}, nil
}

func (k msgServer) MintCreateNftTx(goCtx context.Context, msg *types.MsgMintCreateNftTx) (*types.MsgMintCreateNftTxResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	_ = ctx

	return &types.MsgMintCreateNftTxResponse{}, nil
}

func (k msgServer) MintIssueTx(goCtx context.Context, msg *types.MsgMintIssueTx) (*types.MsgMintIssueTxResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	_ = ctx

	return &types.MsgMintIssueTxResponse{}, nil
}

func (k msgServer) MintBurnTx(goCtx context.Context, msg *types.MsgMintBurnTx) (*types.MsgMintBurnTxResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	_ = ctx

	return &types.MsgMintBurnTxResponse{}, nil
}

// TODO: implement - step1
//   rpc CreateDappProposal(MsgCreateDappProposal) returns (MsgCreateDappProposalResponse);
//   rpc BondDappProposal(MsgBondDappProposal) returns (MsgBondDappProposalResponse);
//   rpc ReclaimDappBondProposal(MsgReclaimDappBondProposal) returns (MsgReclaimDappBondProposalResponse);
//   rpc JoinDappTx(MsgJoinDappTx) returns (MsgJoinDappTxResponse);
//   rpc ExitDapp(MsgExitDapp) returns (MsgExitDappResponse);
//   rpc VoteDappOperatorTx(MsgVoteDappOperatorTx) returns (MsgVoteDappOperatorTxResponse);
//   rpc ExecuteDappTx(MsgExecuteDappTx) returns (MsgExecuteDappTxResponse);
//   rpc DenounceLeaderTx(MsgDenounceLeaderTx) returns (MsgDenounceLeaderTxResponse);
//   rpc TransitionDappTx(MsgTransitionDappTx) returns (MsgTransitionDappTxResponse);
//   rpc ApproveDappTransitionTx(MsgApproveDappTransitionTx) returns (MsgApproveDappTransitionTxResponse);
//   rpc RejectDappTransitionTx(MsgRejectDappTransitionTx) returns (MsgRejectDappTransitionTxResponse);

// TODO: implement - step2
//   rpc PauseDappTx(MsgPauseDappTx) returns (MsgPauseDappTxResponse);
//   rpc UnPauseDappTx(MsgUnPauseDappTx) returns (MsgUnPauseDappTxResponse);
//   rpc ReactivateDappTx(MsgReactivateDappTx) returns (MsgReactivateDappTxResponse);
//   rpc RedeemDappPoolTx(MsgRedeemDappPoolTx) returns (MsgRedeemDappPoolTxResponse);
//   rpc SwapDappPoolTx(MsgSwapDappPoolTx) returns (MsgSwapDappPoolTxResponse);
//   rpc ConvertDappPoolTx(MsgConvertDappPoolTx) returns (MsgConvertDappPoolTxResponse);
//   rpc UpsertDappProposalTx(MsgUpsertDappProposalTx) returns (MsgUpsertDappProposalTxResponse);
//   rpc VoteUpsertDappProposalTx(MsgVoteUpsertDappProposalTx) returns (MsgVoteUpsertDappProposalTxResponse);
//   rpc TransferDappTx(MsgTransferDappTx) returns (MsgTransferDappTxResponse);
//   rpc MintCreateFtTx(MsgMintCreateFtTx) returns (MsgMintCreateFtTxResponse);
//   rpc MintCreateNftTx(MsgMintCreateNftTx) returns (MsgMintCreateNftTxResponse);
//   rpc MintIssueTx(MsgMintIssueTx) returns (MsgMintIssueTxResponse);
//   rpc MintBurnTx(MsgMintBurnTx) returns (MsgMintBurnTxResponse);
