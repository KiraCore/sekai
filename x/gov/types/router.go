package types

import (
	"fmt"

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

func (r ProposalRouter) ApplyProposal(ctx sdk.Context, proposal Content) error {
	h, ok := r.routes[proposal.ProposalType()]
	if !ok {
		panic("invalid proposal type")
	}

	cachedCtx, writeCache := ctx.CacheContext()
	err := h.Apply(cachedCtx, proposal)
	if err == nil {
		writeCache()
	} else { // not halt the chain for proposal execution
		fmt.Println("error applying proposal:", err)
	}
	return err
}

type ProposalHandler interface {
	ProposalType() string
	Apply(ctx sdk.Context, proposal Content) error
}
