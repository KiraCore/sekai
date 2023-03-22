package layer2

import (
	kiratypes "github.com/KiraCore/sekai/types"
	govtypes "github.com/KiraCore/sekai/x/gov/types"
	"github.com/KiraCore/sekai/x/layer2/keeper"
	"github.com/KiraCore/sekai/x/layer2/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type ApplyJoinDappProposalHandler struct {
	keeper keeper.Keeper
}

func NewApplyJoinDappProposalHandler(keeper keeper.Keeper) *ApplyJoinDappProposalHandler {
	return &ApplyJoinDappProposalHandler{
		keeper: keeper,
	}
}

func (a ApplyJoinDappProposalHandler) ProposalType() string {
	return kiratypes.ProposalTypeJoinDapp
}

func (a ApplyJoinDappProposalHandler) IsAllowedAddress(ctx sdk.Context, addr sdk.AccAddress, proposal govtypes.Content) bool {
	p := proposal.(*types.ProposalJoinDapp)

	dapp := a.keeper.GetDapp(ctx, p.DappName)
	if dapp.Name == "" {
		return false
	}

	return a.keeper.IsAllowedAddress(ctx, addr, dapp.Controllers)
}

func (a ApplyJoinDappProposalHandler) AllowedAddresses(ctx sdk.Context, proposal govtypes.Content) []string {
	p := proposal.(*types.ProposalJoinDapp)

	dapp := a.keeper.GetDapp(ctx, p.DappName)
	if dapp.Name == "" {
		return []string{}
	}
	return a.keeper.AllowedAddresses(ctx, dapp.Controllers)
}

func (a ApplyJoinDappProposalHandler) Quorum(ctx sdk.Context, proposal govtypes.Content) uint64 {
	p := proposal.(*types.ProposalJoinDapp)

	dapp := a.keeper.GetDapp(ctx, p.DappName)
	if dapp.Name == "" {
		return 0
	}

	return dapp.VoteQuorum
}

func (a ApplyJoinDappProposalHandler) VotePeriod(ctx sdk.Context, proposal govtypes.Content) uint64 {
	p := proposal.(*types.ProposalJoinDapp)

	dapp := a.keeper.GetDapp(ctx, p.DappName)
	if dapp.Name == "" {
		return 0
	}

	return dapp.VotePeriod
}

func (a ApplyJoinDappProposalHandler) VoteEnactment(ctx sdk.Context, proposal govtypes.Content) uint64 {
	p := proposal.(*types.ProposalJoinDapp)

	dapp := a.keeper.GetDapp(ctx, p.DappName)
	if dapp.Name == "" {
		return 0
	}

	return dapp.VoteEnactment
}

func (a ApplyJoinDappProposalHandler) Apply(ctx sdk.Context, proposalID uint64, proposal govtypes.Content, slash sdk.Dec) error {
	p := proposal.(*types.ProposalJoinDapp)

	a.keeper.SetDappOperator(ctx, types.DappOperator{
		DappName: p.DappName,
		Operator: p.Sender,
		Executor: p.Executor,
		Verifier: p.Verifier,
		Interx:   p.Interx,
		Status:   types.OperatorDeactivatived,
	})
	return nil
}

type ApplyTransitionDappProposalHandler struct {
	keeper keeper.Keeper
}

func NewApplyTransitionDappProposalHandler(keeper keeper.Keeper) *ApplyTransitionDappProposalHandler {
	return &ApplyTransitionDappProposalHandler{
		keeper: keeper,
	}
}

func (a ApplyTransitionDappProposalHandler) ProposalType() string {
	return kiratypes.ProposalTypeTransitionDapp
}

func (a ApplyTransitionDappProposalHandler) IsAllowedAddress(ctx sdk.Context, addr sdk.AccAddress, proposal govtypes.Content) bool {
	p := proposal.(*types.ProposalTransitionDapp)

	dapp := a.keeper.GetDapp(ctx, p.DappName)
	if dapp.Name == "" {
		return false
	}
	// TODO: probably, dapp transition will need to use raw Messages
	return a.keeper.IsAllowedAddress(ctx, addr, dapp.Controllers)
}

func (a ApplyTransitionDappProposalHandler) AllowedAddresses(ctx sdk.Context, proposal govtypes.Content) []string {
	p := proposal.(*types.ProposalTransitionDapp)

	operators := a.keeper.GetDappOperators(ctx, p.DappName)
	addrs := []string{}
	for _, operator := range operators {
		addrs = append(addrs, operator.Operator)
	}
	return addrs
}

func (a ApplyTransitionDappProposalHandler) Quorum(ctx sdk.Context, proposal govtypes.Content) uint64 {
	p := proposal.(*types.ProposalTransitionDapp)

	dapp := a.keeper.GetDapp(ctx, p.DappName)
	if dapp.Name == "" {
		return 0
	}

	// TODO: probably, dapp transition will need to use raw Messages
	return dapp.VerifiersMin
}

func (a ApplyTransitionDappProposalHandler) VotePeriod(ctx sdk.Context, proposal govtypes.Content) uint64 {
	p := proposal.(*types.ProposalTransitionDapp)

	dapp := a.keeper.GetDapp(ctx, p.DappName)
	if dapp.Name == "" {
		return 0
	}

	// TODO: probably, dapp transition will need to use raw Messages
	return dapp.VotePeriod
}

func (a ApplyTransitionDappProposalHandler) VoteEnactment(ctx sdk.Context, proposal govtypes.Content) uint64 {
	p := proposal.(*types.ProposalTransitionDapp)

	dapp := a.keeper.GetDapp(ctx, p.DappName)
	if dapp.Name == "" {
		return 0
	}

	// TODO: probably, dapp transition will need to use raw Messages
	return dapp.VoteEnactment
}

func (a ApplyTransitionDappProposalHandler) Apply(ctx sdk.Context, proposalID uint64, proposal govtypes.Content, slash sdk.Dec) error {
	p := proposal.(*types.ProposalTransitionDapp)

	dapp := a.keeper.GetDapp(ctx, p.DappName)
	if dapp.Name == "" {
		return types.ErrDappDoesNotExist
	}

	dapp.StatusHash = p.StatusHash
	a.keeper.SetDapp(ctx, dapp)
	// TODO: probably, dapp transition will need to use raw Messages
	return nil
}

type ApplyUpsertDappProposalHandler struct {
	keeper keeper.Keeper
}

func NewApplyUpsertDappProposalHandler(keeper keeper.Keeper) *ApplyUpsertDappProposalHandler {
	return &ApplyUpsertDappProposalHandler{
		keeper: keeper,
	}
}

func (a ApplyUpsertDappProposalHandler) ProposalType() string {
	return kiratypes.ProposalTypeUpsertDapp
}

func (a ApplyUpsertDappProposalHandler) IsAllowedAddress(ctx sdk.Context, addr sdk.AccAddress, proposal govtypes.Content) bool {
	p := proposal.(*types.ProposalUpsertDapp)

	dapp := a.keeper.GetDapp(ctx, p.Dapp.Name)
	if dapp.Name == "" {
		return false
	}

	return a.keeper.IsAllowedAddress(ctx, addr, dapp.Controllers)
}

func (a ApplyUpsertDappProposalHandler) AllowedAddresses(ctx sdk.Context, proposal govtypes.Content) []string {
	p := proposal.(*types.ProposalUpsertDapp)

	dapp := a.keeper.GetDapp(ctx, p.Dapp.Name)
	if dapp.Name == "" {
		return []string{}
	}

	return a.keeper.AllowedAddresses(ctx, dapp.Controllers)
}

func (a ApplyUpsertDappProposalHandler) Quorum(ctx sdk.Context, proposal govtypes.Content) uint64 {
	p := proposal.(*types.ProposalUpsertDapp)

	dapp := a.keeper.GetDapp(ctx, p.Dapp.Name)
	if dapp.Name == "" {
		return 0
	}

	return dapp.VoteQuorum
}

func (a ApplyUpsertDappProposalHandler) VotePeriod(ctx sdk.Context, proposal govtypes.Content) uint64 {
	p := proposal.(*types.ProposalUpsertDapp)

	dapp := a.keeper.GetDapp(ctx, p.Dapp.Name)
	if dapp.Name == "" {
		return 0
	}

	return dapp.VotePeriod
}

func (a ApplyUpsertDappProposalHandler) VoteEnactment(ctx sdk.Context, proposal govtypes.Content) uint64 {
	p := proposal.(*types.ProposalUpsertDapp)

	dapp := a.keeper.GetDapp(ctx, p.Dapp.Name)
	if dapp.Name == "" {
		return 0
	}

	return dapp.VoteEnactment
}

func (a ApplyUpsertDappProposalHandler) Apply(ctx sdk.Context, proposalID uint64, proposal govtypes.Content, slash sdk.Dec) error {
	p := proposal.(*types.ProposalUpsertDapp)

	dapp := a.keeper.GetDapp(ctx, p.Dapp.Name)
	if dapp.Name == "" {
		return types.ErrDappDoesNotExist
	}

	dapp.Name = p.Dapp.Name
	dapp.Denom = p.Dapp.Denom
	dapp.Description = p.Dapp.Description
	dapp.Website = p.Dapp.Website
	dapp.Logo = p.Dapp.Logo
	dapp.Social = p.Dapp.Social
	dapp.Docs = p.Dapp.Docs
	dapp.Controllers = p.Dapp.Controllers
	dapp.Bin = p.Dapp.Bin
	dapp.Pool = p.Dapp.Pool
	dapp.Issurance = p.Dapp.Issurance
	dapp.UpdateTimeMax = p.Dapp.UpdateTimeMax
	dapp.ExecutorsMin = p.Dapp.ExecutorsMin
	dapp.ExecutorsMax = p.Dapp.ExecutorsMax
	dapp.VerifiersMin = p.Dapp.VerifiersMin
	dapp.TotalBond = p.Dapp.TotalBond
	dapp.CreationTime = p.Dapp.CreationTime
	dapp.StatusHash = p.Dapp.StatusHash
	dapp.Status = p.Dapp.Status
	dapp.VoteQuorum = p.Dapp.VoteQuorum
	dapp.VotePeriod = p.Dapp.VotePeriod
	dapp.VoteEnactment = p.Dapp.VoteEnactment

	a.keeper.SetDapp(ctx, p.Dapp)
	return nil
}
