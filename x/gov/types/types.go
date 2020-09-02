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

// AddToBlacklist adds permission to whitelist.
func (p *Permissions) AddToBlacklist(perm PermValue) error {
	if p.IsWhitelisted(perm) {
		return fmt.Errorf("permission is already whitelisted")
	}

	p.Blacklist = append(p.Blacklist, uint32(perm))
	return nil
}

func NewNetworkActor(
	addr types.AccAddress,
	roles Roles,
	status uint32,
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
		0,
		nil,
		NewPermissions(nil, nil),
		0,
	)
}