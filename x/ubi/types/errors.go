package types

import "github.com/cosmos/cosmos-sdk/types/errors"

// ubi module errors
var (
	ErrUBIRecordNotFound        = errors.Register(ModuleName, 1, "ubi record not found")
	ErrUbiSumOverflowsHardcap   = errors.Register(ModuleName, 2, "ubi sum overflows hardcap")
	ErrUBIRecordDoesNotExists   = errors.Register(ModuleName, 3, "ubi record does not exist")
	ErrSpendingPoolDoesNotExist = errors.Register(ModuleName, 4, "spending pool does not exist")
)
