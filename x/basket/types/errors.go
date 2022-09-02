package types

import "github.com/cosmos/cosmos-sdk/types/errors"

// basket module errors
var (
	ErrBasketNotFound = errors.Register(ModuleName, 1, "basket not found")
)
