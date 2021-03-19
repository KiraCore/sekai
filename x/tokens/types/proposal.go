package types

import (
	"github.com/KiraCore/sekai/x/gov/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	_ types.Content = &ProposalUpsertTokenAlias{}
	_ types.Content = &ProposalUpsertTokenRates{}
)

const (
	ProposalTypeUpsertTokenAlias       = "UpsertTokenAlias"
	ProposalTypeUpsertTokenRates       = "UpsertTokenRates"
	ProposalTypeTokensWhiteBlackChange = "TokensWhiteBlackChange"
)

func NewProposalUpsertTokenAlias(
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

func (m *ProposalUpsertTokenAlias) VotePermission() types.PermValue {
	return types.PermVoteUpsertTokenAliasProposal
}

func NewProposalUpsertTokenRates(denom string, rate sdk.Dec, feePayments bool) *ProposalUpsertTokenRates {
	return &ProposalUpsertTokenRates{Denom: denom, Rate: rate, FeePayments: feePayments}
}

func (m *ProposalUpsertTokenRates) ProposalType() string {
	return ProposalTypeUpsertTokenRates
}

func (m *ProposalUpsertTokenRates) VotePermission() types.PermValue {
	return types.PermVoteUpsertTokenRateProposal
}

func NewProposalTokensWhiteBlackChange(isBlacklist, isAdd bool, tokens []string) *ProposalTokensWhiteBlackChange {
	return &ProposalTokensWhiteBlackChange{isBlacklist, isAdd, tokens}
}

func (m *ProposalTokensWhiteBlackChange) ProposalType() string {
	return ProposalTypeTokensWhiteBlackChange
}

func (m *ProposalTokensWhiteBlackChange) VotePermission() types.PermValue {
	return types.PermVoteTokensWhiteBlackChangeProposal
}
