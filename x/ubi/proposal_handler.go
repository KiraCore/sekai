package ubi

import (
	"fmt"

	kiratypes "github.com/KiraCore/sekai/types"
	"github.com/KiraCore/sekai/x/gov/types"
	"github.com/KiraCore/sekai/x/ubi/keeper"
	ubitypes "github.com/KiraCore/sekai/x/ubi/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type ApplyUpsertUBIProposalHandler struct {
	keeper keeper.Keeper
	gk     ubitypes.CustomGovKeeper
	sk     ubitypes.SpendingKeeper
}

func NewApplyUpsertUBIProposalHandler(keeper keeper.Keeper, gk ubitypes.CustomGovKeeper, sk ubitypes.SpendingKeeper) *ApplyUpsertUBIProposalHandler {
	return &ApplyUpsertUBIProposalHandler{
		keeper: keeper,
		gk:     gk,
		sk:     sk,
	}
}

func (a ApplyUpsertUBIProposalHandler) ProposalType() string {
	return kiratypes.ProposalTypeUpsertUBI
}

func (a ApplyUpsertUBIProposalHandler) Apply(ctx sdk.Context, proposalID uint64, proposal types.Content, slash sdk.Dec) error {
	p := proposal.(*ubitypes.UpsertUBIProposal)

	spendingPool := a.sk.GetSpendingPool(ctx, p.Pool)
	if spendingPool == nil {
		return ubitypes.ErrSpendingPoolDoesNotExist
	}

	yearSeconds := uint64(31556952)
	hardcap := a.gk.GetNetworkProperties(ctx).UbiHardcap
	allRecords := a.keeper.GetUBIRecords(ctx)
	ubiSum := uint64(0)
	for _, record := range allRecords {
		ubiSum += record.Amount * yearSeconds / record.Period
	}

	fmt.Println("ubiSum", ubiSum)
	fmt.Println("ubiSum+p.Amount*yearSeconds/p.Period", ubiSum+p.Amount*yearSeconds/p.Period)
	fmt.Println("hardcap", hardcap)
	// fail if sum of all ((float)amount / period) * 31556952 for all UBI records is greater than ubi-hard-cap.
	if ubiSum+p.Amount*yearSeconds/p.Period > hardcap {
		return ubitypes.ErrUbiSumOverflowsHardcap
	}

	a.keeper.SetUBIRecord(ctx, ubitypes.UBIRecord{
		Name:              p.Name,
		DistributionStart: p.DistributionStart,
		DistributionEnd:   p.DistributionEnd,
		DistributionLast:  p.DistributionStart,
		Amount:            p.Amount,
		Period:            p.Period,
		Pool:              p.Pool,
	})
	return nil
}

type ApplyRemoveUBIProposalHandler struct {
	keeper keeper.Keeper
}

func NewApplyRemoveUBIProposalHandler(keeper keeper.Keeper) *ApplyRemoveUBIProposalHandler {
	return &ApplyRemoveUBIProposalHandler{keeper: keeper}
}

func (a ApplyRemoveUBIProposalHandler) ProposalType() string {
	return kiratypes.ProposalTypeRemoveUBI
}

func (a ApplyRemoveUBIProposalHandler) Apply(ctx sdk.Context, proposalID uint64, proposal types.Content, slash sdk.Dec) error {
	p := proposal.(*ubitypes.RemoveUBIProposal)
	return a.keeper.DeleteUBIRecord(ctx, p.UbiName)
}
