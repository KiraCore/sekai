package types

import "github.com/cosmos/cosmos-sdk/types/errors"

var (
	ErrSetPermissions = errors.Register(ModuleName, 2, "error setting permissions")
)
