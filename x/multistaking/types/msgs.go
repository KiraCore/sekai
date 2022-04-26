package types

import (
	"github.com/KiraCore/sekai/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	_ sdk.Msg = &MsgUpsertStakingPool{}
)

func NewMsgUpsertStakingPool(
	moniker string,
	valKey sdk.ValAddress,
	pubKey cryptotypes.PubKey,
) (*MsgUpsertStakingPool, error) {
	return &MsgUpsertStakingPool{}, nil
}

func (m *MsgUpsertStakingPool) Route() string {
	return ModuleName
}

func (m *MsgUpsertStakingPool) Type() string {
	return types.MsgTypeUpsertStakingPool
}

func (m *MsgUpsertStakingPool) ValidateBasic() error {
	return nil
}

func (m *MsgUpsertStakingPool) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m *MsgUpsertStakingPool) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{
		sender,
	}
}

var (
	_ sdk.Msg = &MsgDelegate{}
)

func NewMsgDelegate(
	moniker string,
	valKey sdk.ValAddress,
	pubKey cryptotypes.PubKey,
) (*MsgDelegate, error) {
	return &MsgDelegate{}, nil
}

func (m *MsgDelegate) Route() string {
	return ModuleName
}

func (m *MsgDelegate) Type() string {
	return types.MsgTypeDelegate
}

func (m *MsgDelegate) ValidateBasic() error {
	return nil
}

func (m *MsgDelegate) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m *MsgDelegate) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{
		sender,
	}
}

var (
	_ sdk.Msg = &MsgUndelegate{}
)

func NewMsgUndelegate(
	moniker string,
	valKey sdk.ValAddress,
	pubKey cryptotypes.PubKey,
) (*MsgUndelegate, error) {
	return &MsgUndelegate{}, nil
}

func (m *MsgUndelegate) Route() string {
	return ModuleName
}

func (m *MsgUndelegate) Type() string {
	return types.MsgTypeUndelegate
}

func (m *MsgUndelegate) ValidateBasic() error {
	return nil
}

func (m *MsgUndelegate) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m *MsgUndelegate) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{
		sender,
	}
}

var (
	_ sdk.Msg = &MsgClaimRewards{}
)

func NewMsgClaimRewards(
	moniker string,
	valKey sdk.ValAddress,
	pubKey cryptotypes.PubKey,
) (*MsgClaimRewards, error) {
	return &MsgClaimRewards{}, nil
}

func (m *MsgClaimRewards) Route() string {
	return ModuleName
}

func (m *MsgClaimRewards) Type() string {
	return types.MsgTypeClaimRewards
}

func (m *MsgClaimRewards) ValidateBasic() error {
	return nil
}

func (m *MsgClaimRewards) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m *MsgClaimRewards) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{
		sender,
	}
}

var (
	_ sdk.Msg = &MsgClaimRewards{}
)

func NewMsgClaimUndelegation(
	moniker string,
	valKey sdk.ValAddress,
	pubKey cryptotypes.PubKey,
) (*MsgClaimUndelegation, error) {
	return &MsgClaimUndelegation{}, nil
}

func (m *MsgClaimUndelegation) Route() string {
	return ModuleName
}

func (m *MsgClaimUndelegation) Type() string {
	return types.MsgTypeClaimUndelegation
}

func (m *MsgClaimUndelegation) ValidateBasic() error {
	return nil
}

func (m *MsgClaimUndelegation) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m *MsgClaimUndelegation) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{
		sender,
	}
}
