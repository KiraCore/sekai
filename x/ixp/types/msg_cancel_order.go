package types

import (
	"github.com/asaskevich/govalidator"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/pkg/errors"
)

var _ sdk.Msg = &MsgCancelOrder{}

// NewMsgCancelOrder create a new message to cancel order
func NewMsgCancelOrder(orderID string, curator sdk.AccAddress) (*MsgCancelOrder, error) {
	return &MsgCancelOrder{
		OrderID: orderID,
		Curator: curator,
	}, nil
}

// Route returns route name associated with the message
func (message *MsgCancelOrder) Route() string { return ModuleName }

// Type returns transaction type in string
func (message *MsgCancelOrder) Type() string { return CancelOrderTransaction }

// ValidateBasic do basic validation for message by type
func (message *MsgCancelOrder) ValidateBasic() error {
	var _, Error = govalidator.ValidateStruct(message)
	if Error != nil {
		return errors.Wrap(IncorrectMessageCode, Error.Error())
	}
	return nil
}

// GetSignBytes returns to sign bytes for this message
func (message *MsgCancelOrder) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(message))
}

// GetSigners return signers to sign this message before broadcast
func (message *MsgCancelOrder) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{message.Curator}
}
