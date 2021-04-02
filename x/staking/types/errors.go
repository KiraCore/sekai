package types

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/types/errors"
)

var ErrInvalidMonikerLength = fmt.Errorf("invalid moniker length (max 64 bytes)")
var ErrInvalidWebsiteLength = fmt.Errorf("invalid website length (max 64 bytes)")
var ErrInvalidSocialLength = fmt.Errorf("invalid social length (max 64 bytes)")
var ErrInvalidIdentityLength = fmt.Errorf("invalid identity length (max 64 bytes)")
var ErrValidatorInactive = fmt.Errorf("validator is inactive")
var ErrValidatorPaused = fmt.Errorf("validator is paused")

var (
	ErrNetworkActorNotFound    = errors.Register(ModuleName, 2, "network actor not found")
	ErrNotEnoughPermissions    = errors.Register(ModuleName, 3, "not enough permissions")
	ErrValidatorAlreadyClaimed = errors.Register(ModuleName, 4, "validator already claimed")
)
