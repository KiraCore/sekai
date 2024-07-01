package types

import (
	govtypes "github.com/KiraCore/sekai/x/gov/types"
	spendingtypes "github.com/KiraCore/sekai/x/spending/types"
	tokenstypes "github.com/KiraCore/sekai/x/tokens/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type CustomGovKeeper interface {
	GetNetworkProperties(ctx sdk.Context) *govtypes.NetworkProperties
}

type BankKeeper interface {
	MintCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
}

type SpendingKeeper interface {
	GetSpendingPool(ctx sdk.Context, name string) *spendingtypes.SpendingPool
	DepositSpendingPoolFromModule(ctx sdk.Context, moduleName, poolName string, amounts sdk.Coins) error
}

type DistrKeeper interface {
	InflationPossible(ctx sdk.Context) bool
}

// TokensKeeper defines expected interface needed from tokens keeper
type TokensKeeper interface {
	MintCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
	GetTokenInfo(ctx sdk.Context, denom string) *tokenstypes.TokenInfo
}
