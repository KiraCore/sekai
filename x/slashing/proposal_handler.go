package slashing

import (
	kiratypes "github.com/KiraCore/sekai/types"
	govtypes "github.com/KiraCore/sekai/x/gov/types"
	"github.com/KiraCore/sekai/x/slashing/keeper"
	"github.com/KiraCore/sekai/x/slashing/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type ApplyResetWholeValidatorRankProposalHandler struct {
	keeper keeper.Keeper
}

func NewApplyResetWholeValidatorRankProposalHandler(keeper keeper.Keeper) *ApplyResetWholeValidatorRankProposalHandler {
	return &ApplyResetWholeValidatorRankProposalHandler{
		keeper: keeper,
	}
}

func (a ApplyResetWholeValidatorRankProposalHandler) ProposalType() string {
	return kiratypes.ProposalTypeResetWholeValidatorRank
}

func (a ApplyResetWholeValidatorRankProposalHandler) Apply(ctx sdk.Context, proposalID uint64, proposal govtypes.Content, slash sdk.Dec) error {
	_ = proposal.(*types.ProposalResetWholeValidatorRank)

	return a.keeper.ResetWholeValidatorRank(ctx)
}

type ApplySlashValidatorProposalHandler struct {
	keeper keeper.Keeper
}

func NewApplySlashValidatorProposalHandler(keeper keeper.Keeper) *ApplySlashValidatorProposalHandler {
	return &ApplySlashValidatorProposalHandler{
		keeper: keeper,
	}
}

func (a ApplySlashValidatorProposalHandler) ProposalType() string {
	return kiratypes.ProposalTypeSlashValidator
}

func (a ApplySlashValidatorProposalHandler) Apply(ctx sdk.Context, proposalID uint64, proposal govtypes.Content, slash sdk.Dec) error {
	content := proposal.(*types.ProposalSlashValidator)

	a.keeper.SlashStakingPool(ctx, content, slash)
	return nil
}
