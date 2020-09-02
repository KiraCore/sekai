package types

import "github.com/cosmos/cosmos-sdk/types/errors"

var (
	ErrSetPermissions             = errors.Register(ModuleName, 2, "error setting permissions")
	ErrEmptyProposerAccAddress    = errors.Register(ModuleName, 3, "empty proposer key")
	ErrEmptyPermissionsAccAddress = errors.Register(ModuleName, 4, "empty address to set the permissions")
	ErrNotEnoughPermissions       = errors.Register(ModuleName, 5, "not enough permissions")
)
