package keeper

import (
	"context"
	customgovtypes "github.com/KiraCore/sekai/x/gov/types"
	"github.com/KiraCore/sekai/x/upgrade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	"time"
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

	proposalID, err := m.CreateAndSaveProposalWithContent(ctx, &types.ProposalSoftwareUpgrade{
		Resources:            msg.Resources,
		MinHaltTime:          msg.MinHaltTime,
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

func (k msgServer) CreateAndSaveProposalWithContent(ctx sdk.Context, content customgovtypes.Content) (uint64, error) {
	blockTime := ctx.BlockTime()
	proposalID, err := k.cgk.GetNextProposalID(ctx)
	if err != nil {
		return 0, err
	}

	properties := k.cgk.GetNetworkProperties(ctx)

	proposal, err := customgovtypes.NewProposal(
		proposalID,
		content,
		blockTime,
		blockTime.Add(time.Minute*time.Duration(properties.ProposalEndTime)),
		blockTime.Add(time.Minute*time.Duration(properties.ProposalEnactmentTime)),
	)

	k.cgk.SaveProposal(ctx, proposal)
	k.cgk.AddToActiveProposals(ctx, proposal)

	return proposalID, nil
}
