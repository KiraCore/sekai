package createOrder

import (
	"encoding/hex"
	"golang.org/x/crypto/blake2b"
	"strconv"
	"strings"
	"time"

	"github.com/KiraCore/cosmos-sdk/codec"
	sdk "github.com/KiraCore/cosmos-sdk/types"
	"github.com/KiraCore/sekai/types"
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

var lastOrderIndex uint32 = 0

// This is the definitions of the lens of the shortened hashes
var numberOfBytes = 4
var numberOfCharacters = 2 * numberOfBytes

func (k Keeper) CreateOrder(ctx sdk.Context, orderBookID string, orderType uint8, amount int64, limitPrice int64) {

	var limitOrder = types.NewLimitOrder()

	limitOrder.OrderBookID = orderBookID
	limitOrder.OrderType = orderType
	limitOrder.Amount = amount
	limitOrder.LimitPrice = limitPrice

	// Expiry Time Logic

	now := time.Now()
	unix := now.Unix()
	limitOrder.ExpiryTime = unix

	// ID Generation Algorithm
	hashOfIndex := blake2b.Sum256([]byte(orderBookID))
	hashInStringOfIndex := hex.EncodeToString(hashOfIndex[:])
	idHashInStringOfIndex := hashInStringOfIndex[len(hashInStringOfIndex) - numberOfCharacters:]

	orderTypeAsString := strconv.Itoa(int(orderType))
	hashOfType := blake2b.Sum256([]byte(orderTypeAsString))
	hashInStringOfType := hex.EncodeToString(hashOfType[:])
	idHashInStringOfType := hashInStringOfType[len(hashInStringOfType) - numberOfCharacters:]

	limitPriceAsString := strconv.Itoa(int(limitPrice))
	hashOfPrice := blake2b.Sum256([]byte(limitPriceAsString))
	hashInStringOfPrice := hex.EncodeToString(hashOfPrice[:])
	idHashInStringOfPrice := hashInStringOfPrice[len(hashInStringOfPrice) - numberOfCharacters:]

	var ID strings.Builder

	ID.WriteString(idHashInStringOfIndex)
	ID.WriteString(idHashInStringOfType)
	ID.WriteString(idHashInStringOfPrice)

	// Storage Logic
	store := ctx.KVStore(k.storeKey)
	bz := store.Get([]byte("order_meta"))

	var metaData []meta

	if len(bz) == 0 {
		lastOrderIndex = 0
	} else {
		var isSlotEmpty = 0

		k.cdc.MustUnmarshalBinaryBare(bz, &metaData)

		bz := store.Get([]byte("last_order_index"))
		k.cdc.MustUnmarshalBinaryBare(bz, &lastOrderIndex)

		// Need to get list of all Indices, assuming the list is called listOfIndices
		for indexInListOfIndices, elementInListOfIndices := range metaData {
			if uint32(indexInListOfIndices) != elementInListOfIndices.Index {
				lastOrderIndex = uint32(indexInListOfIndices)
				isSlotEmpty = 1
				break
			}
		}

		// It will come to this loop if none of the slots are empty
		if isSlotEmpty != 0 {
			lastOrderIndex = uint32(len(metaData)) + 1
		}
	}

	// Hashing and adding the lastOrderBookIndex to the ID
	lenOfLastOrderIndex := strconv.Itoa(len(strconv.Itoa(int(lastOrderIndex))))
	hashOfLenOfLastOrderIndex := blake2b.Sum256([]byte(lenOfLastOrderIndex))
	hashInStringOfLenOfLastOrderIndexLarge := hex.EncodeToString(hashOfLenOfLastOrderIndex[:])
	hashInStringOfLenOfLastOrderIndex := hashInStringOfLenOfLastOrderIndexLarge[len(hashInStringOfLenOfLastOrderIndexLarge) - numberOfCharacters:]

	ID.WriteString(hashInStringOfLenOfLastOrderIndex)

	id := ID.String()
	limitOrder.ID = id
	limitOrder.Index = lastOrderIndex

	store.Set([]byte(id), k.cdc.MustMarshalBinaryBare(limitOrder))
	store.Set([]byte("last_order_book_index"), k.cdc.MustMarshalBinaryBare(lastOrderIndex))

	// To sort metadata
	var newMetaData []meta

	if len(metaData) == 0 {
		newMetaData = append(newMetaData, newMeta(id, lastOrderIndex))
	} else {
		var appendedFlag = 0

		for _, elementInListOfIndices := range metaData {
			if lastOrderIndex != elementInListOfIndices.Index {
				newMetaData = append(newMetaData, elementInListOfIndices)
			} else {
				appendedFlag = 1

				newMetaData = append(newMetaData, newMeta(id, lastOrderIndex))
				newMetaData = append(newMetaData, elementInListOfIndices)
			}
		}

		if appendedFlag == 0 {
			newMetaData = append(newMetaData, newMeta(id, lastOrderIndex))
		}
	}

	store.Set([]byte("order_meta"), k.cdc.MustMarshalBinaryBare(newMetaData))
}