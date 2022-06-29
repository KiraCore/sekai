package types

import (
	"github.com/KiraCore/sekai/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewMsgCreateCustody(addr sdk.AccAddress, custodySettings CustodySettings) *MsgCreteCustodyRecord {
	return &MsgCreteCustodyRecord{addr, custodySettings}
}

func (m *MsgCreteCustodyRecord) Route() string {
	return ModuleName
}

func (m *MsgCreteCustodyRecord) Type() string {
	return types.MsgTypeCreateCustody
}

func (m *MsgCreteCustodyRecord) ValidateBasic() error {
	return nil
}

func (m *MsgCreteCustodyRecord) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m *MsgCreteCustodyRecord) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{
		m.Address,
	}
}
