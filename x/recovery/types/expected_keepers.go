package types

import (
	collectivestypes "github.com/KiraCore/sekai/x/collectives/types"
	custodytypes "github.com/KiraCore/sekai/x/custody/types"
	govtypes "github.com/KiraCore/sekai/x/gov/types"
	multistakingtypes "github.com/KiraCore/sekai/x/multistaking/types"
	spendingtypes "github.com/KiraCore/sekai/x/spending/types"
	stakingtypes "github.com/KiraCore/sekai/x/staking/types"
	tokenstypes "github.com/KiraCore/sekai/x/tokens/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	auth "github.com/cosmos/cosmos-sdk/x/auth/types"
)

// AccountKeeper expected account keeper
type AccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) auth.AccountI
	SetAccount(ctx sdk.Context, acc auth.AccountI)
}

type BankKeeper interface {
	GetAllBalances(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	GetSupply(ctx sdk.Context, denom string) sdk.Coin
	SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
	MintCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
	BurnCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
}

// StakingKeeper expected staking keeper
type StakingKeeper interface {
	GetValidator(sdk.Context, sdk.ValAddress) (stakingtypes.Validator, error) // get a particular validator by operator address
	AddValidator(ctx sdk.Context, validator stakingtypes.Validator)
	RemoveValidator(ctx sdk.Context, validator stakingtypes.Validator)
}

// GovKeeper expected governance keeper
type GovKeeper interface {
	GetNetworkProperties(ctx sdk.Context) *govtypes.NetworkProperties

	SaveCouncilor(ctx sdk.Context, councilor govtypes.Councilor)
	DeleteCouncilor(ctx sdk.Context, councilor govtypes.Councilor)
	GetCouncilor(ctx sdk.Context, address sdk.AccAddress) (govtypes.Councilor, bool)

	GetIdRecordsByAddress(ctx sdk.Context, address sdk.AccAddress) []govtypes.IdentityRecord
	SetIdentityRecord(ctx sdk.Context, record govtypes.IdentityRecord)
	DeleteIdentityRecordById(ctx sdk.Context, recordId uint64)

	GetIdentityRecordIdByAddressKey(ctx sdk.Context, address sdk.AccAddress, key string) uint64
	GetIdentityRecordById(ctx sdk.Context, recordId uint64) *govtypes.IdentityRecord
	GetIdRecordsByAddressAndKeys(ctx sdk.Context, address sdk.AccAddress, keys []string) ([]govtypes.IdentityRecord, error)

	SetIdentityRecordsVerifyRequest(ctx sdk.Context, request govtypes.IdentityRecordsVerify)
	DeleteIdRecordsVerifyRequest(ctx sdk.Context, requestId uint64)
	GetIdRecordsVerifyRequestsByRequester(ctx sdk.Context, requester sdk.AccAddress) []govtypes.IdentityRecordsVerify
	GetIdRecordsVerifyRequestsByApprover(ctx sdk.Context, requester sdk.AccAddress) []govtypes.IdentityRecordsVerify

	SaveNetworkActor(ctx sdk.Context, actor govtypes.NetworkActor)
	DeleteNetworkActor(ctx sdk.Context, actor govtypes.NetworkActor)
	GetNetworkActorByAddress(ctx sdk.Context, address sdk.AccAddress) (govtypes.NetworkActor, bool)
	SetWhitelistAddressPermKey(ctx sdk.Context, actor govtypes.NetworkActor, perm govtypes.PermValue)
	UnassignRoleFromActor(ctx sdk.Context, actor govtypes.NetworkActor, role uint64)
	AssignRoleToActor(ctx sdk.Context, actor govtypes.NetworkActor, role uint64)
	DeleteWhitelistAddressPermKey(ctx sdk.Context, actor govtypes.NetworkActor, perm govtypes.PermValue)

	SaveVote(ctx sdk.Context, vote govtypes.Vote)
	DeleteVote(ctx sdk.Context, vote govtypes.Vote)
	GetVote(ctx sdk.Context, proposalID uint64, address sdk.AccAddress) (govtypes.Vote, bool)
	GetProposals(ctx sdk.Context) ([]govtypes.Proposal, error)
	SaveProposal(ctx sdk.Context, proposal govtypes.Proposal)
}

type MultiStakingKeeper interface {
	GetCompoundInfoByAddress(ctx sdk.Context, addr string) multistakingtypes.CompoundInfo
	SetCompoundInfo(ctx sdk.Context, info multistakingtypes.CompoundInfo)
	RemoveCompoundInfo(ctx sdk.Context, info multistakingtypes.CompoundInfo)

	SetPoolDelegator(ctx sdk.Context, poolId uint64, delegator sdk.AccAddress)
	RemovePoolDelegator(ctx sdk.Context, poolId uint64, delegator sdk.AccAddress)
	IsPoolDelegator(ctx sdk.Context, poolId uint64, delegator sdk.AccAddress) bool

	SetDelegatorRewards(ctx sdk.Context, delegator sdk.AccAddress, rewards sdk.Coins)
	RemoveDelegatorRewards(ctx sdk.Context, delegator sdk.AccAddress)
	GetDelegatorRewards(ctx sdk.Context, delegator sdk.AccAddress) sdk.Coins

	GetAllStakingPools(ctx sdk.Context) []multistakingtypes.StakingPool
	GetStakingPoolByValidator(ctx sdk.Context, validator string) (pool multistakingtypes.StakingPool, found bool)
	SetStakingPool(ctx sdk.Context, pool multistakingtypes.StakingPool)
	RemoveStakingPool(ctx sdk.Context, pool multistakingtypes.StakingPool)
}

type CollectivesKeeper interface {
	SetCollectiveContributer(ctx sdk.Context, cc collectivestypes.CollectiveContributor)
	DeleteCollectiveContributer(ctx sdk.Context, name, address string)
	GetAllCollectiveContributers(ctx sdk.Context) []collectivestypes.CollectiveContributor
}

type SpendingKeeper interface {
	GetAllSpendingPools(ctx sdk.Context) []spendingtypes.SpendingPool
	SetClaimInfo(ctx sdk.Context, claimInfo spendingtypes.ClaimInfo)
	RemoveClaimInfo(ctx sdk.Context, claimInfo spendingtypes.ClaimInfo)
	GetClaimInfo(ctx sdk.Context, poolName string, address sdk.AccAddress) *spendingtypes.ClaimInfo
}

type CustodyKeeper interface {
	GetCustodyInfoByAddress(ctx sdk.Context, address sdk.AccAddress) *custodytypes.CustodySettings
	SetCustodyRecord(ctx sdk.Context, record custodytypes.CustodyRecord)
	DropCustodyRecord(ctx sdk.Context, address sdk.AccAddress)

	GetCustodyCustodiansByAddress(ctx sdk.Context, address sdk.AccAddress) *custodytypes.CustodyCustodianList
	AddToCustodyCustodians(ctx sdk.Context, record custodytypes.CustodyCustodiansRecord)
	DropCustodyCustodiansByAddress(ctx sdk.Context, address sdk.AccAddress)

	GetCustodyWhiteListByAddress(ctx sdk.Context, address sdk.AccAddress) *custodytypes.CustodyWhiteList
	AddToCustodyWhiteList(ctx sdk.Context, record custodytypes.CustodyWhiteListRecord)
	DropCustodyWhiteListByAddress(ctx sdk.Context, address sdk.AccAddress)

	GetCustodyLimitsByAddress(ctx sdk.Context, address sdk.AccAddress) *custodytypes.CustodyLimits
	AddToCustodyLimits(ctx sdk.Context, record custodytypes.CustodyLimitRecord)
	DropCustodyLimitsByAddress(ctx sdk.Context, address sdk.AccAddress)

	GetCustodyLimitsStatusByAddress(ctx sdk.Context, address sdk.AccAddress) *custodytypes.CustodyStatuses
	AddToCustodyLimitsStatus(ctx sdk.Context, record custodytypes.CustodyLimitStatusRecord)
	DropCustodyLimitsStatus(ctx sdk.Context, addr sdk.AccAddress)

	AddToCustodyPool(ctx sdk.Context, record custodytypes.CustodyPool)
	GetCustodyPoolByAddress(ctx sdk.Context, address sdk.AccAddress) *custodytypes.TransactionPool
	DropCustodyPool(ctx sdk.Context, addr sdk.AccAddress)
}

// TokensKeeper defines expected interface needed from tokens keeper
type TokensKeeper interface {
	MintCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
	BurnCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
	GetTokenInfo(ctx sdk.Context, denom string) *tokenstypes.TokenInfo
}
