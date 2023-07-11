package types

import (
	"github.com/KiraCore/sekai/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.Msg = &MsgDisableBasketDeposits{}

// NewMsgDisableBasketDeposits returns an instance of MsgDisableBasketDeposits
func NewMsgDisableBasketDeposits(proposer sdk.AccAddress, basketId uint64, disabled bool) *MsgDisableBasketDeposits {
	return &MsgDisableBasketDeposits{
		Sender:   proposer.String(),
		BasketId: basketId,
		Disabled: disabled,
	}
}

// Route returns route
func (m *MsgDisableBasketDeposits) Route() string {
	return ModuleName
}

// Type returns return message type
func (m *MsgDisableBasketDeposits) Type() string {
	return types.MsgTypeDisableBasketWithdraws
}

// ValidateBasic returns basic validation result
func (m *MsgDisableBasketDeposits) ValidateBasic() error {
	return nil
}

// GetSignBytes returns to sign bytes
func (m *MsgDisableBasketDeposits) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

// GetSigners returns signers
func (m *MsgDisableBasketDeposits) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

var _ sdk.Msg = &MsgDisableBasketWithdraws{}

// NewMsgDisableBasketWithdraws returns an instance of MsgDisableBasketWithdraws
func NewMsgDisableBasketWithdraws(proposer sdk.AccAddress, basketId uint64, disabled bool) *MsgDisableBasketWithdraws {
	return &MsgDisableBasketWithdraws{
		Sender:   proposer.String(),
		BasketId: basketId,
		Disabled: disabled,
	}
}

// Route returns route
func (m *MsgDisableBasketWithdraws) Route() string {
	return ModuleName
}

// Type returns return message type
func (m *MsgDisableBasketWithdraws) Type() string {
	return types.MsgTypeDisableBasketWithdraws
}

// ValidateBasic returns basic validation result
func (m *MsgDisableBasketWithdraws) ValidateBasic() error {
	return nil
}

// GetSignBytes returns to sign bytes
func (m *MsgDisableBasketWithdraws) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

// GetSigners returns signers
func (m *MsgDisableBasketWithdraws) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

var _ sdk.Msg = &MsgDisableBasketSwaps{}

// NewMsgDisableBasketSwaps returns an instance of MsgDisableBasketSwaps
func NewMsgDisableBasketSwaps(proposer sdk.AccAddress, basketId uint64, disabled bool) *MsgDisableBasketSwaps {
	return &MsgDisableBasketSwaps{
		Sender:   proposer.String(),
		BasketId: basketId,
		Disabled: disabled,
	}
}

// Route returns route
func (m *MsgDisableBasketSwaps) Route() string {
	return ModuleName
}

// Type returns return message type
func (m *MsgDisableBasketSwaps) Type() string {
	return types.MsgTypeDisableBasketSwaps
}

// ValidateBasic returns basic validation result
func (m *MsgDisableBasketSwaps) ValidateBasic() error {
	return nil
}

// GetSignBytes returns to sign bytes
func (m *MsgDisableBasketSwaps) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

// GetSigners returns signers
func (m *MsgDisableBasketSwaps) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

var _ sdk.Msg = &MsgBasketTokenMint{}

// NewMsgBasketTokenMint returns an instance of MsgBasketTokenMint
func NewMsgBasketTokenMint(proposer sdk.AccAddress, basketId uint64, deposit sdk.Coins) *MsgBasketTokenMint {
	return &MsgBasketTokenMint{
		Sender:   proposer.String(),
		BasketId: basketId,
		Deposit:  deposit,
	}
}

// Route returns route
func (m *MsgBasketTokenMint) Route() string {
	return ModuleName
}

// Type returns return message type
func (m *MsgBasketTokenMint) Type() string {
	return types.MsgTypeBasketTokenMint
}

// ValidateBasic returns basic validation result
func (m *MsgBasketTokenMint) ValidateBasic() error {
	return nil
}

// GetSignBytes returns to sign bytes
func (m *MsgBasketTokenMint) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

// GetSigners returns signers
func (m *MsgBasketTokenMint) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

var _ sdk.Msg = &MsgBasketTokenBurn{}

// NewMsgBasketTokenBurn returns an instance of MsgBasketTokenBurn
func NewMsgBasketTokenBurn(proposer sdk.AccAddress, basketId uint64, burnAmount sdk.Coin) *MsgBasketTokenBurn {
	return &MsgBasketTokenBurn{
		Sender:     proposer.String(),
		BasketId:   basketId,
		BurnAmount: burnAmount,
	}
}

// Route returns route
func (m *MsgBasketTokenBurn) Route() string {
	return ModuleName
}

// Type returns return message type
func (m *MsgBasketTokenBurn) Type() string {
	return types.MsgTypeBasketTokenBurn
}

// ValidateBasic returns basic validation result
func (m *MsgBasketTokenBurn) ValidateBasic() error {
	return nil
}

// GetSignBytes returns to sign bytes
func (m *MsgBasketTokenBurn) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

// GetSigners returns signers
func (m *MsgBasketTokenBurn) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

var _ sdk.Msg = &MsgBasketTokenSwap{}

// NewMsgBasketTokenSwap returns an instance of MsgBasketTokenSwap
func NewMsgBasketTokenSwap(proposer sdk.AccAddress, basketId uint64, pairs []SwapPair) *MsgBasketTokenSwap {
	return &MsgBasketTokenSwap{
		Sender:   proposer.String(),
		BasketId: basketId,
		Pairs:    pairs,
	}
}

// Route returns route
func (m *MsgBasketTokenSwap) Route() string {
	return ModuleName
}

// Type returns return message type
func (m *MsgBasketTokenSwap) Type() string {
	return types.MsgTypeBasketTokenSwap
}

// ValidateBasic returns basic validation result
func (m *MsgBasketTokenSwap) ValidateBasic() error {
	return nil
}

// GetSignBytes returns to sign bytes
func (m *MsgBasketTokenSwap) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

// GetSigners returns signers
func (m *MsgBasketTokenSwap) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

var _ sdk.Msg = &MsgBasketClaimRewards{}

// NewMsgBasketClaimRewards returns an instance of MsgBasketClaimRewards
func NewMsgBasketClaimRewards(proposer sdk.AccAddress, basketTokens sdk.Coins) *MsgBasketClaimRewards {
	return &MsgBasketClaimRewards{
		Sender:       proposer.String(),
		BasketTokens: basketTokens,
	}
}

// Route returns route
func (m *MsgBasketClaimRewards) Route() string {
	return ModuleName
}

// Type returns return message type
func (m *MsgBasketClaimRewards) Type() string {
	return types.MsgTypeBasketClaimRewards
}

// ValidateBasic returns basic validation result
func (m *MsgBasketClaimRewards) ValidateBasic() error {
	return nil
}

// GetSignBytes returns to sign bytes
func (m *MsgBasketClaimRewards) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

// GetSigners returns signers
func (m *MsgBasketClaimRewards) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}
