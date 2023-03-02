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

func (a ApplyCreateBasketProposalHandler) Apply(ctx sdk.Context, proposalID uint64, proposal govtypes.Content, slash sdk.Dec) error {
	p := proposal.(*types.ProposalCreateBasket)
	return a.keeper.CreateBasket(ctx, p.Basket)
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

func (a ApplyEditBasketProposalHandler) Apply(ctx sdk.Context, proposalID uint64, proposal govtypes.Content, slash sdk.Dec) error {
	p := proposal.(*types.ProposalEditBasket)

	return a.keeper.EditBasket(ctx, p.Basket)
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

func (a ApplyBasketWithdrawSurplusProposalHandler) Apply(ctx sdk.Context, proposalID uint64, proposal govtypes.Content, slash sdk.Dec) error {
	p := proposal.(*types.ProposalBasketWithdrawSurplus)
	return a.keeper.BasketWithdrawSurplus(ctx, *p)
}
