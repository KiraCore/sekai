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

func NewMsgAddToCustodyCustodians(addr sdk.AccAddress, newAddr []sdk.AccAddress) *MsgAddToCustodyCustodians {
	return &MsgAddToCustodyCustodians{addr, newAddr}
}

func (m *MsgAddToCustodyCustodians) Route() string {
	return ModuleName
}

func (m *MsgAddToCustodyCustodians) Type() string {
	return types.MsgTypeAddToCustodyWhiteList
}

func (m *MsgAddToCustodyCustodians) ValidateBasic() error {
	return nil
}

func (m *MsgAddToCustodyCustodians) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m *MsgAddToCustodyCustodians) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{
		m.Address,
	}
}

func NewMsgRemoveFromCustodyCustodians(addr sdk.AccAddress, newAddr sdk.AccAddress) *MsgRemoveFromCustodyCustodians {
	return &MsgRemoveFromCustodyCustodians{addr, newAddr}
}

func (m *MsgRemoveFromCustodyCustodians) Route() string {
	return ModuleName
}

func (m *MsgRemoveFromCustodyCustodians) Type() string {
	return types.MsgTypeRemoveFromCustodyWhiteList
}

func (m *MsgRemoveFromCustodyCustodians) ValidateBasic() error {
	return nil
}

func (m *MsgRemoveFromCustodyCustodians) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m *MsgRemoveFromCustodyCustodians) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{
		m.Address,
	}
}

func NewMsgDropCustodyCustodians(addr sdk.AccAddress) *MsgDropCustodyCustodians {
	return &MsgDropCustodyCustodians{addr}
}

func (m *MsgDropCustodyCustodians) Route() string {
	return ModuleName
}

func (m *MsgDropCustodyCustodians) Type() string {
	return types.MsgTypeDropCustodyWhiteList
}

func (m *MsgDropCustodyCustodians) ValidateBasic() error {
	return nil
}

func (m *MsgDropCustodyCustodians) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m *MsgDropCustodyCustodians) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{
		m.Address,
	}
}

func NewMsgAddToCustodyWhiteList(addr sdk.AccAddress, newAddr []sdk.AccAddress) *MsgAddToCustodyWhiteList {
	return &MsgAddToCustodyWhiteList{addr, newAddr}
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

func NewMsgRemoveFromCustodyWhiteList(addr sdk.AccAddress, newAddr sdk.AccAddress) *MsgRemoveFromCustodyWhiteList {
	return &MsgRemoveFromCustodyWhiteList{addr, newAddr}
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

func NewMsgAddToCustodyLimits(addr sdk.AccAddress, denom string, amount uint64, limit string) *MsgAddToCustodyLimits {
	return &MsgAddToCustodyLimits{addr, denom, amount, limit}
}

func (m *MsgAddToCustodyLimits) Route() string {
	return ModuleName
}

func (m *MsgAddToCustodyLimits) Type() string {
	return types.MsgTypeAddToCustodyWhiteList
}

func (m *MsgAddToCustodyLimits) ValidateBasic() error {
	return nil
}

func (m *MsgAddToCustodyLimits) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m *MsgAddToCustodyLimits) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{
		m.Address,
	}
}

func NewMsgRemoveFromCustodyLimits(addr sdk.AccAddress, denom string) *MsgRemoveFromCustodyLimits {
	return &MsgRemoveFromCustodyLimits{addr, denom}
}

func (m *MsgRemoveFromCustodyLimits) Route() string {
	return ModuleName
}

func (m *MsgRemoveFromCustodyLimits) Type() string {
	return types.MsgTypeRemoveFromCustodyWhiteList
}

func (m *MsgRemoveFromCustodyLimits) ValidateBasic() error {
	return nil
}

func (m *MsgRemoveFromCustodyLimits) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m *MsgRemoveFromCustodyLimits) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{
		m.Address,
	}
}

func NewMsgDropCustodyLimits(addr sdk.AccAddress) *MsgDropCustodyLimits {
	return &MsgDropCustodyLimits{addr}
}

func (m *MsgDropCustodyLimits) Route() string {
	return ModuleName
}

func (m *MsgDropCustodyLimits) Type() string {
	return types.MsgTypeDropCustodyWhiteList
}

func (m *MsgDropCustodyLimits) ValidateBasic() error {
	return nil
}

func (m *MsgDropCustodyLimits) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m *MsgDropCustodyLimits) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{
		m.Address,
	}
}

func NewMsgAddToCustodyLimitsStatus(addr sdk.AccAddress, denom string, amount uint64) *MsgAddToCustodyLimitsStatus {
	return &MsgAddToCustodyLimitsStatus{addr, denom, amount}
}

func (m *MsgAddToCustodyLimitsStatus) Route() string {
	return ModuleName
}

func (m *MsgAddToCustodyLimitsStatus) Type() string {
	return types.MsgTypeDropCustodyWhiteList
}

func (m *MsgAddToCustodyLimitsStatus) ValidateBasic() error {
	return nil
}

func (m *MsgAddToCustodyLimitsStatus) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m *MsgAddToCustodyLimitsStatus) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{
		m.Address,
	}
}
