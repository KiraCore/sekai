package upgrade

import (
	"github.com/KiraCore/sekai/x/gov/types"
	"github.com/KiraCore/sekai/x/upgrade/keeper"
	upgradetypes "github.com/KiraCore/sekai/x/upgrade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type ApplySoftwareUpgradeProposalHandler struct {
	keeper keeper.Keeper
}

func NewApplySoftwareUpgradeProposalHandler(keeper keeper.Keeper) *ApplySoftwareUpgradeProposalHandler {
	return &ApplySoftwareUpgradeProposalHandler{
		keeper: keeper,
	}
}

func (a ApplySoftwareUpgradeProposalHandler) ProposalType() string {
	return upgradetypes.ProposalTypeSoftwareUpgrade
}

func (a ApplySoftwareUpgradeProposalHandler) Apply(ctx sdk.Context, proposal types.Content) {
	p := proposal.(*upgradetypes.ProposalSoftwareUpgrade)

	plan := upgradetypes.NewUpgradePlan(p.MinHaltTime, p.MaxEnrolmentDuration, p.Resources.Id, p.RollbackChecksum)
	a.keeper.SaveUpgradePlan(ctx, plan)
}
