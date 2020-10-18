package types

import (
	"errors"
	"math"
	"strconv"

	customgovtypes "github.com/KiraCore/sekai/x/gov/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Msg types
const (
	UpsertTokenRate = customgovtypes.UpsertTokenRate
)

var (
	_ sdk.Msg = &MsgUpsertTokenRate{}
)

// NewMsgUpsertTokenRate returns an instance of MsgUpserTokenRate
func NewMsgUpsertTokenRate(
	proposer sdk.AccAddress,
	denom string,
	rate string,
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
	return UpsertTokenRate
}

// ValidateBasic returns basic validation result
func (m *MsgUpsertTokenRate) ValidateBasic() error {
	rateFloat, err := strconv.ParseFloat(m.Rate, 64)
	if err != nil {
		return err
	}

	if rateFloat <= 0 {
		return errors.New("rate should be positive")
	}

	if rateFloat >= RateMaximum {
		return errors.New("rate is larger than maximum")
	}

	rateRawFloat := rateFloat * RateDecimalDenominator
	if math.Floor(rateRawFloat) != rateRawFloat {
		return errors.New("decimal is bigger than maximum decimal")
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
