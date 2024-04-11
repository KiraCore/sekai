package types

import (
	"github.com/KiraCore/sekai/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewMsgCreateSpendingPool(
	name string,
	claimStart uint64,
	claimEnd uint64,
	rates sdk.DecCoins,
	voteQuorum sdk.Dec,
	votePeriod uint64,
	voteEnactment uint64,
	owners PermInfo,
	beneficiaries WeightedPermInfo,
	sender sdk.AccAddress,
	dynamicRate bool,
	dynamicRatePeriod uint64,
) *MsgCreateSpendingPool {
	return &MsgCreateSpendingPool{
		Name:              name,
		ClaimStart:        claimStart,
		ClaimEnd:          claimEnd,
		Rates:             rates,
		VoteQuorum:        voteQuorum,
		VotePeriod:        votePeriod,
		VoteEnactment:     voteEnactment,
		Owners:            owners,
		Beneficiaries:     beneficiaries,
		Sender:            sender.String(),
		DynamicRate:       dynamicRate,
		DynamicRatePeriod: dynamicRatePeriod,
	}
}

func (m *MsgCreateSpendingPool) Route() string {
	return ModuleName
}

func (m *MsgCreateSpendingPool) Type() string {
	return types.MsgTypeCreateSpendingPool
}

func (m *MsgCreateSpendingPool) ValidateBasic() error {
	if m.Sender == "" {
		return ErrEmptyProposerAccAddress
	}
	for _, beneficiary := range m.Beneficiaries.Accounts {
		if beneficiary.Weight.IsNil() || beneficiary.Weight.IsZero() {
			return ErrEmptyWeightBeneficiary
		}
	}
	for _, beneficiary := range m.Beneficiaries.Roles {
		if beneficiary.Weight.IsNil() || beneficiary.Weight.IsZero() {
			return ErrEmptyWeightBeneficiary
		}
	}
	return nil
}

func (m *MsgCreateSpendingPool) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m *MsgCreateSpendingPool) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{
		addr,
	}
}

func NewMsgDepositSpendingPool(
	name string,
	amount sdk.Coins,
	sender sdk.AccAddress,
) *MsgDepositSpendingPool {
	return &MsgDepositSpendingPool{
		PoolName: name,
		Amount:   amount,
		Sender:   sender.String(),
	}
}

func (m *MsgDepositSpendingPool) Route() string {
	return ModuleName
}

func (m *MsgDepositSpendingPool) Type() string {
	return types.MsgTypeDepositSpendingPool
}

func (m *MsgDepositSpendingPool) ValidateBasic() error {
	if m.Sender == "" {
		return ErrEmptyProposerAccAddress
	}
	return nil
}

func (m *MsgDepositSpendingPool) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m *MsgDepositSpendingPool) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{
		addr,
	}
}

func NewMsgRegisterSpendingPoolBeneficiary(
	name string,
	sender sdk.AccAddress,
) *MsgRegisterSpendingPoolBeneficiary {
	return &MsgRegisterSpendingPoolBeneficiary{
		PoolName: name,
		Sender:   sender.String(),
	}
}

func (m *MsgRegisterSpendingPoolBeneficiary) Route() string {
	return ModuleName
}

func (m *MsgRegisterSpendingPoolBeneficiary) Type() string {
	return types.MsgTypeRegisterSpendingPoolBeneficiary
}

func (m *MsgRegisterSpendingPoolBeneficiary) ValidateBasic() error {
	if m.Sender == "" {
		return ErrEmptyProposerAccAddress
	}
	return nil
}

func (m *MsgRegisterSpendingPoolBeneficiary) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m *MsgRegisterSpendingPoolBeneficiary) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{
		addr,
	}
}

func NewMsgClaimSpendingPool(
	name string,
	sender sdk.AccAddress,
) *MsgClaimSpendingPool {
	return &MsgClaimSpendingPool{
		PoolName: name,
		Sender:   sender.String(),
	}
}

func (m *MsgClaimSpendingPool) Route() string {
	return ModuleName
}

func (m *MsgClaimSpendingPool) Type() string {
	return types.MsgTypeClaimSpendingPool
}

func (m *MsgClaimSpendingPool) ValidateBasic() error {
	if m.Sender == "" {
		return ErrEmptyProposerAccAddress
	}
	return nil
}

func (m *MsgClaimSpendingPool) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m *MsgClaimSpendingPool) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{
		addr,
	}
}
