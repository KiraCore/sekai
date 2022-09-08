package keeper

import (
	"github.com/KiraCore/sekai/x/basket/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) MintBasketToken(ctx sdk.Context, msg *types.MsgBasketTokenMint) error {
	// check if basket is available
	basket, err := k.GetBasketById(ctx, msg.BasketId)
	if err != nil {
		return err
	}

	if !basket.MintsDisabled {
		return types.ErrMintsDisabledBasket
	}

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return err
	}

	err = k.bk.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, msg.Deposit)
	if err != nil {
		return err
	}

	rates, _ := basket.RatesAndIndexes()

	basketTokenAmount := sdk.ZeroDec()
	for _, token := range msg.Deposit {
		rate, ok := rates[token.Denom]
		if !ok {
			return types.ErrInvalidBasketDepositDenom
		}
		basketTokenAmount = basketTokenAmount.Add(token.Amount.ToDec().Mul(rate))
	}

	basketCoin := sdk.NewCoin(basket.GetBasketDenom(), basketTokenAmount.RoundInt())
	basketCoins := sdk.Coins{basketCoin}
	err = k.bk.MintCoins(ctx, types.ModuleName, basketCoins)
	if err != nil {
		return err
	}
	err = k.bk.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sender, basketCoins)
	if err != nil {
		return err
	}

	basket, err = basket.IncreaseBasketTokens(msg.Deposit)
	if err != nil {
		return err
	}
	k.SetBasket(ctx, basket)
	return nil
}

func (k Keeper) BurnBasketToken(ctx sdk.Context, msg *types.MsgBasketTokenBurn) error {
	// check if basket is available
	basket, err := k.GetBasketById(ctx, msg.BasketId)
	if err != nil {
		return err
	}

	if !basket.BurnsDisabled {
		return types.ErrBurnsDisabledBasket
	}

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return err
	}

	burnCoins := sdk.Coins{msg.BurnAmount}
	err = k.bk.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, burnCoins)
	if err != nil {
		return err
	}

	err = k.bk.BurnCoins(ctx, types.ModuleName, burnCoins)
	if err != nil {
		return err
	}

	if msg.BurnAmount.Denom != basket.GetBasketDenom() {
		return types.ErrInvalidBasketDenom
	}

	supply := k.bk.GetSupply(ctx, msg.BurnAmount.Denom)
	portion := msg.BurnAmount.Amount.ToDec().Quo(supply.Amount.ToDec())

	withdrawCoins := sdk.Coins{}
	for _, token := range basket.Tokens {
		if !token.Withdraws {
			continue
		}
		withdrawAmount := token.Amount.ToDec().Mul(portion).RoundInt()
		if withdrawAmount.IsPositive() {
			withdrawCoins = withdrawCoins.Add(sdk.NewCoin(token.Denom, withdrawAmount))
		}
	}

	if withdrawCoins.IsZero() {
		return types.ErrNotAbleToWithdrawAnyTokens
	}

	err = k.bk.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sender, withdrawCoins)
	if err != nil {
		return err
	}

	basket, err = basket.DecreaseBasketTokens(withdrawCoins)
	if err != nil {
		return err
	}
	k.SetBasket(ctx, basket)
	return nil
}

func (k Keeper) BasketSwap(ctx sdk.Context, msg *types.MsgBasketTokenSwap) error {
	// check if basket is available
	basket, err := k.GetBasketById(ctx, msg.BasketId)
	if err != nil {
		return err
	}

	if !basket.SwapsDisabled {
		return types.ErrSwapsDisabledBasket
	}

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return err
	}

	inCoins := sdk.Coins{msg.InAmount}
	err = k.bk.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, inCoins)
	if err != nil {
		return err
	}

	rates, indexes := basket.RatesAndIndexes()

	inRate, ok := rates[msg.InAmount.Denom]
	if !ok {
		return types.ErrInvalidBasketDepositDenom
	}

	outRate, ok := rates[msg.OutToken]
	if !ok {
		return types.ErrInvalidBasketWithdrawDenom
	}

	inTokenIndex := indexes[msg.InAmount.Denom]
	if basket.Tokens[inTokenIndex].Swaps {
		return types.ErrSwapsDisabledForInToken
	}

	outTokenIndex := indexes[msg.OutToken]
	if basket.Tokens[outTokenIndex].Swaps {
		return types.ErrSwapsDisabledForOutToken
	}

	outAmount := msg.InAmount.Amount.ToDec().Mul(inRate).Quo(outRate).Mul(sdk.OneDec().Sub(basket.SwapFee)).RoundInt()
	if outAmount.IsZero() {
		return types.ErrNotAbleToWithdrawAnyTokens
	}

	outCoins := sdk.Coins{sdk.NewCoin(msg.OutToken, outAmount)}
	err = k.bk.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sender, outCoins)
	if err != nil {
		return err
	}

	basket, err = basket.IncreaseBasketTokens(sdk.Coins{msg.InAmount})
	if err != nil {
		return err
	}

	basket, err = basket.DecreaseBasketTokens(outCoins)
	if err != nil {
		return err
	}

	k.SetBasket(ctx, basket)
	return nil
}
