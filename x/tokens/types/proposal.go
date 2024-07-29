package types

import (
	"errors"

	"cosmossdk.io/math"
	kiratypes "github.com/KiraCore/sekai/types"
	"github.com/KiraCore/sekai/x/gov/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	_ types.Content = &ProposalUpsertTokenInfo{}
	_ types.Content = &ProposalTokensWhiteBlackChange{}
)

func NewUpsertTokenInfosProposal(
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
) *ProposalUpsertTokenInfo {
	return &ProposalUpsertTokenInfo{
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

func (m *ProposalUpsertTokenInfo) ProposalType() string {
	return kiratypes.ProposalTypeUpsertTokenInfos
}

func (m *ProposalUpsertTokenInfo) ProposalPermission() types.PermValue {
	return types.PermCreateUpsertTokenInfoProposal
}

func (m *ProposalUpsertTokenInfo) VotePermission() types.PermValue {
	return types.PermVoteUpsertTokenInfoProposal
}

// ValidateBasic returns basic validation
func (m *ProposalUpsertTokenInfo) ValidateBasic() error {
	if m.StakeCap.LT(sdk.NewDec(0)) { // not positive
		return errors.New("reward cap should be positive")
	}

	if m.StakeCap.GT(sdk.OneDec()) { // more than 1
		return errors.New("reward cap not be more than 100%")
	}
	return nil
}

func NewTokensWhiteBlackChangeProposal(isBlacklist, isAdd bool, tokens []string) *ProposalTokensWhiteBlackChange {
	return &ProposalTokensWhiteBlackChange{isBlacklist, isAdd, tokens}
}

func (m *ProposalTokensWhiteBlackChange) ProposalType() string {
	return kiratypes.ProposalTypeTokensWhiteBlackChange
}

func (m *ProposalTokensWhiteBlackChange) ProposalPermission() types.PermValue {
	return types.PermCreateTokensWhiteBlackChangeProposal
}

func (m *ProposalTokensWhiteBlackChange) VotePermission() types.PermValue {
	return types.PermVoteTokensWhiteBlackChangeProposal
}

// ValidateBasic returns basic validation
func (m *ProposalTokensWhiteBlackChange) ValidateBasic() error {
	return nil
}
