package types

import (
	"github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrInvalidUpgradeTime    = errors.Register(ModuleName, 1, "invalid upgrade time")
	ErrInvalidGenesisVersion = errors.Register(ModuleName, 2, "invalid genesis version")
)
