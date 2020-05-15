package createOrderBook

import (
	"encoding/hex"
	"strconv"

	sdk "github.com/KiraCore/cosmos-sdk/types"
	"github.com/KiraCore/cosmos-sdk/codec"
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


	var curatorHash = blake2b.Sum256(curator)

	var id = hex.EncodeToString(curatorHash[len(curatorHash) - 4:])


	store := ctx.KVStore(k.storeKey)

	store.Set([]byte(id), k.cdc.MustMarshalBinaryBare(orderbook))
}