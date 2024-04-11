package spending

import (
	kiratypes "github.com/KiraCore/sekai/types"
	govtypes "github.com/KiraCore/sekai/x/gov/types"
	"github.com/KiraCore/sekai/x/spending/keeper"
	"github.com/KiraCore/sekai/x/spending/types"
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

func (a ApplyUpdateSpendingPoolProposalHandler) IsAllowedAddress(ctx sdk.Context, addr sdk.AccAddress, proposal govtypes.Content) bool {
	p := proposal.(*types.UpdateSpendingPoolProposal)

	pool := a.keeper.GetSpendingPool(ctx, p.Name)
	if pool == nil {
		return false
	}

	return a.keeper.IsAllowedAddress(ctx, addr, *pool.Owners)
}

func (a ApplyUpdateSpendingPoolProposalHandler) AllowedAddresses(ctx sdk.Context, proposal govtypes.Content) []string {
	p := proposal.(*types.UpdateSpendingPoolProposal)

	pool := a.keeper.GetSpendingPool(ctx, p.Name)
	if pool == nil {
		return []string{}
	}

	return a.keeper.AllowedAddresses(ctx, *pool.Owners)
}

func (a ApplyUpdateSpendingPoolProposalHandler) Quorum(ctx sdk.Context, proposal govtypes.Content) sdk.Dec {
	p := proposal.(*types.UpdateSpendingPoolProposal)

	pool := a.keeper.GetSpendingPool(ctx, p.Name)
	if pool == nil {
		return sdk.ZeroDec()
	}

	return pool.VoteQuorum
}

func (a ApplyUpdateSpendingPoolProposalHandler) VotePeriod(ctx sdk.Context, proposal govtypes.Content) uint64 {
	p := proposal.(*types.UpdateSpendingPoolProposal)

	pool := a.keeper.GetSpendingPool(ctx, p.Name)
	if pool == nil {
		return 0
	}

	return pool.VotePeriod
}

func (a ApplyUpdateSpendingPoolProposalHandler) VoteEnactment(ctx sdk.Context, proposal govtypes.Content) uint64 {
	p := proposal.(*types.UpdateSpendingPoolProposal)

	pool := a.keeper.GetSpendingPool(ctx, p.Name)
	if pool == nil {
		return 0
	}

	return pool.VoteEnactment
}

func (a ApplyUpdateSpendingPoolProposalHandler) Apply(ctx sdk.Context, proposalID uint64, proposal govtypes.Content, slash sdk.Dec) error {
	p := proposal.(*spendingtypes.UpdateSpendingPoolProposal)

	pool := a.keeper.GetSpendingPool(ctx, p.Name)
	if pool == nil {
		return types.ErrPoolDoesNotExist
	}

	a.keeper.SetSpendingPool(ctx, types.SpendingPool{
		Name:              p.Name,
		ClaimStart:        p.ClaimStart,
		ClaimEnd:          p.ClaimEnd,
		Rates:             p.Rates,
		VoteQuorum:        p.VoteQuorum,
		VotePeriod:        p.VotePeriod,
		VoteEnactment:     p.VoteEnactment,
		Owners:            &p.Owners,
		Beneficiaries:     &p.Beneficiaries,
		Balances:          pool.Balances,
		DynamicRate:       p.DynamicRate,
		DynamicRatePeriod: p.DynamicRatePeriod,
	})

	return nil
}

type ApplySpendingPoolDistributionProposalHandler struct {
	keeper keeper.Keeper
	gk     types.CustomGovKeeper
}

func NewApplySpendingPoolDistributionProposalHandler(keeper keeper.Keeper, gk types.CustomGovKeeper) *ApplySpendingPoolDistributionProposalHandler {
	return &ApplySpendingPoolDistributionProposalHandler{
		keeper: keeper,
		gk:     gk,
	}
}

func (a ApplySpendingPoolDistributionProposalHandler) ProposalType() string {
	return kiratypes.ProposalTypeSpendingPoolDistribution
}

func (a ApplySpendingPoolDistributionProposalHandler) IsAllowedAddress(ctx sdk.Context, addr sdk.AccAddress, proposal govtypes.Content) bool {
	p := proposal.(*types.SpendingPoolDistributionProposal)

	pool := a.keeper.GetSpendingPool(ctx, p.PoolName)
	if pool == nil {
		return false
	}

	return a.keeper.IsAllowedAddress(ctx, addr, *pool.Owners)
}

func (a ApplySpendingPoolDistributionProposalHandler) AllowedAddresses(ctx sdk.Context, proposal govtypes.Content) []string {
	p := proposal.(*types.SpendingPoolDistributionProposal)

	pool := a.keeper.GetSpendingPool(ctx, p.PoolName)
	if pool == nil {
		return []string{}
	}

	return a.keeper.AllowedAddresses(ctx, *pool.Owners)
}

func (a ApplySpendingPoolDistributionProposalHandler) Quorum(ctx sdk.Context, proposal govtypes.Content) sdk.Dec {
	p := proposal.(*types.SpendingPoolDistributionProposal)

	pool := a.keeper.GetSpendingPool(ctx, p.PoolName)
	if pool == nil {
		return sdk.ZeroDec()
	}

	return pool.VoteQuorum
}

func (a ApplySpendingPoolDistributionProposalHandler) VotePeriod(ctx sdk.Context, proposal govtypes.Content) uint64 {
	p := proposal.(*types.SpendingPoolDistributionProposal)

	pool := a.keeper.GetSpendingPool(ctx, p.PoolName)
	if pool == nil {
		return 0
	}

	return pool.VotePeriod
}

func (a ApplySpendingPoolDistributionProposalHandler) VoteEnactment(ctx sdk.Context, proposal govtypes.Content) uint64 {
	p := proposal.(*types.SpendingPoolDistributionProposal)

	pool := a.keeper.GetSpendingPool(ctx, p.PoolName)
	if pool == nil {
		return 0
	}

	return pool.VoteEnactment
}

func (a ApplySpendingPoolDistributionProposalHandler) Apply(ctx sdk.Context, proposalID uint64, proposal govtypes.Content, slash sdk.Dec) error {
	p := proposal.(*spendingtypes.SpendingPoolDistributionProposal)
	_ = p

	pool := a.keeper.GetSpendingPool(ctx, p.PoolName)
	duplicateMap := map[string]bool{}
	var beneficiaries []spendingtypes.WeightedAccount

	for _, acc := range pool.Beneficiaries.Accounts {
		if _, ok := duplicateMap[acc.Account]; !ok {
			duplicateMap[acc.Account] = true
			beneficiaries = append(beneficiaries, acc)
		}
	}
	for _, role := range pool.Beneficiaries.Roles {
		actorIter := a.gk.GetNetworkActorsByRole(ctx, role.Role)

		for ; actorIter.Valid(); actorIter.Next() {
			if _, ok := duplicateMap[sdk.AccAddress(actorIter.Value()).String()]; !ok {
				duplicateMap[sdk.AccAddress(actorIter.Value()).String()] = true
				beneficiaries = append(beneficiaries, spendingtypes.WeightedAccount{
					Account: sdk.AccAddress(actorIter.Value()).String(),
					Weight:  role.Weight,
				})
			}
		}
	}

	for _, beneficiary := range beneficiaries {
		beneficiaryAcc, err := sdk.AccAddressFromBech32(beneficiary.Account)
		if err != nil {
			return err
		}

		err = a.keeper.ClaimSpendingPool(ctx, p.PoolName, beneficiaryAcc)
		if err != nil {
			return err
		}
	}

	return nil
}

type ApplySpendingPoolWithdrawProposalHandler struct {
	keeper keeper.Keeper
	bk     types.BankKeeper
}

func NewApplySpendingPoolWithdrawProposalHandler(keeper keeper.Keeper, bk types.BankKeeper) *ApplySpendingPoolWithdrawProposalHandler {
	return &ApplySpendingPoolWithdrawProposalHandler{
		keeper: keeper,
		bk:     bk,
	}
}

func (a ApplySpendingPoolWithdrawProposalHandler) ProposalType() string {
	return kiratypes.ProposalTypeSpendingPoolWithdraw
}

func (a ApplySpendingPoolWithdrawProposalHandler) IsAllowedAddress(ctx sdk.Context, addr sdk.AccAddress, proposal govtypes.Content) bool {
	p := proposal.(*types.SpendingPoolWithdrawProposal)

	pool := a.keeper.GetSpendingPool(ctx, p.PoolName)
	if pool == nil {
		return false
	}

	return a.keeper.IsAllowedAddress(ctx, addr, *pool.Owners)
}

func (a ApplySpendingPoolWithdrawProposalHandler) AllowedAddresses(ctx sdk.Context, proposal govtypes.Content) []string {
	p := proposal.(*types.SpendingPoolWithdrawProposal)

	pool := a.keeper.GetSpendingPool(ctx, p.PoolName)
	if pool == nil {
		return []string{}
	}

	return a.keeper.AllowedAddresses(ctx, *pool.Owners)
}

func (a ApplySpendingPoolWithdrawProposalHandler) Quorum(ctx sdk.Context, proposal govtypes.Content) sdk.Dec {
	p := proposal.(*types.SpendingPoolWithdrawProposal)

	pool := a.keeper.GetSpendingPool(ctx, p.PoolName)
	if pool == nil {
		return sdk.ZeroDec()
	}

	return pool.VoteQuorum
}

func (a ApplySpendingPoolWithdrawProposalHandler) VotePeriod(ctx sdk.Context, proposal govtypes.Content) uint64 {
	p := proposal.(*types.SpendingPoolWithdrawProposal)

	pool := a.keeper.GetSpendingPool(ctx, p.PoolName)
	if pool == nil {
		return 0
	}

	return pool.VotePeriod
}

func (a ApplySpendingPoolWithdrawProposalHandler) VoteEnactment(ctx sdk.Context, proposal govtypes.Content) uint64 {
	p := proposal.(*types.SpendingPoolWithdrawProposal)

	pool := a.keeper.GetSpendingPool(ctx, p.PoolName)
	if pool == nil {
		return 0
	}

	return pool.VoteEnactment
}

func (a ApplySpendingPoolWithdrawProposalHandler) Apply(ctx sdk.Context, proposalID uint64, proposal govtypes.Content, slash sdk.Dec) error {
	p := proposal.(*spendingtypes.SpendingPoolWithdrawProposal)

	pool := a.keeper.GetSpendingPool(ctx, p.PoolName)
	if pool == nil {
		return types.ErrPoolDoesNotExist
	}

	for _, beneficiary := range p.Beneficiaries {
		beneficiaryAcc, err := sdk.AccAddressFromBech32(beneficiary)
		if err != nil {
			return err
		}

		if !a.keeper.IsAllowedBeneficiary(ctx, beneficiaryAcc, *pool.Beneficiaries) {
			return types.ErrNotPoolBeneficiary
		}

		err = a.bk.SendCoinsFromModuleToAccount(ctx, types.ModuleName, beneficiaryAcc, p.Amounts)
		if err != nil {
			return err
		}

		// update pool to reduce pool's balance
		pool.Balances = sdk.Coins(pool.Balances).Sub(sdk.Coins(p.Amounts)...)
	}

	a.keeper.SetSpendingPool(ctx, *pool)
	return nil
}
