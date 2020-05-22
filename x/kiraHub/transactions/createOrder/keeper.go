package createOrder

import (
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

	// Storage Logic

}