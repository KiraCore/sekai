package signerkey

import (
	"time"

	"github.com/KiraCore/cosmos-sdk/codec"
	sdk "github.com/KiraCore/cosmos-sdk/types"
	"github.com/KiraCore/sekai/types"
)

// Keeper is an interface to keep signer keys
type Keeper struct {
	cdc      *codec.Codec // The wire codec for binary encoding/decoding.
	storeKey sdk.StoreKey // Unexposed key to access store from sdk.Context
}

// GetSignerKeys return SignerKeys by a curator
func (k Keeper) GetSignerKeys(ctx sdk.Context, curator sdk.AccAddress) []types.SignerKey {

	var signerKeys []types.SignerKey
	var queryOutput = []types.SignerKey{}

	store := ctx.KVStore(k.storeKey)
	bz := store.Get([]byte("signer_keys")) // TODO: should use iterator instead of this

	k.cdc.MustUnmarshalBinaryBare(bz, &signerKeys)

	for _, signerKey := range signerKeys {
		if signerKey.Curator.Equals(curator) {
			queryOutput = append(queryOutput, signerKey)
		}
	}

	return queryOutput
}

// NewKeeper is a utility to create a keeper
func NewKeeper(cdc *codec.Codec, storeKey sdk.StoreKey) Keeper {
	return Keeper{
		cdc:      cdc,
		storeKey: storeKey,
	}
}

var lastOrderIndex uint32 = 0

// This is the definitions of the lens of the shortened hashes
var numberOfBytes = 4
var numberOfCharacters = 2 * numberOfBytes

// CreateSignerKey create signer key and put it into the keeper
func (k Keeper) CreateSignerKey(ctx sdk.Context,
	pubKey [4096]byte,
	keyType types.SignerKeyType,
	Permissions []int,
	curator sdk.AccAddress) {

	var signerKeys []types.SignerKey
	now := time.Now()
	unix := now.Unix() // TODO: this won't work as every node has little time differece in unix

	var signerKey = types.NewSignerKey(pubKey, keyType, unix, true, Permissions, curator)

	// Storage Logic
	store := ctx.KVStore(k.storeKey)
	bz := store.Get([]byte("signer_keys"))

	k.cdc.MustUnmarshalBinaryBare(bz, &signerKeys)
	signerKeys = append(signerKeys, signerKey)

	store.Set([]byte("signer_keys"), k.cdc.MustMarshalBinaryBare(signerKeys))
}

// TODO: should add deleteSignerKey and disableSignerKey after discussion
