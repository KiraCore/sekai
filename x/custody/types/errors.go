package types

import (
	"github.com/cosmos/cosmos-sdk/types/errors"
)

// errors
var (
	ErrNoWhiteLists        = errors.Register(ModuleName, 2, "error getting whitelist")
	ErrNoWhiteListsElement = errors.Register(ModuleName, 3, "error getting whitelist element")
	ErrNotInWhiteList      = errors.Register(ModuleName, 4, "recipient not in the whitelist")
	ErrNotInLimits         = errors.Register(ModuleName, 5, "denom limit has been reached")
	ErrWrongKey            = errors.Register(ModuleName, 6, "wrong key")
	ErrNotEnoughReward     = errors.Register(ModuleName, 7, "not enough custody reward")
	ErrWrongTargetAddr     = errors.Register(ModuleName, 8, "wrong target address")
)
