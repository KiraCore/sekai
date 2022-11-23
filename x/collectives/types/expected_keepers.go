package types

import (
	govtypes "github.com/KiraCore/sekai/x/gov/types"
	multistakingtypes "github.com/KiraCore/sekai/x/multistaking/types"
	spendingtypes "github.com/KiraCore/sekai/x/spending/types"
	stakingtypes "github.com/KiraCore/sekai/x/staking/types"
	tokenstypes "github.com/KiraCore/sekai/x/tokens/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

type CustomGovKeeper interface {
	GetNetworkProperties(ctx sdk.Context) *govtypes.NetworkProperties
}

type BankKeeper interface {
	GetSupply(ctx sdk.Context, denom string) sdk.Coin
	GetAllBalances(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	MintCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
	BurnCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	SendCoinsFromModuleToModule(ctx sdk.Context, senderModule, recipientModule string, amt sdk.Coins) error
}

type AccountKeeper interface {
	GetModuleAccount(ctx sdk.Context, moduleName string) authtypes.ModuleAccountI
}

type SpendingKeeper interface {
	GetSpendingPool(ctx sdk.Context, name string) *spendingtypes.SpendingPool
	DepositSpendingPoolFromModule(ctx sdk.Context, moduleName, poolName string, amount sdk.Coin) error
	DepositSpendingPoolFromAccount(ctx sdk.Context, addr sdk.AccAddress, poolName string, amount sdk.Coin) error
}

type StakingKeeper interface {
	GetValidatorByConsAddr(ctx sdk.Context, consAddr sdk.ConsAddress) (stakingtypes.Validator, error)
	GetValidator(sdk.Context, sdk.ValAddress) (stakingtypes.Validator, error)
}

type MultiStakingKeeper interface {
	GetStakingPoolByValidator(ctx sdk.Context, validator string) (pool multistakingtypes.StakingPool, found bool)
	IncreasePoolRewards(ctx sdk.Context, pool multistakingtypes.StakingPool, rewards sdk.Coins)
	RegisterDelegator(ctx sdk.Context, delegator sdk.AccAddress)
	ClaimRewards(ctx sdk.Context, delegator sdk.AccAddress) sdk.Coins
	ClaimRewardsFromModule(ctx sdk.Context, module string) sdk.Coins
}

// TokensKeeper defines expected interface needed to get token rate
type TokensKeeper interface {
	GetTokenRate(ctx sdk.Context, denom string) *tokenstypes.TokenRate
}
