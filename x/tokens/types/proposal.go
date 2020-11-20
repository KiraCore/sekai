package types

import "github.com/KiraCore/sekai/x/gov/types"

const ProposalTypeUpsertTokenAlias = "UpsertTokenAlias"

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
