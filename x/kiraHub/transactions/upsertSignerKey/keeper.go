package signerkey

import (
	"errors"
	"strings"
	"time"

	"github.com/KiraCore/cosmos-sdk/codec"
	sdk "github.com/KiraCore/cosmos-sdk/types"
	"github.com/KiraCore/sekai/types"
)

// KeySignerKeys describes the key where to save in KVStore
const KeySignerKeys = "signer_keys"

// KeyPubKeyCurator describes the owner of each pubKey
const KeyPubKeyCurator = "pub_key_curator"

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
	bz := store.Get([]byte(KeySignerKeys)) // TODO: should use iterator instead of this

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

// UpsertSignerKey create signer key and put it into the keeper
func (k Keeper) UpsertSignerKey(ctx sdk.Context,
	pubKey string,
	keyType types.SignerKeyType,
	Permissions []int,
	curator sdk.AccAddress) error {

	var newSignerKeys []types.SignerKey
	// TODO: expiry key should be entered from a user or set automatically?
	// for now, set it to last block's time + 10 days
	// TODO: order should use block time instead of current timestamp from local computer
	unix := ctx.BlockHeader().Time.Unix() + time.Hour.Milliseconds()*24*10

	var signerKey = types.NewSignerKey(pubKey, keyType, unix, true, Permissions, curator)

	// Storage Logic
	store := ctx.KVStore(k.storeKey)
	bz := store.Get([]byte(KeySignerKeys))

	var signerKeys []types.SignerKey
	k.cdc.MustUnmarshalBinaryBare(bz, &signerKeys)
	// TODO: we need to create index that will help us quickly identify keys belonging to specific user.
	// TODO: must add a check to make sure that 2 accounts can't have the same sub-key
	// TODO: navigating around whole signer keys is inefficient, should update it to efficient and make it by sender
	for _, sk := range signerKeys {
		if strings.Compare(sk.PubKey, pubKey) == 0 {
			if keyType == sk.KeyType {
				return errors.New("keyType shouldn't be different for same pub key")
			}
			if sk.Curator.Equals(curator) {
				return errors.New("this key is owned by another curator already")
			}
			newSignerKeys = append(newSignerKeys, signerKey)
			store.Set([]byte(KeyPubKeyCurator), []byte(signerKey.PubKey))
		} else if sk.ExpiryTime > unix {
			newSignerKeys = append(newSignerKeys, sk)
		}
	}
	// TODO: easily query if sub-key x belongs to account y

	store.Set([]byte(KeySignerKeys), k.cdc.MustMarshalBinaryBare(newSignerKeys))
	// TODO: should add test for creating / updating after v0.0.5 release.
	return nil
}

// TODO: should add deleteSignerKey after discussion but this should create another directory under transactions folder?
