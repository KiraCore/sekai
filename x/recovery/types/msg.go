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
func (msg MsgRegisterRecoverySecret) Type() string  { return types.MsgTypeRegisterRecoverySecret }
func (msg MsgRegisterRecoverySecret) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Address)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr.Bytes()}
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
func NewMsgRotateRecoveryAddress(feePayer, addr, recovery, proof string) *MsgRotateRecoveryAddress {
	return &MsgRotateRecoveryAddress{
		FeePayer: feePayer,
		Address:  addr,
		Recovery: recovery,
		Proof:    proof,
	}
}

func (msg MsgRotateRecoveryAddress) Route() string { return RouterKey }
func (msg MsgRotateRecoveryAddress) Type() string  { return types.MsgTypeRotateRecoveryAddress }
func (msg MsgRotateRecoveryAddress) GetSigners() []sdk.AccAddress {
	feePayer, err := sdk.AccAddressFromBech32(msg.FeePayer)
	if err != nil {
		panic(err)
	}
	addr, err := sdk.AccAddressFromBech32(msg.Address)
	if err != nil {
		panic(err)
	}

	if msg.Address == msg.FeePayer {
		return []sdk.AccAddress{addr.Bytes()}
	}
	return []sdk.AccAddress{feePayer.Bytes(), addr.Bytes()}
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
func (msg MsgIssueRecoveryTokens) Type() string  { return types.MsgTypeIssueRecoveryTokens }
func (msg MsgIssueRecoveryTokens) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Address)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr.Bytes()}
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
func NewMsgBurnRecoveryTokens(sender sdk.AccAddress, rrCoin sdk.Coin) *MsgBurnRecoveryTokens {
	return &MsgBurnRecoveryTokens{
		Address: sender.String(),
		RrCoin:  rrCoin,
	}
}

func (msg MsgBurnRecoveryTokens) Route() string { return RouterKey }
func (msg MsgBurnRecoveryTokens) Type() string  { return types.MsgTypeBurnRecoveryTokens }
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

// verify interface at compile time
var _ sdk.Msg = &MsgRegisterRRTokenHolder{}

// NewMsgRegisterRRTokenHolder creates a new MsgRegisterRRTokenHolder instance
func NewMsgRegisterRRTokenHolder(sender sdk.AccAddress) *MsgRegisterRRTokenHolder {
	return &MsgRegisterRRTokenHolder{
		Holder: sender.String(),
	}
}

func (msg MsgRegisterRRTokenHolder) Route() string { return RouterKey }
func (msg MsgRegisterRRTokenHolder) Type() string  { return types.MsgTypeRegisterRRTokenHolder }
func (msg MsgRegisterRRTokenHolder) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Holder)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

// GetSignBytes gets the bytes for the message signer to sign on
func (msg MsgRegisterRRTokenHolder) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic validity check for the AnteHandler
func (msg MsgRegisterRRTokenHolder) ValidateBasic() error {
	if msg.Holder == "" {
		return ErrInvalidAccAddress
	}

	return nil
}

// verify interface at compile time
var _ sdk.Msg = &MsgClaimRRHolderRewards{}

// NewMsgClaimRRHolderRewards creates a new MsgClaimRRHolderRewards instance
func NewMsgClaimRRHolderRewards(sender sdk.AccAddress) *MsgClaimRRHolderRewards {
	return &MsgClaimRRHolderRewards{
		Sender: sender.String(),
	}
}

func (msg MsgClaimRRHolderRewards) Route() string { return RouterKey }
func (msg MsgClaimRRHolderRewards) Type() string  { return types.MsgTypeClaimRRHolderRewards }
func (msg MsgClaimRRHolderRewards) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

// GetSignBytes gets the bytes for the message signer to sign on
func (msg MsgClaimRRHolderRewards) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic validity check for the AnteHandler
func (msg MsgClaimRRHolderRewards) ValidateBasic() error {
	if msg.Sender == "" {
		return ErrInvalidAccAddress
	}

	return nil
}

// verify interface at compile time
var _ sdk.Msg = &MsgRotateValidatorByHalfRRTokenHolder{}

// NewMsgRotateValidatorByHalfRRTokenHolder creates a new MsgRotateValidatorByHalfRRTokenHolder instance
//nolint:interfacer
func NewMsgRotateValidatorByHalfRRTokenHolder(rrHolder, addr, recovery string) *MsgRotateValidatorByHalfRRTokenHolder {
	return &MsgRotateValidatorByHalfRRTokenHolder{
		RrHolder: rrHolder,
		Address:  addr,
		Recovery: recovery,
	}
}

func (msg MsgRotateValidatorByHalfRRTokenHolder) Route() string { return RouterKey }
func (msg MsgRotateValidatorByHalfRRTokenHolder) Type() string {
	return types.MsgTypeRotateValidatorByHalfRRTokenHolder
}
func (msg MsgRotateValidatorByHalfRRTokenHolder) GetSigners() []sdk.AccAddress {
	rrHolder, err := sdk.AccAddressFromBech32(msg.RrHolder)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{rrHolder.Bytes()}
}

// GetSignBytes gets the bytes for the message signer to sign on
func (msg MsgRotateValidatorByHalfRRTokenHolder) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic validity check for the AnteHandler
func (msg MsgRotateValidatorByHalfRRTokenHolder) ValidateBasic() error {
	if msg.Address == "" {
		return ErrInvalidAccAddress
	}

	return nil
}
