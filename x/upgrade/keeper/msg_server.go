package keeper

import (
	"context"
	"time"

	customgovtypes "github.com/KiraCore/sekai/x/gov/types"
	"github.com/KiraCore/sekai/x/upgrade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
)

type msgServer struct {
	keeper Keeper
	cgk    types.CustomGovKeeper
}

// NewMsgServerImpl returns an implementation of the upgrade MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper, cgk types.CustomGovKeeper) types.MsgServer {
	return &msgServer{
		keeper: keeper,
		cgk:    cgk,
	}
}

func (m msgServer) ProposalSoftwareUpgrade(goCtx context.Context, msg *types.MsgProposalSoftwareUpgradeRequest) (*types.MsgProposalSoftwareUpgradeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	isAllowed := m.cgk.CheckIfAllowedPermission(ctx, msg.Proposer, customgovtypes.PermCreateUpsertTokenRateProposal)
	if !isAllowed {
		return nil, errors.Wrap(customgovtypes.ErrNotEnoughPermissions, customgovtypes.PermCreateUpsertTokenRateProposal.String())
	}

	proposalID, err := m.CreateAndSaveProposalWithContent(ctx, msg.Memo, &types.ProposalSoftwareUpgrade{
		Name:                 msg.Name,
		Resources:            msg.Resources,
		Height:               msg.Height,
		MinUpgradeTime:       msg.MinUpgradeTime,
		OldChainId:           msg.OldChainId,
		NewChainId:           msg.NewChainId,
		RollbackChecksum:     msg.RollbackChecksum,
		MaxEnrolmentDuration: msg.MaxEnrolmentDuration,
		Memo:                 msg.Memo,
	})

	return &types.MsgProposalSoftwareUpgradeResponse{
		ProposalID: proposalID,
	}, err
}

func (k msgServer) CreateAndSaveProposalWithContent(ctx sdk.Context, description string, content customgovtypes.Content) (uint64, error) {
	blockTime := ctx.BlockTime()
	proposalID := k.cgk.GetNextProposalIDAndIncrement(ctx)
	properties := k.cgk.GetNetworkProperties(ctx)

	proposal, err := customgovtypes.NewProposal(
		proposalID,
		content,
		blockTime,
		blockTime.Add(time.Second*time.Duration(properties.ProposalEndTime)),
		blockTime.Add(time.Second*time.Duration(properties.ProposalEndTime)+
			time.Second*time.Duration(properties.ProposalEnactmentTime),
		),
		ctx.BlockHeight()+int64(properties.MinProposalEndBlocks),
		ctx.BlockHeight()+int64(properties.MinProposalEndBlocks+properties.MinProposalEnactmentBlocks),
		description,
	)
	if err != nil {
		return 0, err
	}

	k.cgk.SaveProposal(ctx, proposal)
	k.cgk.AddToActiveProposals(ctx, proposal)

	return proposalID, nil
}
