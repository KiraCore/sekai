package types

import (
	"github.com/KiraCore/sekai/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.Msg = &MsgSetExecutionFee{}

func NewMsgSetExecutionFee(
	name string,
	transactionType string,
	executionFee uint64,
	failureFee uint64,
	timeout uint64,
	defaultParams uint64,
	proposer sdk.AccAddress,
) *MsgSetExecutionFee {
	return &MsgSetExecutionFee{
		Name:              name,
		TransactionType:   transactionType,
		ExecutionFee:      executionFee,
		FailureFee:        failureFee,
		Timeout:           timeout,
		DefaultParameters: defaultParams,
		Proposer:          proposer,
	}
}

func (m *MsgSetExecutionFee) Route() string {
	return ModuleName
}

func (m *MsgSetExecutionFee) Type() string {
	return types.MsgTypeSetExecutionFee
}

func (m *MsgSetExecutionFee) ValidateBasic() error {
	if m.Proposer.Empty() {
		return ErrEmptyProposerAccAddress
	}

	return nil
}

func (m *MsgSetExecutionFee) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m *MsgSetExecutionFee) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{
		m.Proposer,
	}
}
