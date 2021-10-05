package types

import (
	govtypes "github.com/KiraCore/sekai/x/gov/types"
	"github.com/KiraCore/sekai/x/tokens/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// BankKeeper defines the expected interface needed to retrieve account balances.
type BankKeeper interface {
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
}

// TokensKeeper defines expected interface needed to get token rate
type TokensKeeper interface {
	GetTokenRate(ctx sdk.Context, denom string) *types.TokenRate
}

// CustomGovKeeper defines the expected interface contract the tokens module requires
type CustomGovKeeper interface {
	GetExecutionFee(ctx sdk.Context, txType string) *govtypes.ExecutionFee
}
