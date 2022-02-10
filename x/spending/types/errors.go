package types

import "github.com/cosmos/cosmos-sdk/types/errors"

// spending module errors
var (
	ErrEmptyProposerAccAddress   = errors.Register(ModuleName, 1, "empty proposer address")
	ErrAlreadyRegisteredPoolName = errors.Register(ModuleName, 2, "already registered spending pool name")
	ErrPoolDoesNotExist          = errors.Register(ModuleName, 3, "pool does not exist")
)
