package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewTokenInfo generates a new token rate struct.
func NewTokenInfo(
	denom string,
	feeRate sdk.Dec,
	feePayments bool,
	stakeCap sdk.Dec,
	stakeMin sdk.Int,
	stakeToken bool,
	invalidated bool,
	symbol string,
	name string,
	icon string,
	decimals uint32,
) TokenInfo {
	return TokenInfo{
		Denom:       denom,
		FeeRate:     feeRate,
		FeePayments: feePayments,
		StakeCap:    stakeCap,
		StakeMin:    stakeMin,
		StakeToken:  stakeToken,
		Invalidated: invalidated,
		Symbol:      symbol,
		Name:        name,
		Icon:        icon,
		Decimals:    decimals,
	}
}
