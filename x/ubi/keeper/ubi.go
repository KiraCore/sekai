package keeper

import (
	"fmt"
	"strings"

	errorsmod "cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/KiraCore/sekai/x/ubi/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
)

func (k Keeper) GetUBIRecordByName(ctx sdk.Context, name string) *types.UBIRecord {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.PrefixKeyUBIRecord))
	bz := prefixStore.Get([]byte(name))
	if bz == nil {
		return nil
	}

	rate := new(types.UBIRecord)
	k.cdc.MustUnmarshal(bz, rate)

	return rate
}

func (k Keeper) GetUBIRecords(ctx sdk.Context) []types.UBIRecord {
	var records []types.UBIRecord

	// get iterator for token rates
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte(types.PrefixKeyUBIRecord))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		name := strings.TrimPrefix(string(iterator.Key()), string(types.PrefixKeyUBIRecord))
		record := k.GetUBIRecordByName(ctx, name)
		if record != nil {
			records = append(records, *record)
		}
	}
	return records
}

func (k Keeper) SetUBIRecord(ctx sdk.Context, record types.UBIRecord) {
	store := ctx.KVStore(k.storeKey)
	key := append([]byte(types.PrefixKeyUBIRecord), []byte(record.Name)...)
	store.Set(key, k.cdc.MustMarshal(&record))
}

func (k Keeper) DeleteUBIRecord(ctx sdk.Context, name string) error {
	store := ctx.KVStore(k.storeKey)
	key := append([]byte(types.PrefixKeyUBIRecord), []byte(name)...)
	if !store.Has(key) {
		return errorsmod.Wrap(types.ErrUBIRecordDoesNotExists, fmt.Sprintf("ubi record does not exist: %s", name))
	}

	store.Delete(key)
	return nil
}

func (k Keeper) ProcessUBIRecord(ctx sdk.Context, record types.UBIRecord) error {
	if !k.dk.InflationPossible(ctx) {
		return nil
	}
	currUnixTimestamp := uint64(ctx.BlockTime().Unix())
	record.DistributionLast = currUnixTimestamp
	k.SetUBIRecord(ctx, record)

	amount := sdk.NewInt(int64(record.Amount)).Mul(sdk.NewInt(1000_000))

	defaultDenom := k.DefaultDenom(ctx)
	// if dynamic ubi record, mint only missing amount
	if record.Dynamic {
		spendingPool := k.sk.GetSpendingPool(ctx, record.Pool)
		if spendingPool == nil {
			return types.ErrSpendingPoolDoesNotExist
		}

		defaultDenomBalance := sdk.Coins(spendingPool.Balances).AmountOf(defaultDenom)

		if amount.LTE(defaultDenomBalance) {
			return nil
		}
		amount = amount.Sub(defaultDenomBalance)
	}

	coin := sdk.NewCoin(defaultDenom, amount)
	err := k.tk.MintCoins(ctx, minttypes.ModuleName, sdk.NewCoins(coin))
	if err != nil {
		return err
	}

	return k.sk.DepositSpendingPoolFromModule(ctx, minttypes.ModuleName, record.Pool, sdk.Coins{coin})
}
