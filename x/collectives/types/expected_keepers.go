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
	SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
}

type SpendingKeeper interface {
	GetSpendingPool(ctx sdk.Context, name string) *spendingtypes.SpendingPool
	DepositSpendingPoolFromAccount(ctx sdk.Context, addr sdk.AccAddress, poolName string, amounts sdk.Coins) error
}

type MultiStakingKeeper interface {
	RegisterDelegator(ctx sdk.Context, delegator sdk.AccAddress)
	ClaimRewards(ctx sdk.Context, delegator sdk.AccAddress) sdk.Coins
}

// TokensKeeper defines expected interface needed to get token rate
type TokensKeeper interface {
	GetTokenInfo(ctx sdk.Context, denom string) *tokenstypes.TokenInfo
}
