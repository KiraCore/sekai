package types

import (
	"errors"

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
	rate sdk.Dec,
	feePayments bool,
	stakeCap sdk.Dec,
	stakeMin sdk.Int,
	stakeToken bool,
	isInvalidated bool,
	symbol string,
	name string,
	icon string,
	decimals uint32,
) *ProposalUpsertTokenInfo {
	return &ProposalUpsertTokenInfo{
		Denom:       denom,
		Rate:        rate,
		FeePayments: feePayments,
		StakeCap:    stakeCap,
		StakeMin:    stakeMin,
		StakeToken:  stakeToken,
		Invalidated: isInvalidated,
		Symbol:      symbol,
		Name:        name,
		Icon:        icon,
		Decimals:    decimals,
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
