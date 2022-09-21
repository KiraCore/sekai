package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (b Basket) GetBasketDenom() string {
	return fmt.Sprintf("b%d/%s", b.Id, b.Suffix)
}

func (b Basket) RatesAndIndexes() (map[string]sdk.Dec, map[string]int) {
	rates := make(map[string]sdk.Dec)
	indexes := make(map[string]int)
	for index, token := range b.Tokens {
		rates[token.Denom] = token.Weight
		indexes[token.Denom] = index
	}
	return rates, indexes
}

func (b Basket) IncreaseBasketTokens(coins sdk.Coins) (Basket, error) {
	rates, indexes := b.RatesAndIndexes()

	for _, token := range coins {
		_, ok := rates[token.Denom]
		if !ok {
			return b, ErrInvalidBasketDepositDenom
		}
		b.Tokens[indexes[token.Denom]].Amount = b.Tokens[indexes[token.Denom]].Amount.Add(token.Amount)
	}
	return b, nil
}

func (b Basket) DecreaseBasketTokens(coins sdk.Coins) (Basket, error) {
	rates, indexes := b.RatesAndIndexes()

	for _, token := range coins {
		_, ok := rates[token.Denom]
		if !ok {
			return b, ErrInvalidBasketDepositDenom
		}
		b.Tokens[indexes[token.Denom]].Amount = b.Tokens[indexes[token.Denom]].Amount.Sub(token.Amount)
		if b.Tokens[indexes[token.Denom]].Amount.IsNegative() {
			return b, ErrInsufficientBasketDepositToken
		}
	}
	return b, nil
}

func (b Basket) DenomExists(checkTokens []string) bool {
	for _, token := range b.Tokens {
		for _, checkToken := range checkTokens {
			if token.Denom == checkToken {
				return true
			}
		}
	}
	return false
}

func (b Basket) ValidateTokensCap() error {
	totalTokens := sdk.ZeroDec()
	for _, token := range b.Tokens {
		totalTokens = totalTokens.Add(token.Amount.ToDec().Mul(token.Weight))
	}

	for _, token := range b.Tokens {
		if token.Amount.ToDec().Mul(token.Weight).GTE(totalTokens.Mul(b.TokensCap)) {
			return sdkerrors.Wrap(ErrTokenExceedingCap, fmt.Sprintf("denom=%s", token.Denom))
		}
	}
	return nil
}

func (b Basket) AverageDisbalance() sdk.Dec {
	if len(b.Tokens) == 0 {
		return sdk.ZeroDec()
	}

	totalVal := sdk.ZeroDec()
	for _, token := range b.Tokens {
		totalVal = totalVal.Add(token.Weight.Mul(token.Amount.ToDec()))
	}
	averageVal := totalVal.Quo(sdk.NewDec(int64(len(b.Tokens))))
	totalDisbalance := sdk.ZeroDec()
	for _, token := range b.Tokens {
		disbalance := averageVal.Sub(token.Weight.Mul(token.Amount.ToDec())).Quo(averageVal)
		totalDisbalance = totalDisbalance.Add(disbalance.Abs())
	}
	averageDisbalance := totalDisbalance.Quo(sdk.NewDec(int64(len(b.Tokens))))
	return averageDisbalance
}

func (b Basket) SlippageFee(oldDisbalance sdk.Dec) sdk.Dec {
	disbalance := b.AverageDisbalance()
	disbalanceDiff := disbalance.Sub(oldDisbalance)
	if b.SlipppageFeeMin.GT(disbalanceDiff) {
		return b.SlipppageFeeMin
	}
	return disbalanceDiff
}
