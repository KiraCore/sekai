package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewTokenInfo generates a new token rate struct.
func NewTokenInfo(
	denom string,
	feeRate sdk.Dec,
	feeEnabled bool,
	stakeCap sdk.Dec,
	stakeMin sdk.Int,
	stakeEnabled bool,
	inactive bool,
	symbol string,
	name string,
	icon string,
	decimals uint32,
) TokenInfo {
	return TokenInfo{
		Denom:        denom,
		FeeRate:      feeRate,
		FeeEnabled:   feeEnabled,
		StakeCap:     stakeCap,
		StakeMin:     stakeMin,
		StakeEnabled: stakeEnabled,
		Inactive:     inactive,
		Symbol:       symbol,
		Name:         name,
		Icon:         icon,
		Decimals:     decimals,
	}
}
