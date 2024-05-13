package collectives

import (
	kiratypes "github.com/KiraCore/sekai/types"
	"github.com/KiraCore/sekai/x/collectives/keeper"
	"github.com/KiraCore/sekai/x/collectives/types"
	govtypes "github.com/KiraCore/sekai/x/gov/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type ApplyCollectiveSendDonationProposalHandler struct {
	keeper keeper.Keeper
}

func NewApplyCollectiveSendDonationProposalHandler(keeper keeper.Keeper) *ApplyCollectiveSendDonationProposalHandler {
	return &ApplyCollectiveSendDonationProposalHandler{
		keeper: keeper,
	}
}

func (a ApplyCollectiveSendDonationProposalHandler) ProposalType() string {
	return kiratypes.ProposalTypeCollectiveSendDonation
}

func (a ApplyCollectiveSendDonationProposalHandler) IsAllowedAddress(ctx sdk.Context, addr sdk.AccAddress, proposal govtypes.Content) bool {
	p := proposal.(*types.ProposalCollectiveSendDonation)

	collective := a.keeper.GetCollective(ctx, p.Name)
	if collective.Name == "" {
		return false
	}

	return a.keeper.IsAllowedAddress(ctx, addr, collective.OwnersWhitelist)
}

func (a ApplyCollectiveSendDonationProposalHandler) AllowedAddresses(ctx sdk.Context, proposal govtypes.Content) []string {
	p := proposal.(*types.ProposalCollectiveSendDonation)

	collective := a.keeper.GetCollective(ctx, p.Name)
	if collective.Name == "" {
		return []string{}
	}

	return a.keeper.AllowedAddresses(ctx, collective.OwnersWhitelist)
}

func (a ApplyCollectiveSendDonationProposalHandler) Quorum(ctx sdk.Context, proposal govtypes.Content) sdk.Dec {
	p := proposal.(*types.ProposalCollectiveSendDonation)

	collective := a.keeper.GetCollective(ctx, p.Name)
	if collective.Name == "" {
		return sdk.ZeroDec()
	}

	return collective.VoteQuorum
}

func (a ApplyCollectiveSendDonationProposalHandler) VotePeriod(ctx sdk.Context, proposal govtypes.Content) uint64 {
	p := proposal.(*types.ProposalCollectiveSendDonation)

	collective := a.keeper.GetCollective(ctx, p.Name)
	if collective.Name == "" {
		return 0
	}

	return collective.VotePeriod
}

func (a ApplyCollectiveSendDonationProposalHandler) VoteEnactment(ctx sdk.Context, proposal govtypes.Content) uint64 {
	p := proposal.(*types.ProposalCollectiveSendDonation)

	collective := a.keeper.GetCollective(ctx, p.Name)
	if collective.Name == "" {
		return 0
	}

	return collective.VoteEnactment
}

func (a ApplyCollectiveSendDonationProposalHandler) Apply(ctx sdk.Context, proposalID uint64, proposal govtypes.Content, slash sdk.Dec) error {
	p := proposal.(*types.ProposalCollectiveSendDonation)

	addr, err := sdk.AccAddressFromBech32(p.Address)
	if err != nil {
		return err
	}
	return a.keeper.SendDonation(ctx, p.Name, addr, p.Amounts)
}

type ApplyCollectiveUpdateProposalHandler struct {
	keeper keeper.Keeper
}

func NewApplyCollectiveUpdateProposalHandler(keeper keeper.Keeper) *ApplyCollectiveUpdateProposalHandler {
	return &ApplyCollectiveUpdateProposalHandler{
		keeper: keeper,
	}
}

func (a ApplyCollectiveUpdateProposalHandler) ProposalType() string {
	return kiratypes.ProposalTypeCollectiveUpdate
}

func (a ApplyCollectiveUpdateProposalHandler) IsAllowedAddress(ctx sdk.Context, addr sdk.AccAddress, proposal govtypes.Content) bool {
	p := proposal.(*types.ProposalCollectiveUpdate)

	collective := a.keeper.GetCollective(ctx, p.Name)
	if collective.Name == "" {
		return false
	}

	return a.keeper.IsAllowedAddress(ctx, addr, collective.OwnersWhitelist)
}

func (a ApplyCollectiveUpdateProposalHandler) AllowedAddresses(ctx sdk.Context, proposal govtypes.Content) []string {
	p := proposal.(*types.ProposalCollectiveUpdate)

	collective := a.keeper.GetCollective(ctx, p.Name)
	if collective.Name == "" {
		return []string{}
	}

	return a.keeper.AllowedAddresses(ctx, collective.OwnersWhitelist)
}

func (a ApplyCollectiveUpdateProposalHandler) Quorum(ctx sdk.Context, proposal govtypes.Content) sdk.Dec {
	p := proposal.(*types.ProposalCollectiveUpdate)

	collective := a.keeper.GetCollective(ctx, p.Name)
	if collective.Name == "" {
		return sdk.ZeroDec()
	}

	return collective.VoteQuorum
}

func (a ApplyCollectiveUpdateProposalHandler) VotePeriod(ctx sdk.Context, proposal govtypes.Content) uint64 {
	p := proposal.(*types.ProposalCollectiveUpdate)

	collective := a.keeper.GetCollective(ctx, p.Name)
	if collective.Name == "" {
		return 0
	}

	return collective.VotePeriod
}

func (a ApplyCollectiveUpdateProposalHandler) VoteEnactment(ctx sdk.Context, proposal govtypes.Content) uint64 {
	p := proposal.(*types.ProposalCollectiveUpdate)

	collective := a.keeper.GetCollective(ctx, p.Name)
	if collective.Name == "" {
		return 0
	}

	return collective.VoteEnactment
}

func (a ApplyCollectiveUpdateProposalHandler) Apply(ctx sdk.Context, proposalID uint64, proposal govtypes.Content, slash sdk.Dec) error {
	p := proposal.(*types.ProposalCollectiveUpdate)

	collective := a.keeper.GetCollective(ctx, p.Name)
	if collective.Name == "" {
		return types.ErrCollectiveDoesNotExist
	}

	collective.Description = p.Description
	collective.Status = p.Status
	collective.DepositWhitelist = p.DepositWhitelist
	collective.OwnersWhitelist = p.OwnersWhitelist
	collective.SpendingPools = p.SpendingPools
	collective.ClaimStart = p.ClaimStart
	collective.ClaimPeriod = p.ClaimPeriod
	collective.ClaimEnd = p.ClaimEnd
	collective.VoteQuorum = p.VoteQuorum
	collective.VotePeriod = p.VotePeriod
	collective.VoteEnactment = p.VoteEnactment

	a.keeper.SetCollective(ctx, collective)
	return nil
}

type ApplyCollectiveRemoveProposalHandler struct {
	keeper keeper.Keeper
}

func NewApplyCollectiveRemoveProposalHandler(keeper keeper.Keeper) *ApplyCollectiveRemoveProposalHandler {
	return &ApplyCollectiveRemoveProposalHandler{
		keeper: keeper,
	}
}

func (a ApplyCollectiveRemoveProposalHandler) ProposalType() string {
	return kiratypes.ProposalTypeCollectiveRemove
}

func (a ApplyCollectiveRemoveProposalHandler) IsAllowedAddress(ctx sdk.Context, addr sdk.AccAddress, proposal govtypes.Content) bool {
	p := proposal.(*types.ProposalCollectiveRemove)

	collective := a.keeper.GetCollective(ctx, p.Name)
	if collective.Name == "" {
		return false
	}

	return a.keeper.IsAllowedAddress(ctx, addr, collective.OwnersWhitelist)
}

func (a ApplyCollectiveRemoveProposalHandler) AllowedAddresses(ctx sdk.Context, proposal govtypes.Content) []string {
	p := proposal.(*types.ProposalCollectiveRemove)

	collective := a.keeper.GetCollective(ctx, p.Name)
	if collective.Name == "" {
		return []string{}
	}

	return a.keeper.AllowedAddresses(ctx, collective.OwnersWhitelist)
}

func (a ApplyCollectiveRemoveProposalHandler) Quorum(ctx sdk.Context, proposal govtypes.Content) sdk.Dec {
	p := proposal.(*types.ProposalCollectiveRemove)

	collective := a.keeper.GetCollective(ctx, p.Name)
	if collective.Name == "" {
		return sdk.ZeroDec()
	}

	return collective.VoteQuorum
}

func (a ApplyCollectiveRemoveProposalHandler) VotePeriod(ctx sdk.Context, proposal govtypes.Content) uint64 {
	p := proposal.(*types.ProposalCollectiveRemove)

	collective := a.keeper.GetCollective(ctx, p.Name)
	if collective.Name == "" {
		return 0
	}

	return collective.VotePeriod
}

func (a ApplyCollectiveRemoveProposalHandler) VoteEnactment(ctx sdk.Context, proposal govtypes.Content) uint64 {
	p := proposal.(*types.ProposalCollectiveRemove)

	collective := a.keeper.GetCollective(ctx, p.Name)
	if collective.Name == "" {
		return 0
	}

	return collective.VoteEnactment
}

func (a ApplyCollectiveRemoveProposalHandler) Apply(ctx sdk.Context, proposalID uint64, proposal govtypes.Content, slash sdk.Dec) error {
	p := proposal.(*types.ProposalCollectiveRemove)
	collective := a.keeper.GetCollective(ctx, p.Name)
	if collective.Name == "" {
		return types.ErrCollectiveDoesNotExist
	}

	a.keeper.ExecuteCollectiveRemove(ctx, collective)
	return nil
}
