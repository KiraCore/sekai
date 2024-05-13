package types

import (
	"fmt"
	"time"

	errorsmod "cosmossdk.io/errors"
	kiratypes "github.com/KiraCore/sekai/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/gogo/protobuf/proto"
)

// NewProposal creates a new proposal
func NewProposal(
	proposalID uint64,
	title, description string,
	content Content,
	submitTime time.Time,
	votingEndTime time.Time,
	enactmentEndTime time.Time,
	minVotingEndBlockHeight int64,
	minEnactmentEndBlockHeight int64,
) (Proposal, error) {
	msg, ok := content.(proto.Message)
	if !ok {
		return Proposal{}, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("%T does not implement proto.Message", content))
	}

	any, err := codectypes.NewAnyWithValue(msg)
	if err != nil {
		return Proposal{}, err
	}

	return Proposal{
		ProposalId:                 proposalID,
		Title:                      title,
		Description:                description,
		SubmitTime:                 submitTime,
		VotingEndTime:              votingEndTime,
		EnactmentEndTime:           enactmentEndTime,
		MinVotingEndBlockHeight:    minVotingEndBlockHeight,
		MinEnactmentEndBlockHeight: minEnactmentEndBlockHeight,
		Content:                    any,
		Result:                     Pending,
	}, nil
}

// GetContent returns the proposal Content
func (p Proposal) GetContent() Content {
	content, ok := p.Content.GetCachedValue().(Content)
	if !ok {
		return nil
	}
	return content
}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (p Proposal) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	var content Content
	return unpacker.UnpackAny(p.Content, &content)
}

var _ Content = &WhitelistAccountPermissionProposal{}

// NewWhitelistAccountPermissionProposal creates a new assign permission proposal
func NewWhitelistAccountPermissionProposal(
	address types.AccAddress,
	permission PermValue,
) Content {
	return &WhitelistAccountPermissionProposal{
		Address:    address,
		Permission: permission,
	}
}

// ProposalType returns proposal's type
func (m *WhitelistAccountPermissionProposal) ProposalType() string {
	return kiratypes.ProposalTypeWhitelistAccountPermission
}

// ValidateBasic returns basic validation
func (m *WhitelistAccountPermissionProposal) ValidateBasic() error {
	if m.Address.Empty() {
		return ErrEmptyPermissionsAccAddress
	}
	return nil
}

func (m *WhitelistAccountPermissionProposal) ProposalPermission() PermValue {
	return PermWhitelistAccountPermissionProposal
}

func (m *WhitelistAccountPermissionProposal) VotePermission() PermValue {
	return PermVoteWhitelistAccountPermissionProposal
}

var _ Content = &BlacklistAccountPermissionProposal{}

// NewBlacklistAccountPermissionProposal creates a new assign permission proposal
func NewBlacklistAccountPermissionProposal(
	address types.AccAddress,
	permission PermValue,
) Content {
	return &BlacklistAccountPermissionProposal{
		Address:    address,
		Permission: permission,
	}
}

// ProposalType returns proposal's type
func (m *BlacklistAccountPermissionProposal) ProposalType() string {
	return kiratypes.ProposalTypeBlacklistAccountPermission
}

// ValidateBasic returns basic validation
func (m *BlacklistAccountPermissionProposal) ValidateBasic() error {
	if m.Address.Empty() {
		return ErrEmptyPermissionsAccAddress
	}
	return nil
}

func (m *BlacklistAccountPermissionProposal) ProposalPermission() PermValue {
	return PermBlacklistAccountPermissionProposal
}

func (m *BlacklistAccountPermissionProposal) VotePermission() PermValue {
	return PermVoteBlacklistAccountPermissionProposal
}

var _ Content = &RemoveWhitelistedAccountPermissionProposal{}

// NewRemoveWhitelistedAccountPermissionProposal creates a new assign permission proposal
func NewRemoveWhitelistedAccountPermissionProposal(
	address types.AccAddress,
	permission PermValue,
) Content {
	return &RemoveWhitelistedAccountPermissionProposal{
		Address:    address,
		Permission: permission,
	}
}

// ProposalType returns proposal's type
func (m *RemoveWhitelistedAccountPermissionProposal) ProposalType() string {
	return kiratypes.ProposalTypeRemoveWhitelistedAccountPermission
}

// ValidateBasic returns basic validation
func (m *RemoveWhitelistedAccountPermissionProposal) ValidateBasic() error {
	if m.Address.Empty() {
		return ErrEmptyPermissionsAccAddress
	}
	return nil
}

func (m *RemoveWhitelistedAccountPermissionProposal) ProposalPermission() PermValue {
	return PermRemoveWhitelistedAccountPermissionProposal
}

func (m *RemoveWhitelistedAccountPermissionProposal) VotePermission() PermValue {
	return PermVoteRemoveWhitelistedAccountPermissionProposal
}

var _ Content = &RemoveBlacklistedAccountPermissionProposal{}

// NewRemoveBlacklistedAccountPermissionProposal creates a new assign permission proposal
func NewRemoveBlacklistedAccountPermissionProposal(
	address types.AccAddress,
	permission PermValue,
) Content {
	return &RemoveBlacklistedAccountPermissionProposal{
		Address:    address,
		Permission: permission,
	}
}

// ProposalType returns proposal's type
func (m *RemoveBlacklistedAccountPermissionProposal) ProposalType() string {
	return kiratypes.ProposalTypeRemoveBlacklistedAccountPermission
}

// ValidateBasic returns basic validation
func (m *RemoveBlacklistedAccountPermissionProposal) ValidateBasic() error {
	if m.Address.Empty() {
		return ErrEmptyPermissionsAccAddress
	}
	return nil
}

func (m *RemoveBlacklistedAccountPermissionProposal) ProposalPermission() PermValue {
	return PermRemoveBlacklistedAccountPermissionProposal
}

func (m *RemoveBlacklistedAccountPermissionProposal) VotePermission() PermValue {
	return PermVoteRemoveBlacklistedAccountPermissionProposal
}

var _ Content = &RemoveBlacklistedAccountPermissionProposal{}

// NewAssignRoleToAccountProposal creates a new assign permission proposal
func NewAssignRoleToAccountProposal(
	address types.AccAddress,
	roleIdentifier string,
) Content {
	return &AssignRoleToAccountProposal{
		Address:        address,
		RoleIdentifier: roleIdentifier,
	}
}

// ProposalType returns proposal's type
func (m *AssignRoleToAccountProposal) ProposalType() string {
	return kiratypes.ProposalTypeAssignRoleToAccount
}

// ValidateBasic returns basic validation
func (m *AssignRoleToAccountProposal) ValidateBasic() error {
	if m.Address.Empty() {
		return ErrEmptyPermissionsAccAddress
	}
	return nil
}

func (m *AssignRoleToAccountProposal) ProposalPermission() PermValue {
	return PermAssignRoleToAccountProposal
}

func (m *AssignRoleToAccountProposal) VotePermission() PermValue {
	return PermVoteAssignRoleToAccountProposal
}

// NewUnassignRoleFromAccountProposal creates a new assign permission proposal
func NewUnassignRoleFromAccountProposal(
	address types.AccAddress,
	roleIdentifier string,
) Content {
	return &UnassignRoleFromAccountProposal{
		Address:        address,
		RoleIdentifier: roleIdentifier,
	}
}

// ProposalType returns proposal's type
func (m *UnassignRoleFromAccountProposal) ProposalType() string {
	return kiratypes.ProposalTypeUnassignRoleFromAccount
}

// ValidateBasic returns basic validation
func (m *UnassignRoleFromAccountProposal) ValidateBasic() error {
	if m.Address.Empty() {
		return ErrEmptyPermissionsAccAddress
	}
	return nil
}

func (m *UnassignRoleFromAccountProposal) ProposalPermission() PermValue {
	return PermUnassignRoleFromAccountProposal
}

func (m *UnassignRoleFromAccountProposal) VotePermission() PermValue {
	return PermVoteUnassignRoleFromAccountProposal
}

// NewSetNetworkPropertyProposal creates a new set network property proposal
func NewSetNetworkPropertyProposal(
	property NetworkProperty,
	value NetworkPropertyValue,
) Content {
	return &SetNetworkPropertyProposal{
		NetworkProperty: property,
		Value:           value,
	}
}

// ProposalType returns proposal's type
func (m *SetNetworkPropertyProposal) ProposalType() string {
	return kiratypes.ProposalTypeSetNetworkProperty
}

func (m *SetNetworkPropertyProposal) ProposalPermission() PermValue {
	return PermCreateSetNetworkPropertyProposal
}

// VotePermission returns permission to vote on this proposal
func (m *SetNetworkPropertyProposal) VotePermission() PermValue {
	return PermVoteSetNetworkPropertyProposal
}

// ValidateBasic returns basic validation
func (m *SetNetworkPropertyProposal) ValidateBasic() error {
	switch m.NetworkProperty {
	case MinTxFee,
		MaxTxFee,
		VoteQuorum,
		MinimumProposalEndTime,
		ProposalEnactmentTime,
		EnableForeignFeePayments,
		MischanceRankDecreaseAmount,
		MischanceConfidence,
		MaxMischance,
		InactiveRankDecreasePercent,
		PoorNetworkMaxBankSend,
		MinValidators,
		UnjailMaxTime,
		EnableTokenWhitelist,
		EnableTokenBlacklist,
		MinIdentityApprovalTip,
		UniqueIdentityKeys,
		UbiHardcap,
		ValidatorsFeeShare,
		InflationRate,
		InflationPeriod,
		UnstakingPeriod,
		MaxDelegators,
		MinDelegationPushout,
		SlashingPeriod,
		MaxJailedPercentage,
		MaxSlashingPercentage,
		MinCustodyReward,
		MaxCustodyBufferSize,
		MaxCustodyTxSize,
		AbstentionRankDecreaseAmount,
		MaxAbstention,
		MinCollectiveBond,
		MinCollectiveBondingTime,
		MaxCollectiveOutputs,
		MinCollectiveClaimPeriod,
		ValidatorRecoveryBond,
		MaxAnnualInflation,
		MinDappBond,
		MaxDappBond,
		DappBondDuration:
		return nil
	default:
		return ErrInvalidNetworkProperty
	}
}

func NewUpsertDataRegistryProposal(key, hash, reference, encoding string, size uint64) Content {
	return &UpsertDataRegistryProposal{
		Key:       key,
		Hash:      hash,
		Reference: reference,
		Encoding:  encoding,
		Size_:     size,
	}
}

func (m *UpsertDataRegistryProposal) ProposalType() string {
	return kiratypes.ProposalTypeUpsertDataRegistry
}

func (m *UpsertDataRegistryProposal) ProposalPermission() PermValue {
	return PermCreateUpsertDataRegistryProposal
}

func (m *UpsertDataRegistryProposal) VotePermission() PermValue {
	return PermVoteUpsertDataRegistryProposal
}

// ValidateBasic returns basic validation
func (m *UpsertDataRegistryProposal) ValidateBasic() error {
	return nil
}

func NewSetPoorNetworkMessagesProposal(msgs []string) Content {
	return &SetPoorNetworkMessagesProposal{
		Messages: msgs,
	}
}

func (m *SetPoorNetworkMessagesProposal) ProposalType() string {
	return kiratypes.ProposalTypeSetPoorNetworkMessages
}

func (m *SetPoorNetworkMessagesProposal) ProposalPermission() PermValue {
	return PermCreateSetPoorNetworkMessagesProposal
}

func (m *SetPoorNetworkMessagesProposal) VotePermission() PermValue {
	return PermVoteSetPoorNetworkMessagesProposal
}

// ValidateBasic returns basic validation
func (m *SetPoorNetworkMessagesProposal) ValidateBasic() error {
	return nil
}

func NewCreateRoleProposal(sid, description string, whitelist []PermValue, blacklist []PermValue) Content {
	return &CreateRoleProposal{
		RoleSid:                sid,
		RoleDescription:        description,
		WhitelistedPermissions: whitelist,
		BlacklistedPermissions: blacklist,
	}
}

func (m *CreateRoleProposal) ProposalType() string {
	return kiratypes.ProposalTypeCreateRole
}

func (m *CreateRoleProposal) ProposalPermission() PermValue {
	return PermCreateRoleProposal
}

func (m *CreateRoleProposal) VotePermission() PermValue {
	return PermVoteCreateRoleProposal
}

// ValidateBasic returns basic validation
func (m *CreateRoleProposal) ValidateBasic() error {
	if m.RoleSid == "" {
		return ErrInvalidRoleIdentifier
	}
	if len(m.WhitelistedPermissions) == 0 && len(m.BlacklistedPermissions) == 0 {
		return ErrEmptyPermissions
	}

	return nil
}

func NewRemoveRoleProposal(roleIdentifier string) Content {
	return &RemoveRoleProposal{
		RoleSid: roleIdentifier,
	}
}

func (m *RemoveRoleProposal) ProposalType() string {
	return kiratypes.ProposalTypeRemoveRole
}

func (m *RemoveRoleProposal) ProposalPermission() PermValue {
	return PermRemoveRoleProposal
}

func (m *RemoveRoleProposal) VotePermission() PermValue {
	return PermVoteRemoveRoleProposal
}

// ValidateBasic returns basic validation
func (m *RemoveRoleProposal) ValidateBasic() error {
	if m.RoleSid == "" {
		return ErrInvalidRoleIdentifier
	}
	return nil
}

func NewWhitelistRolePermissionProposal(roleIdentifier string, permission PermValue) Content {
	return &WhitelistRolePermissionProposal{
		RoleIdentifier: roleIdentifier,
		Permission:     permission,
	}
}

func (m *WhitelistRolePermissionProposal) ProposalType() string {
	return kiratypes.ProposalTypeWhitelistRolePermission
}

func (m *WhitelistRolePermissionProposal) ProposalPermission() PermValue {
	return PermWhitelistRolePermissionProposal
}

func (m *WhitelistRolePermissionProposal) VotePermission() PermValue {
	return PermVoteWhitelistRolePermissionProposal
}

// ValidateBasic returns basic validation
func (m *WhitelistRolePermissionProposal) ValidateBasic() error {
	if m.RoleIdentifier == "" {
		return ErrInvalidRoleIdentifier
	}
	return nil
}

func NewBlacklistRolePermissionProposal(roleIdentifier string, permission PermValue) Content {
	return &BlacklistRolePermissionProposal{
		RoleIdentifier: roleIdentifier,
		Permission:     permission,
	}
}

func (m *BlacklistRolePermissionProposal) ProposalType() string {
	return kiratypes.ProposalTypeBlacklistRolePermission
}

func (m *BlacklistRolePermissionProposal) ProposalPermission() PermValue {
	return PermBlacklistRolePermissionProposal
}

func (m *BlacklistRolePermissionProposal) VotePermission() PermValue {
	return PermVoteBlacklistRolePermissionProposal
}

// ValidateBasic returns basic validation
func (m *BlacklistRolePermissionProposal) ValidateBasic() error {
	if m.RoleIdentifier == "" {
		return ErrInvalidRoleIdentifier
	}
	return nil
}

func NewRemoveWhitelistedRolePermissionProposal(roleSid string, permission PermValue) Content {
	return &RemoveWhitelistedRolePermissionProposal{
		RoleSid:    roleSid,
		Permission: permission,
	}
}

func (m *RemoveWhitelistedRolePermissionProposal) ProposalType() string {
	return kiratypes.ProposalTypeRemoveWhitelistedRolePermission
}

func (m *RemoveWhitelistedRolePermissionProposal) ProposalPermission() PermValue {
	return PermRemoveWhitelistedRolePermissionProposal
}

func (m *RemoveWhitelistedRolePermissionProposal) VotePermission() PermValue {
	return PermVoteRemoveWhitelistedRolePermissionProposal
}

// ValidateBasic returns basic validation
func (m *RemoveWhitelistedRolePermissionProposal) ValidateBasic() error {
	if m.RoleSid == "" {
		return ErrInvalidRoleIdentifier
	}
	return nil
}

func NewRemoveBlacklistedRolePermissionProposal(roleSid string, permission PermValue) Content {
	return &RemoveBlacklistedRolePermissionProposal{
		RoleSid:    roleSid,
		Permission: permission,
	}
}

func (m *RemoveBlacklistedRolePermissionProposal) ProposalType() string {
	return kiratypes.ProposalTypeRemoveBlacklistedRolePermission
}

func (m *RemoveBlacklistedRolePermissionProposal) ProposalPermission() PermValue {
	return PermRemoveBlacklistedRolePermissionProposal
}

func (m *RemoveBlacklistedRolePermissionProposal) VotePermission() PermValue {
	return PermVoteRemoveBlacklistedRolePermissionProposal
}

// ValidateBasic returns basic validation
func (m *RemoveBlacklistedRolePermissionProposal) ValidateBasic() error {
	if m.RoleSid == "" {
		return ErrInvalidRoleIdentifier
	}

	return nil
}

func NewSetProposalDurationsProposal(typeofProposals []string, durations []uint64) Content {
	return &SetProposalDurationsProposal{
		TypeofProposals:   typeofProposals,
		ProposalDurations: durations,
	}
}

func (m *SetProposalDurationsProposal) ProposalType() string {
	return kiratypes.ProposalTypeSetProposalDurations
}

func (m *SetProposalDurationsProposal) ProposalPermission() PermValue {
	return PermCreateSetProposalDurationProposal
}

func (m *SetProposalDurationsProposal) VotePermission() PermValue {
	return PermVoteSetProposalDurationProposal
}

// ValidateBasic returns basic validation
func (m *SetProposalDurationsProposal) ValidateBasic() error {
	if len(m.TypeofProposals) == 0 {
		return fmt.Errorf("at least one proposal type should be set")
	}
	if len(m.TypeofProposals) != len(m.ProposalDurations) {
		return fmt.Errorf("the length of proposal types and durations should be equal")
	}
	for _, pt := range m.TypeofProposals {
		if pt == "" {
			return fmt.Errorf("empty proposal type is not allowed")
		}
	}
	for _, pd := range m.ProposalDurations {
		if pd == 0 {
			return fmt.Errorf("zero proposal duration is not allowed")
		}
	}
	return nil
}

func NewResetWholeCouncilorRankProposal(proposer sdk.AccAddress) *ProposalResetWholeCouncilorRank {
	return &ProposalResetWholeCouncilorRank{
		Proposer: proposer,
	}
}

func (m *ProposalResetWholeCouncilorRank) ProposalType() string {
	return kiratypes.ProposalTypeResetWholeCouncilorRank
}

func (m *ProposalResetWholeCouncilorRank) ProposalPermission() PermValue {
	return PermCreateResetWholeCouncilorRankProposal
}

func (m *ProposalResetWholeCouncilorRank) VotePermission() PermValue {
	return PermVoteResetWholeCouncilorRankProposal
}

// ValidateBasic returns basic validation
func (m *ProposalResetWholeCouncilorRank) ValidateBasic() error {
	return nil
}

func NewJailCouncilorProposal(proposer sdk.AccAddress, description string, councilors []string) *ProposalJailCouncilor {
	return &ProposalJailCouncilor{
		Proposer:    proposer,
		Description: description,
		Councilors:  councilors,
	}
}

func (m *ProposalJailCouncilor) ProposalType() string {
	return kiratypes.ProposalTypeJailCouncilor
}

func (m *ProposalJailCouncilor) ProposalPermission() PermValue {
	return PermCreateJailCouncilorProposal
}

func (m *ProposalJailCouncilor) VotePermission() PermValue {
	return PermVoteJailCouncilorProposal
}

// ValidateBasic returns basic validation
func (m *ProposalJailCouncilor) ValidateBasic() error {
	return nil
}

func NewSetExecutionFeesProposal(proposer sdk.AccAddress, description string, executionFees []ExecutionFee) *ProposalSetExecutionFees {
	return &ProposalSetExecutionFees{
		Proposer:      proposer,
		Description:   description,
		ExecutionFees: executionFees,
	}
}

func (m *ProposalSetExecutionFees) ProposalType() string {
	return kiratypes.ProposalTypeSetExecutionFees
}

func (m *ProposalSetExecutionFees) ProposalPermission() PermValue {
	return PermCreateSetExecutionFeesProposal
}

func (m *ProposalSetExecutionFees) VotePermission() PermValue {
	return PermVoteSetExecutionFeesProposal
}

// ValidateBasic returns basic validation
func (m *ProposalSetExecutionFees) ValidateBasic() error {
	return nil
}
