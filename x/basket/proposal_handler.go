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

	_, err := a.keeper.GetBasketById(ctx, p.BasketId)
	if err != nil {
		return err
	}

	// TODO: implement
	// a.keeper.SetBasket(ctx, p.Basket)
	return nil
}
