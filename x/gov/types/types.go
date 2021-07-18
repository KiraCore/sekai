package types

import (
	"fmt"

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

func NewCouncilor(
	moniker string,
	address types.AccAddress,
) Councilor {
	return Councilor{
		Moniker: moniker,
		Address: address,
	}
}
