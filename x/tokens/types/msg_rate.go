package types

import (
	"errors"

	"github.com/KiraCore/sekai/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	_ sdk.Msg = &MsgUpsertTokenRate{}
)

// NewMsgUpsertTokenRate returns an instance of MsgUpserTokenRate
func NewMsgUpsertTokenRate(
	proposer sdk.AccAddress,
	denom string,
	rate sdk.Dec,
	feePayments bool,
) *MsgUpsertTokenRate {
	return &MsgUpsertTokenRate{
		Proposer:    proposer,
		Denom:       denom,
		Rate:        rate,
		FeePayments: feePayments,
	}
}

// Route returns route
func (m *MsgUpsertTokenRate) Route() string {
	return ModuleName
}

// Type returns return message type
func (m *MsgUpsertTokenRate) Type() string {
	return types.MsgTypeUpsertTokenRate
}

// ValidateBasic returns basic validation result
func (m *MsgUpsertTokenRate) ValidateBasic() error {
	if m.Denom == "ukex" {
		return errors.New("bond denom rate is read-only")
	}

	if m.Rate.LTE(sdk.NewDec(0)) { // not positive
		return errors.New("rate should be positive")
	}

	return nil
}

// GetSignBytes returns to sign bytes
func (m *MsgUpsertTokenRate) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

// GetSigners returns signers
func (m *MsgUpsertTokenRate) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{
		m.Proposer,
	}
}
