package types

import (
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/types"
)

// NewPermissions generates a new permissions struct.
func NewPermissions(whitelist []PermValue, blacklist []PermValue) *Permissions {
	var b []uint32
	for _, bv := range blacklist {
		b = append(b, uint32(bv))
	}

	var w []uint32
	for _, wv := range whitelist {
		w = append(w, uint32(wv))
	}

	return &Permissions{
		Blacklist: b,
		Whitelist: w,
	}
}

// IsBlacklisted returns if the perm is blacklisted or not.
func (p *Permissions) IsBlacklisted(perm PermValue) bool {
	for _, bPerm := range p.Blacklist {
		if bPerm == uint32(perm) {
			return true
		}
	}

	return false
}

// IsWhitelisted returns if the perm is whitelisted or not.
func (p *Permissions) IsWhitelisted(perm PermValue) bool {
	for _, bPerm := range p.Whitelist {
		if bPerm == uint32(perm) {
			return true
		}
	}

	return false
}

// AddToWhitelist adds permission to whitelist.
func (p *Permissions) AddToWhitelist(perm PermValue) error {
	if p.IsBlacklisted(perm) {
		return fmt.Errorf("permission is already blacklisted")
	}

	p.Whitelist = append(p.Whitelist, uint32(perm))
	return nil
}

// AddToBlacklist adds permission to blacklist. It fails if the permission is whitelisted.
func (p *Permissions) AddToBlacklist(perm PermValue) error {
	if p.IsWhitelisted(perm) {
		return fmt.Errorf("permission is already whitelisted")
	}

	p.Blacklist = append(p.Blacklist, uint32(perm))
	return nil
}

// RemoveFromWhitelist removes permission from whitelist. It fails if permission is not
// whitelisted.
func (m *Permissions) RemoveFromWhitelist(perm PermValue) error {
	if !m.IsWhitelisted(perm) {
		return fmt.Errorf("permission is not whitelisted")
	}

	for i, v := range m.Whitelist {
		if v == uint32(perm) {
			m.Whitelist = append(m.Whitelist[:i], m.Whitelist[i+1:]...)
			return nil
		}
	}

	return nil
}

func (m *Permissions) RemoveFromBlacklist(perm PermValue) error {
	if !m.IsBlacklisted(perm) {
		return fmt.Errorf("permission is not blacklisted")
	}

	for i, v := range m.Blacklist {
		if v == uint32(perm) {
			m.Blacklist = append(m.Blacklist[:i], m.Blacklist[i+1:]...)
			return nil
		}
	}

	return nil
}

func NewNetworkActor(
	addr types.AccAddress,
	roles Roles,
	status ActorStatus,
	votes []uint32,
	perm *Permissions,
	skin uint64,
) NetworkActor {
	return NetworkActor{
		Address:     addr,
		Roles:       roles,
		Status:      status,
		Votes:       votes,
		Permissions: perm,
		Skin:        skin,
	}
}

func (m *NetworkActor) HasRole(role Role) bool {
	for _, r := range m.Roles {
		if r == uint64(role) {
			return true
		}
	}
	return false
}

func (m *NetworkActor) SetRole(role Role) {
	if !m.HasRole(role) {
		m.Roles = append(m.Roles, uint64(role))
	}
}

func (m *NetworkActor) RemoveRole(role Role) {
	for i, r := range m.Roles {
		if r == uint64(role) {
			m.Roles = append(m.Roles[:i], m.Roles[i+1:]...)
			return
		}
	}
}

func (m *NetworkActor) IsActive() bool {
	return m.Status == Active
}

// NewDefaultActor returns a default actor with:
// - The provided addr.
// - Roles set to nil
// - Status set to 0
// - Votes set to nil
// - Empty permissions
// - Skin set to 0
func NewDefaultActor(addr types.AccAddress) NetworkActor {
	return NewNetworkActor(
		addr,
		nil,
		Undefined,
		nil,
		NewPermissions(nil, nil),
		0,
	)
}

func NewCouncilor(
	moniker string,
	website string,
	social string,
	identity string,
	address types.AccAddress,
) Councilor {
	return Councilor{
		Moniker:  moniker,
		Website:  website,
		Social:   social,
		Identity: identity,
		Address:  address,
	}
}

func NewProposalAssignPermission(address types.AccAddress, permission PermValue, submitTime time.Time) ProposalAssignPermission {
	return ProposalAssignPermission{
		Address:    address,
		Permission: uint32(permission),
		SubmitTime: submitTime,
	}
}

func NewVote(proposalID uint64, addr types.AccAddress, option VoteOption) Vote {
	return Vote{
		ProposalId: proposalID,
		Voter:      addr,
		Option:     option,
	}
}
