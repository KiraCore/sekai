package types

import "github.com/cosmos/cosmos-sdk/types/errors"

// distributor module errors
var (
	ErrdistributorRecordNotFound              = errors.Register(ModuleName, 1, "distributor record not found")
	ErrdistributorSumOverflowsHardcap         = errors.Register(ModuleName, 2, "distributor sum overflows hardcap")
	ErrdistributorRecordDoesNotExists         = errors.Register(ModuleName, 3, "distributor record does not exist")
	ErrdistributorOnlyAllowedOnBondDenomPools = errors.Register(ModuleName, 4, "distributor is only allowed for bond denom pools")
	ErrSpendingPoolDoesNotExist               = errors.Register(ModuleName, 5, "spending pool does not exist")
)
