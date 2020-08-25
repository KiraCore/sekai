package types

import (
	"github.com/asaskevich/govalidator"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/pkg/errors"
)

var _ sdk.Msg = &MsgCancelOrder{}

func NewMsgCancelOrder(orderID string) (*MsgCancelOrder, error) {
	return &MsgCancelOrder{
		OrderID: orderID,
	}, nil
}

func (message MsgCancelOrder) Route() string { return ModuleName }
func (message MsgCancelOrder) Type() string  { return CreateOrderTransaction }

func (message MsgCancelOrder) ValidateBasic() error {
	var _, Error = govalidator.ValidateStruct(message)
	if Error != nil {
		return errors.Wrap(IncorrectMessageCode, Error.Error())
	}
	return nil
}

func (message MsgCancelOrder) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(message))
}

func (message MsgCancelOrder) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{message.Curator}
}
