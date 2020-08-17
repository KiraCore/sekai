package keeper

import (
	"encoding/hex"
	"math/rand"
	"strconv"
	"strings"
	"time"

	sdk "github.com/KiraCore/cosmos-sdk/types"
	"github.com/KiraCore/sekai/x/kiraHub/types"
	"golang.org/x/crypto/blake2b"
)

func (k Keeper) GetOrders(ctx sdk.Context, order_book_id string, maxOrders uint32, minAmount uint32) []types.LimitOrder {

	var metaData []orderMeta
	var queryOutput = []types.LimitOrder{}
	var order types.LimitOrder

	store := ctx.KVStore(k.storeKey)
	bz := store.Get([]byte("limit_order_meta"))

	k.cdc.MustUnmarshalBinaryBare(bz, &metaData)

	for _, elementInListOfIndices := range metaData {
		if elementInListOfIndices.OrderBookID == order_book_id {
			bz := store.Get([]byte(elementInListOfIndices.OrderID))
			k.cdc.MustUnmarshalBinaryBare(bz, &order)

			queryOutput = append(queryOutput, order)
		}
	}

	return queryOutput
}

type orderMeta struct {
	OrderBookID string
	OrderID     string
	Index       uint32
}

func newOrderMeta(orderBookID string, orderID string, index uint32) orderMeta {
	return orderMeta{
		OrderBookID: orderBookID,
		OrderID:     orderID,
		Index:       index,
	}
}

var lastOrderIndex uint32 = 0

func (k Keeper) CreateOrder(ctx sdk.Context, orderBookID string, orderType types.LimitOrderType, amount int64, limitPrice int64, expiryTime int64, curator sdk.AccAddress) {

	//var orderBook = createOrderBook.NewKeeper(k.cdc, k.storeKey).GetOrderBookByID(ctx, orderBookID)

	// Validation Check
	//if string(orderBook[0].Curator) != string(curator) {
	//	return
	//}

	var limitOrder = types.NewLimitOrder()

	limitOrder.OrderBookID = orderBookID
	limitOrder.OrderType = orderType
	limitOrder.Amount = amount
	limitOrder.LimitPrice = limitPrice
	limitOrder.Curator = curator
	limitOrder.ExpiryTime = expiryTime

	// ID Generation Algorithm
	hashOfIndex := blake2b.Sum256([]byte(orderBookID))
	hashInStringOfIndex := hex.EncodeToString(hashOfIndex[:])
	idHashInStringOfIndex := hashInStringOfIndex[len(hashInStringOfIndex)-numberOfCharacters:]

	orderTypeAsString := strconv.Itoa(int(orderType))
	hashOfType := blake2b.Sum256([]byte(orderTypeAsString))
	hashInStringOfType := hex.EncodeToString(hashOfType[:])
	idHashInStringOfType := hashInStringOfType[len(hashInStringOfType)-numberOfCharacters:]

	limitPriceAsString := strconv.Itoa(int(limitPrice))
	hashOfPrice := blake2b.Sum256([]byte(limitPriceAsString))
	hashInStringOfPrice := hex.EncodeToString(hashOfPrice[:])
	idHashInStringOfPrice := hashInStringOfPrice[len(hashInStringOfPrice)-numberOfCharacters:]

	var ID strings.Builder

	ID.WriteString(idHashInStringOfIndex)
	ID.WriteString(idHashInStringOfType)
	ID.WriteString(idHashInStringOfPrice)

	// Storage Logic
	store := ctx.KVStore(k.storeKey)
	bz := store.Get([]byte("limit_order_meta"))

	var metaData []orderMeta

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
	hashInStringOfLenOfLastOrderIndex := hashInStringOfLenOfLastOrderIndexLarge[len(hashInStringOfLenOfLastOrderIndexLarge)-numberOfCharacters:]

	ID.WriteString(hashInStringOfLenOfLastOrderIndex)

	id := ID.String()
	//limitOrder.ID = id
	limitOrder.Index = lastOrderIndex

	store.Set([]byte(id), k.cdc.MustMarshalBinaryBare(limitOrder))
	store.Set([]byte("last_order_index"), k.cdc.MustMarshalBinaryBare(lastOrderIndex))

	// To sort metadata
	var newMetaData []orderMeta

	if len(metaData) == 0 {
		newMetaData = append(newMetaData, newOrderMeta(orderBookID, id, lastOrderIndex))
	} else {
		var appendedFlag = 0

		for _, elementInListOfIndices := range metaData {
			if lastOrderIndex != elementInListOfIndices.Index {
				newMetaData = append(newMetaData, elementInListOfIndices)
			} else {
				appendedFlag = 1

				newMetaData = append(newMetaData, newOrderMeta(orderBookID, id, lastOrderIndex))
				newMetaData = append(newMetaData, elementInListOfIndices)
			}
		}

		if appendedFlag == 0 {
			newMetaData = append(newMetaData, newOrderMeta(id, id, lastOrderIndex))
		}
	}

	store.Set([]byte("limit_order_meta"), k.cdc.MustMarshalBinaryBare(newMetaData))
}

func (k Keeper) cancelOrder(ctx sdk.Context, orderID string) {
	// Load Order
	var order types.LimitOrder

	store := ctx.KVStore(k.storeKey)
	bz := store.Get([]byte(orderID))
	k.cdc.MustUnmarshalBinaryBare(bz, &order)

	// Cancel Order
	order.IsCancelled = true

	// Store Order
	store.Set([]byte(orderID), k.cdc.MustMarshalBinaryBare(order))
}

func (k Keeper) handleOrders(ctx sdk.Context, orderBookID string) {

	// Loading Limit Orders
	var metaData []orderMeta
	var limitBuy []types.LimitOrder
	var limitSell []types.LimitOrder
	var orderBooks []types.OrderBook

	store := ctx.KVStore(k.storeKey)
	bz := store.Get([]byte("limit_order_meta"))
	k.cdc.MustUnmarshalBinaryBare(bz, &metaData)

	bz = store.Get([]byte("limit_buy"))
	k.cdc.MustUnmarshalBinaryBare(bz, &limitBuy)

	bz = store.Get([]byte("limit_sell"))
	k.cdc.MustUnmarshalBinaryBare(bz, &limitSell)

	if limitBuy == nil && limitSell == nil {
		for _, elementInListOfIndices := range metaData {

			var order types.LimitOrder
			var orderBook types.OrderBook

			bz := store.Get([]byte(elementInListOfIndices.OrderID))
			k.cdc.MustUnmarshalBinaryBare(bz, &order)

			if order.OrderType == 1 {
				limitBuy = append(limitBuy, order)
			} else if order.OrderType == 2 {
				limitSell = append(limitSell, order)
			}

			if len(orderBooks) == 0 {
				bz := store.Get([]byte(elementInListOfIndices.OrderBookID))
				k.cdc.MustUnmarshalBinaryBare(bz, &orderBook)

			} else {
				var retrievedFlag = 0

				for _, orderbook := range orderBooks {
					if orderbook.ID == elementInListOfIndices.OrderBookID {
						retrievedFlag = 1
					}
				}

				if retrievedFlag == 0 {
					bz := store.Get([]byte(elementInListOfIndices.OrderBookID))
					k.cdc.MustUnmarshalBinaryBare(bz, &orderBook)
				}
			}

			orderBooks = append(orderBooks, orderBook)
		}
	} else {
		for _, elementInListOfIndices := range metaData {
			var orderBook types.OrderBook

			if len(orderBooks) == 0 {
				bz := store.Get([]byte(elementInListOfIndices.OrderBookID))
				k.cdc.MustUnmarshalBinaryBare(bz, &orderBook)

			} else {
				var retrievedFlag = 0

				for _, orderbook := range orderBooks {
					if orderbook.ID == elementInListOfIndices.OrderBookID {
						retrievedFlag = 1
					}
				}

				if retrievedFlag == 0 {
					bz := store.Get([]byte(elementInListOfIndices.OrderBookID))
					k.cdc.MustUnmarshalBinaryBare(bz, &orderBook)
				}
			}

			orderBooks = append(orderBooks, orderBook)
		}
	}

	// Remove Cancelled & Expired
	for i, elementInListOfIndices := range limitBuy {
		if time.Now().Unix() > elementInListOfIndices.ExpiryTime || elementInListOfIndices.IsCancelled == true {
			limitBuy = append(limitBuy[:i], limitBuy[i+1:]...)
		}
	}

	for i, elementInListOfIndices := range limitSell {
		if time.Now().Unix() > elementInListOfIndices.ExpiryTime || elementInListOfIndices.IsCancelled == true {
			limitSell = append(limitSell[:i], limitSell[i+1:]...)
		}
	}

	// Order By Tx Fee
	limitBuy = mergesort(limitBuy, "limitPrice")
	limitSell = mergesort(limitSell, "limitPrice")

	// Assign ID
	for _, elementInListOfIndices := range metaData {
		for _, buy := range limitBuy {
			if elementInListOfIndices.Index == buy.Index {
				buy.ID = elementInListOfIndices.OrderID
			}
		}

		for _, sell := range limitSell {
			if elementInListOfIndices.Index == sell.Index {
				sell.ID = elementInListOfIndices.OrderID
			}
		}
	}

	var matchBuy []types.LimitOrder
	var matchSell []types.LimitOrder

	// Find orders that increase liquidity
	for index, elementInListOfIndices := range metaData {

		var order types.LimitOrder

		bz := store.Get([]byte(elementInListOfIndices.OrderID))
		k.cdc.MustUnmarshalBinaryBare(bz, &order)

		liquidityAdder(order, matchBuy, matchSell, limitBuy, limitSell, index, 0)
	}

	// Generate Seed
	blockHeader := ctx.BlockHeader().LastBlockId.Hash
	blockIDHex := hex.EncodeToString(blockHeader[:])
	blockIDInt, _ := strconv.Atoi(blockIDHex[:])

	rand.Seed(int64(blockIDInt))

	// Randomize Orders
	newBuy := fisheryatesShuffle(matchBuy)
	newSell := fisheryatesShuffle(matchSell)

	// Pick Orders

	terminate := 0

	for terminate != 1 {

		var buy types.LimitOrder
		buyF := 0

		var sell types.LimitOrder
		sellF := 0

		if len(newBuy) > 1 {
			buy = newSell[0]
			buyF = 1
			if len(newBuy) > 1 {
				newBuy = newBuy[1:]
			}
		} else {
			buyF = 0
		}

		if len(newSell) >= 1 {
			sell = newSell[0]
			sellF = 1
			if len(newSell) > 1 {
				newSell = newSell[1:]
			}
		} else {
			sellF = 0
		}

		if buyF == 0 && sellF == 0 {
			terminate = 1
		}

		// New Orders Matched
		if buyF == 1 && sellF == 1 && buy.LimitPrice > sell.LimitPrice {
			if buy.OrderBookID == sell.OrderBookID {

				if buy.Amount > sell.Amount {
					buy.Amount -= sell.Amount
					matchPayout(sell.Curator, buy.Curator, buy.LimitPrice, sell.Amount)
				} else if buy.Amount < sell.Amount {
					sell.Amount -= buy.Amount
					matchPayout(sell.Curator, buy.Curator, buy.LimitPrice, buy.Amount)
				}

			} else {
				var buyOrderBook = k.GetOrderBookByID(ctx, buy.OrderBookID)
				var sellOrderBook = k.GetOrderBookByID(ctx, sell.OrderBookID)

				if buyOrderBook[0].Base == sellOrderBook[0].Base && buyOrderBook[0].Quote == sellOrderBook[0].Quote {
					// Matching
					if buy.Amount > sell.Amount {
						buy.Amount -= sell.Amount
						matchPayout(sell.Curator, buy.Curator, buy.LimitPrice, sell.Amount)
					} else if buy.Amount < sell.Amount {
						sell.Amount -= buy.Amount
						matchPayout(sell.Curator, buy.Curator, buy.LimitPrice, buy.Amount)
					}
				}
			}
		} else {

			// Match With State
			if sellF == 1 {
				for index, stateBuy := range limitBuy {

					if stateBuy.LimitPrice > sell.LimitPrice {
						if stateBuy.OrderBookID == sell.OrderBookID {

							// Order Matched
							if stateBuy.Amount > sell.Amount {
								stateBuy.Amount -= sell.Amount
								matchPayout(sell.Curator, stateBuy.Curator, stateBuy.LimitPrice, sell.Amount)

								break
							} else if stateBuy.Amount < sell.Amount {
								sell.Amount -= stateBuy.Amount
								matchPayout(sell.Curator, stateBuy.Curator, stateBuy.LimitPrice, stateBuy.Amount)

								limitBuy = append(limitBuy[:index], limitBuy[index+1:]...)

								continue
							}

						} else {
							var buyOrderBook = k.GetOrderBookByID(ctx, buy.OrderBookID)
							var sellOrderBook = k.GetOrderBookByID(ctx, sell.OrderBookID)

							if buyOrderBook[0].Base == sellOrderBook[0].Base && buyOrderBook[0].Quote == sellOrderBook[0].Quote {

								// Order Matched
								if stateBuy.Amount > sell.Amount {
									stateBuy.Amount -= sell.Amount
									matchPayout(sell.Curator, stateBuy.Curator, stateBuy.LimitPrice, sell.Amount)

									break
								} else if stateBuy.Amount < sell.Amount {
									sell.Amount -= stateBuy.Amount
									matchPayout(sell.Curator, stateBuy.Curator, stateBuy.LimitPrice, stateBuy.Amount)

									limitBuy = append(limitBuy[:index], limitBuy[index+1:]...)

									continue
								}
							}
						}
					}
				}
			}

			if buyF == 1 {
				for index, stateSell := range limitSell {

					if stateSell.LimitPrice < buy.LimitPrice {
						if stateSell.OrderBookID == buy.OrderBookID {

							// Order Matched
							if stateSell.Amount > buy.Amount {
								stateSell.Amount -= buy.Amount
								matchPayout(stateSell.Curator, buy.Curator, buy.LimitPrice, buy.Amount)

								break
							} else if stateSell.Amount < buy.Amount {
								buy.Amount -= stateSell.Amount
								matchPayout(stateSell.Curator, buy.Curator, buy.LimitPrice, stateSell.Amount)

								limitSell = append(limitSell[:index], limitSell[index+1:]...)

								continue
							}

						} else {
							var buyOrderBook = k.GetOrderBookByID(ctx, buy.OrderBookID)
							var sellOrderBook = k.GetOrderBookByID(ctx, sell.OrderBookID)

							if buyOrderBook[0].Base == sellOrderBook[0].Base && buyOrderBook[0].Quote == sellOrderBook[0].Quote {

								// Order Matched
								if stateSell.Amount > buy.Amount {
									stateSell.Amount -= buy.Amount
									matchPayout(stateSell.Curator, buy.Curator, buy.LimitPrice, buy.Amount)

									break
								} else if stateSell.Amount < buy.Amount {
									buy.Amount -= stateSell.Amount
									matchPayout(stateSell.Curator, buy.Curator, buy.LimitPrice, stateSell.Amount)

									limitSell = append(limitSell[:index], limitSell[index+1:]...)

									continue
								}
							}
						}
					}
				}
			}
		}

		if terminate != 1 {
			for index, buy := range newBuy {
				liquidityAdder(buy, newBuy, newSell, limitBuy, limitSell, index, 1)
			}

			for index, sell := range newSell {
				liquidityAdder(sell, newBuy, newSell, limitBuy, limitSell, index, 1)
			}
		}
	}

	// Persist the limitBuy and limitSell
	store.Set([]byte("limit_buy"), k.cdc.MustMarshalBinaryBare(limitBuy))
	store.Set([]byte("limit_sell"), k.cdc.MustMarshalBinaryBare(limitSell))
}

func liquidityAdder(order types.LimitOrder, matchBuy []types.LimitOrder, matchSell []types.LimitOrder, limitBuy []types.LimitOrder, limitSell []types.LimitOrder, index int, funcType int) {
	if order.OrderType == 1 {
		if order.LimitPrice < findMin(limitSell) {
			if funcType == 0 {
				matchBuy = append(matchBuy, order)
			}
		} else {
			limitBuy = append(limitBuy, order)
			if funcType == 1 {
				matchBuy = append(matchBuy[:index], matchBuy[index+1:]...)
			}
		}
	} else if order.OrderType == 2 {
		if order.LimitPrice > findMax(limitBuy) {
			if funcType == 0 {
				matchSell = append(matchSell, order)
			}
		} else {
			limitSell = append(limitSell, order)
			if funcType == 1 {
				matchSell = append(matchSell[:index], matchSell[index+1:]...)
			}
		}
	}

	return
}

func merge(orderList []types.LimitOrder, middle int, sortBy string) {
	var helper = orderList

	helperLeft := 0
	helperRight := middle
	current := 0
	high := len(orderList) - 1

	switch sortBy {
	case "limitPrice":
		for helperLeft <= middle-1 && helperRight <= high {
			if helper[helperLeft].LimitPrice <= helper[helperRight].LimitPrice {
				orderList[current] = helper[helperLeft]
				helperLeft++
			} else {
				orderList[current] = helper[helperRight]
				helperRight++
			}
			current++
		}
	case "index":
		for helperLeft <= middle-1 && helperRight <= high {
			if helper[helperLeft].Index <= helper[helperRight].Index {
				orderList[current] = helper[helperLeft]
				helperLeft++
			} else {
				orderList[current] = helper[helperRight]
				helperRight++
			}
			current++
		}
	}

	for helperLeft <= middle-1 {
		orderList[current] = helper[helperLeft]
		current++
		helperLeft++
	}
}

func mergesort(orderList []types.LimitOrder, sortBy string) []types.LimitOrder {
	if len(orderList) > 1 {
		middle := len(orderList) / 2
		mergesort(orderList[:middle], sortBy)
		mergesort(orderList[middle:], sortBy)
		merge(orderList, middle, sortBy)
	}

	return orderList
}

// Use this instead of mergeSort when Orders exceed a million in number
//func parallelMergeSort(s []int) []int {
//	len := len(s)
//
//	if len > 1 {
//		middle := len / 2
//
//		var wg sync.WaitGroup
//		wg.Add(2)
//
//		go func() {
//			defer wg.Done()
//			parallelMergeSort(s[:middle])
//		}()
//
//		go func() {
//			defer wg.Done()
//			parallelMergeSort(s[middle:])
//		}()
//
//		wg.Wait()
//		parallelMerge(s, middle)
//	}
//
//	return s
//}

func findMax(list []types.LimitOrder) int64 {
	var max int64 = list[0].LimitPrice
	for _, value := range list {
		if max < value.LimitPrice {
			max = value.LimitPrice
		}
	}
	return max
}

func findMin(list []types.LimitOrder) int64 {
	var min int64 = list[0].LimitPrice
	for _, value := range list {
		if min > value.LimitPrice {
			min = value.LimitPrice
		}
	}
	return min
}

func fisheryatesShuffle(list []types.LimitOrder) []types.LimitOrder {
	N := len(list)

	for i := 0; i < N; i++ {
		// choose index uniformly in [i, N-1]
		r := i + rand.Intn(N-i)
		list[r], list[i] = list[i], list[r]
	}

	return list
}

func matchPayout(seller, buyer sdk.AccAddress, price, amount int64) {
	// Seller needs to pay buyer the price*amount as a tx

	// Return error if buyer doesnt have enough funds or if theres any other problem in the sellers end
}
