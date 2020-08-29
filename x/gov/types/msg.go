package types

import sdk "github.com/cosmos/cosmos-sdk/types"

var _ sdk.Msg = &MsgWhitelistPermissions{}

func (m *MsgWhitelistPermissions) Route() string {
	return ModuleName
}

func (m *MsgWhitelistPermissions) Type() string {
	panic("implement me")
}

func (m *MsgWhitelistPermissions) ValidateBasic() error {
	panic("implement me")
}

func (m *MsgWhitelistPermissions) GetSignBytes() []byte {
	panic("implement me")
}

func (m *MsgWhitelistPermissions) GetSigners() []sdk.AccAddress {
	panic("implement me")
}

