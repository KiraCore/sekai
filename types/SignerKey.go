package types

import (
	"fmt"
	"strings"

	sdk "github.com/KiraCore/cosmos-sdk/types"
)

// SignerKeyType describes the enum for signer key types
type SignerKeyType int

// define signer key types
const (
	Secp256k1 SignerKeyType = iota
	Ed25519
)

func (e SignerKeyType) String() string {
	switch e {
	case Secp256k1:
		return "Secp256k1"
	case Ed25519:
		return "Ed25519"
	default:
		return fmt.Sprintf("%d", int(e))
	}
}

// SignerKey describes signer public keys with its status
type SignerKey struct {
	PubKey      string         `json:"pubkey"`      // pubkey - New public key to register (max 4096 Bytes)
	KeyType     SignerKeyType  `json:"type"`        // type - Key type enum (e.g. secp256k1, ed25519)
	ExpiryTime  int64          `json:"expires"`     // expires - UTC time (8 Bytes) when key expires
	Enabled     bool           `json:"enabled"`     // enabled - boolean field defining if key is use
	Permissions []int          `json:"permissions"` // permissions - array of integers defining permissions that the owner assigned to the key (this field does not have to be used in PoC)
	Curator     sdk.AccAddress `json:"curator"`     // the address that sent this public key
}

// NewSignerKey returns SignerKey instance
func NewSignerKey(
	pubKey string,
	keyType SignerKeyType,
	expiryTime int64,
	enabled bool,
	permissions []int,
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

func (o SignerKey) String() string {
	return strings.TrimSpace(fmt.Sprintf(`PubKey: %s, KeyType: %+v, ExpiryTime: %d, Enabled: %t, Permissions: %+v`, o.PubKey, o.KeyType, o.ExpiryTime, o.Enabled, o.Permissions))
}
