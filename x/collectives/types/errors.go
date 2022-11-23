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
	ErrLockPeriodCanOnlyBeIncreased              = errors.Register(ModuleName, 9, "lock period can only be increased")
	ErrLockPeriodCannotExceedOneYear             = errors.Register(ModuleName, 10, "lock period cannot exceed one year")
	ErrBondsLockedOnTheCollective                = errors.Register(ModuleName, 11, "your bonds are locked on the collective")
	ErrInvalidDonationValue                      = errors.Register(ModuleName, 12, "invalid donation value: should be between 0 and 1")
	ErrClaimPeriodLowerThanNetworkConfig         = errors.Register(ModuleName, 13, "claim period is lower than network properties configuration")
)
