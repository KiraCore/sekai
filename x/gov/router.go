package gov

import (
	"github.com/KiraCore/sekai/x/gov/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type ProposalRouter struct {
	routes map[string]ProposalHandler
}

func NewProposalRouter(handlers []ProposalHandler) ProposalRouter {
	routes := make(map[string]ProposalHandler, len(handlers))
	for _, h := range handlers {
		routes[h.ProposalType()] = h
	}

	return ProposalRouter{routes: routes}
}

func (r ProposalRouter) ApplyProposal(ctx sdk.Context, proposal types.Content) {
	h, ok := r.routes[proposal.ProposalType()]
	if !ok {
		panic("invalid proposal type")
	}

	h.Apply(ctx, proposal)
}

type ProposalHandler interface {
	ProposalType() string
	Apply(ctx sdk.Context, proposal types.Content)
}
