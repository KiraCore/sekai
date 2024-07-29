package types

import "github.com/cosmos/cosmos-sdk/types/errors"

// tokens module errors
var (
	ErrTokenInfoNotFound                = errors.Register(ModuleName, 3, "token rate not found")
	ErrTotalRewardsCapExceeds100Percent = errors.Register(ModuleName, 4, "total rewards cap exceeds 100%")
	ErrUnimplementedTxType              = errors.Register(ModuleName, 5, "not implemented tx type")
	ErrCannotExceedTokenCap             = errors.Register(ModuleName, 6, "cannot exceed token cap")
	ErrBondDenomIsReadOnly              = errors.Register(ModuleName, 7, "bond denom rate is read-only")
	ErrTokenNotRegistered               = errors.Register(ModuleName, 8, "token not registered")
	ErrSupplyCapShouldNotBeIncreased    = errors.Register(ModuleName, 9, "supply cap should not be increased")
)
