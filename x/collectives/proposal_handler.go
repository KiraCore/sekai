package collectives

import (
	kiratypes "github.com/KiraCore/sekai/types"
	"github.com/KiraCore/sekai/x/collectives/keeper"
	"github.com/KiraCore/sekai/x/collectives/types"
	govtypes "github.com/KiraCore/sekai/x/gov/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type ApplyCollectiveSendDonationProposalHandler struct {
	keeper keeper.Keeper
}

func NewApplyCollectiveSendDonationProposalHandler(keeper keeper.Keeper) *ApplyCollectiveSendDonationProposalHandler {
	return &ApplyCollectiveSendDonationProposalHandler{
		keeper: keeper,
	}
}

func (a ApplyCollectiveSendDonationProposalHandler) ProposalType() string {
	return kiratypes.ProposalTypeCollectiveSendDonation
}

func (a ApplyCollectiveSendDonationProposalHandler) Apply(ctx sdk.Context, proposalID uint64, proposal govtypes.Content, slash uint64) error {
	p := proposal.(*types.ProposalCollectiveSendDonation)
	return a.keeper.SendDonation(ctx, p.Name, p.Address, p.Amount)
}

type ApplyCollectiveUpdateProposalHandler struct {
	keeper keeper.Keeper
}

func NewApplyCollectiveUpdateProposalHandler(keeper keeper.Keeper) *ApplyCollectiveUpdateProposalHandler {
	return &ApplyCollectiveUpdateProposalHandler{
		keeper: keeper,
	}
}

func (a ApplyCollectiveUpdateProposalHandler) ProposalType() string {
	return kiratypes.ProposalTypeCollectiveUpdate
}

func (a ApplyCollectiveUpdateProposalHandler) Apply(ctx sdk.Context, proposalID uint64, proposal govtypes.Content, slash uint64) error {
	p := proposal.(*types.ProposalCollectiveUpdate)

	a.keeper.SetCollective(ctx, p.Collective)
	return nil
}

type ApplyCollectiveRemoveProposalHandler struct {
	keeper keeper.Keeper
}

func NewApplyCollectiveRemoveProposalHandler(keeper keeper.Keeper) *ApplyCollectiveRemoveProposalHandler {
	return &ApplyCollectiveRemoveProposalHandler{
		keeper: keeper,
	}
}

func (a ApplyCollectiveRemoveProposalHandler) ProposalType() string {
	return kiratypes.ProposalTypeCollectiveRemove
}

func (a ApplyCollectiveRemoveProposalHandler) Apply(ctx sdk.Context, proposalID uint64, proposal govtypes.Content, slash uint64) error {
	p := proposal.(*types.ProposalCollectiveRemove)
	a.keeper.DeleteCollective(ctx, p.Name)
	return nil
}
