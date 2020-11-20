package types

import (
	"github.com/KiraCore/sekai/x/gov/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	_ sdk.Msg       = &MsgProposalUpsertTokenAlias{}
	_ types.Content = &ProposalUpsertTokenAlias{}
)

const MsgProposalUpsertTokenAliasType = "propose-upsert-token-alias"

func NewMsgProposalUpsertTokenAlias(
	proposer sdk.AccAddress,
	symbol string,
	name string,
	icon string,
	decimals uint32,
	denoms []string,
) *MsgProposalUpsertTokenAlias {
	return &MsgProposalUpsertTokenAlias{
		Proposer: proposer,
		Symbol:   symbol,
		Name:     name,
		Icon:     icon,
		Decimals: decimals,
		Denoms:   denoms,
	}
}

func (m *MsgProposalUpsertTokenAlias) Route() string {
	return ModuleName
}

func (m *MsgProposalUpsertTokenAlias) Type() string {
	return MsgProposalUpsertTokenAliasType
}

func (m *MsgProposalUpsertTokenAlias) ValidateBasic() error {
	return nil
}

func (m *MsgProposalUpsertTokenAlias) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m *MsgProposalUpsertTokenAlias) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.Proposer}
}

