package types

import (
	"github.com/asaskevich/govalidator"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/pkg/errors"
)

var _ sdk.Msg = &MsgCreateOrderBook{}

// NewMsgCreateOrderBook create a new message to create an orderbook
func NewMsgCreateOrderBook(
	base string,
	quote string,
	mnemonic string,
	curator sdk.AccAddress,
) (*MsgCreateOrderBook, error) {
	return &MsgCreateOrderBook{
		Base:     base,
		Quote:    quote,
		Mnemonic: mnemonic,
		Curator:  curator,
	}, nil
}

// Route returns route name associated with the message
func (message *MsgCreateOrderBook) Route() string { return ModuleName }

// Type returns transaction type in string
func (message *MsgCreateOrderBook) Type() string { return CreateOrderBookTransaction }

// ValidateBasic do basic validation for message by type
func (message *MsgCreateOrderBook) ValidateBasic() error {
	var _, Error = govalidator.ValidateStruct(message)
	if Error != nil {
		return errors.Wrap(IncorrectMessageCode, Error.Error())
	}
	return nil
}

// GetSignBytes returns to sign bytes for this message
func (message *MsgCreateOrderBook) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(message))
}

// GetSigners return signers to sign this message before broadcast
func (message *MsgCreateOrderBook) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{message.Curator}
}
