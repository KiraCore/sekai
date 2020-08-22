package types

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
