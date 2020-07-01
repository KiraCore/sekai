package kiraHub

import (
	"github.com/KiraCore/cosmos-sdk/codec"
	sdk "github.com/KiraCore/cosmos-sdk/types"
	"github.com/KiraCore/sekai/x/kiraHub/transactions/createOrderBook"
	"github.com/KiraCore/sekai/x/kiraHub/transactions/createOrder"
	)

type Keeper interface {
	getCreateOrderBookKeeper() createOrderBook.Keeper
	getCreateOrderKeeper()     createOrder.Keeper
}

type baseKeeper struct {
	createOrderBookKeeper   createOrderBook.Keeper
	createOrderKeeper 		createOrder.Keeper
}

func NewKeeper(codec *codec.Codec, storeKey sdk.StoreKey) Keeper {
	return baseKeeper{
		createOrderBookKeeper:   createOrderBook.NewKeeper(codec, storeKey),
		createOrderKeeper: 		 createOrder.NewKeeper(codec, storeKey),
	}
}

func (baseKeeper baseKeeper) getCreateOrderBookKeeper() createOrderBook.Keeper { return baseKeeper.createOrderBookKeeper }

func (baseKeeper baseKeeper) getCreateOrderKeeper() createOrder.Keeper { return baseKeeper.createOrderKeeper }






