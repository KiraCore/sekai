package types

import "github.com/cosmos/cosmos-sdk/types/errors"

// tokens module errors
var (
	ErrTokenAliasNotFound = errors.Register(ModuleName, 2, "token alias not found")
	ErrTokenRateNotFound  = errors.Register(ModuleName, 3, "token rate not found")
)
