package keeper

import (
	"fmt"

	"github.com/KiraCore/sekai/x/basket/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

func (k Keeper) MintBasketToken(ctx sdk.Context, msg *types.MsgBasketTokenMint) error {
	// check if basket is available
	basket, err := k.GetBasketById(ctx, msg.BasketId)
	if err != nil {
		return err
	}

	if basket.MintsDisabled {
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

		_, indexes := basket.RatesAndIndexes()
		tokenIndex := indexes[token.Denom]
		if !basket.Tokens[tokenIndex].Deposits {
			return sdkerrors.Wrap(types.ErrDepositsDisabledForToken, fmt.Sprintf("denom=%s", token.Denom))
		}
		basketTokenAmount = basketTokenAmount.Add(token.Amount.ToDec().Mul(rate))
	}

	basketCoin := sdk.NewCoin(basket.GetBasketDenom(), basketTokenAmount.RoundInt())

	if basketCoin.Amount.LT(basket.MintsMin) {
		return types.ErrAmountBelowBaksetMintsMin
	}

	// register action and check mints max
	k.RegisterMintAction(ctx, msg.BasketId, basketCoin.Amount)
	if k.GetLimitsPeriodMintAmount(ctx, msg.BasketId, basket.LimitsPeriod).GT(basket.MintsMax) {
		return types.ErrAmountAboveBaksetMintsMax
	}

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

	err = basket.ValidateTokensCap()
	if err != nil {
		return err
	}

	basket.Amount = basket.Amount.Add(basketCoin.Amount)
	k.SetBasket(ctx, basket)
	return nil
}

func (k Keeper) BurnBasketToken(ctx sdk.Context, msg *types.MsgBasketTokenBurn) error {
	// check if basket is available
	basket, err := k.GetBasketById(ctx, msg.BasketId)
	if err != nil {
		return err
	}

	if basket.BurnsDisabled {
		return types.ErrBurnsDisabledBasket
	}

	if msg.BurnAmount.Amount.LT(basket.BurnsMin) {
		return types.ErrAmountBelowBaksetBurnsMin
	}

	// register action and check burns max
	k.RegisterBurnAction(ctx, msg.BasketId, msg.BurnAmount.Amount)
	if k.GetLimitsPeriodBurnAmount(ctx, msg.BasketId, basket.LimitsPeriod).GT(basket.BurnsMax) {
		return types.ErrAmountAboveBaksetBurnsMax
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

	err = basket.ValidateTokensCap()
	if err != nil {
		return err
	}

	basket.Amount = basket.Amount.Sub(msg.BurnAmount.Amount)
	k.SetBasket(ctx, basket)
	return nil
}

func (k Keeper) BasketSwap(ctx sdk.Context, msg *types.MsgBasketTokenSwap) error {
	// check if basket is available
	basket, err := k.GetBasketById(ctx, msg.BasketId)
	if err != nil {
		return err
	}

	if basket.SwapsDisabled {
		return types.ErrSwapsDisabledBasket
	}

	oldDisbalance := basket.AverageDisbalance()

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
	if !basket.Tokens[inTokenIndex].Swaps {
		return types.ErrSwapsDisabledForInToken
	}

	outTokenIndex := indexes[msg.OutToken]
	if !basket.Tokens[outTokenIndex].Swaps {
		return types.ErrSwapsDisabledForOutToken
	}

	swapValue := msg.InAmount.Amount.ToDec().Mul(inRate).RoundInt()
	if swapValue.LT(basket.SwapsMin) {
		return types.ErrAmountBelowBaksetSwapsMin
	}

	// register action and check swaps max
	k.RegisterSwapAction(ctx, msg.BasketId, swapValue)
	if k.GetLimitsPeriodSwapAmount(ctx, msg.BasketId, basket.LimitsPeriod).GT(basket.SwapsMax) {
		return types.ErrAmountAboveBaksetSwapsMax
	}

	// calculate out amount considering fees and rates
	swapAmount := msg.InAmount.Amount.ToDec().Mul(sdk.OneDec().Sub(basket.SwapFee)).RoundInt()

	// pay network for fee
	feeAmount := msg.InAmount.Amount.Sub(swapAmount)
	if feeAmount.IsPositive() {
		err := k.bk.SendCoinsFromModuleToModule(ctx, types.ModuleName, authtypes.FeeCollectorName, sdk.Coins{sdk.NewCoin(msg.InAmount.Denom, feeAmount)})
		if err != nil {
			return sdkerrors.Wrapf(sdkerrors.ErrInsufficientFunds, err.Error())
		}
	}

	outAmount := swapAmount.ToDec().Mul(inRate).Quo(outRate).RoundInt()
	if outAmount.IsZero() {
		return types.ErrNotAbleToWithdrawAnyTokens
	}

	// increase in tokens
	basket, err = basket.IncreaseBasketTokens(sdk.Coins{sdk.NewCoin(msg.InAmount.Denom, swapAmount)})
	if err != nil {
		return err
	}

	outCoins := sdk.Coins{sdk.NewCoin(msg.OutToken, outAmount)}
	// decrease out tokens
	basket, err = basket.DecreaseBasketTokens(outCoins)
	if err != nil {
		return err
	}

	// calculate slippage fee
	slippageFee := basket.SlippageFee(oldDisbalance)
	finalOutAmount := outAmount.ToDec().Mul(sdk.OneDec().Sub(slippageFee)).RoundInt()
	finalOutCoins := sdk.Coins{sdk.NewCoin(msg.OutToken, finalOutAmount)}
	err = k.bk.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sender, finalOutCoins)
	if err != nil {
		return err
	}

	// increase surplus by slippage fee
	slippageFeeAmount := outAmount.Sub(finalOutAmount)
	basket.Surplus = sdk.Coins(basket.Surplus).Add(sdk.NewCoin(msg.OutToken, slippageFeeAmount))

	err = basket.ValidateTokensCap()
	if err != nil {
		return err
	}
	k.SetBasket(ctx, basket)
	return nil
}

func (k Keeper) BasketWithdrawSurplus(ctx sdk.Context, p types.ProposalBasketWithdrawSurplus) error {
	withdrawTarget, err := sdk.AccAddressFromBech32(p.WithdrawTarget)
	if err != nil {
		return err
	}

	for _, basketId := range p.BasketIds {
		// check if basket is available
		basket, err := k.GetBasketById(ctx, basketId)
		if err != nil {
			return err
		}

		err = k.bk.SendCoinsFromModuleToAccount(ctx, types.ModuleName, withdrawTarget, sdk.Coins(basket.Surplus))
		if err != nil {
			return err
		}

		basket.Surplus = sdk.Coins{}
		k.SetBasket(ctx, basket)
	}

	// withdraw delegation rewards
	delegator := k.ak.GetModuleAccount(ctx, types.ModuleName).GetAddress()
	k.mk.RegisterDelegator(ctx, delegator)
	rewards := k.mk.ClaimRewardsFromModule(ctx, types.ModuleName)
	if rewards.IsAllPositive() {
		err = k.bk.SendCoinsFromModuleToAccount(ctx, types.ModuleName, withdrawTarget, rewards)
		if err != nil {
			return err
		}
	}

	return nil
}

func (k Keeper) RegisterBasketModuleAsDelegator(ctx sdk.Context) error {
	// withdraw delegation rewards
	delegator := k.ak.GetModuleAccount(ctx, types.ModuleName).GetAddress()
	k.mk.RegisterDelegator(ctx, delegator)

	return nil
}