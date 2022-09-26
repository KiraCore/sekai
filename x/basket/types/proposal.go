package types

import (
	kiratypes "github.com/KiraCore/sekai/types"
	"github.com/KiraCore/sekai/x/gov/types"
)

func NewProposalCreateBasket(basket Basket) *ProposalCreateBasket {
	return &ProposalCreateBasket{
		Basket: basket,
	}
}

func (m *ProposalCreateBasket) ProposalType() string {
	return kiratypes.ProposalTypeCreateBasket
}

func (m *ProposalCreateBasket) ProposalPermission() types.PermValue {
	return types.PermCreateBasketProposal
}

func (m *ProposalCreateBasket) VotePermission() types.PermValue {
	return types.PermVoteBasketProposal
}

// ValidateBasic returns basic validation
func (m *ProposalCreateBasket) ValidateBasic() error {
	return nil
}

func NewProposalEditBasket(basket Basket) *ProposalEditBasket {
	return &ProposalEditBasket{
		Basket: basket,
	}
}

func (m *ProposalEditBasket) ProposalType() string {
	return kiratypes.ProposalTypeEditBasket
}

func (m *ProposalEditBasket) ProposalPermission() types.PermValue {
	return types.PermCreateBasketProposal
}

func (m *ProposalEditBasket) VotePermission() types.PermValue {
	return types.PermVoteBasketProposal
}

// ValidateBasic returns basic validation
func (m *ProposalEditBasket) ValidateBasic() error {
	return nil
}

func NewProposalBasketWithdrawSurplus(basketIds []uint64, withdrawTarget string) *ProposalBasketWithdrawSurplus {
	return &ProposalBasketWithdrawSurplus{
		BasketIds:      basketIds,
		WithdrawTarget: withdrawTarget,
	}
}

func (m *ProposalBasketWithdrawSurplus) ProposalType() string {
	return kiratypes.ProposalTypeBasketWithdrawSurplus
}

func (m *ProposalBasketWithdrawSurplus) ProposalPermission() types.PermValue {
	return types.PermCreateBasketProposal
}

func (m *ProposalBasketWithdrawSurplus) VotePermission() types.PermValue {
	return types.PermVoteBasketProposal
}

// ValidateBasic returns basic validation
func (m *ProposalBasketWithdrawSurplus) ValidateBasic() error {
	return nil
}
