package types

// FindTokenIndex find token index from tokens list
func FindTokenIndex(tokens []string, token string) int {
	for index, t := range tokens {
		if t == token {
			return index
		}
	}
	return -1
}

// IsFrozen returns is frozen
func (t TokensWhiteBlack) IsFrozen(denom string, enableTokenBlacklist, enableTokenWhitelist bool) bool {
	if enableTokenBlacklist {
		if FindTokenIndex(t.Blacklisted, denom) >= 0 {
			return true
		}
	}
	if enableTokenWhitelist {
		if FindTokenIndex(t.Whitelisted, denom) < 0 {
			return true
		}
	}
	return false
}
