package types

import "github.com/cosmos/cosmos-sdk/types/errors"

// tokens module errors
var (
	ErrTokenAliasNotFound               = errors.Register(ModuleName, 2, "token alias not found")
	ErrTokenRateNotFound                = errors.Register(ModuleName, 3, "token rate not found")
	ErrTotalRewardsCapExceeds100Percent = errors.Register(ModuleName, 4, "total rewards cap exceeds 100%")
)
