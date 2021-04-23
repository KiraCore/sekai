package staking

import (
	govtypes "github.com/KiraCore/sekai/x/gov/types"
	"github.com/KiraCore/sekai/x/staking/keeper"
	"github.com/KiraCore/sekai/x/staking/types"
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
	return types.ProposalTypeUnjailValidator
}

func (a ApplyUnjailValidatorProposalHandler) Apply(ctx sdk.Context, proposal govtypes.Content) {
	p := proposal.(*types.ProposalUnjailValidator)

	err := a.keeper.Unjail(ctx, sdk.ValAddress(p.Proposer))
	if err != nil {
		panic("error unjailing")
	}
}
