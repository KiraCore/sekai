package types

import "github.com/cosmos/cosmos-sdk/types/errors"

// tokens module errors
var (
	ErrEmptyProposerAccAddress = errors.Register(ModuleName, 1, "empty proposer address")
)
