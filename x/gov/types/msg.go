package types

import (
	"github.com/KiraCore/sekai/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	// Proposal
	_ sdk.Msg = &MsgVoteProposal{}

	// Permissions
	_ sdk.Msg = &MsgWhitelistPermissions{}
	_ sdk.Msg = &MsgBlacklistPermissions{}
	_ sdk.Msg = &MsgProposalAssignPermission{}
	_ sdk.Msg = &MsgProposalUpsertDataRegistry{}
	_ sdk.Msg = &MsgProposalSetPoorNetworkMessages{}

	// Councilor
	_ sdk.Msg = &MsgClaimCouncilor{}

	// Roles
	_ sdk.Msg = &MsgCreateRole{}
	_ sdk.Msg = &MsgAssignRole{}
	_ sdk.Msg = &MsgRemoveRole{}

	_ sdk.Msg = &MsgWhitelistRolePermission{}
	_ sdk.Msg = &MsgBlacklistRolePermission{}
	_ sdk.Msg = &MsgRemoveWhitelistRolePermission{}
	_ sdk.Msg = &MsgRemoveBlacklistRolePermission{}

	_ sdk.Msg = &MsgProposalCreateRole{}
)

func NewMsgWhitelistPermissions(
	proposer, address sdk.AccAddress,
	permission uint32,
) *MsgWhitelistPermissions {
	return &MsgWhitelistPermissions{
		Proposer:   proposer,
		Address:    address,
		Permission: permission,
	}
}

func (m *MsgWhitelistPermissions) Route() string {
	return ModuleName
}

func (m *MsgWhitelistPermissions) Type() string {
	return types.MsgTypeWhitelistPermissions
}

func (m *MsgWhitelistPermissions) ValidateBasic() error {
	if m.Proposer.Empty() {
		return ErrEmptyProposerAccAddress
	}

	if m.Address.Empty() {
		return ErrEmptyPermissionsAccAddress
	}

	return nil
}

func (m *MsgWhitelistPermissions) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m *MsgWhitelistPermissions) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{
		m.Proposer,
	}
}

func NewMsgBlacklistPermissions(
	proposer, address sdk.AccAddress,
	permission uint32,
) *MsgBlacklistPermissions {
	return &MsgBlacklistPermissions{
		Proposer:   proposer,
		Address:    address,
		Permission: permission,
	}
}

func (m *MsgBlacklistPermissions) Route() string {
	return ModuleName
}

func (m *MsgBlacklistPermissions) Type() string {
	return types.MsgTypeBlacklistPermissions
}

func (m *MsgBlacklistPermissions) ValidateBasic() error {
	if m.Proposer.Empty() {
		return ErrEmptyProposerAccAddress
	}

	if m.Address.Empty() {
		return ErrEmptyPermissionsAccAddress
	}

	return nil
}

func (m *MsgBlacklistPermissions) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m *MsgBlacklistPermissions) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{
		m.Proposer,
	}
}

func NewMsgClaimCouncilor(
	moniker string,
	website string,
	social string,
	identity string,
	address sdk.AccAddress,
) *MsgClaimCouncilor {
	return &MsgClaimCouncilor{
		Moniker:  moniker,
		Website:  website,
		Social:   social,
		Identity: identity,
		Address:  address,
	}
}

func (m *MsgClaimCouncilor) Route() string {
	return ModuleName
}

func (m *MsgClaimCouncilor) Type() string {
	return types.MsgTypeClaimCouncilor
}

func (m *MsgClaimCouncilor) ValidateBasic() error {
	if m.Address.Empty() {
		return ErrCouncilorEmptyAddress
	}

	return nil
}

func (m *MsgClaimCouncilor) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m *MsgClaimCouncilor) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{
		m.Address,
	}
}

func NewMsgWhitelistRolePermission(
	proposer sdk.AccAddress,
	role uint32,
	permission uint32,
) *MsgWhitelistRolePermission {
	return &MsgWhitelistRolePermission{Proposer: proposer, Role: role, Permission: permission}
}

func (m *MsgWhitelistRolePermission) Route() string {
	return ModuleName
}

func (m *MsgWhitelistRolePermission) Type() string {
	return types.MsgTypeWhitelistRolePermission
}

func (m *MsgWhitelistRolePermission) ValidateBasic() error {
	if m.Proposer.Empty() {
		return ErrEmptyProposerAccAddress
	}

	return nil
}

func (m *MsgWhitelistRolePermission) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m *MsgWhitelistRolePermission) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{
		m.Proposer,
	}
}

func NewMsgBlacklistRolePermission(
	proposer sdk.AccAddress,
	role uint32,
	permission uint32,
) *MsgBlacklistRolePermission {
	return &MsgBlacklistRolePermission{Proposer: proposer, Role: role, Permission: permission}
}

func (m *MsgBlacklistRolePermission) Route() string {
	return ModuleName
}

func (m *MsgBlacklistRolePermission) Type() string {
	return types.MsgTypeBlacklistRolePermission
}

func (m *MsgBlacklistRolePermission) ValidateBasic() error {
	if m.Proposer.Empty() {
		return ErrEmptyProposerAccAddress
	}

	return nil
}

func (m *MsgBlacklistRolePermission) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m *MsgBlacklistRolePermission) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{
		m.Proposer,
	}
}

func NewMsgRemoveWhitelistRolePermission(
	proposer sdk.AccAddress,
	role uint32,
	permission uint32,
) *MsgRemoveWhitelistRolePermission {
	return &MsgRemoveWhitelistRolePermission{Proposer: proposer, Role: role, Permission: permission}
}

func (m *MsgRemoveWhitelistRolePermission) Route() string {
	return ModuleName
}

func (m *MsgRemoveWhitelistRolePermission) Type() string {
	return types.MsgTypeRemoveWhitelistRolePermission
}

func (m *MsgRemoveWhitelistRolePermission) ValidateBasic() error {
	if m.Proposer.Empty() {
		return ErrEmptyProposerAccAddress
	}

	return nil
}

func (m *MsgRemoveWhitelistRolePermission) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m *MsgRemoveWhitelistRolePermission) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{
		m.Proposer,
	}
}

func NewMsgRemoveBlacklistRolePermission(
	proposer sdk.AccAddress,
	role uint32,
	permission uint32,
) *MsgRemoveBlacklistRolePermission {
	return &MsgRemoveBlacklistRolePermission{Proposer: proposer, Role: role, Permission: permission}
}

func (m *MsgRemoveBlacklistRolePermission) Route() string {
	return ModuleName
}

func (m *MsgRemoveBlacklistRolePermission) Type() string {
	return types.MsgTypeRemoveBlacklistRolePermission
}

func (m *MsgRemoveBlacklistRolePermission) ValidateBasic() error {
	if m.Proposer.Empty() {
		return ErrEmptyProposerAccAddress
	}

	return nil
}

func (m *MsgRemoveBlacklistRolePermission) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m *MsgRemoveBlacklistRolePermission) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{
		m.Proposer,
	}
}

func NewMsgCreateRole(proposer sdk.AccAddress, role uint32) *MsgCreateRole {
	return &MsgCreateRole{Proposer: proposer, Role: role}
}

func (m *MsgCreateRole) Route() string {
	return ModuleName
}

func (m *MsgCreateRole) Type() string {
	return types.MsgTypeCreateRole
}

func (m *MsgCreateRole) ValidateBasic() error {
	if m.Proposer.Empty() {
		return ErrEmptyProposerAccAddress
	}

	return nil
}

func (m *MsgCreateRole) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m *MsgCreateRole) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{
		m.Proposer,
	}
}

func NewMsgAssignRole(proposer, address sdk.AccAddress, role uint32) *MsgAssignRole {
	return &MsgAssignRole{Proposer: proposer, Address: address, Role: role}
}

func (m *MsgAssignRole) Route() string {
	return ModuleName
}

func (m *MsgAssignRole) Type() string {
	return types.MsgTypeAssignRole
}

func (m *MsgAssignRole) ValidateBasic() error {
	if m.Proposer.Empty() {
		return ErrEmptyProposerAccAddress
	}

	if m.Address.Empty() {
		return ErrEmptyPermissionsAccAddress
	}

	return nil
}

func (m *MsgAssignRole) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m *MsgAssignRole) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{
		m.Proposer,
	}
}

func NewMsgRemoveRole(proposer, address sdk.AccAddress, role uint32) *MsgRemoveRole {
	return &MsgRemoveRole{Proposer: proposer, Address: address, Role: role}
}

func (m *MsgRemoveRole) Route() string {
	return ModuleName
}

func (m *MsgRemoveRole) Type() string {
	return types.MsgTypeRemoveRole
}

func (m *MsgRemoveRole) ValidateBasic() error {
	if m.Proposer.Empty() {
		return ErrEmptyProposerAccAddress
	}

	if m.Address.Empty() {
		return ErrEmptyPermissionsAccAddress
	}

	return nil
}

func (m *MsgRemoveRole) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m *MsgRemoveRole) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{
		m.Proposer,
	}
}

func NewMsgProposalAssignPermission(proposer sdk.AccAddress, description string, address sdk.AccAddress, permission PermValue) *MsgProposalAssignPermission {
	return &MsgProposalAssignPermission{
		Proposer:    proposer,
		Description: description,
		Address:     address,
		Permission:  uint32(permission),
	}
}

func (m *MsgProposalAssignPermission) Route() string {
	return ModuleName
}

func (m *MsgProposalAssignPermission) Type() string {
	return types.MsgTypeProposalAssignPermission
}

func (m *MsgProposalAssignPermission) ValidateBasic() error {
	if m.Proposer.Empty() {
		return ErrEmptyProposerAccAddress
	}

	if m.Address.Empty() {
		return ErrEmptyPermissionsAccAddress
	}

	return nil
}

func (m *MsgProposalAssignPermission) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m *MsgProposalAssignPermission) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{
		m.Proposer,
	}
}

func NewMsgProposalSetNetworkProperty(proposer sdk.AccAddress, description string, property NetworkProperty, value uint64) *MsgProposalSetNetworkProperty {
	return &MsgProposalSetNetworkProperty{
		Proposer:        proposer,
		Description:     description,
		NetworkProperty: property,
		Value:           value,
	}
}

func (m *MsgProposalSetNetworkProperty) Route() string {
	return ModuleName
}

func (m *MsgProposalSetNetworkProperty) Type() string {
	return types.MsgTypeProposalSetNetworkProperty
}

func (m *MsgProposalSetNetworkProperty) ValidateBasic() error {
	if m.Proposer.Empty() {
		return ErrEmptyProposerAccAddress
	}

	switch m.NetworkProperty {
	case MinTxFee,
		MaxTxFee,
		VoteQuorum,
		ProposalEndTime,
		ProposalEnactmentTime,
		EnableForeignFeePayments,
		MischanceRankDecreaseAmount,
		InactiveRankDecreasePercent,
		PoorNetworkMaxBankSend,
		MinValidators:
		return nil
	default:
		return ErrInvalidNetworkProperty
	}
}

func (m *MsgProposalSetNetworkProperty) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m *MsgProposalSetNetworkProperty) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{
		m.Proposer,
	}
}

func NewMsgProposalUpsertDataRegistry(proposer sdk.AccAddress, description string, key, hash, reference, encoding string, size uint64) *MsgProposalUpsertDataRegistry {
	return &MsgProposalUpsertDataRegistry{
		Proposer:    proposer,
		Description: description,
		Key:         key,
		Hash:        hash,
		Reference:   reference,
		Encoding:    encoding,
		Size_:       size,
	}
}

func (m *MsgProposalUpsertDataRegistry) Route() string {
	return ModuleName
}

func (m *MsgProposalUpsertDataRegistry) Type() string {
	return types.MsgTypeProposalUpsertDataRegistry
}

func (m *MsgProposalUpsertDataRegistry) ValidateBasic() error {
	if m.Proposer.Empty() {
		return ErrEmptyProposerAccAddress
	}
	return nil
}

func (m *MsgProposalUpsertDataRegistry) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m *MsgProposalUpsertDataRegistry) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{
		m.Proposer,
	}
}

func NewMsgProposalSetPoorNetworkMessages(proposer sdk.AccAddress, description string, messages []string) *MsgProposalSetPoorNetworkMessages {
	return &MsgProposalSetPoorNetworkMessages{
		Proposer:    proposer,
		Description: description,
		Messages:    messages,
	}
}

func (m *MsgProposalSetPoorNetworkMessages) Route() string {
	return ModuleName
}

func (m *MsgProposalSetPoorNetworkMessages) Type() string {
	return types.MsgTypeProposalSetPoorNetworkMessages
}

func (m *MsgProposalSetPoorNetworkMessages) ValidateBasic() error {
	return nil
}

func (m *MsgProposalSetPoorNetworkMessages) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m *MsgProposalSetPoorNetworkMessages) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.Proposer}
}

func NewMsgVoteProposal(proposalID uint64, voter sdk.AccAddress, option VoteOption) *MsgVoteProposal {
	return &MsgVoteProposal{
		ProposalId: proposalID,
		Voter:      voter,
		Option:     option,
	}
}

func (m *MsgVoteProposal) Route() string {
	return ModuleName
}

func (m *MsgVoteProposal) Type() string {
	return types.MsgTypeVoteProposal
}

func (m *MsgVoteProposal) ValidateBasic() error {
	if m.Voter.Empty() {
		return ErrEmptyProposerAccAddress
	}

	return nil
}

func (m *MsgVoteProposal) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m *MsgVoteProposal) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{
		m.Voter,
	}
}

func NewMsgProposalCreateRole(
	proposer sdk.AccAddress,
	description string,
	role Role,
	whitelistPerms []PermValue,
	blacklistPerms []PermValue,
) *MsgProposalCreateRole {
	return &MsgProposalCreateRole{
		Proposer:               proposer,
		Description:            description,
		Role:                   uint32(role),
		WhitelistedPermissions: whitelistPerms,
		BlacklistedPermissions: blacklistPerms,
	}
}

func (m *MsgProposalCreateRole) Route() string {
	return ModuleName
}

func (m *MsgProposalCreateRole) Type() string {
	return types.MsgTypeProposalCreateRole
}

func (m *MsgProposalCreateRole) ValidateBasic() error {
	if m.Proposer.Empty() {
		return ErrEmptyProposerAccAddress
	}
	return nil
}

func (m *MsgProposalCreateRole) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m *MsgProposalCreateRole) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{
		m.Proposer,
	}
}
