package keeper

import (
	"github.com/KiraCore/sekai/x/gov/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Keeper struct {
	cdc      codec.BinaryMarshaler
	storeKey sdk.StoreKey
}

func NewKeeper(storeKey sdk.StoreKey, cdc codec.BinaryMarshaler) Keeper {
	return Keeper{cdc: cdc, storeKey: storeKey}
}

// SetNetworkProperties set network properties on KVStore
func (k Keeper) SetNetworkProperties(ctx sdk.Context, properties *types.NetworkProperties) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixNetworkProperties)
	prefixStore.Set([]byte("property"), k.cdc.MustMarshalBinaryBare(properties))
}

// GetNetworkProperties get network properties from KVStore
func (k Keeper) GetNetworkProperties(ctx sdk.Context) *types.NetworkProperties {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixNetworkProperties)
	bz := prefixStore.Get([]byte("property"))

	properties := new(types.NetworkProperties)
	k.cdc.MustUnmarshalBinaryBare(bz, properties)
	return properties
}

// SetExecutionFee set fee by execution function name
func (k Keeper) SetExecutionFee(ctx sdk.Context, fee *types.ExecutionFee) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixExecutionFee)
	key := []byte(fee.TransactionType)
	prefixStore.Set(key, k.cdc.MustMarshalBinaryBare(fee))
}

// GetExecutionFee get fee from execution function name
func (k Keeper) GetExecutionFee(ctx sdk.Context, txType string) *types.ExecutionFee {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixExecutionFee)
	key := []byte(txType)
	if !prefixStore.Has(key) {
		return nil
	}
	bz := prefixStore.Get([]byte(txType))

	fee := new(types.ExecutionFee)
	k.cdc.MustUnmarshalBinaryBare(bz, fee)
	return fee
}
