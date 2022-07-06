package types

import (
	"github.com/cosmos/cosmos-sdk/types/errors"
)

// errors
var (
	ErrNoWhiteLists        = errors.Register(ModuleName, 2, "error getting whitelist")
	ErrNoWhiteListsElement = errors.Register(ModuleName, 3, "error getting whitelist element")
	ErrNotInWhiteList      = errors.Register(ModuleName, 4, "recipient not in the whitelist")
)
