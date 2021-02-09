package tokens

import (
	"github.com/KiraCore/sekai/x/gov/types"
	"github.com/KiraCore/sekai/x/tokens/keeper"
	tokenstypes "github.com/KiraCore/sekai/x/tokens/types"
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
	return tokenstypes.ProposalTypeUpsertTokenAlias
}

func (a ApplyUpsertTokenAliasProposalHandler) Apply(ctx sdk.Context, proposal types.Content) {
	p := proposal.(*tokenstypes.ProposalUpsertTokenAlias)

	tokenAlians := tokenstypes.NewTokenAlias(p.Symbol, p.Name, p.Icon, p.Decimals, p.Denoms)
	err := a.keeper.UpsertTokenAlias(ctx, *tokenAlians)
	if err != nil {
		panic(err)
	}
}

type ApplyUpsertTokenRatesProposalHandler struct {
	keeper keeper.Keeper
}

func NewApplyUpsertTokenRatesProposalHandler(keeper keeper.Keeper) *ApplyUpsertTokenRatesProposalHandler {
	return &ApplyUpsertTokenRatesProposalHandler{keeper: keeper}
}

func (a ApplyUpsertTokenRatesProposalHandler) ProposalType() string {
	return tokenstypes.ProposalTypeUpsertTokenRates
}

func (a ApplyUpsertTokenRatesProposalHandler) Apply(ctx sdk.Context, proposal types.Content) {
	p := proposal.(*tokenstypes.ProposalUpsertTokenRates)

	tokenAlians := tokenstypes.NewTokenRate(p.Denom, p.Rate, p.FeePayments)
	err := a.keeper.UpsertTokenRate(ctx, *tokenAlians)
	if err != nil {
		panic(err)
	}
}
