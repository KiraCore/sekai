// noalias
// DONTCOVER
package types

import (
	govtypes "github.com/KiraCore/sekai/x/gov/types"
	stakingtypes "github.com/KiraCore/sekai/x/staking/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	auth "github.com/cosmos/cosmos-sdk/x/auth/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// AccountKeeper expected account keeper
type AccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) auth.AccountI
	IterateAccounts(ctx sdk.Context, process func(auth.AccountI) (stop bool))
}

// BankKeeper defines the expected interface needed to retrieve account balances.
type BankKeeper interface {
	GetAllBalances(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin
	LockedCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	SpendableCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
}

// ParamSubspace defines the expected Subspace interfacace
type ParamSubspace interface {
	HasKeyTable() bool
	WithKeyTable(table paramtypes.KeyTable) paramtypes.Subspace
	Get(ctx sdk.Context, key []byte, ptr interface{})
	GetParamSet(ctx sdk.Context, ps paramtypes.ParamSet)
	SetParamSet(ctx sdk.Context, ps paramtypes.ParamSet)
}

// StakingKeeper expected staking keeper
type StakingKeeper interface {
	// iterate through validators by operator address, execute func for each validator
	IterateValidators(sdk.Context,
		func(index int64, validator *stakingtypes.Validator) (stop bool))

	GetValidator(sdk.Context, sdk.ValAddress) (stakingtypes.Validator, error)            // get a particular validator by operator address
	GetValidatorByConsAddr(sdk.Context, sdk.ConsAddress) (stakingtypes.Validator, error) // get a particular validator by consensus address
	GetValidatorSet(ctx sdk.Context) []stakingtypes.Validator                            // get all validator set

	// activate/inactivate the validator and delegators of the validator, specifying offence height, offence power, and slash fraction
	Inactivate(sdk.Context, sdk.ValAddress) error // inactivate a validator
	Activate(sdk.Context, sdk.ValAddress) error   // activate a validator
	Jail(sdk.Context, sdk.ValAddress) error       // jail a validator
	ResetWholeValidatorRank(sdk.Context)          // reset whole validator rank

	// pause/unpause the validator and delegators of the validator, specifying offence height, offence power, and slash fraction
	Pause(sdk.Context, sdk.ValAddress) error   // pause a validator
	Unpause(sdk.Context, sdk.ValAddress) error // unpause a validator

	HandleValidatorSignature(sdk.Context, sdk.ValAddress, bool, int64) error

	// MaxValidators returns the maximum amount of joined validators
	MaxValidators(sdk.Context) uint32

	GetIdRecordsByAddress(sdk.Context, sdk.AccAddress) []govtypes.IdentityRecord
}

// StakingHooks event hooks for staking validator object (noalias)
type StakingHooks interface {
	AfterValidatorCreated(ctx sdk.Context, valAddr sdk.ValAddress)                           // Must be called when a validator is created
	AfterValidatorRemoved(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) // Must be called when a validator is deleted
	AfterValidatorJoined(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress)  // Must be called when a validator is joined
}

// GovKeeper expected governance keeper
type GovKeeper interface {
	GetNetworkProperties(sdk.Context) *govtypes.NetworkProperties // returns network properties
	CheckIfAllowedPermission(ctx sdk.Context, addr sdk.AccAddress, permValue govtypes.PermValue) bool
	GetNextProposalID(ctx sdk.Context) uint64
	SaveProposal(ctx sdk.Context, proposal govtypes.Proposal)
	AddToActiveProposals(ctx sdk.Context, proposal govtypes.Proposal)
}
