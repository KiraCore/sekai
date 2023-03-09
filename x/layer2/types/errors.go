package types

import (
	"github.com/cosmos/cosmos-sdk/types/errors"
)

// layer2 module errors
var (
	ErrDappDoesNotExist = errors.Register(ModuleName, 1, "dapp not found")
)
