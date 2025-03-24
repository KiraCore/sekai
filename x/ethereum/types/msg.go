package types

import (
	"github.com/KiraCore/sekai/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewMsgRelay(addr sdk.AccAddress, data string) *MsgRelay {
	return &MsgRelay{addr, data}
}

func (m *MsgRelay) Route() string {
	return ModuleName
}

func (m *MsgRelay) Type() string {
	return types.MsgTypeCreateCustody
}

func (m *MsgRelay) ValidateBasic() error {
	return nil
}

func (m *MsgRelay) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m *MsgRelay) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{
		m.Address,
	}
}
