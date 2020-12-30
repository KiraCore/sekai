package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewTokenAlias generates a new token alias struct.
func NewTokenAlias(
	symbol string,
	name string,
	icon string,
	decimals uint32,
	denoms []string,
) *TokenAlias {
	return &TokenAlias{
		Symbol:   symbol,
		Name:     name,
		Icon:     icon,
		Decimals: decimals,
		Denoms:   denoms,
	}
}

// NewTokenRate generates a new token rate struct.
func NewTokenRate(
	denom string,
	rate sdk.Dec,
	feePayments bool,
) *TokenRate {
	return &TokenRate{
		Denom:       denom,
		Rate:        rate,
		FeePayments: feePayments,
	}
}
