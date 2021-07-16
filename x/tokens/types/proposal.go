package types

import (
	"github.com/KiraCore/sekai/x/gov/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	_ types.Content = &ProposalUpsertTokenAlias{}
	_ types.Content = &ProposalUpsertTokenRates{}
	_ types.Content = &ProposalTokensWhiteBlackChange{}
)

const (
	ProposalTypeUpsertTokenAlias       = "UpsertTokenAlias"
	ProposalTypeUpsertTokenRates       = "UpsertTokenRates"
	ProposalTypeTokensWhiteBlackChange = "TokensWhiteBlackChange"
)

func NewUpsertTokenAliasProposal(
	symbol string,
	name string,
	icon string,
	decimals uint32,
	denoms []string,
) *ProposalUpsertTokenAlias {
	return &ProposalUpsertTokenAlias{
		Symbol:   symbol,
		Name:     name,
		Icon:     icon,
		Decimals: decimals,
		Denoms:   denoms,
	}
}

func (m *ProposalUpsertTokenAlias) ProposalType() string {
	return ProposalTypeUpsertTokenAlias
}

func (m *ProposalUpsertTokenAlias) ProposalPermission() types.PermValue {
	return types.PermCreateUpsertTokenAliasProposal
}

func (m *ProposalUpsertTokenAlias) VotePermission() types.PermValue {
	return types.PermVoteUpsertTokenAliasProposal
}

// ValidateBasic returns basic validation
func (m *ProposalUpsertTokenAlias) ValidateBasic() error {
	return nil
}

func NewUpsertTokenRatesProposal(denom string, rate sdk.Dec, feePayments bool) *ProposalUpsertTokenRates {
	return &ProposalUpsertTokenRates{Denom: denom, Rate: rate, FeePayments: feePayments}
}

func (m *ProposalUpsertTokenRates) ProposalType() string {
	return ProposalTypeUpsertTokenRates
}

func (m *ProposalUpsertTokenRates) ProposalPermission() types.PermValue {
	return types.PermCreateUpsertTokenRateProposal
}

func (m *ProposalUpsertTokenRates) VotePermission() types.PermValue {
	return types.PermVoteUpsertTokenRateProposal
}

// ValidateBasic returns basic validation
func (m *ProposalUpsertTokenRates) ValidateBasic() error {
	return nil
}

func NewTokensWhiteBlackChangeProposal(isBlacklist, isAdd bool, tokens []string) *ProposalTokensWhiteBlackChange {
	return &ProposalTokensWhiteBlackChange{isBlacklist, isAdd, tokens}
}

func (m *ProposalTokensWhiteBlackChange) ProposalType() string {
	return ProposalTypeTokensWhiteBlackChange
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
