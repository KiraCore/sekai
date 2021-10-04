package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewUpgradePlan(name string, resources []Resource, upgradeTime int64, oldChainId, newChainId string, maxEnrollmentTime int64, rollbackChecksum string, instateUpgrade, rebootRequired, skipHandler bool) Plan {
	return Plan{
		Name:                 name,
		Resources:            resources,
		UpgradeTime:          upgradeTime,
		OldChainId:           oldChainId,
		NewChainId:           newChainId,
		RollbackChecksum:     rollbackChecksum,
		MaxEnrolmentDuration: maxEnrollmentTime,
		InstateUpgrade:       instateUpgrade,
		RebootRequired:       rebootRequired,
		SkipHandler:          skipHandler,
	}
}

func (plan Plan) ShouldExecute(ctx sdk.Context) bool {
	return ctx.BlockTime().Unix() >= plan.UpgradeTime
}
