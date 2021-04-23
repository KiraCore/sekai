package slashing

import (
	"fmt"

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
	return types.ProposalTypeResetWholeValidatorRank
}

func (a ApplyResetWholeValidatorRankProposalHandler) Apply(ctx sdk.Context, proposal govtypes.Content) {
	_ = proposal.(*types.ProposalResetWholeValidatorRank)

	err := a.keeper.ResetWholeValidatorRank(ctx)
	if err != nil {
		panic(fmt.Errorf("error resetting whole validator rank: %+v", err))
	}
}
