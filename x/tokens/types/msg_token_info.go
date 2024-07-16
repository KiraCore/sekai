package types

import (
	"errors"

	"cosmossdk.io/math"
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
	tokenType string,
	feeRate math.LegacyDec,
	feeEnabled bool,
	supply math.Int,
	supplyCap math.Int,
	stakeCap math.LegacyDec,
	stakeMin math.Int,
	stakeEnabled bool,
	inactive bool,
	symbol string,
	name string,
	icon string,
	decimals uint32,
	description string,
	website string,
	social string,
	holders uint64,
	mintingFee math.Int,
	owner string,
	ownerEditDisabled bool,
	nftMetadata string,
	nftHash string,
) *MsgUpsertTokenInfo {
	return &MsgUpsertTokenInfo{
		Proposer:          proposer,
		Denom:             denom,
		TokenType:         tokenType,
		FeeRate:           feeRate,
		FeeEnabled:        feeEnabled,
		Supply:            supply,
		SupplyCap:         supplyCap,
		StakeCap:          stakeCap,
		StakeMin:          stakeMin,
		StakeEnabled:      stakeEnabled,
		Inactive:          inactive,
		Symbol:            symbol,
		Name:              name,
		Icon:              icon,
		Decimals:          decimals,
		Description:       description,
		Website:           website,
		Social:            social,
		Holders:           holders,
		MintingFee:        mintingFee,
		Owner:             owner,
		OwnerEditDisabled: ownerEditDisabled,
		NftMetadata:       nftMetadata,
		NftHash:           nftHash,
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

	if !m.FeeRate.IsPositive() { // not positive
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
