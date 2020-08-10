package kiraHub

import (
	"github.com/KiraCore/cosmos-sdk/codec"
	sdk "github.com/KiraCore/cosmos-sdk/types"
	"github.com/KiraCore/sekai/x/kiraHub/transactions/createOrder"
	"github.com/KiraCore/sekai/x/kiraHub/transactions/createOrderBook"
	signerkey "github.com/KiraCore/sekai/x/kiraHub/transactions/upsertSignerKey"
)

type Keeper interface {
	getCreateOrderBookKeeper() createOrderBook.Keeper
	getCreateOrderKeeper() createOrder.Keeper
	getUpsertSignerKeyKeeper() signerkey.Keeper
}

type baseKeeper struct {
	createOrderBookKeeper createOrderBook.Keeper
	createOrderKeeper     createOrder.Keeper
	upsertSignerKeyKeeper signerkey.Keeper
}

func NewKeeper(codec *codec.Codec, storeKey sdk.StoreKey) Keeper {
	return baseKeeper{
		createOrderBookKeeper: createOrderBook.NewKeeper(codec, storeKey),
		createOrderKeeper:     createOrder.NewKeeper(codec, storeKey),
		upsertSignerKeyKeeper: signerkey.NewKeeper(codec, storeKey),
	}
}

func (baseKeeper baseKeeper) getCreateOrderBookKeeper() createOrderBook.Keeper {
	return baseKeeper.createOrderBookKeeper
}

func (baseKeeper baseKeeper) getCreateOrderKeeper() createOrder.Keeper {
	return baseKeeper.createOrderKeeper
}

func (baseKeeper baseKeeper) getUpsertSignerKeyKeeper() signerkey.Keeper {
	return baseKeeper.upsertSignerKeyKeeper
}
