package createOrderBook

import (
	"encoding/hex"

	"github.com/KiraCore/cosmos-sdk/codec"
	sdk "github.com/KiraCore/cosmos-sdk/types"
	"github.com/KiraCore/sekai/types"

	"golang.org/x/crypto/blake2b"
)

type Keeper struct {
	cdc *codec.Codec // The wire codec for binary encoding/decoding.
	storeKey sdk.StoreKey // Unexposed key to access store from sdk.Context
}

func NewKeeper(cdc *codec.Codec, storeKey sdk.StoreKey) Keeper {
	return Keeper{
		cdc:        cdc,
		storeKey:   storeKey,
	}
}

func (k Keeper) CreateOrderBook(ctx sdk.Context, quote sdk.Coins, base sdk.Coins, curator sdk.AccAddress, mnemonic string) {
	var orderbook = types.OrderBook{}

	orderbook.Quote = quote
	orderbook.Base = base
	orderbook.Curator = curator
	orderbook.Mnemonic = mnemonic

	// This is the definitions of the lens of the shortened hashes
	numberOfBytes := 4
	numberOfCharacters := 2*numberOfBytes

	// Creating the hashes of the parts of the ID
	hashOfCurator := blake2b.Sum256([]byte(curator))
	hashInStringOfCurator := hex.EncodeToString(hashOfCurator[:])
	idHashInStringOfCurator := hashInStringOfCurator[len(hashInStringOfCurator) - numberOfCharacters:]

	store := ctx.KVStore(k.storeKey)

	store.Set([]byte(id), k.cdc.MustMarshalBinaryBare(orderbook))
}