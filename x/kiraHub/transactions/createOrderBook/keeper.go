package createOrderBook

import (
	"encoding/hex"
	"strings"

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

	// ARJUN CHANGE THIS TO THE DYNAMIC INDEX PULLED FROM THE KVSTORE
	// var last_order_book_index = 1

	// This is the definitions of the lens of the shortened hashes
	var numberOfBytes = 4
	var numberOfCharacters = 2 * numberOfBytes

	// Creating the hashes of the parts of the ID
	hashOfCurator := blake2b.Sum256([]byte(curator))
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
	// Still need to add the functionalities of last_order_book_index

	id := ID.String()

	store := ctx.KVStore(k.storeKey)

	store.Set([]byte(id), k.cdc.MustMarshalBinaryBare(orderbook))
}