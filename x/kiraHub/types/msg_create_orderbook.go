package types

import (
	"github.com/asaskevich/govalidator"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/pkg/errors"
)

var _ sdk.Msg = &MsgCreateOrderBook{}

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

func (message MsgCreateOrderBook) Route() string { return ModuleName }
func (message MsgCreateOrderBook) Type() string  { return CreateOrderBookTransaction }

func (message MsgCreateOrderBook) ValidateBasic() error {
	var _, Error = govalidator.ValidateStruct(message)
	if Error != nil {
		return errors.Wrap(IncorrectMessageCode, Error.Error())
	}
	return nil
}

func (message MsgCreateOrderBook) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(message))
}

func (message MsgCreateOrderBook) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{message.Curator}
}
