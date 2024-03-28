package types

import (
	"fmt"
	"strconv"
	"strings"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
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

func (b Basket) DerivativeBasket() bool {
	for _, token := range b.Tokens {
		// all the tokens should start with `v%d/` to be a staking derivative basket
		split := strings.Split(token.Denom, "/")
		if len(split) == 1 {
			return false
		}
		if len(split[0]) < 2 {
			return false
		}
		if !strings.HasPrefix(split[0], "v") {
			return false
		}
		_, err := strconv.Atoi(split[0][1:])
		if err != nil {
			return false
		}
	}
	return true
}

func (b Basket) ValidateTokensCap() error {
	totalTokens := sdk.ZeroDec()
	for _, token := range b.Tokens {
		totalTokens = totalTokens.Add(sdk.NewDecFromInt(token.Amount).Mul(token.Weight))
	}

	for _, token := range b.Tokens {
		if sdk.NewDecFromInt(token.Amount).Mul(token.Weight).GT(totalTokens.Mul(b.TokensCap)) {
			return errorsmod.Wrap(ErrTokenExceedingCap, fmt.Sprintf("denom=%s", token.Denom))
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
		totalVal = totalVal.Add(token.Weight.Mul(sdk.NewDecFromInt(token.Amount)))
	}
	averageVal := totalVal.Quo(sdk.NewDec(int64(len(b.Tokens))))
	totalDisbalance := sdk.ZeroDec()
	for _, token := range b.Tokens {
		disbalance := averageVal.Sub(token.Weight.Mul(sdk.NewDecFromInt(token.Amount))).Quo(averageVal)
		totalDisbalance = totalDisbalance.Add(disbalance.Abs())
	}
	averageDisbalance := totalDisbalance.Quo(sdk.NewDec(int64(len(b.Tokens))))
	return averageDisbalance
}

func (b Basket) SlippageFee(oldDisbalance sdk.Dec) sdk.Dec {
	disbalance := b.AverageDisbalance()
	disbalanceDiff := disbalance.Sub(oldDisbalance)
	// do not charge anything when disbalanceDiff is negative
	if disbalanceDiff.LT(sdk.ZeroDec()) {
		return sdk.ZeroDec()
	}
	if b.SlipppageFeeMin.GT(disbalanceDiff) {
		return b.SlipppageFeeMin
	}
	return disbalanceDiff
}
