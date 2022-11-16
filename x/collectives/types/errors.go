package types

import (
	"github.com/cosmos/cosmos-sdk/types/errors"
)

// collectives module errors
var (
	ErrCollectiveDoesNotExist   = errors.Register(ModuleName, 1, "collective not found")
	ErrCollectiveAlreadyExists  = errors.Register(ModuleName, 2, "collective already exists")
	ErrNotCollectiveContributer = errors.Register(ModuleName, 3, "not a collective contributer")
	ErrDonationLocked           = errors.Register(ModuleName, 4, "donation is in lock")
)
