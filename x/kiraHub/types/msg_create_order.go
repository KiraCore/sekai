package types

import (
	sdk "github.com/KiraCore/cosmos-sdk/types"
	"github.com/asaskevich/govalidator"
	"github.com/pkg/errors"
)

var _ sdk.Msg = &MsgCreateOrder{}

func NewMsgCreateOrder(
	orderBookID string,
	orderType LimitOrderType,
	amount int64,
	limitPrice int64,
	curator sdk.AccAddress,
) (*MsgCreateOrder, error) {

	return &MsgCreateOrder{
		OrderBookID: orderBookID,
		OrderType:   orderType,
		Amount:      amount,
		LimitPrice:  limitPrice,
		Curator:     curator,
	}, nil
}

func (message MsgCreateOrder) Route() string { return ModuleName }
func (message MsgCreateOrder) Type() string  { return CreateOrderTransaction }

func (message MsgCreateOrder) ValidateBasic() error {
	var _, Error = govalidator.ValidateStruct(message)
	if Error != nil {
		return errors.Wrap(IncorrectMessageCode, Error.Error())
	}
	return nil
}

func (message MsgCreateOrder) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(message))
}

func (message MsgCreateOrder) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{message.Curator}
}
