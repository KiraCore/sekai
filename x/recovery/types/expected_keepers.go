package types

import (
	collectivestypes "github.com/KiraCore/sekai/x/collectives/types"
	govtypes "github.com/KiraCore/sekai/x/gov/types"
	multistakingtypes "github.com/KiraCore/sekai/x/multistaking/types"
	stakingtypes "github.com/KiraCore/sekai/x/staking/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	auth "github.com/cosmos/cosmos-sdk/x/auth/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// AccountKeeper expected account keeper
type AccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) auth.AccountI
	SetAccount(ctx sdk.Context, acc auth.AccountI)
	IterateAccounts(ctx sdk.Context, process func(auth.AccountI) (stop bool))
}

type BankKeeper interface {
	GetAllBalances(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
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
	SaveCouncilor(ctx sdk.Context, councilor govtypes.Councilor)
	DeleteCouncilor(ctx sdk.Context, councilor govtypes.Councilor)
	GetCouncilor(ctx sdk.Context, address sdk.AccAddress) (govtypes.Councilor, bool)

	GetIdRecordsByAddress(ctx sdk.Context, address sdk.AccAddress) []govtypes.IdentityRecord
	SetIdentityRecord(ctx sdk.Context, record govtypes.IdentityRecord)
	DeleteIdentityRecordById(ctx sdk.Context, recordId uint64)

	SetIdentityRecordsVerifyRequest(ctx sdk.Context, request govtypes.IdentityRecordsVerify)
	DeleteIdRecordsVerifyRequest(ctx sdk.Context, requestId uint64)
	GetIdRecordsVerifyRequestsByRequester(ctx sdk.Context, requester sdk.AccAddress) []govtypes.IdentityRecordsVerify
	GetIdRecordsVerifyRequestsByApprover(ctx sdk.Context, requester sdk.AccAddress) []govtypes.IdentityRecordsVerify

	SaveNetworkActor(ctx sdk.Context, actor govtypes.NetworkActor)
	DeleteNetworkActor(ctx sdk.Context, actor govtypes.NetworkActor)
	GetNetworkActorByAddress(ctx sdk.Context, address sdk.AccAddress) (govtypes.NetworkActor, bool)
	SetWhitelistAddressPermKey(ctx sdk.Context, actor govtypes.NetworkActor, perm govtypes.PermValue)
	RemoveRoleFromActor(ctx sdk.Context, actor govtypes.NetworkActor, role uint64)
	AssignRoleToActor(ctx sdk.Context, actor govtypes.NetworkActor, role uint64)
	DeleteWhitelistAddressPermKey(ctx sdk.Context, actor govtypes.NetworkActor, perm govtypes.PermValue)

	SaveVote(ctx sdk.Context, vote govtypes.Vote)
	DeleteVote(ctx sdk.Context, vote govtypes.Vote)
	GetVote(ctx sdk.Context, proposalID uint64, address sdk.AccAddress) (govtypes.Vote, bool)
	GetProposals(ctx sdk.Context) ([]govtypes.Proposal, error)
}

type MultiStakingKeeper interface {
	GetStakingPoolByValidator(ctx sdk.Context, validator string) (pool multistakingtypes.StakingPool, found bool)
	IncreasePoolRewards(ctx sdk.Context, pool multistakingtypes.StakingPool, rewards sdk.Coins)
	GetAllStakingPools(ctx sdk.Context) []multistakingtypes.StakingPool
	SlashStakingPool(ctx sdk.Context, validator string, slash sdk.Dec)
}

type CollectivesKeeper interface {
	SetCollectiveContributer(ctx sdk.Context, cc collectivestypes.CollectiveContributor)
	DeleteCollectiveContributer(ctx sdk.Context, name, address string)
	GetAllCollectiveContributers(ctx sdk.Context) []collectivestypes.CollectiveContributor
}
