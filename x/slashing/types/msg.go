package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// slashing message types
const (
	TypeMsgActivate = "activate"
)

// verify interface at compile time
var _ sdk.Msg = &MsgActivate{}

// NewMsgActivate creates a new MsgActivate instance
//nolint:interfacer
func NewMsgActivate(validatorAddr sdk.ValAddress) *MsgActivate {
	return &MsgActivate{
		ValidatorAddr: validatorAddr.String(),
	}
}

func (msg MsgActivate) Route() string { return RouterKey }
func (msg MsgActivate) Type() string  { return TypeMsgActivate }
func (msg MsgActivate) GetSigners() []sdk.AccAddress {
	valAddr, err := sdk.ValAddressFromBech32(msg.ValidatorAddr)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{valAddr.Bytes()}
}

// GetSignBytes gets the bytes for the message signer to sign on
func (msg MsgActivate) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic validity check for the AnteHandler
func (msg MsgActivate) ValidateBasic() error {
	if msg.ValidatorAddr == "" {
		return ErrBadValidatorAddr
	}

	return nil
}
