package types

import (
	"errors"

	appparams "github.com/KiraCore/sekai/app/params"
	"github.com/KiraCore/sekai/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	_ sdk.Msg = &MsgUpsertTokenInfo{}
)

// NewMsgUpsertTokenInfo returns an instance of MsgUpserTokenInfo
func NewMsgUpsertTokenInfo(
	proposer sdk.AccAddress,
	denom string,
	rate sdk.Dec,
	feeEnabled bool,
	stakeCap sdk.Dec,
	stakeMin sdk.Int,
	stakeEnabled bool,
	inactive bool,
	symbol string,
	name string,
	icon string,
	decimals uint32,
) *MsgUpsertTokenInfo {
	return &MsgUpsertTokenInfo{
		Proposer:     proposer,
		Denom:        denom,
		Rate:         rate,
		FeeEnabled:   feeEnabled,
		StakeCap:     stakeCap,
		StakeMin:     stakeMin,
		StakeEnabled: stakeEnabled,
		Inactive:     inactive,
		Symbol:       symbol,
		Name:         name,
		Icon:         icon,
		Decimals:     decimals,
	}
}

// Route returns route
func (m *MsgUpsertTokenInfo) Route() string {
	return ModuleName
}

// Type returns return message type
func (m *MsgUpsertTokenInfo) Type() string {
	return types.MsgTypeUpsertTokenInfo
}

// ValidateBasic returns basic validation result
func (m *MsgUpsertTokenInfo) ValidateBasic() error {
	if m.Denom == appparams.DefaultDenom {
		return errors.New("bond denom rate is read-only")
	}

	if m.Rate.LTE(sdk.NewDec(0)) { // not positive
		return errors.New("rate should be positive")
	}

	if m.StakeCap.LT(sdk.NewDec(0)) { // not positive
		return errors.New("reward cap should be positive")
	}

	if m.StakeCap.GT(sdk.OneDec()) { // more than 1
		return errors.New("reward cap not be more than 100%")
	}

	return nil
}

// GetSignBytes returns to sign bytes
func (m *MsgUpsertTokenInfo) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

// GetSigners returns signers
func (m *MsgUpsertTokenInfo) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{
		m.Proposer,
	}
}
