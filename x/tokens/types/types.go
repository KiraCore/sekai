package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewTokenAlias generates a new token alias struct.
func NewTokenAlias(
	expiration uint32,
	enactment uint32,
	allowedVoteTypes []VoteType,
	symbol string,
	name string,
	icon string,
	decimals uint32,
	denoms []string,
	status ProposalStatus,
) *TokenAlias {
	return &TokenAlias{
		Expiration:       expiration,
		Enactment:        enactment,
		AllowedVoteTypes: allowedVoteTypes,
		Symbol:           symbol,
		Name:             name,
		Icon:             icon,
		Decimals:         decimals,
		Denoms:           denoms,
		Status:           status,
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
