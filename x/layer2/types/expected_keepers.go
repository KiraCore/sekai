package types

import (
	govtypes "github.com/KiraCore/sekai/x/gov/types"
	spendingtypes "github.com/KiraCore/sekai/x/spending/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type CustomGovKeeper interface {
	GetNetworkProperties(ctx sdk.Context) *govtypes.NetworkProperties
}

type BankKeeper interface {
	MintCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
}

type SpendingKeeper interface {
	CreateSpendingPool(ctx sdk.Context, pool spendingtypes.SpendingPool) error
	DepositSpendingPoolFromModule(ctx sdk.Context, moduleName, poolName string, amount sdk.Coins) error
}
