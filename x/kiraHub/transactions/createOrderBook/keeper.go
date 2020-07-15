package createOrderBook

import (
	"encoding/hex"
	"strconv"
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

func NewKeeper(cdc *codec.Codec, storeKey sdk.StoreKey) Keeper {
	return Keeper{
		cdc:        cdc,
		storeKey:   storeKey,
	}
}

type meta struct {
	ID string
	Index uint32
}

func newMeta(id string, index uint32) meta {
	return meta{
		ID: id,
		Index: index,
	}
}

var lastOrderBookIndex uint32 = 0

// This is the definitions of the lens of the shortened hashes
var numberOfBytes = 4
var numberOfCharacters = 2 * numberOfBytes

func (k Keeper) CreateOrderBook(ctx sdk.Context, quote string, base string, curator sdk.AccAddress, mnemonic string) {
	var orderbook = types.NewOrderBook()

	fmt.Println("Last Block ID: ", ctx.BlockHeader().LastBlockId)

	orderbook.Quote = quote
	orderbook.Base = base
	orderbook.Curator = curator
	orderbook.Mnemonic = mnemonic

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

	var ID strings.Builder

	ID.WriteString(idHashInStringOfCurator)
	ID.WriteString(idHashInStringOfBase)
	ID.WriteString(idHashInStringOfQuote)

	store := ctx.KVStore(k.storeKey)
	bz := store.Get([]byte("order_book_meta"))

	var metaData []meta

	if len(bz) == 0 {
		lastOrderBookIndex = 0
	} else {
		var isSlotEmpty = 0

		k.cdc.MustUnmarshalBinaryBare(bz, &metaData)

		bz := store.Get([]byte("last_order_book_index"))
		k.cdc.MustUnmarshalBinaryBare(bz, &lastOrderBookIndex)

		// Need to get list of all Indices, assuming the list is called listOfIndices
		for indexInListOfIndices, elementInListOfIndices := range metaData {
			if uint32(indexInListOfIndices) != elementInListOfIndices.Index {
				lastOrderBookIndex = uint32(indexInListOfIndices)
				isSlotEmpty = 1
				break
			}
		}

		// It will come to this loop if none of the slots are empty
		if isSlotEmpty != 0 {
			lastOrderBookIndex = uint32(len(metaData)) + 1
		}
		// lastOrderBookIndex = uint32(len(bz)) + 1
	}

	// Hashing and adding the lastOrderBookIndex to the ID
	lenOfLastOrderBookIndex := strconv.Itoa(len(strconv.Itoa(int(lastOrderBookIndex))))
	hashOfLenOfLastOrderBookIndex := blake2b.Sum256([]byte(lenOfLastOrderBookIndex))
	hashInStringOfLenOfLastOrderBookIndexLarge := hex.EncodeToString(hashOfLenOfLastOrderBookIndex[:])
	hashInStringOfLenOfLastOrderBookIndex := hashInStringOfLenOfLastOrderBookIndexLarge[len(hashInStringOfLenOfLastOrderBookIndexLarge) - numberOfCharacters:]

	ID.WriteString(hashInStringOfLenOfLastOrderBookIndex)

	id := ID.String()
	orderbook.ID = id
	orderbook.Index = lastOrderBookIndex

	store.Set([]byte(id), k.cdc.MustMarshalBinaryBare(orderbook))
	store.Set([]byte("last_order_book_index"), k.cdc.MustMarshalBinaryBare(lastOrderBookIndex))

	// To sort metadata
	var newMetaData []meta

	if len(metaData) == 0 {
		newMetaData = append(newMetaData, newMeta(id, lastOrderBookIndex))
	} else {
		var appendedFlag = 0

		for _, elementInListOfIndices := range metaData {
			if lastOrderBookIndex != elementInListOfIndices.Index {
				newMetaData = append(newMetaData, elementInListOfIndices)
			} else {
				appendedFlag = 1

				newMetaData = append(newMetaData, newMeta(id, lastOrderBookIndex))
				newMetaData = append(newMetaData, elementInListOfIndices)
			}
		}

		if appendedFlag == 0 {
			newMetaData = append(newMetaData, newMeta(id, lastOrderBookIndex))
		}
	}

	store.Set([]byte("order_book_meta"), k.cdc.MustMarshalBinaryBare(newMetaData))

}

func (k Keeper) GetOrderBookByID(ctx sdk.Context, id string) []types.OrderBook {

	store := ctx.KVStore(k.storeKey)
	bz := store.Get([]byte(id))

	var orderbook types.OrderBook
	k.cdc.MustUnmarshalBinaryBare(bz, &orderbook)

	var orderbooksQueried = []types.OrderBook{orderbook}
	return orderbooksQueried
}

func (k Keeper) GetOrderBookByIndex(ctx sdk.Context, index uint32) []types.OrderBook {

	store := ctx.KVStore(k.storeKey)

	var orderbook types.OrderBook
	var orderbooksQueried = []types.OrderBook{}
	var metaData []meta

	lenOfLastOrderBookIndex := strconv.Itoa(len(strconv.Itoa(int(index))))
	hashOfLenOfLastOrderBookIndex := blake2b.Sum256([]byte(lenOfLastOrderBookIndex))
	hashInStringOfLenOfLastOrderBookIndexLarge := hex.EncodeToString(hashOfLenOfLastOrderBookIndex[:])
	hashInStringOfLenOfLastOrderBookIndex := hashInStringOfLenOfLastOrderBookIndexLarge[len(hashInStringOfLenOfLastOrderBookIndexLarge) - numberOfCharacters:]


	bz := store.Get([]byte("order_book_meta"))
	k.cdc.MustUnmarshalBinaryBare(bz, &metaData)

	for _, element := range metaData {

		// Matching
		if hashInStringOfLenOfLastOrderBookIndex == element.ID[3 * numberOfCharacters: 4 * numberOfCharacters] {
			bz := store.Get([]byte(element.ID))
			k.cdc.MustUnmarshalBinaryBare(bz, &orderbook)
			orderbooksQueried = append(orderbooksQueried, orderbook)
		}
	}

	return orderbooksQueried
}

func (k Keeper) GetOrderBookByBase(ctx sdk.Context, base string) []types.OrderBook {

	store := ctx.KVStore(k.storeKey)

	var orderbook types.OrderBook
	var orderbooksQueried = []types.OrderBook{}
	var metaData []meta

	hashOfBase := blake2b.Sum256([]byte(base))
	hashInStringOfBase := hex.EncodeToString(hashOfBase[:])
	idHashInStringOfBase := hashInStringOfBase[len(hashInStringOfBase) - numberOfCharacters:]

	bz := store.Get([]byte("order_book_meta"))
	k.cdc.MustUnmarshalBinaryBare(bz, &metaData)

	for _, element := range metaData {

		// Matching
		if idHashInStringOfBase == element.ID[numberOfCharacters: 2 * numberOfCharacters] {
			bz := store.Get([]byte(element.ID))
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
	var metaData []meta

	hashOfQuote := blake2b.Sum256([]byte(quote))
	hashInStringOfQuote := hex.EncodeToString(hashOfQuote[:])
	idHashInStringOfQuote := hashInStringOfQuote[len(hashInStringOfQuote) - numberOfCharacters:]

	bz := store.Get([]byte("order_book_meta"))
	k.cdc.MustUnmarshalBinaryBare(bz, &metaData)

	for _, element := range metaData {

		// Matching
		if idHashInStringOfQuote == element.ID[2 * numberOfCharacters: 3 * numberOfCharacters] {
			bz := store.Get([]byte(element.ID))
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
	var metaData []meta

	hashOfBase := blake2b.Sum256([]byte(base))
	hashInStringOfBase := hex.EncodeToString(hashOfBase[:])
	idHashInStringOfBase := hashInStringOfBase[len(hashInStringOfBase) - numberOfCharacters:]

	hashOfQuote := blake2b.Sum256([]byte(quote))
	hashInStringOfQuote := hex.EncodeToString(hashOfQuote[:])
	idHashInStringOfQuote := hashInStringOfQuote[len(hashInStringOfQuote) - numberOfCharacters:]

	bz := store.Get([]byte("order_book_meta"))
	k.cdc.MustUnmarshalBinaryBare(bz, &metaData)

	for _, element := range metaData {

		// Matching
		if idHashInStringOfBase == element.ID[numberOfCharacters: 2 * numberOfCharacters] &&
			idHashInStringOfQuote == element.ID[2 * numberOfCharacters: 3 * numberOfCharacters] {
			bz := store.Get([]byte(element.ID))
			k.cdc.MustUnmarshalBinaryBare(bz, &orderbook)
			orderbooksQueried = append(orderbooksQueried, orderbook)
		}
	}

	return orderbooksQueried
}

func (k Keeper) GetOrderBookByCurator(ctx sdk.Context, curatorString string) []types.OrderBook {

	store := ctx.KVStore(k.storeKey)

	var orderbook types.OrderBook
	var orderbooksQueried = []types.OrderBook{}
	var metaData []meta

	var curator, _ = sdk.AccAddressFromBech32(curatorString)

	hashOfCurator := blake2b.Sum256(curator)
	hashInStringOfCurator := hex.EncodeToString(hashOfCurator[:])
	idHashInStringOfCurator := hashInStringOfCurator[len(hashInStringOfCurator) - numberOfCharacters:]

	bz := store.Get([]byte("order_book_meta"))
	k.cdc.MustUnmarshalBinaryBare(bz, &metaData)

	for _, element := range metaData {

		// Matching
		if idHashInStringOfCurator == element.ID[0:numberOfCharacters] {
			bz := store.Get([]byte(element.ID))
			k.cdc.MustUnmarshalBinaryBare(bz, &orderbook)
			orderbooksQueried = append(orderbooksQueried, orderbook)
		}
	}

	return orderbooksQueried
}
