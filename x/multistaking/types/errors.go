package types


import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrNotValidatorOwner    = errors.Register(ModuleName, 2, "executor is not validator owner")
)