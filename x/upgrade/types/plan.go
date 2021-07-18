package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewUpgradePlan(name string, resources []Resource, height, minUpgradeTime, maxEnrollmentTime int64, rollbackChecksum string, instateUpgrade bool) Plan {
	return Plan{
		Name:                 name,
		Resources:            resources,
		Height:               height,
		MinUpgradeTime:       minUpgradeTime,
		RollbackChecksum:     rollbackChecksum,
		MaxEnrolmentDuration: maxEnrollmentTime,
		InstateUpgrade:       instateUpgrade,
	}
}

func (plan Plan) ShouldExecute(ctx sdk.Context) bool {
	return ctx.BlockHeight() == plan.Height || ctx.BlockHeight() > plan.MinUpgradeTime
}
