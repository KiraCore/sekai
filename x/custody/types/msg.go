package types

import (
	errorsmod "cosmossdk.io/errors"
	"github.com/KiraCore/sekai/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func NewMsgCreateCustody(addr sdk.AccAddress, custodySettings CustodySettings, oldKey, newKey, next, target string) *MsgCreateCustodyRecord {
	return &MsgCreateCustodyRecord{addr, custodySettings, oldKey, newKey, next, target}
}

func (m *MsgCreateCustodyRecord) Route() string {
	return ModuleName
}

func (m *MsgCreateCustodyRecord) Type() string {
	return types.MsgTypeCreateCustody
}

func (m *MsgCreateCustodyRecord) ValidateBasic() error {
	return nil
}

func (m *MsgCreateCustodyRecord) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m *MsgCreateCustodyRecord) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{
		m.Address,
	}
}

func NewMsgDisableCustody(addr sdk.AccAddress, oldKey, newKey, next, target string) *MsgDisableCustodyRecord {
	return &MsgDisableCustodyRecord{addr, oldKey, newKey, next, target}
}

func (m *MsgDisableCustodyRecord) Route() string {
	return ModuleName
}

func (m *MsgDisableCustodyRecord) Type() string {
	return types.MsgTypeDisableCustody
}

func (m *MsgDisableCustodyRecord) ValidateBasic() error {
	return nil
}

func (m *MsgDisableCustodyRecord) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m *MsgDisableCustodyRecord) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{
		m.Address,
	}
}

func NewMsgDropCustody(addr sdk.AccAddress, oldKey, target string) *MsgDropCustodyRecord {
	return &MsgDropCustodyRecord{addr, oldKey, target}
}

func (m *MsgDropCustodyRecord) Route() string {
	return ModuleName
}

func (m *MsgDropCustodyRecord) Type() string {
	return types.MsgTypeDropCustody
}

func (m *MsgDropCustodyRecord) ValidateBasic() error {
	return nil
}

func (m *MsgDropCustodyRecord) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m *MsgDropCustodyRecord) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{
		m.Address,
	}
}

func NewMsgApproveCustodyTransaction(from sdk.AccAddress, to sdk.AccAddress, hash string) *MsgApproveCustodyTransaction {
	return &MsgApproveCustodyTransaction{from, to, hash}
}

func (m *MsgApproveCustodyTransaction) Route() string {
	return ModuleName
}

func (m *MsgApproveCustodyTransaction) Type() string {
	return types.MsgTypeAddToCustodyCustodians
}

func (m *MsgApproveCustodyTransaction) ValidateBasic() error {
	return nil
}

func (m *MsgApproveCustodyTransaction) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m *MsgApproveCustodyTransaction) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{
		m.FromAddress,
	}
}

func NewMsgDeclineCustodyTransaction(from sdk.AccAddress, to sdk.AccAddress, hash string) *MsgDeclineCustodyTransaction {
	return &MsgDeclineCustodyTransaction{from, to, hash}
}

func (m *MsgDeclineCustodyTransaction) Route() string {
	return ModuleName
}

func (m *MsgDeclineCustodyTransaction) Type() string {
	return types.MsgTypeAddToCustodyCustodians
}

func (m *MsgDeclineCustodyTransaction) ValidateBasic() error {
	return nil
}

func (m *MsgDeclineCustodyTransaction) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m *MsgDeclineCustodyTransaction) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{
		m.FromAddress,
	}
}

func NewMsgPasswordConfirmTransaction(from sdk.AccAddress, sender sdk.AccAddress, hash string, password string) *MsgPasswordConfirmTransaction {
	return &MsgPasswordConfirmTransaction{from, sender, hash, password}
}

func (m *MsgPasswordConfirmTransaction) Route() string {
	return ModuleName
}

func (m *MsgPasswordConfirmTransaction) Type() string {
	return types.MsgPasswordConfirmTransaction
}

func (m *MsgPasswordConfirmTransaction) ValidateBasic() error {
	return nil
}

func (m *MsgPasswordConfirmTransaction) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m *MsgPasswordConfirmTransaction) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{
		m.FromAddress,
	}
}

func NewMsgAddToCustodyCustodians(addr sdk.AccAddress, newAddr []sdk.AccAddress, oldKey, newKey, next, target string) *MsgAddToCustodyCustodians {
	return &MsgAddToCustodyCustodians{addr, newAddr, oldKey, newKey, next, target}
}

func (m *MsgAddToCustodyCustodians) Route() string {
	return ModuleName
}

func (m *MsgAddToCustodyCustodians) Type() string {
	return types.MsgTypeAddToCustodyCustodians
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

func NewMsgRemoveFromCustodyCustodians(addr sdk.AccAddress, newAddr sdk.AccAddress, oldKey, newKey, next, target string) *MsgRemoveFromCustodyCustodians {
	return &MsgRemoveFromCustodyCustodians{addr, newAddr, oldKey, newKey, next, target}
}

func (m *MsgRemoveFromCustodyCustodians) Route() string {
	return ModuleName
}

func (m *MsgRemoveFromCustodyCustodians) Type() string {
	return types.MsgTypeRemoveFromCustodyCustodians
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

func NewMsgDropCustodyCustodians(addr sdk.AccAddress, oldKey, newKey, next, target string) *MsgDropCustodyCustodians {
	return &MsgDropCustodyCustodians{addr, oldKey, newKey, next, target}
}

func (m *MsgDropCustodyCustodians) Route() string {
	return ModuleName
}

func (m *MsgDropCustodyCustodians) Type() string {
	return types.MsgTypeDropCustodyCustodians
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

func NewMsgAddToCustodyWhiteList(addr sdk.AccAddress, newAddr []sdk.AccAddress, oldKey, newKey, next, target string) *MsgAddToCustodyWhiteList {
	return &MsgAddToCustodyWhiteList{addr, newAddr, oldKey, newKey, next, target}
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

func NewMsgRemoveFromCustodyWhiteList(addr sdk.AccAddress, newAddr sdk.AccAddress, oldKey, newKey, next, target string) *MsgRemoveFromCustodyWhiteList {
	return &MsgRemoveFromCustodyWhiteList{addr, newAddr, oldKey, newKey, next, target}
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

func NewMsgDropCustodyWhiteList(addr sdk.AccAddress, oldKey, newKey, next, target string) *MsgDropCustodyWhiteList {
	return &MsgDropCustodyWhiteList{addr, oldKey, newKey, next, target}
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

func NewMsgAddToCustodyLimits(addr sdk.AccAddress, denom string, amount uint64, limit, oldKey, newKey, next, target string) *MsgAddToCustodyLimits {
	return &MsgAddToCustodyLimits{addr, denom, amount, limit, oldKey, newKey, next, target}
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

func NewMsgRemoveFromCustodyLimits(addr sdk.AccAddress, denom, oldKey, newKey, next, target string) *MsgRemoveFromCustodyLimits {
	return &MsgRemoveFromCustodyLimits{addr, denom, oldKey, newKey, next, target}
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

func NewMsgDropCustodyLimits(addr sdk.AccAddress, oldKey, newKey, next, target string) *MsgDropCustodyLimits {
	return &MsgDropCustodyLimits{addr, oldKey, newKey, next, target}
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

func NewMsgSend(fromAddr, toAddr sdk.AccAddress, amount sdk.Coins, password string, reward sdk.Coins) *MsgSend {
	return &MsgSend{FromAddress: fromAddr.String(), ToAddress: toAddr.String(), Amount: amount, Password: password, Reward: reward}
}

func (m MsgSend) Route() string {
	return ModuleName
}

func (m MsgSend) Type() string {
	return types.MsgTypeSend
}

func (m MsgSend) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.FromAddress)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid sender address (%s)", err)
	}

	_, err = sdk.AccAddressFromBech32(m.ToAddress)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid recipient address (%s)", err)
	}

	if !m.Amount.IsValid() {
		return errorsmod.Wrap(sdkerrors.ErrInvalidCoins, m.Amount.String())
	}

	if !m.Amount.IsAllPositive() {
		return errorsmod.Wrap(sdkerrors.ErrInvalidCoins, m.Amount.String())
	}

	return nil
}

func (m MsgSend) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

func (m MsgSend) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(m.FromAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}
