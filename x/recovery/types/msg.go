package types

import (
	"github.com/KiraCore/sekai/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// verify interface at compile time
var _ sdk.Msg = &MsgRegisterRecoverySecret{}

// NewMsgRegisterRecoverySecret creates a new MsgRegisterRecoverySecret instance
//nolint:interfacer
func NewMsgRegisterRecoverySecret(addr, challenge, nonce, proof string) *MsgRegisterRecoverySecret {
	return &MsgRegisterRecoverySecret{
		Address:   addr,
		Challenge: challenge,
		Nonce:     nonce,
		Proof:     proof,
	}
}

func (msg MsgRegisterRecoverySecret) Route() string { return RouterKey }
func (msg MsgRegisterRecoverySecret) Type() string  { return types.MsgTypeActivate }
func (msg MsgRegisterRecoverySecret) GetSigners() []sdk.AccAddress {
	valAddr, err := sdk.AccAddressFromBech32(msg.Address)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{valAddr.Bytes()}
}

// GetSignBytes gets the bytes for the message signer to sign on
func (msg MsgRegisterRecoverySecret) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic validity check for the AnteHandler
func (msg MsgRegisterRecoverySecret) ValidateBasic() error {
	if msg.Address == "" {
		return ErrInvalidAccAddress
	}

	return nil
}

// verify interface at compile time
var _ sdk.Msg = &MsgRotateRecoveryAddress{}

// NewMsgRotateRecoveryAddress creates a new MsgRotateRecoveryAddress instance
//nolint:interfacer
func NewMsgRotateRecoveryAddress(addr, recovery, proof string) *MsgRotateRecoveryAddress {
	return &MsgRotateRecoveryAddress{
		Address:  addr,
		Recovery: recovery,
		Proof:    proof,
	}
}

func (msg MsgRotateRecoveryAddress) Route() string { return RouterKey }
func (msg MsgRotateRecoveryAddress) Type() string  { return types.MsgTypeUnpause }
func (msg MsgRotateRecoveryAddress) GetSigners() []sdk.AccAddress {
	valAddr, err := sdk.AccAddressFromBech32(msg.Address)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{valAddr.Bytes()}
}

// GetSignBytes gets the bytes for the message signer to sign on
func (msg MsgRotateRecoveryAddress) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic validity check for the AnteHandler
func (msg MsgRotateRecoveryAddress) ValidateBasic() error {
	if msg.Address == "" {
		return ErrInvalidAccAddress
	}

	return nil
}

// verify interface at compile time
var _ sdk.Msg = &MsgIssueRecoveryTokens{}

// NewMsgIssueRecoveryTokens creates a new MsgIssueRecoveryTokens instance
//nolint:interfacer
func NewMsgIssueRecoveryTokens(addr string) *MsgIssueRecoveryTokens {
	return &MsgIssueRecoveryTokens{
		Address: addr,
	}
}

func (msg MsgIssueRecoveryTokens) Route() string { return RouterKey }
func (msg MsgIssueRecoveryTokens) Type() string  { return types.MsgTypePause }
func (msg MsgIssueRecoveryTokens) GetSigners() []sdk.AccAddress {
	valAddr, err := sdk.AccAddressFromBech32(msg.Address)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{valAddr.Bytes()}
}

// GetSignBytes gets the bytes for the message signer to sign on
func (msg MsgIssueRecoveryTokens) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic validity check for the AnteHandler
func (msg MsgIssueRecoveryTokens) ValidateBasic() error {
	if msg.Address == "" {
		return ErrInvalidAccAddress
	}

	return nil
}

// verify interface at compile time
var _ sdk.Msg = &MsgBurnRecoveryTokens{}

// NewMsgBurnRecoveryTokens creates a new MsgBurnRecoveryTokens instance
func NewMsgBurnRecoveryTokens(sender sdk.AccAddress) *MsgBurnRecoveryTokens {
	return &MsgBurnRecoveryTokens{
		Address: sender.String(),
	}
}

func (msg MsgBurnRecoveryTokens) Route() string { return RouterKey }
func (msg MsgBurnRecoveryTokens) Type() string  { return types.MsgTypePause }
func (msg MsgBurnRecoveryTokens) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Address)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

// GetSignBytes gets the bytes for the message signer to sign on
func (msg MsgBurnRecoveryTokens) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic validity check for the AnteHandler
func (msg MsgBurnRecoveryTokens) ValidateBasic() error {
	if msg.Address == "" {
		return ErrInvalidAccAddress
	}

	return nil
}
