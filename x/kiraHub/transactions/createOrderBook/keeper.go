package createOrderBook

import (
	"encoding/hex"
	"strings"
	"fmt"

	"github.com/KiraCore/cosmos-sdk/codec"
	sdk "github.com/KiraCore/cosmos-sdk/types"
	"github.com/KiraCore/sekai/types"

	"golang.org/x/crypto/blake2b"
)

type Keeper struct {
	cdc *codec.Codec // The wire codec for binary encoding/decoding.
	storeKey sdk.StoreKey // Unexposed key to access store from sdk.Context
}

var last_order_book_index = 0

// This is the definitions of the lens of the shortened hashes
var numberOfBytes = 4
var numberOfCharacters = 2 * numberOfBytes


func NewKeeper(cdc *codec.Codec, storeKey sdk.StoreKey) Keeper {
	return Keeper{
		cdc:        cdc,
		storeKey:   storeKey,
	}
}

func (k Keeper) CreateOrderBook(ctx sdk.Context, quote string, base string, curator sdk.AccAddress, mnemonic string) {
	var orderbook = types.OrderBook{}

	orderbook.Quote = quote
	orderbook.Base = base
	orderbook.Curator = curator
	orderbook.Mnemonic = mnemonic

	// ARJUN CHANGE THIS TO THE DYNAMIC INDEX PULLED FROM THE KVSTORE
	// var last_order_book_index = 1



	// Creating the hashes of the parts of the ID
	hashOfCurator := blake2b.Sum256(curator)
	hashInStringOfCurator := hex.EncodeToString(hashOfCurator[:])
	idHashInStringOfCurator := hashInStringOfCurator[len(hashInStringOfCurator) - numberOfCharacters:]

	hashOfBase := blake2b.Sum256([]byte(base))
	hashInStringOfBase := hex.EncodeToString(hashOfBase[:])
	idHashInStringOfBase := hashInStringOfBase[len(hashInStringOfBase) - numberOfCharacters:]

	hashOfQuote := blake2b.Sum256([]byte(quote))
	hashInStringOfQuote := hex.EncodeToString(hashOfQuote[:])
	idHashInStringOfQuote := hashInStringOfQuote[len(hashInStringOfQuote) - numberOfCharacters:]

	idHashInStringOfIndex = fmt.Sprintf("%x", len(last_order_book_index))
	var ID strings.Builder

	ID.WriteString(idHashInStringOfCurator)
	ID.WriteString(idHashInStringOfBase)
	ID.WriteString(idHashInStringOfQuote)
	ID.WriteString(idHashInStringOfIndex)
	// Still need to add the functionalities of last_order_book_index

	id := ID.String()
	orderbook.ID = id

	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(id), k.cdc.MustMarshalBinaryBare(orderbook))

	var idsArray []string

	bz := store.Get([]byte("ids"))

	if len(bz) != 0 {
		k.cdc.MustUnmarshalBinaryBare(bz, &idsArray)
	}

	idsArray = append(idsArray, id)
	store.Set([]byte("ids"), k.cdc.MustMarshalBinaryBare(idsArray))

}

func (k Keeper) GetOrderBookByID(ctx sdk.Context, id string) []types.OrderBook {

	store := ctx.KVStore(k.storeKey)
	bz := store.Get([]byte(id))

	var orderbook types.OrderBook
	k.cdc.MustUnmarshalBinaryBare(bz, &orderbook)

	var orderbooksQueried = []types.OrderBook{orderbook}
	return orderbooksQueried
}

func (k Keeper) GetOrderBookByBase(ctx sdk.Context, base string) []types.OrderBook {

	store := ctx.KVStore(k.storeKey)

	var orderbook types.OrderBook
	var orderbooksQueried = []types.OrderBook{}
	var idsArray []string

	hashOfBase := blake2b.Sum256([]byte(base))
	hashInStringOfBase := hex.EncodeToString(hashOfBase[:])
	idHashInStringOfBase := hashInStringOfBase[len(hashInStringOfBase) - numberOfCharacters:]

	bz := store.Get([]byte("ids"))
	k.cdc.MustUnmarshalBinaryBare(bz, &idsArray)

	for index, id := range idsArray {

		// Matching
		if idHashInStringOfBase == id[numberOfCharacters: 2 * numberOfCharacters] {
			bz := store.Get([]byte(id))
			k.cdc.MustUnmarshalBinaryBare(bz, &orderbook)
			orderbooksQueried = append(orderbooksQueried, orderbook)
		}
	}

	return orderbooksQueried
}

func (k Keeper) GetOrderBookByQuote(ctx sdk.Context, quote string) []types.OrderBook {

	store := ctx.KVStore(k.storeKey)

	var orderbook types.OrderBook
	var orderbooksQueried = []types.OrderBook{}
	var idsArray []string

	hashOfQuote := blake2b.Sum256([]byte(quote))
	hashInStringOfQuote := hex.EncodeToString(hashOfQuote[:])
	idHashInStringOfQuote := hashInStringOfQuote[len(hashInStringOfQuote) - numberOfCharacters:]

	bz := store.Get([]byte("ids"))
	k.cdc.MustUnmarshalBinaryBare(bz, &idsArray)

	for index, id := range idsArray {

		// Matching
		if idHashInStringOfQuote == id[2 * numberOfCharacters: 3 * numberOfCharacters] {
			bz := store.Get([]byte(id))
			k.cdc.MustUnmarshalBinaryBare(bz, &orderbook)
			orderbooksQueried = append(orderbooksQueried, orderbook)
		}
	}

	return orderbooksQueried
}

func (k Keeper) GetOrderBookByTP(ctx sdk.Context, base string, quote string) []types.OrderBook {

	store := ctx.KVStore(k.storeKey)

	var orderbook types.OrderBook
	var orderbooksQueried = []types.OrderBook{}
	var idsArray []string

	hashOfBase := blake2b.Sum256([]byte(base))
	hashInStringOfBase := hex.EncodeToString(hashOfBase[:])
	idHashInStringOfBase := hashInStringOfBase[len(hashInStringOfBase) - numberOfCharacters:]

	hashOfQuote := blake2b.Sum256([]byte(quote))
	hashInStringOfQuote := hex.EncodeToString(hashOfQuote[:])
	idHashInStringOfQuote := hashInStringOfQuote[len(hashInStringOfQuote) - numberOfCharacters:]

	bz := store.Get([]byte("ids"))
	k.cdc.MustUnmarshalBinaryBare(bz, &idsArray)

	for index, id := range idsArray {

		// Matching
		if idHashInStringOfBase == id[numberOfCharacters: 2 * numberOfCharacters] &&
			idHashInStringOfQuote == id[2 * numberOfCharacters: 3 * numberOfCharacters] {
			bz := store.Get([]byte(id))
			k.cdc.MustUnmarshalBinaryBare(bz, &orderbook)
			orderbooksQueried = append(orderbooksQueried, orderbook)
		}
	}

	return orderbooksQueried
}

func (k Keeper) GetOrderBookByCurator(ctx sdk.Context, curator string) []types.OrderBook {

	store := ctx.KVStore(k.storeKey)

	var orderbook types.OrderBook
	var orderbooksQueried = []types.OrderBook{}
	var idsArray []string

	hashOfCurator := blake2b.Sum256([]byte(curator))
	hashInStringOfCurator := hex.EncodeToString(hashOfCurator[:])
	idHashInStringOfCurator := hashInStringOfCurator[len(hashInStringOfCurator) - numberOfCharacters:]

	bz := store.Get([]byte("ids"))
	k.cdc.MustUnmarshalBinaryBare(bz, &idsArray)

	for index, id := range idsArray {

		// Matching
		if idHashInStringOfCurator == id[0:numberOfCharacters] {
			bz := store.Get([]byte(id))
			k.cdc.MustUnmarshalBinaryBare(bz, &orderbook)
			orderbooksQueried = append(orderbooksQueried, orderbook)
		}
	}

	return orderbooksQueried
}