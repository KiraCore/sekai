package types

import (
	"github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrNotValidatorOwner                = errors.Register(ModuleName, 1, "executor is not validator owner")
	ErrStakingPoolNotFound              = errors.Register(ModuleName, 2, "staking pool not found")
	ErrUndelegationNotFound             = errors.Register(ModuleName, 3, "undelegation not found")
	ErrNotEnoughTimePassed              = errors.Register(ModuleName, 4, "not enough time passed")
	ErrNotAllowedStakingToken           = errors.Register(ModuleName, 5, "not allowed staking token")
	ErrDenomStakingMinTokensNotReached  = errors.Register(ModuleName, 6, "denom staking minimum amount not reached")
	ErrNotActiveValidator               = errors.Register(ModuleName, 7, "not an active validator")
	ErrMaxDelegatorsReached             = errors.Register(ModuleName, 8, "max delegators reached for the pool")
	ErrActionNotSupportedForSlashedPool = errors.Register(ModuleName, 9, "action not supported for slashed pool")
	ErrInsufficientTotalStakingTokens   = errors.Register(ModuleName, 10, "insufficient total staking tokens")
)
