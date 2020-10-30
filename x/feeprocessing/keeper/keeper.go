package keeper

import (
	"encoding/json"

	"github.com/KiraCore/sekai/x/feeprocessing/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Keeper manages module's storage
type Keeper struct {
	cdc      codec.BinaryMarshaler
	storeKey sdk.StoreKey
	bk       types.BankKeeper
	tk       types.TokensKeeper
}

// NewKeeper returns new instance of a keeper
func NewKeeper(storeKey sdk.StoreKey, cdc codec.BinaryMarshaler, bk types.BankKeeper, tk types.TokensKeeper) Keeper {
	return Keeper{cdc: cdc, storeKey: storeKey, bk: bk, tk: tk}
}

// GetSenderCoinsHistory returns fee payment history of an address
func (k Keeper) GetSenderCoinsHistory(ctx sdk.Context, senderAddr sdk.AccAddress) sdk.Coins {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyFeePaymentHistory)
	bz := prefixStore.Get([]byte(senderAddr))
	coins := []sdk.Coin{}
	json.Unmarshal(bz, &coins)
	return coins
}

// SetSenderCoinsHistory set fee payment history of an address
func (k Keeper) SetSenderCoinsHistory(ctx sdk.Context, senderAddr sdk.AccAddress, coins sdk.Coins) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyFeePaymentHistory)
	bz, _ := json.Marshal(coins)
	prefixStore.Set([]byte(senderAddr), bz)
}

// SendCoinsFromModuleToAccount is a wrapper of bank keeper's SendCoinsFromModuleToAccount
func (k Keeper) SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error {
	recipientSentCoins := k.GetSenderCoinsHistory(ctx, recipientAddr)
	paybackCoins := sdk.Coins{}
	totalAmount := int64(0)
	filledAmount := int64(0)

	for _, coin := range amt {
		rate := k.tk.GetTokenRate(ctx, coin.Denom)
		if rate != nil {
			totalAmount += int64(rate.Rate) * coin.Amount.Int64()
		}
	}

	for _, coin := range recipientSentCoins {
		rate := k.tk.GetTokenRate(ctx, coin.Denom)
		if rate == nil {
			continue
		}
		fillAmt := int64(rate.Rate) * coin.Amount.Int64()
		if fillAmt > totalAmount-filledAmount {
			coinAmt := (totalAmount - filledAmount) / int64(rate.Rate)
			paybackCoins.Add(sdk.NewInt64Coin(coin.Denom, coinAmt))
			filledAmount = totalAmount
		} else {
			filledAmount += fillAmt
			paybackCoins.Add(coin)
		}
	}

	k.SetSenderCoinsHistory(ctx, recipientAddr, recipientSentCoins.Sub(paybackCoins))
	return k.bk.SendCoinsFromModuleToAccount(ctx, senderModule, recipientAddr, paybackCoins)
}

// SendCoinsFromAccountToModule is a wrapper of bank keeper's SendCoinsFromAccountToModule
func (k Keeper) SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error {
	coins := k.GetSenderCoinsHistory(ctx, senderAddr)
	k.SetSenderCoinsHistory(ctx, senderAddr, coins.Add(amt...))
	return k.bk.SendCoinsFromAccountToModule(ctx, senderAddr, recipientModule, amt)
}
