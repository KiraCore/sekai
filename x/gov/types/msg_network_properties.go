package types

import sdk "github.com/cosmos/cosmos-sdk/types"

var _ sdk.Msg = &MsgSetNetworkProperties{}

func NewMsgSetNetworkProperties(
	proposer sdk.AccAddress,
	properties *NetworkProperties,
) *MsgSetNetworkProperties {
	return &MsgSetNetworkProperties{
		Proposer:          proposer,
		NetworkProperties: properties,
	}
}

func (m *MsgSetNetworkProperties) Route() string {
	return ModuleName
}

func (m *MsgSetNetworkProperties) Type() string {
	return WhitelistPermissions
}

func (m *MsgSetNetworkProperties) ValidateBasic() error {
	if m.Proposer.Empty() {
		return ErrEmptyProposerAccAddress
	}

	return nil
}

func (m *MsgSetNetworkProperties) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m *MsgSetNetworkProperties) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{
		m.Proposer,
	}
}
