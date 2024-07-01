package types

import (
	"cosmossdk.io/math"
)

// NewTokenInfo generates a new token rate struct.
func NewTokenInfo(
	denom string,
	tokenType string,
	feeRate math.LegacyDec,
	feeEnabled bool,
	supply math.Int,
	supplyCap math.Int,
	stakeCap math.LegacyDec,
	stakeMin math.Int,
	stakeEnabled bool,
	inactive bool,
	symbol string,
	name string,
	icon string,
	decimals uint32,
	description string,
	website string,
	social string,
	holders uint64,
	mintingFee math.Int,
	owner string,
	ownerEditDisabled bool,
	nftMetadata string,
	nftHash string,
) TokenInfo {
	if tokenType == "adr43" {
		decimals = 0
	}
	if tokenType == "adr20" {
		nftHash = ""
		nftMetadata = ""
	}
	return TokenInfo{
		Denom:             denom,
		TokenType:         tokenType,
		FeeRate:           feeRate,
		FeeEnabled:        feeEnabled,
		Supply:            supply,
		SupplyCap:         supplyCap,
		StakeCap:          stakeCap,
		StakeMin:          stakeMin,
		StakeEnabled:      stakeEnabled,
		Inactive:          inactive,
		Symbol:            symbol,
		Name:              name,
		Icon:              icon,
		Decimals:          decimals,
		Description:       description,
		Website:           website,
		Social:            social,
		Holders:           holders,
		MintingFee:        mintingFee,
		Owner:             owner,
		OwnerEditDisabled: ownerEditDisabled,
		NftMetadata:       nftMetadata,
		NftHash:           nftHash,
	}
}
