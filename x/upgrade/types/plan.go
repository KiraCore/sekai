package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewUpgradePlan(name string, resources []Resource, upgradeTime, maxEnrollmentTime int64, rollbackChecksum string, instateUpgrade bool) Plan {
	return Plan{
		Name:                 name,
		Resources:            resources,
		UpgradeTime:          upgradeTime,
		RollbackChecksum:     rollbackChecksum,
		MaxEnrolmentDuration: maxEnrollmentTime,
		InstateUpgrade:       instateUpgrade,
	}
}

func (plan Plan) ShouldExecute(ctx sdk.Context) bool {
	return ctx.BlockTime().Unix() >= plan.UpgradeTime
}
