package types

import (
	"github.com/cosmos/cosmos-sdk/types/errors"
)

// collectives module errors
var (
	ErrCollectiveDoesNotExist                    = errors.Register(ModuleName, 1, "collective not found")
	ErrCollectiveAlreadyExists                   = errors.Register(ModuleName, 2, "collective already exists")
	ErrNotCollectiveContributer                  = errors.Register(ModuleName, 3, "not a collective contributer")
	ErrDonationLocked                            = errors.Register(ModuleName, 4, "donation is in lock")
	ErrTotalSpendingPoolWeightShouldBeOne        = errors.Register(ModuleName, 5, "total spending pool weight should be one")
	ErrNumberOfSpendingPoolsBiggerThanMaxOutputs = errors.Register(ModuleName, 6, "number of spending pools is bigger than max outputs")
	InitialBondLowerThanTenPercentOfMinimumBond  = errors.Register(ModuleName, 7, "initial bond is lower than 10%% of minimum bond")
	ErrNotWhitelistedForCollectiveDeposit        = errors.Register(ModuleName, 8, "not whitelisted for collective deposit")
)
