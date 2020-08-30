package types

import sdk "github.com/cosmos/cosmos-sdk/types"

var _ sdk.Msg = &MsgWhitelistPermissions{}

func NewMsgWhitelistPermissions(
	proposer, address sdk.AccAddress,
	permission uint32,
) *MsgWhitelistPermissions {
	return &MsgWhitelistPermissions{
		Proposer:   proposer,
		Address:    address,
		Permission: permission,
	}
}

func (m *MsgWhitelistPermissions) Route() string {
	return ModuleName
}

func (m *MsgWhitelistPermissions) Type() string {
	return WhitelistPermissions
}

func (m *MsgWhitelistPermissions) ValidateBasic() error {
	if m.Proposer.Empty() {
		return ErrEmptyProposerAccAddress
	}

	if m.Address.Empty() {
		return ErrEmptyPermissionsAccAddress
	}

	return nil
}

func (m *MsgWhitelistPermissions) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m *MsgWhitelistPermissions) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{
		m.Proposer,
	}
}
