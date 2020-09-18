package types

import (
	"github.com/asaskevich/govalidator"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/pkg/errors"
)

var _ sdk.Msg = &MsgCreateOrder{}

// NewMsgCreateOrder create a new message to create an order
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

// Route returns route name associated with the message
func (message *MsgCreateOrder) Route() string { return ModuleName }

// Type returns transaction type in string
func (message *MsgCreateOrder) Type() string { return CreateOrderTransaction }

// ValidateBasic do basic validation for message by type
func (message *MsgCreateOrder) ValidateBasic() error {
	var _, Error = govalidator.ValidateStruct(message)
	if Error != nil {
		return errors.Wrap(IncorrectMessageCode, Error.Error())
	}
	return nil
}

// GetSignBytes returns to sign bytes for this message
func (message *MsgCreateOrder) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(message))
}

// GetSigners return signers to sign this message before broadcast
func (message *MsgCreateOrder) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{message.Curator}
}
