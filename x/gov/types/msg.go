package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Msg types
const (
	WhitelistPermissions = "whitelist-permissions"
	BlacklistPermissions = "blacklist-permissions"

	ClaimCouncilor = "claim-councilor"

	CreateRole = "create-role"
	AssignRole = "assign-role"
	RemoveRole = "remove-role"

	WhitelistRolePermission       = "whitelist-role-permission"
	BlacklistRolePermission       = "blacklist-role-permission"
	RemoveWhitelistRolePermission = "remove-whitelist-role-permission"
	RemoveBlacklistRolePermission = "remove-blacklist-role-permission"
)

var (
	_ sdk.Msg = &MsgWhitelistPermissions{}
	_ sdk.Msg = &MsgBlacklistPermissions{}

	_ sdk.Msg = &MsgClaimCouncilor{}

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
	return WhitelistPermissions
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
	return BlacklistPermissions
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
	return ClaimCouncilor
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
	return WhitelistRolePermission
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
	return BlacklistRolePermission
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
	return RemoveWhitelistRolePermission
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
	return RemoveBlacklistRolePermission
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
	return CreateRole
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
	return AssignRole
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
	return RemoveRole
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
