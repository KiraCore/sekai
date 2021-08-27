package types

import (
	"fmt"

	"github.com/KiraCore/sekai/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	proto "github.com/gogo/protobuf/proto"
	"gopkg.in/yaml.v2"
)

var (
	// Proposal
	_ sdk.Msg = &MsgVoteProposal{}

	// Permissions
	_ sdk.Msg = &MsgWhitelistPermissions{}
	_ sdk.Msg = &MsgBlacklistPermissions{}
	_ sdk.Msg = &MsgSubmitProposal{}

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
	address sdk.AccAddress,
) *MsgClaimCouncilor {
	return &MsgClaimCouncilor{
		Moniker: moniker,
		Address: address,
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

func NewMsgRegisterIdentityRecords(address sdk.AccAddress, infos []IdentityInfoEntry) *MsgRegisterIdentityRecords {
	return &MsgRegisterIdentityRecords{
		Address: address,
		Infos:   infos,
	}
}

func (m *MsgRegisterIdentityRecords) Route() string {
	return ModuleName
}

func (m *MsgRegisterIdentityRecords) Type() string {
	return types.MsgTypeRegisterIdentityRecords
}

func (m *MsgRegisterIdentityRecords) ValidateBasic() error {
	if m.Address.Empty() {
		return ErrEmptyProposerAccAddress
	}
	if len(m.Infos) == 0 {
		return ErrEmptyInfos
	}
	return nil
}

func (m *MsgRegisterIdentityRecords) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m *MsgRegisterIdentityRecords) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{
		m.Address,
	}
}

func NewMsgDeleteIdentityRecords(address sdk.AccAddress, keys []string) *MsgDeleteIdentityRecords {
	return &MsgDeleteIdentityRecords{
		Address: address,
		Keys:    keys,
	}
}

func (m *MsgDeleteIdentityRecords) Route() string {
	return ModuleName
}

func (m *MsgDeleteIdentityRecords) Type() string {
	return types.MsgTypeEditIdentityRecord
}

func (m *MsgDeleteIdentityRecords) ValidateBasic() error {
	if m.Address.Empty() {
		return ErrEmptyProposerAccAddress
	}
	return nil
}

func (m *MsgDeleteIdentityRecords) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m *MsgDeleteIdentityRecords) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{
		m.Address,
	}
}

func NewMsgRequestIdentityRecordsVerify(address, verifier sdk.AccAddress, recordIds []uint64, tip sdk.Coin) *MsgRequestIdentityRecordsVerify {
	return &MsgRequestIdentityRecordsVerify{
		Address:   address,
		Verifier:  verifier,
		RecordIds: recordIds,
		Tip:       tip,
	}
}

func (m *MsgRequestIdentityRecordsVerify) Route() string {
	return ModuleName
}

func (m *MsgRequestIdentityRecordsVerify) Type() string {
	return types.MsgTypeRequestIdentityRecordsVerify
}

func (m *MsgRequestIdentityRecordsVerify) ValidateBasic() error {
	if m.Address.Empty() {
		return ErrEmptyProposerAccAddress
	}
	if m.Verifier.Empty() {
		return ErrEmptyVerifierAccAddress
	}
	if !m.Tip.IsValid() {
		return ErrInvalidTip
	}
	if len(m.RecordIds) == 0 {
		return ErrInvalidRecordIds
	}
	return nil
}

func (m *MsgRequestIdentityRecordsVerify) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m *MsgRequestIdentityRecordsVerify) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{
		m.Address,
	}
}

func NewMsgHandleIdentityRecordsVerifyRequest(verifier sdk.AccAddress, requestId uint64, isApprove bool) *MsgHandleIdentityRecordsVerifyRequest {
	return &MsgHandleIdentityRecordsVerifyRequest{
		Verifier:        verifier,
		VerifyRequestId: requestId,
		Yes:             isApprove,
	}
}

func (m *MsgHandleIdentityRecordsVerifyRequest) Route() string {
	return ModuleName
}

func (m *MsgHandleIdentityRecordsVerifyRequest) Type() string {
	return types.MsgTypeHandleIdentityRecordsVerifyRequest
}

func (m *MsgHandleIdentityRecordsVerifyRequest) ValidateBasic() error {
	if m.Verifier.Empty() {
		return ErrEmptyVerifierAccAddress
	}
	if m.VerifyRequestId == 0 {
		return ErrInvalidVerifyRequestId
	}
	return nil
}

func (m *MsgHandleIdentityRecordsVerifyRequest) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m *MsgHandleIdentityRecordsVerifyRequest) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{
		m.Verifier,
	}
}

func NewMsgCancelIdentityRecordsVerifyRequest(executor sdk.AccAddress, verifyRequestId uint64) *MsgCancelIdentityRecordsVerifyRequest {
	return &MsgCancelIdentityRecordsVerifyRequest{
		Executor:        executor,
		VerifyRequestId: verifyRequestId,
	}
}

func (m *MsgCancelIdentityRecordsVerifyRequest) Route() string {
	return ModuleName
}

func (m *MsgCancelIdentityRecordsVerifyRequest) Type() string {
	return types.MsgTypeCancelIdentityRecordsVerifyRequest
}

func (m *MsgCancelIdentityRecordsVerifyRequest) ValidateBasic() error {
	if m.Executor.Empty() {
		return ErrEmptyProposerAccAddress
	}
	if m.VerifyRequestId == 0 {
		return ErrInvalidVerifyRequestId
	}
	return nil
}

func (m *MsgCancelIdentityRecordsVerifyRequest) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m *MsgCancelIdentityRecordsVerifyRequest) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{
		m.Executor,
	}
}

// NewMsgSubmitProposal creates a new MsgSubmitProposal.
//nolint:interfacer
func NewMsgSubmitProposal(proposer sdk.AccAddress, title, description string, content Content) (*MsgSubmitProposal, error) {
	m := &MsgSubmitProposal{
		Proposer:    proposer,
		Title:       title,
		Description: description,
	}
	err := m.SetContent(content)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (m *MsgSubmitProposal) GetProposer() sdk.AccAddress {
	return m.Proposer
}

func (m *MsgSubmitProposal) GetContent() Content {
	content, ok := m.Content.GetCachedValue().(Content)
	if !ok {
		return nil
	}
	return content
}

func (m *MsgSubmitProposal) SetContent(content Content) error {
	msg, ok := content.(proto.Message)
	if !ok {
		return fmt.Errorf("can't proto marshal %T", msg)
	}
	any, err := codectypes.NewAnyWithValue(msg)
	if err != nil {
		return err
	}
	m.Content = any
	return nil
}

// Route implements Msg
func (m MsgSubmitProposal) Route() string { return RouterKey }

// Type implements Msg
func (m MsgSubmitProposal) Type() string { return types.MsgTypeSubmitProposal }

// ValidateBasic implements Msg
func (m MsgSubmitProposal) ValidateBasic() error {
	if m.Proposer.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, m.Proposer.String())
	}

	content := m.GetContent()
	if content == nil {
		return sdkerrors.Wrap(ErrInvalidProposalContent, "missing content")
	}
	if err := content.ValidateBasic(); err != nil {
		return err
	}

	return nil
}

// GetSignBytes implements Msg
func (m MsgSubmitProposal) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&m)
	return sdk.MustSortJSON(bz)
}

// GetSigners implements Msg
func (m MsgSubmitProposal) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.Proposer}
}

// String implements the Stringer interface
func (m MsgSubmitProposal) String() string {
	out, _ := yaml.Marshal(m)
	return string(out)
}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (m MsgSubmitProposal) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	var content Content
	return unpacker.UnpackAny(m.Content, &content)
}
