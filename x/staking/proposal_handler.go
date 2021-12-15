package staking

import (
	"fmt"
	"time"

	kiratypes "github.com/KiraCore/sekai/types"
	govtypes "github.com/KiraCore/sekai/x/gov/types"
	"github.com/KiraCore/sekai/x/staking/keeper"
	"github.com/KiraCore/sekai/x/staking/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type ApplyUnjailValidatorProposalHandler struct {
	keeper    keeper.Keeper
	govkeeper types.GovKeeper
}

func NewApplyUnjailValidatorProposalHandler(keeper keeper.Keeper, govkeeper types.GovKeeper) *ApplyUnjailValidatorProposalHandler {
	return &ApplyUnjailValidatorProposalHandler{
		keeper:    keeper,
		govkeeper: govkeeper,
	}
}

func (a ApplyUnjailValidatorProposalHandler) ProposalType() string {
	return kiratypes.ProposalTypeUnjailValidator
}

func (a ApplyUnjailValidatorProposalHandler) Apply(ctx sdk.Context, proposalID uint64, proposal govtypes.Content) error {
	p := proposal.(*types.ProposalUnjailValidator)

	valAddr, err := sdk.ValAddressFromBech32(p.ValAddr)
	if err != nil {
		return err
	}

	validator, err := a.keeper.GetValidator(ctx, valAddr)
	if err != nil {
		return err
	}

	if !validator.IsJailed() {
		return fmt.Errorf("validator is not jailed")
	}

	networkProperties := a.govkeeper.GetNetworkProperties(ctx)
	maxUnjailingTime := networkProperties.JailMaxTime

	info, found := a.keeper.GetValidatorJailInfo(ctx, validator.ValKey)
	if !found {
		return fmt.Errorf("validator jailing info not found")
	}

	if info.Time.Add(time.Duration(maxUnjailingTime) * time.Second).Before(ctx.BlockTime()) {
		return fmt.Errorf("time to unjail passed")
	}

	return a.keeper.Unjail(ctx, valAddr)
}
