package types

import (
	"github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrNotValidatorOwner   = errors.Register(ModuleName, 2, "executor is not validator owner")
	ErrStakingPoolNotFound = errors.Register(ModuleName, 3, "staking pool not found")
)
