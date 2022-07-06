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

func NewMsgAddToCustodyWhiteList(addr sdk.AccAddress, new_addr sdk.AccAddress) *MsgAddToCustodyWhiteList {
	return &MsgAddToCustodyWhiteList{addr, new_addr}
}

func (m *MsgAddToCustodyWhiteList) Route() string {
	return ModuleName
}

func (m *MsgAddToCustodyWhiteList) Type() string {
	return types.MsgTypeAddToCustodyWhiteList
}

func (m *MsgAddToCustodyWhiteList) ValidateBasic() error {
	return nil
}

func (m *MsgAddToCustodyWhiteList) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m *MsgAddToCustodyWhiteList) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{
		m.Address,
	}
}

func NewMsgRemoveFromCustodyWhiteList(addr sdk.AccAddress, new_addr sdk.AccAddress) *MsgRemoveFromCustodyWhiteList {
	return &MsgRemoveFromCustodyWhiteList{addr, new_addr}
}

func (m *MsgRemoveFromCustodyWhiteList) Route() string {
	return ModuleName
}

func (m *MsgRemoveFromCustodyWhiteList) Type() string {
	return types.MsgTypeRemoveFromCustodyWhiteList
}

func (m *MsgRemoveFromCustodyWhiteList) ValidateBasic() error {
	return nil
}

func (m *MsgRemoveFromCustodyWhiteList) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m *MsgRemoveFromCustodyWhiteList) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{
		m.Address,
	}
}

func NewMsgDropCustodyWhiteList(addr sdk.AccAddress) *MsgDropCustodyWhiteList {
	return &MsgDropCustodyWhiteList{addr}
}

func (m *MsgDropCustodyWhiteList) Route() string {
	return ModuleName
}

func (m *MsgDropCustodyWhiteList) Type() string {
	return types.MsgTypeDropCustodyWhiteList
}

func (m *MsgDropCustodyWhiteList) ValidateBasic() error {
	return nil
}

func (m *MsgDropCustodyWhiteList) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m *MsgDropCustodyWhiteList) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{
		m.Address,
	}
}
