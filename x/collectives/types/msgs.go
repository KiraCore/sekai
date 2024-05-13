package types

import (
	"github.com/KiraCore/sekai/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.Msg = &MsgCreateCollective{}

// NewMsgCreateCollective returns an instance of MsgCreateCollective
func NewMsgCreateCollective(
	proposer sdk.AccAddress,
	name, description string,
	bonds sdk.Coins,
	depositWhitelist DepositWhitelist,
	ownersWhitelist OwnersWhitelist,
	weightedSpendingPool []WeightedSpendingPool,
	claimStart, claimPeriod, claimEnd uint64,
	voteQuorum sdk.Dec,
	votePeriod, voteEnactment uint64,
) *MsgCreateCollective {
	return &MsgCreateCollective{
		Sender:           proposer.String(),
		Name:             name,
		Description:      description,
		Bonds:            bonds,
		DepositWhitelist: depositWhitelist,
		OwnersWhitelist:  ownersWhitelist,
		SpendingPools:    weightedSpendingPool,
		ClaimStart:       claimStart,
		ClaimPeriod:      claimPeriod,
		ClaimEnd:         claimEnd,
		VoteQuorum:       voteQuorum,
		VotePeriod:       votePeriod,
		VoteEnactment:    voteEnactment,
	}
}

// Route returns route
func (m *MsgCreateCollective) Route() string {
	return ModuleName
}

// Type returns return message type
func (m *MsgCreateCollective) Type() string {
	return types.MsgTypeCreateCollective
}

// ValidateBasic returns basic validation result
func (m *MsgCreateCollective) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		return err
	}

	totalWeight := sdk.ZeroDec()
	for _, wpool := range m.SpendingPools {
		totalWeight = totalWeight.Add(wpool.Weight)
	}
	if !totalWeight.Equal(sdk.OneDec()) {
		return ErrTotalSpendingPoolWeightShouldBeOne
	}

	if m.Name == "" {
		return ErrEmptyCollectiveName
	}
	if sdk.Coins(m.Bonds).Empty() {
		return ErrEmptyCollectiveBonds
	}
	if len(m.OwnersWhitelist.Accounts) == 0 && len(m.OwnersWhitelist.Roles) == 0 {
		return ErrEmptyOwnersList
	}
	if len(m.SpendingPools) == 0 {
		return ErrEmptySpendingPools
	}
	if m.ClaimPeriod == 0 {
		return ErrInvalidClaimPeriod
	}
	if m.VoteQuorum.IsZero() {
		return ErrInvalidVoteQuorum
	}
	if m.VotePeriod == 0 {
		return ErrInvalidVotePeriod
	}

	return nil
}

// GetSignBytes returns to sign bytes
func (m *MsgCreateCollective) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

// GetSigners returns signers
func (m *MsgCreateCollective) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

var _ sdk.Msg = &MsgBondCollective{}

// NewMsgBondCollective returns an instance of MsgBondCollective
func NewMsgBondCollective(
	proposer sdk.AccAddress,
	name string,
	bonds sdk.Coins,
) *MsgBondCollective {
	return &MsgBondCollective{
		Sender: proposer.String(),
		Name:   name,
		Bonds:  bonds,
	}
}

// Route returns route
func (m *MsgBondCollective) Route() string {
	return ModuleName
}

// Type returns return message type
func (m *MsgBondCollective) Type() string {
	return types.MsgTypeBondCollective
}

// ValidateBasic returns basic validation result
func (m *MsgBondCollective) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		return err
	}

	if m.Name == "" {
		return ErrEmptyCollectiveName
	}
	if sdk.Coins(m.Bonds).Empty() {
		return ErrEmptyCollectiveBonds
	}

	return nil
}

// GetSignBytes returns to sign bytes
func (m *MsgBondCollective) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

// GetSigners returns signers
func (m *MsgBondCollective) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

var _ sdk.Msg = &MsgDonateCollective{}

// NewMsgDonateCollective returns an instance of MsgDonateCollective
func NewMsgDonateCollective(
	proposer sdk.AccAddress,
	name string,
	locking uint64,
	donation sdk.Dec,
	donationLock bool,
) *MsgDonateCollective {
	return &MsgDonateCollective{
		Sender:       proposer.String(),
		Name:         name,
		Locking:      locking,
		Donation:     donation,
		DonationLock: donationLock,
	}
}

// Route returns route
func (m *MsgDonateCollective) Route() string {
	return ModuleName
}

// Type returns return message type
func (m *MsgDonateCollective) Type() string {
	return types.MsgTypeDonateCollective
}

// ValidateBasic returns basic validation result
func (m *MsgDonateCollective) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		return err
	}

	if m.Name == "" {
		return ErrEmptyCollectiveName
	}
	// In addition to the locking period, whitelisted contributors should be able to configure their individual intended donation to the collective,
	// that is a percentage of rewards (value between 0 and 1) that should be deposited and controlled by the collective.
	if m.Donation.LT(sdk.ZeroDec()) || m.Donation.GT(sdk.OneDec()) {
		return ErrInvalidDonationValue
	}
	return nil
}

// GetSignBytes returns to sign bytes
func (m *MsgDonateCollective) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

// GetSigners returns signers
func (m *MsgDonateCollective) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

var _ sdk.Msg = &MsgWithdrawCollective{}

// NewMsgWithdrawCollective returns an instance of MsgWithdrawCollective
func NewMsgWithdrawCollective(proposer sdk.AccAddress, name string) *MsgWithdrawCollective {
	return &MsgWithdrawCollective{
		Sender: proposer.String(),
		Name:   name,
	}
}

// Route returns route
func (m *MsgWithdrawCollective) Route() string {
	return ModuleName
}

// Type returns return message type
func (m *MsgWithdrawCollective) Type() string {
	return types.MsgTypeWithdrawCollective
}

// ValidateBasic returns basic validation result
func (m *MsgWithdrawCollective) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		return err
	}

	if m.Name == "" {
		return ErrEmptyCollectiveName
	}

	return nil
}

// GetSignBytes returns to sign bytes
func (m *MsgWithdrawCollective) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

// GetSigners returns signers
func (m *MsgWithdrawCollective) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}
