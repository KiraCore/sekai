package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewUpgradePlan(minHaltTime, maxEnrollmentTime int64, name, rollbackChecksum string) Plan {
	return Plan{
		MinHaltTime:          minHaltTime,
		RollbackChecksum:     rollbackChecksum,
		MaxEnrolmentDuration: maxEnrollmentTime,
		Name:                 name,
	}
}

func (plan Plan) ShouldExecute(ctx sdk.Context) bool {
	return ctx.BlockHeight() > plan.MinHaltTime
}
