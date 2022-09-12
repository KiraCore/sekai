package basket

import (
	kiratypes "github.com/KiraCore/sekai/types"
	"github.com/KiraCore/sekai/x/basket/keeper"
	"github.com/KiraCore/sekai/x/basket/types"
	govtypes "github.com/KiraCore/sekai/x/gov/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type ApplyCreateBasketProposalHandler struct {
	keeper keeper.Keeper
}

func NewApplyCreateBasketProposalHandler(keeper keeper.Keeper) *ApplyCreateBasketProposalHandler {
	return &ApplyCreateBasketProposalHandler{
		keeper: keeper,
	}
}

func (a ApplyCreateBasketProposalHandler) ProposalType() string {
	return kiratypes.ProposalTypeCreateBasket
}

func (a ApplyCreateBasketProposalHandler) Apply(ctx sdk.Context, proposalID uint64, proposal govtypes.Content, slash uint64) error {
	p := proposal.(*types.ProposalCreateBasket)

	// TODO: create a new basket id
	// TODO: check suffix is duplicated
	// TODO: surpluse should be zero
	// TODO: use FlagSlippageFeeMin
	// TODO: use FlagTokensCap
	// TODO: use FlagLimitsPeriod
	// TODO: ensure basket.Tokens[i].Amount is zero
	// TODO: ensure denom is not empty
	// TODO: ensure weights are not zero for a denom
	// TODO: ensure denoms not duplicate within tokens list

	a.keeper.SetBasket(ctx, p.Basket)
	return nil
}

type ApplyEditBasketProposalHandler struct {
	keeper keeper.Keeper
}

func NewApplyEditBasketProposalHandler(keeper keeper.Keeper) *ApplyEditBasketProposalHandler {
	return &ApplyEditBasketProposalHandler{
		keeper: keeper,
	}
}

func (a ApplyEditBasketProposalHandler) ProposalType() string {
	return kiratypes.ProposalTypeEditBasket
}

func (a ApplyEditBasketProposalHandler) Apply(ctx sdk.Context, proposalID uint64, proposal govtypes.Content, slash uint64) error {
	p := proposal.(*types.ProposalEditBasket)

	_, err := a.keeper.GetBasketById(ctx, p.Basket.Id)
	if err != nil {
		return err
	}

	// TODO: check suffix is not changed
	// TODO: check id existance when editing
	// TODO: use previous surplus
	// TODO: basket tokens removal consideration
	// TODO: ensure basket.Tokens[i].Amount is not used
	// TODO: ensure denom is not empty
	// TODO: ensure weights are not zero for a denom
	// TODO: ensure denoms not duplicate within tokens list

	a.keeper.SetBasket(ctx, p.Basket)
	return nil
}

type ApplyBasketWithdrawSurplusProposalHandler struct {
	keeper keeper.Keeper
}

func NewApplyBasketWithdrawSurplusProposalHandler(keeper keeper.Keeper) *ApplyBasketWithdrawSurplusProposalHandler {
	return &ApplyBasketWithdrawSurplusProposalHandler{
		keeper: keeper,
	}
}

func (a ApplyBasketWithdrawSurplusProposalHandler) ProposalType() string {
	return kiratypes.ProposalTypeBasketWithdrawSurplus
}

func (a ApplyBasketWithdrawSurplusProposalHandler) Apply(ctx sdk.Context, proposalID uint64, proposal govtypes.Content, slash uint64) error {
	p := proposal.(*types.ProposalBasketWithdrawSurplus)
	return a.keeper.BasketWithdrawSurplus(ctx, *p)
}
