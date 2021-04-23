package types

import (
	kiratypes "github.com/KiraCore/sekai/types"
	"github.com/KiraCore/sekai/x/gov/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	_ sdk.Msg       = &MsgProposalSoftwareUpgradeRequest{}
	_ types.Content = &ProposalSoftwareUpgrade{}
)

func NewMsgProposalSoftwareUpgradeRequest(
	proposer sdk.AccAddress,
	id string, git string, checkout string, minHaltTime int64,
	oldChainId, newChainId, rollBackMemo, checkSum string, maxEnrollmentDuration int64, memo string) *MsgProposalSoftwareUpgradeRequest {
	return &MsgProposalSoftwareUpgradeRequest{
		Resources: &Resource{
			Id:       id,
			Git:      git,
			Checkout: checkout,
			Checksum: checkSum,
		},
		MinHaltTime:          minHaltTime,
		OldChainId:           oldChainId,
		NewChainId:           newChainId,
		RollbackChecksum:     rollBackMemo,
		MaxEnrolmentDuration: maxEnrollmentDuration,
		Memo:                 memo,
		Proposer:             proposer,
	}
}

func (m *MsgProposalSoftwareUpgradeRequest) Route() string {
	return ModuleName
}

func (m *MsgProposalSoftwareUpgradeRequest) Type() string {
	return kiratypes.MsgProposalSoftwareUpgrade
}

func (m *MsgProposalSoftwareUpgradeRequest) ValidateBasic() error {
	return nil
}

func (m *MsgProposalSoftwareUpgradeRequest) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m *MsgProposalSoftwareUpgradeRequest) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.Proposer}
}
