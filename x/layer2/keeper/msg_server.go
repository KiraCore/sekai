package keeper

import (
	"context"

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
	_ = ctx

	// TODO: permission check PermCreateDappProposalWithoutBond or Bonds check (1% of properties.MinDappBond)
	// TODO: send initial bond to module account
	// TODO: create dapp object

	return &types.MsgCreateDappProposalResponse{}, nil
}

func (k msgServer) BondDappProposal(goCtx context.Context, msg *types.MsgBondDappProposal) (*types.MsgBondDappProposalResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	_ = ctx

	// Since the value of tokens that are required to successfully pass a
	// proposal is very large, we must make it possible for many people to
	// collaborate in the dApp proposal submission. We should start with the
	// deployer having to commit a **minimum of 1%** of the `min_dapp_bond`
	// while the remaining 99% of the bond can be supplied by ANY other user.

	return &types.MsgBondDappProposalResponse{}, nil
}

func (k msgServer) ReclaimDappBondProposal(goCtx context.Context, msg *types.MsgReclaimDappBondProposal) (*types.MsgReclaimDappBondProposalResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	_ = ctx

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
