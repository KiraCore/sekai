package createOrderBook

import (
	sdk "github.com/KiraCore/cosmos-sdk/types"
	"github.com/KiraCore/cosmos-sdk/types/errors"
	"github.com/asaskevich/govalidator"

	constants "github.com/KiraCore/sekai/x/kiraHub/constants"
)

type Message struct {
	Base string		   	   `json:"base"  yaml:"base"  valid:"required~Base is required"`
	Quote string		   `json:"quote" yaml:"quote" valid:"required~Quote is required"`
	Mnemonic string 	   `json:"mnemonic" yaml:"mnemonic" valid:"required~Mnemonic is required"`
	Curator sdk.AccAddress `json:"curator"  yaml:"curator" valid:"required~Curator is required"`
}

func (message Message) Reset() {
	panic("implement me")
}

func (message Message) String() string {
	panic("implement me")
}

func (message Message) ProtoMessage() {
	panic("implement me")
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
	return []sdk.AccAddress{message.Curator}
}