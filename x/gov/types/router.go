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

func (r ProposalRouter) IsAllowedAddressDynamicProposal(ctx sdk.Context, addr sdk.AccAddress, proposal Content) bool {
	h, ok := r.routes[proposal.ProposalType()]
	if !ok {
		panic("invalid proposal type")
	}

	dh, ok := h.(DynamicVoterProposalHandler)
	if !ok {
		return false
	}
	return dh.IsAllowedAddress(ctx, addr, proposal)
}

func (r ProposalRouter) AllowedAddressesDynamicProposal(ctx sdk.Context, proposal Content) []string {
	h, ok := r.routes[proposal.ProposalType()]
	if !ok {
		panic("invalid proposal type")
	}

	dh, ok := h.(DynamicVoterProposalHandler)
	if !ok {
		return []string{}
	}
	return dh.AllowedAddresses(ctx, proposal)
}

func (r ProposalRouter) QuorumDynamicProposal(ctx sdk.Context, proposal Content) sdk.Dec {
	h, ok := r.routes[proposal.ProposalType()]
	if !ok {
		panic("invalid proposal type")
	}

	dh, ok := h.(DynamicVoterProposalHandler)
	if !ok {
		return sdk.ZeroDec()
	}
	return dh.Quorum(ctx, proposal)
}

func (r ProposalRouter) VotePeriodDynamicProposal(ctx sdk.Context, proposal Content) uint64 {
	h, ok := r.routes[proposal.ProposalType()]
	if !ok {
		panic("invalid proposal type")
	}

	dh, ok := h.(DynamicVoterProposalHandler)
	if !ok {
		return 0
	}
	return dh.VotePeriod(ctx, proposal)
}

func (r ProposalRouter) EnactmentPeriodDynamicProposal(ctx sdk.Context, proposal Content) uint64 {
	h, ok := r.routes[proposal.ProposalType()]
	if !ok {
		panic("invalid proposal type")
	}

	dh, ok := h.(DynamicVoterProposalHandler)
	if !ok {
		return 0
	}
	return dh.VoteEnactment(ctx, proposal)
}

func (r ProposalRouter) ApplyProposal(ctx sdk.Context, proposalID uint64, proposal Content, slash sdk.Dec) error {
	h, ok := r.routes[proposal.ProposalType()]
	if !ok {
		panic("invalid proposal type")
	}

	cachedCtx, writeCache := ctx.CacheContext()
	err := h.Apply(cachedCtx, proposalID, proposal, slash)
	if err == nil {
		writeCache()
	} else { // not halt the chain for proposal execution
		fmt.Println("error applying proposal:", err)
	}
	return err
}

type ProposalHandler interface {
	ProposalType() string
	Apply(ctx sdk.Context, proposalID uint64, proposal Content, slash sdk.Dec) error
}

type DynamicVoterProposalHandler interface {
	ProposalType() string
	IsAllowedAddress(ctx sdk.Context, addr sdk.AccAddress, proposal Content) bool
	Quorum(ctx sdk.Context, proposal Content) sdk.Dec
	VotePeriod(ctx sdk.Context, proposal Content) uint64
	VoteEnactment(ctx sdk.Context, proposal Content) uint64
	AllowedAddresses(ctx sdk.Context, proposal Content) []string
	Apply(ctx sdk.Context, proposalID uint64, proposal Content, slash sdk.Dec) error
}
