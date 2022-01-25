package tokens

import (
	kiratypes "github.com/KiraCore/sekai/types"
	"github.com/KiraCore/sekai/x/gov/types"
	"github.com/KiraCore/sekai/x/spending/keeper"
	tokenstypes "github.com/KiraCore/sekai/x/spending/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type ApplyUpsertTokenAliasProposalHandler struct {
	keeper keeper.Keeper
}

func NewApplyUpsertTokenAliasProposalHandler(keeper keeper.Keeper) *ApplyUpsertTokenAliasProposalHandler {
	return &ApplyUpsertTokenAliasProposalHandler{
		keeper: keeper,
	}
}

func (a ApplyUpsertTokenAliasProposalHandler) ProposalType() string {
	return kiratypes.ProposalTypeUpsertTokenAlias
}

func (a ApplyUpsertTokenAliasProposalHandler) Apply(ctx sdk.Context, proposalID uint64, proposal types.Content) error {
	p := proposal.(*tokenstypes.ProposalUpsertTokenAlias)

	tokenAlians := tokenstypes.NewTokenAlias(p.Symbol, p.Name, p.Icon, p.Decimals, p.Denoms)
	return a.keeper.UpsertTokenAlias(ctx, *tokenAlians)
}
