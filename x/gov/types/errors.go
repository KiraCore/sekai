package types

import "github.com/cosmos/cosmos-sdk/types/errors"

// errors
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
	ErrRoleNotAssigned            = errors.Register(ModuleName, 14, "role not assigned")
	ErrCouncilorNotFound          = errors.Register(ModuleName, 15, "councilor not found")
	ErrUserIsNotCouncilor         = errors.Register(ModuleName, 16, "user is not councilor")
	ErrProposalDoesNotExist       = errors.Register(ModuleName, 17, "proposal does not exist")
	ErrActorIsNotActive           = errors.Register(ModuleName, 18, "actor is not active")
	ErrInvalidNetworkProperty     = errors.Register(ModuleName, 19, "invalid network property")
)
