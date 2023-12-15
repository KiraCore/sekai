package types

import (
	"fmt"

	"github.com/KiraCore/sekai/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	_ sdk.Msg = &MsgUpsertStakingPool{}
)

func NewMsgUpsertStakingPool(sender, validator string, enabled bool, commission sdk.Dec) *MsgUpsertStakingPool {
	return &MsgUpsertStakingPool{
		Sender:     sender,
		Validator:  validator,
		Enabled:    enabled,
		Commission: commission,
	}
}

func (m *MsgUpsertStakingPool) Route() string {
	return ModuleName
}

func (m *MsgUpsertStakingPool) Type() string {
	return types.MsgTypeUpsertStakingPool
}

func (m *MsgUpsertStakingPool) ValidateBasic() error {
	if m.Commission.LT(sdk.NewDecWithPrec(1, 2)) { // 1%
		return fmt.Errorf("commission should not be less than 1%%")
	}
	if m.Commission.GT(sdk.NewDecWithPrec(5, 1)) { // 50%
		return fmt.Errorf("commission should not be more than 50%%")
	}
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

func NewMsgDelegate(delegator string, validator string, coins sdk.Coins) *MsgDelegate {
	return &MsgDelegate{
		DelegatorAddress: delegator,
		ValidatorAddress: validator,
		Amounts:          coins,
	}
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
	sender, err := sdk.AccAddressFromBech32(m.DelegatorAddress)
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

func NewMsgUndelegate(delegator string, validator string, coins sdk.Coins) *MsgUndelegate {
	return &MsgUndelegate{
		DelegatorAddress: delegator,
		ValidatorAddress: validator,
		Amounts:          coins,
	}
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
	sender, err := sdk.AccAddressFromBech32(m.DelegatorAddress)
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

func NewMsgClaimRewards(sender string) *MsgClaimRewards {
	return &MsgClaimRewards{
		Sender: sender,
	}
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
	_ sdk.Msg = &MsgClaimUndelegation{}
)

func NewMsgClaimUndelegation(sender string, undelegationId uint64) *MsgClaimUndelegation {
	return &MsgClaimUndelegation{
		Sender:         sender,
		UndelegationId: undelegationId,
	}
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

var (
	_ sdk.Msg = &MsgSetCompoundInfo{}
)

func NewMsgSetCompoundInfo(sender string, allDenom bool, denoms []string) *MsgSetCompoundInfo {
	return &MsgSetCompoundInfo{
		Sender:         sender,
		AllDenom:       allDenom,
		CompoundDenoms: denoms,
	}
}

func (m *MsgSetCompoundInfo) Route() string {
	return ModuleName
}

func (m *MsgSetCompoundInfo) Type() string {
	return types.MsgTypeSetCompoundInfo
}

func (m *MsgSetCompoundInfo) ValidateBasic() error {
	return nil
}

func (m *MsgSetCompoundInfo) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m *MsgSetCompoundInfo) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{
		sender,
	}
}

var (
	_ sdk.Msg = &MsgSetCompoundInfo{}
)

func NewMsgRegisterDelegator(sender string) *MsgRegisterDelegator {
	return &MsgRegisterDelegator{
		Delegator: sender,
	}
}

func (m *MsgRegisterDelegator) Route() string {
	return ModuleName
}

func (m *MsgRegisterDelegator) Type() string {
	return types.MsgTypeRegisterDelegator
}

func (m *MsgRegisterDelegator) ValidateBasic() error {
	return nil
}

func (m *MsgRegisterDelegator) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m *MsgRegisterDelegator) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(m.Delegator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{
		sender,
	}
}

var (
	_ sdk.Msg = &MsgClaimMaturedUndelegations{}
)

func NewMsgClaimMaturedUndelegations(sender string) *MsgClaimMaturedUndelegations {
	return &MsgClaimMaturedUndelegations{
		Sender: sender,
	}
}

func (m *MsgClaimMaturedUndelegations) Route() string {
	return ModuleName
}

func (m *MsgClaimMaturedUndelegations) Type() string {
	return types.MsgTypeClaimMaturedUndelegations
}

func (m *MsgClaimMaturedUndelegations) ValidateBasic() error {
	return nil
}

func (m *MsgClaimMaturedUndelegations) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m *MsgClaimMaturedUndelegations) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}
