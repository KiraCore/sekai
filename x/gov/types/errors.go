package types

import "github.com/cosmos/cosmos-sdk/types/errors"

var (
	ErrSetPermissions             = errors.Register(ModuleName, 2, "error setting permissions")
	ErrEmptyProposerAccAddress    = errors.Register(ModuleName, 3, "empty proposer key")
	ErrEmptyPermissionsAccAddress = errors.Register(ModuleName, 4, "empty address to set the permissions")
	ErrNotEnoughPermissions       = errors.Register(ModuleName, 5, "not enough permissions")
	ErrCouncilorEmptyAddress      = errors.Register(ModuleName, 6, "empty councilor address")
	ErrRoleDoesNotExist           = errors.Register(ModuleName, 7, "role does not exist")
	ErrRoleExist                  = errors.Register(ModuleName, 12, "role already exist")
	ErrWhitelisting               = errors.Register(ModuleName, 8, "error adding to whitelist")
	ErrBlacklisting               = errors.Register(ModuleName, 9, "error adding to blacklist")
	ErrRemovingWhitelist          = errors.Register(ModuleName, 10, "error removing from whitelist")
	ErrRemovingBlacklist          = errors.Register(ModuleName, 11, "error removing from blacklist")
	ErrRoleAlreadyAssigned        = errors.Register(ModuleName, 13, "role already assigned")
)
