package types

import (
	kiratypes "github.com/KiraCore/sekai/types"
	"github.com/KiraCore/sekai/x/gov/types"
)

var (
	_ types.Content = &ProposalUpsertTokenAlias{}
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
	return kiratypes.ProposalTypeUpsertTokenAlias
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
