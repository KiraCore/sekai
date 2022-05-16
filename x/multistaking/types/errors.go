package types

import (
	"github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrNotValidatorOwner    = errors.Register(ModuleName, 2, "executor is not validator owner")
	ErrStakingPoolNotFound  = errors.Register(ModuleName, 3, "staking pool not found")
	ErrUndelegationNotFound = errors.Register(ModuleName, 4, "undelegation not found")
	ErrNotEnoughTimePassed  = errors.Register(ModuleName, 5, "not enough time passed")
)
