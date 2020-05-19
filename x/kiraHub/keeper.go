package kiraHub

import (
	"github.com/KiraCore/cosmos-sdk/codec"
	sdk "github.com/KiraCore/cosmos-sdk/types"
	"github.com/KiraCore/sekai/x/kiraHub/transactions/createOrderBook"
)

type Keeper interface {
	getCreateOrderBookKeeper() createOrderBook.Keeper
}

type baseKeeper struct {
	createOrderBookKeeper   createOrderBook.Keeper
}

func NewKeeper(codec *codec.Codec, storeKey sdk.StoreKey) Keeper {
	return baseKeeper{
		createOrderBookKeeper:   createOrderBook.NewKeeper(codec, storeKey),
	}
}

func (baseKeeper baseKeeper) getCreateOrderBookKeeper() createOrderBook.Keeper { return baseKeeper.createOrderBookKeeper }






