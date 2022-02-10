package tokens

import (
	kiratypes "github.com/KiraCore/sekai/types"
	"github.com/KiraCore/sekai/x/gov/types"
	"github.com/KiraCore/sekai/x/spending/keeper"
	spendingtypes "github.com/KiraCore/sekai/x/spending/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type ApplyUpdateSpendingPoolProposalHandler struct {
	keeper keeper.Keeper
}

func NewApplyUpdateSpendingPoolProposalHandler(keeper keeper.Keeper) *ApplyUpdateSpendingPoolProposalHandler {
	return &ApplyUpdateSpendingPoolProposalHandler{
		keeper: keeper,
	}
}

func (a ApplyUpdateSpendingPoolProposalHandler) ProposalType() string {
	return kiratypes.ProposalTypeUpdateSpendingPool
}

func (a ApplyUpdateSpendingPoolProposalHandler) Apply(ctx sdk.Context, proposalID uint64, proposal types.Content) error {
	p := proposal.(*spendingtypes.ProposalUpsertTokenAlias)

	tokenAlians := spendingtypes.NewTokenAlias(p.Symbol, p.Name, p.Icon, p.Decimals, p.Denoms)
	return a.keeper.UpsertTokenAlias(ctx, *tokenAlians)
}

type ApplySpendingPoolDistributionProposalHandler struct {
	keeper keeper.Keeper
}

func NewApplyUpdateSpendingPoolProposalHandler(keeper keeper.Keeper) *ApplySpendingPoolDistributionProposalHandler {
	return &ApplySpendingPoolDistributionProposalHandler{
		keeper: keeper,
	}
}

func (a ApplySpendingPoolDistributionProposalHandler) ProposalType() string {
	return kiratypes.ProposalTypeSpendingPoolDistribution
}

func (a ApplySpendingPoolDistributionProposalHandler) Apply(ctx sdk.Context, proposalID uint64, proposal types.Content) error {
	p := proposal.(*spendingtypes.ProposalUpsertTokenAlias)

	tokenAlians := spendingtypes.NewTokenAlias(p.Symbol, p.Name, p.Icon, p.Decimals, p.Denoms)
	return a.keeper.UpsertTokenAlias(ctx, *tokenAlians)
}

type ApplySpendingPoolWithdrawProposalHandler struct {
	keeper keeper.Keeper
}

func NewApplyUpdateSpendingPoolProposalHandler(keeper keeper.Keeper) *ApplySpendingPoolWithdrawProposalHandler {
	return &ApplySpendingPoolWithdrawProposalHandler{
		keeper: keeper,
	}
}

func (a ApplySpendingPoolWithdrawProposalHandler) ProposalType() string {
	return kiratypes.ProposalTypeSpendingPoolWithdraw
}

func (a ApplySpendingPoolWithdrawProposalHandler) Apply(ctx sdk.Context, proposalID uint64, proposal types.Content) error {
	p := proposal.(*spendingtypes.ProposalUpsertTokenAlias)

	tokenAlians := spendingtypes.NewTokenAlias(p.Symbol, p.Name, p.Icon, p.Decimals, p.Denoms)
	return a.keeper.UpsertTokenAlias(ctx, *tokenAlians)
}
