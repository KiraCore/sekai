package signerkey

import (
	"errors"
	"strings"
	"time"

	"github.com/KiraCore/cosmos-sdk/codec"
	sdk "github.com/KiraCore/cosmos-sdk/types"
	"github.com/KiraCore/sekai/types"
)

// PrefixKeySignerKeys describes the key where to save in KVStore
const PrefixKeySignerKeys = "signer_keys"

// PrefixKeyPubKeyCurator describes the owner of each pubKey
const PrefixKeyPubKeyCurator = "pub_key_curator"

// Keeper is an interface to keep signer keys
type Keeper struct {
	cdc      *codec.Codec // The wire codec for binary encoding/decoding.
	storeKey sdk.StoreKey // Unexposed key to access store from sdk.Context
}

// GetSignerKeys return SignerKeys by a curator
func (k Keeper) GetSignerKeys(ctx sdk.Context, curator sdk.AccAddress) []types.SignerKey {

	var signerKeys []types.SignerKey

	store := ctx.KVStore(k.storeKey)
	curatorStoreID := append([]byte(PrefixKeySignerKeys), []byte(curator)...)

	if store.Has(curatorStoreID) {
		bz := store.Get(curatorStoreID)
		k.cdc.MustUnmarshalBinaryBare(bz, &signerKeys)
	}

	return signerKeys
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
	// TODO: createOrder/createOrderBook should use block time instead of current timestamp from local computer
	unix := ctx.BlockHeader().Time.Unix() + time.Hour.Milliseconds()*24*10

	var signerKey = types.NewSignerKey(pubKey, keyType, unix, true, Permissions, curator)

	var signerKeys []types.SignerKey
	// Storage Logic
	store := ctx.KVStore(k.storeKey)
	curatorStoreID := append([]byte(PrefixKeySignerKeys), []byte(curator)...)

	if store.Has(curatorStoreID) {
		bz := store.Get(curatorStoreID)
		k.cdc.MustUnmarshalBinaryBare(bz, &signerKeys)
	}

	pubKeyStoreID := append([]byte(PrefixKeyPubKeyCurator), []byte(pubKey)...)

	for _, sk := range signerKeys {
		if strings.Compare(sk.PubKey, pubKey) == 0 {
			if keyType == sk.KeyType {
				return errors.New("keyType shouldn't be different for same pub key")
			}
			newSignerKeys = append(newSignerKeys, signerKey)
		} else if sk.ExpiryTime > unix {
			newSignerKeys = append(newSignerKeys, sk)
		} else { // Delete pubKey curator when it is expired
			if !store.Has(pubKeyStoreID) {
				return errors.New("pubKey to curator mapping is not set properly at the time of expired key cleanup")
			}
			store.Delete(pubKeyStoreID)
		}
	}
	if !store.Has([]byte(pubKey)) { // when pubKey's owner does not exist
		newSignerKeys = append(newSignerKeys, signerKey)
		// Set pubKey curator when pubKey is newly added
		store.Set(pubKeyStoreID, curator)
	} else {
		originCurator := store.Get(pubKeyStoreID)
		if !curator.Equals(sdk.AccAddress(originCurator)) {
			return errors.New("this key is owned by another curator already")
		}
	}

	store.Set(curatorStoreID, k.cdc.MustMarshalBinaryBare(newSignerKeys))
	return nil
}

// TODO: should add test for creating / updating after v0.0.5 release.
// TODO: should add deleteSignerKey after discussion but this should create another directory under transactions folder?
