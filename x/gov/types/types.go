package types

import "github.com/cosmos/cosmos-sdk/types"

// NewPermissions generates a new permissions struct.
func NewPermissions(whitelist []PermValue, blacklist []PermValue) Permissions {
	var b []uint32
	for _, bv := range blacklist {
		b = append(b, uint32(bv))
	}

	var w []uint32
	for _, wv := range whitelist {
		w = append(w, uint32(wv))
	}

	return Permissions{
		Blacklist: b,
		Whitelist: w,
	}
}

// IsBlacklisted returns if the perm is blacklisted or not.
func (p Permissions) IsBlacklisted(perm PermValue) bool {
	for _, bPerm := range p.Blacklist {
		if bPerm == uint32(perm) {
			return true
		}
	}

	return false
}

// IsWhitelisted returns if the perm is whitelisted or not.
func (p Permissions) IsWhitelisted(perm PermValue) bool {
	for _, bPerm := range p.Whitelist {
		if bPerm == uint32(perm) {
			return true
		}
	}

	return false
}

func NewNetworkActor(
	addr types.AccAddress,
	roles Roles,
	status uint32,
	votes []uint32,
	perm Permissions,
	skin uint64,
) NetworkActor {
	return NetworkActor{
		Address:     addr,
		Roles:       roles,
		Status:      status,
		Votes:       votes,
		Permissions: &perm,
		Skin:        skin,
	}
}

func (m NetworkActor) HasPermissionFor(perm PermValue) bool {
	if m.Permissions.IsWhitelisted(perm) {
		return true
	}

	return false
}
