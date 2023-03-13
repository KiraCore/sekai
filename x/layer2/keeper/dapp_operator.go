package keeper

import (
	"github.com/KiraCore/sekai/x/layer2/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) SetDappOperator(ctx sdk.Context, operator types.DappOperator) {
	bz := k.cdc.MustMarshal(&operator)
	store := ctx.KVStore(k.storeKey)
	store.Set(types.DappOperatorKey(operator.DappName, operator.Operator), bz)
}

func (k Keeper) DeleteDappOperator(ctx sdk.Context, name, address string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.DappOperatorKey(name, address))
}

func (k Keeper) GetDappOperator(ctx sdk.Context, name string, user string) types.DappOperator {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.DappOperatorKey(name, user))
	if bz == nil {
		return types.DappOperator{}
	}

	operator := types.DappOperator{}
	k.cdc.MustUnmarshal(bz, &operator)
	return operator
}

func (k Keeper) GetDappOperators(ctx sdk.Context, name string) []types.DappOperator {
	store := ctx.KVStore(k.storeKey)

	operators := []types.DappOperator{}
	it := sdk.KVStorePrefixIterator(store, append([]byte(types.PrefixDappOperatorKey), name...))
	defer it.Close()

	for ; it.Valid(); it.Next() {
		operator := types.DappOperator{}
		k.cdc.MustUnmarshal(it.Value(), &operator)
		operators = append(operators, operator)
	}
	return operators
}

func (k Keeper) GetAllDappOperators(ctx sdk.Context) []types.DappOperator {
	store := ctx.KVStore(k.storeKey)

	operators := []types.DappOperator{}
	it := sdk.KVStorePrefixIterator(store, []byte(types.PrefixDappOperatorKey))
	defer it.Close()

	for ; it.Valid(); it.Next() {
		operator := types.DappOperator{}
		k.cdc.MustUnmarshal(it.Value(), &operator)
		operators = append(operators, operator)
	}
	return operators
}
