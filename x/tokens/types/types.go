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
	invalidated bool,
) *TokenAlias {
	return &TokenAlias{
		Symbol:      symbol,
		Name:        name,
		Icon:        icon,
		Decimals:    decimals,
		Denoms:      denoms,
		Invalidated: invalidated,
	}
}

// NewTokenRate generates a new token rate struct.
func NewTokenRate(
	denom string,
	feeRate sdk.Dec,
	feePayments bool,
	stakeCap sdk.Dec,
	stakeMin sdk.Int,
	stakeToken bool,
	invalidated bool,
) *TokenRate {
	return &TokenRate{
		Denom:       denom,
		FeeRate:     feeRate,
		FeePayments: feePayments,
		StakeCap:    stakeCap,
		StakeMin:    stakeMin,
		StakeToken:  stakeToken,
		Invalidated: invalidated,
	}
}
