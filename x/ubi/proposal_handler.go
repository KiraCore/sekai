package ubi

import (
	kiratypes "github.com/KiraCore/sekai/types"
	"github.com/KiraCore/sekai/x/gov/types"
	"github.com/KiraCore/sekai/x/ubi/keeper"
	ubitypes "github.com/KiraCore/sekai/x/ubi/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type ApplyUpsertUBIProposalHandler struct {
	keeper keeper.Keeper
}

func NewUpsertUBIProposalHandler(keeper keeper.Keeper) *ApplyUpsertUBIProposalHandler {
	return &ApplyUpsertUBIProposalHandler{
		keeper: keeper,
	}
}

func (a ApplyUpsertUBIProposalHandler) ProposalType() string {
	return kiratypes.ProposalTypeUpsertUBI
}

func (a ApplyUpsertUBIProposalHandler) Apply(ctx sdk.Context, proposalID uint64, proposal types.Content) error {
	p := proposal.(*ubitypes.UpsertUBIProposal)

	// TODO: The proposal should fail if sum of all ((float)amount / period) * 31556952 for all UBI records is greater than ubi-hard-cap.

	_ = p
	return nil
}

type ApplyRemoveUBIProposalHandler struct {
	keeper keeper.Keeper
}

func NewApplyRemoveUBIProposalHandler(keeper keeper.Keeper) *ApplyRemoveUBIProposalHandler {
	return &ApplyRemoveUBIProposalHandler{keeper: keeper}
}

func (a ApplyRemoveUBIProposalHandler) ProposalType() string {
	return kiratypes.ProposalTypeRemoveUBI
}

func (a ApplyRemoveUBIProposalHandler) Apply(ctx sdk.Context, proposalID uint64, proposal types.Content) error {
	p := proposal.(*ubitypes.RemoveUBIProposal)
	_ = p
	return nil
}
