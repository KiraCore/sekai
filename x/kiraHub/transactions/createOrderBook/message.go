package createOrderBook

import (
	sdk "github.com/KiraCore/cosmos-sdk/types"
	errors "github.com/KiraCore/cosmos-sdk/types/errors"
	"github.com/asaskevich/govalidator"

	constants "github.com/KiraCore/sekai/x/kiraHub/constants"
)

type Message struct {
	Base string		   	   `json:"base"  yaml:"base"  valid:"required~base"`
	Quote string		   `json:"quote" yaml:"quote" valid:"required~quote"`
	Mnemonic string 	   `json:"mnemonic" yaml:"mnemonic" valid:"required~mnemonic"`
	Curator string 		   `json:"curator"  yaml:"curator" valid:"required~curator"`
}

var _ sdk.Msg = Message{}

func (message Message) Route() string { return constants.ModuleName }
func (message Message) Type() string  { return constants.CreateOrderBookTransaction }

func (message Message) ValidateBasic() error {
	var _, Error = govalidator.ValidateStruct(message)
	if Error != nil {
		return errors.Wrap(constants.IncorrectMessageCode, Error.Error())
	}
	return nil
}

func (message Message) GetSignBytes() []byte {
	return sdk.MustSortJSON(PackageCodec.MustMarshalJSON(message))
}

func (message Message) GetSigners() []sdk.AccAddress {
	var curator, _ = sdk.AccAddressFromBech32(message.Curator)
	return []sdk.AccAddress{curator}
}