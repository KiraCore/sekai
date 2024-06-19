package types

import (
	govtypes "github.com/KiraCore/sekai/x/gov/types"
	multistakingtypes "github.com/KiraCore/sekai/x/multistaking/types"
	tokenstypes "github.com/KiraCore/sekai/x/tokens/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

type CustomGovKeeper interface {
	GetNetworkProperties(ctx sdk.Context) *govtypes.NetworkProperties
}

type BankKeeper interface {
	GetSupply(ctx sdk.Context, denom string) sdk.Coin
	MintCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
	BurnCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	SendCoinsFromModuleToModule(ctx sdk.Context, senderModule, recipientModule string, amt sdk.Coins) error
}

type AccountKeeper interface {
	GetModuleAccount(ctx sdk.Context, moduleName string) authtypes.ModuleAccountI
}

type MultiStakingKeeper interface {
	RegisterDelegator(ctx sdk.Context, delegator sdk.AccAddress)
	ClaimRewards(ctx sdk.Context, delegator sdk.AccAddress) sdk.Coins
	ClaimRewardsFromModule(ctx sdk.Context, module string) sdk.Coins
}

// MultistakingHooks event hooks for multistaking
type MultistakingHooks interface {
	AfterUpsertStakingPool(ctx sdk.Context, valAddr sdk.ValAddress, pool multistakingtypes.StakingPool) // Must be called when a upsert staking pool
	AfterSlashStakingPool(ctx sdk.Context, valAddr sdk.ValAddress, pool multistakingtypes.StakingPool, slash sdk.Dec)
}

// TokensKeeper defines expected interface needed from tokens keeper
type TokensKeeper interface {
	MintCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
	BurnCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
	GetTokenInfo(ctx sdk.Context, denom string) *tokenstypes.TokenInfo
	GetAllTokenInfos(ctx sdk.Context) []tokenstypes.TokenInfo
}
