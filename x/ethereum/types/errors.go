package types

import (
	"github.com/cosmos/cosmos-sdk/types/errors"
)

// errors
var (
	ErrEthTxNotValid = errors.Register(ModuleName, 2, "ETH tx not valid")
)
