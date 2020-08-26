package types

import sdk "github.com/cosmos/cosmos-sdk/types"

// NewSignerKey returns SignerKey instance
func NewSignerKey(
	pubKey string,
	keyType SignerKeyType,
	expiryTime int64,
	enabled bool,
	permissions []int64,
	curator sdk.AccAddress) SignerKey {
	return SignerKey{
		PubKey:      pubKey,
		KeyType:     keyType,
		ExpiryTime:  expiryTime,
		Enabled:     enabled,
		Permissions: permissions,
		Curator:     curator,
	}
}
