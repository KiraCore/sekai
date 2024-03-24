package types

import (
	"github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrWrongBridgeAddr = errors.Register(ModuleName, 8, "wrong bridge address")
)
