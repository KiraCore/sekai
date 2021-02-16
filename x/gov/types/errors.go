package types

import (
	"github.com/cosmos/cosmos-sdk/types/errors"
)

// errors
var (
	ErrSetPermissions             = errors.Register(ModuleName, 2, "error setting permissions")
	ErrEmptyProposerAccAddress    = errors.Register(ModuleName, 3, "empty proposer key")
	ErrEmptyPermissionsAccAddress = errors.Register(ModuleName, 4, "empty address to set the permissions")
	ErrNotEnoughPermissions       = errors.Register(ModuleName, 5, "not enough permissions")
	ErrCouncilorEmptyAddress      = errors.Register(ModuleName, 6, "empty councilor address")
	ErrRoleDoesNotExist           = errors.Register(ModuleName, 7, "role does not exist")
	ErrWhitelisting               = errors.Register(ModuleName, 8, "error adding to whitelist")
	ErrRoleExist                  = errors.Register(ModuleName, 12, "role already exist")
	ErrRoleAlreadyAssigned        = errors.Register(ModuleName, 13, "role already assigned")
	ErrRoleNotAssigned            = errors.Register(ModuleName, 14, "role not assigned")
	ErrCouncilorNotFound          = errors.Register(ModuleName, 15, "councilor not found")
	ErrProposalDoesNotExist       = errors.Register(ModuleName, 17, "proposal does not exist")
	ErrActorIsNotActive           = errors.Register(ModuleName, 18, "actor is not active")
	ErrInvalidNetworkProperty     = errors.Register(ModuleName, 19, "invalid network property")
	ErrFeeNotExist                = errors.Register(ModuleName, 20, "fee does not exist")
	ErrPoorNetworkMsgsNotSet      = errors.Register(ModuleName, 21, "poor network messages not set")
	ErrGettingProposals           = errors.Register(ModuleName, 23, "error getting proposals")
	ErrGettingProposalVotes       = errors.Register(ModuleName, 24, "error getting votes for proposal")
	ErrVotingTimeEnded            = errors.Register(ModuleName, 25, "voting time has ended")
)
