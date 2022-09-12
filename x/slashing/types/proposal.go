package types

import (
	"time"

	kiratypes "github.com/KiraCore/sekai/types"
	"github.com/KiraCore/sekai/x/gov/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewResetWholeValidatorRankProposal(proposer sdk.AccAddress) *ProposalResetWholeValidatorRank {
	return &ProposalResetWholeValidatorRank{
		Proposer: proposer,
	}
}

func (m *ProposalResetWholeValidatorRank) ProposalType() string {
	return kiratypes.ProposalTypeResetWholeValidatorRank
}

func (m *ProposalResetWholeValidatorRank) ProposalPermission() types.PermValue {
	return types.PermCreateResetWholeValidatorRankProposal
}

func (m *ProposalResetWholeValidatorRank) VotePermission() types.PermValue {
	return types.PermVoteResetWholeValidatorRankProposal
}

// ValidateBasic returns basic validation
func (m *ProposalResetWholeValidatorRank) ValidateBasic() error {
	return nil
}

func NewSlashValidatorProposal(
	offender string,
	stakingPoolId uint64,
	misbehaviourTime time.Time,
	misbehaviourType string,
	jailPercentage uint64,
	colluders []string,
	refutation string,
) *ProposalSlashValidator {
	return &ProposalSlashValidator{
		Offender:         offender,
		StakingPoolId:    stakingPoolId,
		MisbehaviourTime: misbehaviourTime,
		MisbehaviourType: misbehaviourType,
		JailPercentage:   jailPercentage,
		Colluders:        colluders,
		Refutation:       refutation,
	}
}

func (m *ProposalSlashValidator) ProposalType() string {
	return kiratypes.ProposalTypeSlashValidator
}

func (m *ProposalSlashValidator) ProposalPermission() types.PermValue {
	return types.PermCreateSlashValidatorProposal
}

func (m *ProposalSlashValidator) VotePermission() types.PermValue {
	return types.PermVoteSlashValidatorProposal
}

// ValidateBasic returns basic validation
func (m *ProposalSlashValidator) ValidateBasic() error {
	return nil
}
