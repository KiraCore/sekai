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
	// id string, git string, checkout string, checksum string,
	name string,
	resources []Resource,
	height, minUpgradeTime int64,
	oldChainId, newChainId, rollBackMemo string,
	maxEnrollmentDuration int64, memo string,
	instateUpgrade bool) *MsgProposalSoftwareUpgradeRequest {
	return &MsgProposalSoftwareUpgradeRequest{
		Name:                 name,
		Resources:            resources,
		Height:               height,
		MinUpgradeTime:       minUpgradeTime,
		OldChainId:           oldChainId,
		NewChainId:           newChainId,
		RollbackChecksum:     rollBackMemo,
		MaxEnrolmentDuration: maxEnrollmentDuration,
		Memo:                 memo,
		InstateUpgrade:       instateUpgrade,
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
