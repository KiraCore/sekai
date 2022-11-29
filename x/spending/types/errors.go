package types

import "github.com/cosmos/cosmos-sdk/types/errors"

// spending module errors
var (
	ErrEmptyProposerAccAddress   = errors.Register(ModuleName, 1, "empty proposer address")
	ErrAlreadyRegisteredPoolName = errors.Register(ModuleName, 2, "already registered spending pool name")
	ErrPoolDoesNotExist          = errors.Register(ModuleName, 3, "pool does not exist")
	ErrNotPoolOwner              = errors.Register(ModuleName, 4, "not a pool owner")
	ErrNotPoolBeneficiary        = errors.Register(ModuleName, 5, "not a pool beneiciary")
	ErrInvalidSpendingPoolName   = errors.Register(ModuleName, 6, "invalid spending pool name")
	ErrNotRegisteredForRewards   = errors.Register(ModuleName, 7, "not registered for rewards")
	ErrNoMoreRewardsToClaim      = errors.Register(ModuleName, 8, "no more rewards to claim")
	ErrInvalidProposalExists     = errors.Register(ModuleName, 9, "invalid proposal exists")
	ErrEmptyWeightBeneficiary    = errors.Register(ModuleName, 10, "empty weight beneficiary not allowed")
)
