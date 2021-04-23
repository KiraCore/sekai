package types

func NewUpgradePlan(minHaltTime, maxEnrollmentTime int64, name, rollbackChecksum string) Plan {
	return Plan{
		MinHaltTime:          minHaltTime,
		RollbackChecksum:     rollbackChecksum,
		MaxEnrolmentDuration: maxEnrollmentTime,
		Name:                 name,
	}
}
