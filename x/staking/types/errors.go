package types

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrInvalidMonikerLength = fmt.Errorf("invalid moniker length (max 32 bytes)")
	ErrValidatorInactive    = fmt.Errorf("validator is inactive")
	ErrValidatorPaused      = fmt.Errorf("validator is paused")
	ErrValidatorJailed      = fmt.Errorf("validator is jailed")
	ErrValidatorActive      = fmt.Errorf("validator is active")
)

var (
	ErrNetworkActorNotFound    = errors.Register(ModuleName, 2, "network actor not found")
	ErrNotEnoughPermissions    = errors.Register(ModuleName, 3, "not enough permissions")
	ErrValidatorAlreadyClaimed = errors.Register(ModuleName, 4, "validator already claimed")
	ErrValidatorMonikerExists  = errors.Register(ModuleName, 5, "validator moniker exists")
)
