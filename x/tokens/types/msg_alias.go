package types

import (
	"github.com/KiraCore/sekai/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	_ sdk.Msg = &MsgUpsertTokenAlias{}
)

// NewMsgUpsertTokenAlias returns an instance of MsgUpsertTokenAlias
func NewMsgUpsertTokenAlias(
	proposer sdk.AccAddress,
	symbol string,
	name string,
	icon string,
	decimals uint32,
	denoms []string,
) *MsgUpsertTokenAlias {
	return &MsgUpsertTokenAlias{
		Proposer: proposer,
		Symbol:   symbol,
		Name:     name,
		Icon:     icon,
		Decimals: decimals,
		Denoms:   denoms,
	}
}

// Route returns route
func (m *MsgUpsertTokenAlias) Route() string {
	return ModuleName
}

// Type returns return message type
func (m *MsgUpsertTokenAlias) Type() string {
	return types.MsgTypeUpsertTokenAlias
}

// ValidateBasic returns basic validation result
func (m *MsgUpsertTokenAlias) ValidateBasic() error {
	return nil
}

// GetSignBytes returns to sign bytes
func (m *MsgUpsertTokenAlias) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

// GetSigners returns signers
func (m *MsgUpsertTokenAlias) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{
		m.Proposer,
	}
}
