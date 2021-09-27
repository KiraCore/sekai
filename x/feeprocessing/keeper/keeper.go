package keeper

import (
	"bytes"
	"encoding/json"

	kiratypes "github.com/KiraCore/sekai/types"
	"github.com/KiraCore/sekai/x/feeprocessing/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

// Keeper manages module's storage
type Keeper struct {
	cdc      codec.BinaryCodec
	storeKey sdk.StoreKey
	bk       types.BankKeeper
	tk       types.TokensKeeper
	cgk      types.CustomGovKeeper
}

// NewKeeper returns new instance of a keeper
func NewKeeper(storeKey sdk.StoreKey, cdc codec.BinaryCodec, bk types.BankKeeper, tk types.TokensKeeper, cgk types.CustomGovKeeper) Keeper {
	return Keeper{
		cdc,
		storeKey,
		bk,
		tk,
		cgk,
	}
}

// BondDenom returns the denom that is basically used for fee payment
func (k Keeper) BondDenom(ctx sdk.Context) string {
	return "ukex"
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
	totalAmount := sdk.NewDec(0)
	filledAmount := sdk.NewDec(0)

	for _, coin := range amt {
		rate := k.tk.GetTokenRate(ctx, coin.Denom)
		if rate != nil {
			totalAmount = totalAmount.Add(rate.Rate.Mul(coin.Amount.ToDec()))
		}
	}

	for _, coin := range recipientSentCoins {
		rate := k.tk.GetTokenRate(ctx, coin.Denom)
		if rate == nil {
			continue
		}
		toFillAmt := totalAmount.Sub(filledAmount)
		fillAmt := rate.Rate.Mul(coin.Amount.ToDec())
		if fillAmt.GT(toFillAmt) {
			// we don't pay back full amount if there's remainder in div operation
			coinAmt := toFillAmt.BigInt().Div(toFillAmt.BigInt(), rate.Rate.BigInt())
			if coinAmt.Int64() > 0 {
				paybackCoins = paybackCoins.Add(sdk.NewInt64Coin(coin.Denom, coinAmt.Int64()))
				filledAmount = filledAmount.Add(rate.Rate.Mul(sdk.NewDec(coinAmt.Int64())))
			}
		} else {
			filledAmount = filledAmount.Add(fillAmt)
			paybackCoins = paybackCoins.Add(coin)
		}
		if totalAmount.Equal(filledAmount) {
			break
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

// GetExecutionsStatus returns array of executions status registered on that block
func (k Keeper) GetExecutionsStatus(ctx sdk.Context) []types.ExecutionStatus {
	executions := []types.ExecutionStatus{}
	store := ctx.KVStore(k.storeKey)
	if store.Has(types.KeyExecutionStatus) {
		bz := store.Get(types.KeyExecutionStatus)
		json.Unmarshal(bz, &executions)
	}
	return executions
}

// AddExecutionStart add execution on executions status array
func (k Keeper) AddExecutionStart(ctx sdk.Context, msg sdk.Msg) {
	executions := k.GetExecutionsStatus(ctx)
	executions = append(executions, types.ExecutionStatus{
		MsgType:  kiratypes.MsgType(msg),
		FeePayer: msg.GetSigners()[0],
		Success:  false,
	})
	bz, _ := json.Marshal(executions)
	store := ctx.KVStore(k.storeKey)
	store.Set(types.KeyExecutionStatus, bz)
}

// SetExecutionStatusSuccess set statatus of particular message to true
func (k Keeper) SetExecutionStatusSuccess(ctx sdk.Context, msg sdk.Msg) {
	executions := k.GetExecutionsStatus(ctx)
	for i, exec := range executions {
		// when execution message is same as param and success is false, just set success flag to be true and break
		if exec.MsgType == kiratypes.MsgType(msg) && bytes.Equal(exec.FeePayer, msg.GetSigners()[0]) && exec.Success == false {
			executions[i].Success = true
			break
		}
	}
	bz, _ := json.Marshal(executions)
	store := ctx.KVStore(k.storeKey)
	store.Set(types.KeyExecutionStatus, bz)
}

// ProcessExecutionFeeReturn process the executions fee return and clear it up
func (k Keeper) ProcessExecutionFeeReturn(ctx sdk.Context) {
	executions := k.GetExecutionsStatus(ctx)
	for _, exec := range executions {
		fee := k.cgk.GetExecutionFee(ctx, exec.MsgType)
		if fee != nil {
			amount := int64(0)
			if exec.Success && fee.ExecutionFee < fee.FailureFee {
				amount = int64(fee.FailureFee - fee.ExecutionFee)
			}
			if !exec.Success && fee.FailureFee < fee.ExecutionFee {
				amount = int64(fee.ExecutionFee - fee.FailureFee)
			}
			if amount > 0 {
				// handle extra fee based on handler result
				bondDenom := k.BondDenom(ctx)
				fees := sdk.Coins{sdk.NewInt64Coin(bondDenom, amount)}
				k.SendCoinsFromModuleToAccount(ctx, authtypes.FeeCollectorName, exec.FeePayer, fees)
			}
		}
	}
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.KeyExecutionStatus)
}
