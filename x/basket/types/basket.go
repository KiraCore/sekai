package types

import (
	"fmt"

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
