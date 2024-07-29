package types

import (
	stakingtypes "github.com/KiraCore/sekai/x/staking/types"
	tokenstypes "github.com/KiraCore/sekai/x/tokens/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MultistakingHooks event hooks for multistaking
type MultistakingHooks interface {
	AfterUpsertStakingPool(ctx sdk.Context, valAddr sdk.ValAddress, pool StakingPool) // Must be called when a upsert staking pool
	AfterSlashStakingPool(ctx sdk.Context, valAddr sdk.ValAddress, pool StakingPool, slash sdk.Dec)
}

// StakingKeeper expected staking keeper
type StakingKeeper interface {
	DefaultDenom(sdk.Context) string
	GetValidator(sdk.Context, sdk.ValAddress) (stakingtypes.Validator, error)
}

type DistributorKeeper interface {
	SetFeesTreasury(ctx sdk.Context, coins sdk.Coins)
	GetFeesTreasury(ctx sdk.Context) sdk.Coins
}

type BankKeeper interface {
	GetAllBalances(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	MintCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
	BurnCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	SendCoinsFromModuleToModule(ctx sdk.Context, senderModule, recipientModule string, amt sdk.Coins) error
}

// TokensKeeper defines expected interface needed from tokens keeper
type TokensKeeper interface {
	MintCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
	GetTokenInfo(ctx sdk.Context, denom string) *tokenstypes.TokenInfo
}
