package createOrder

import (
	sdk "github.com/KiraCore/cosmos-sdk/types"
	"github.com/KiraCore/cosmos-sdk/types/errors"
	constants "github.com/KiraCore/sekai/x/kiraHub/constants"
	"github.com/asaskevich/govalidator"
)

type Message struct {
	OrderBookID string         `json:"order_book_id" yaml:"order_book_id"  valid:"required~OrderBookID is required"`
	OrderType   uint8          `json:"order_type" yaml:"order_type" valid:"required~OrderType is required"`
	Amount      int64          `json:"amount" yaml:"amount" valid:"required~Amount is required"`
	LimitPrice  int64          `json:"limit_price"  yaml:"limit_price" valid:"required~Limit Price is required"`
	Curator     sdk.AccAddress `json:"curator"  yaml:"curator" valid:"required~Curator is required"`
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
func (message Message) Type() string  { return constants.CreateOrderTransaction }

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
