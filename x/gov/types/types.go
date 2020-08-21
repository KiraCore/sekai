package types

func (p Permissions) IsBlacklisted(perm PermValue) bool {
	for _, bPerm := range p.Blacklist {
		if bPerm == uint32(perm) {
			return true
		}
	}

	return false
}
