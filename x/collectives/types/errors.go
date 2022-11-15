package types

import (
	"github.com/cosmos/cosmos-sdk/types/errors"
)

// collectives module errors
var (
	ErrCollectiveDoesNotExist = errors.Register(ModuleName, 1, "collective not found")
)
