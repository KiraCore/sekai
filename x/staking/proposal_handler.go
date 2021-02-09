package staking

import (
	"github.com/KiraCore/sekai/x/gov/types"
	"github.com/KiraCore/sekai/x/staking/keeper"
	types3 "github.com/KiraCore/sekai/x/staking/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type ApplyUnjailValidatorProposalHandler struct {
	keeper keeper.Keeper
}

func NewApplyUnjailValidatorProposalHandler(keeper keeper.Keeper) *ApplyUnjailValidatorProposalHandler {
	return &ApplyUnjailValidatorProposalHandler{
		keeper: keeper,
	}
}

func (a ApplyUnjailValidatorProposalHandler) ProposalType() string {
	return types3.ProposalTypeUnjailValidator
}

func (a ApplyUnjailValidatorProposalHandler) Apply(ctx sdk.Context, proposal types.Content) {
	p := proposal.(*types3.ProposalUnjailValidator)

	err := a.keeper.Unjail(ctx, sdk.ValAddress(p.Proposer))
	if err != nil {
		panic("error unjailing")
	}
}
