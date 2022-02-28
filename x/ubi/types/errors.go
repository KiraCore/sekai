package types

import "github.com/cosmos/cosmos-sdk/types/errors"

// ubi module errors
var (
	ErrUBIRecordNotFound = errors.Register(ModuleName, 2, "ubi record not found")
)
