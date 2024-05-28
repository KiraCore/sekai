package types

import "github.com/cosmos/cosmos-sdk/types/errors"

// tokens module errors
var (
	ErrTokenInfoNotFound                = errors.Register(ModuleName, 3, "token rate not found")
	ErrTotalRewardsCapExceeds100Percent = errors.Register(ModuleName, 4, "total rewards cap exceeds 100%")
	ErrUnimplementedTxType              = errors.Register(ModuleName, 5, "not implemented tx type")
)
