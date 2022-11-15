package types

import (
	"github.com/KiraCore/sekai/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.Msg = &MsgCreateCollective{}

// MsgCreateCollective returns an instance of MsgCreateCollective
func NewMMsgCreateCollective(proposer sdk.AccAddress) *MsgCreateCollective {
	return &MsgCreateCollective{
		Sender: proposer.String(),
	}
}

// Route returns route
func (m *MsgCreateCollective) Route() string {
	return ModuleName
}

// Type returns return message type
func (m *MsgCreateCollective) Type() string {
	return types.MsgTypeCreateCollective
}

// ValidateBasic returns basic validation result
func (m *MsgCreateCollective) ValidateBasic() error {
	return nil
}

// GetSignBytes returns to sign bytes
func (m *MsgCreateCollective) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

// GetSigners returns signers
func (m *MsgCreateCollective) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

var _ sdk.Msg = &MsgBondCollective{}

// NewMsgBondCollective returns an instance of MsgBondCollective
func NewMsgBondCollective(proposer sdk.AccAddress) *MsgBondCollective {
	return &MsgBondCollective{
		Sender: proposer.String(),
	}
}

// Route returns route
func (m *MsgBondCollective) Route() string {
	return ModuleName
}

// Type returns return message type
func (m *MsgBondCollective) Type() string {
	return types.MsgTypeBondCollective
}

// ValidateBasic returns basic validation result
func (m *MsgBondCollective) ValidateBasic() error {
	return nil
}

// GetSignBytes returns to sign bytes
func (m *MsgBondCollective) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

// GetSigners returns signers
func (m *MsgBondCollective) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

var _ sdk.Msg = &MsgDonateCollective{}

// NewMsgDonateCollective returns an instance of MsgDonateCollective
func NewMsgDonateCollective(proposer sdk.AccAddress) *MsgDonateCollective {
	return &MsgDonateCollective{
		Sender: proposer.String(),
	}
}

// Route returns route
func (m *MsgDonateCollective) Route() string {
	return ModuleName
}

// Type returns return message type
func (m *MsgDonateCollective) Type() string {
	return types.MsgTypeDonateCollective
}

// ValidateBasic returns basic validation result
func (m *MsgDonateCollective) ValidateBasic() error {
	return nil
}

// GetSignBytes returns to sign bytes
func (m *MsgDonateCollective) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

// GetSigners returns signers
func (m *MsgDonateCollective) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

var _ sdk.Msg = &MsgWithdrawCollective{}

// NewMsgWithdrawCollective returns an instance of MsgWithdrawCollective
func NewMsgWithdrawCollective(proposer sdk.AccAddress) *MsgWithdrawCollective {
	return &MsgWithdrawCollective{
		Sender: proposer.String(),
	}
}

// Route returns route
func (m *MsgWithdrawCollective) Route() string {
	return ModuleName
}

// Type returns return message type
func (m *MsgWithdrawCollective) Type() string {
	return types.MsgTypeWithdrawCollective
}

// ValidateBasic returns basic validation result
func (m *MsgWithdrawCollective) ValidateBasic() error {
	return nil
}

// GetSignBytes returns to sign bytes
func (m *MsgWithdrawCollective) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

// GetSigners returns signers
func (m *MsgWithdrawCollective) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}
