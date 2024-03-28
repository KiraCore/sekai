package keeper

import (
	errorsmod "cosmossdk.io/errors"
	"github.com/KiraCore/sekai/x/recovery/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) GetRecoveryRecord(ctx sdk.Context, address string) (types.RecoveryRecord, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.RecoveryRecordKey(address))
	if bz == nil {
		return types.RecoveryRecord{}, errorsmod.Wrapf(types.ErrRecoveryRecordDoesNotExist, "RecoveryRecord: %s does not exist", address)
	}
	record := types.RecoveryRecord{}
	k.cdc.MustUnmarshal(bz, &record)
	return record, nil
}

func (k Keeper) GetRecoveryAddressFromChallenge(ctx sdk.Context, challenge string) string {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.RecoveryChallengeKey(challenge))
	if bz == nil {
		return ""
	}
	return string(bz)
}

func (k Keeper) SetRecoveryRecord(ctx sdk.Context, record types.RecoveryRecord) {
	bz := k.cdc.MustMarshal(&record)
	store := ctx.KVStore(k.storeKey)
	store.Set(types.RecoveryRecordKey(record.Address), bz)
	store.Set(types.RecoveryChallengeKey(record.Challenge), []byte(record.Address))
}

func (k Keeper) DeleteRecoveryRecord(ctx sdk.Context, record types.RecoveryRecord) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.RecoveryRecordKey(record.Address))
	store.Delete(types.RecoveryChallengeKey(record.Challenge))
}

func (k Keeper) GetAllRecoveryRecords(ctx sdk.Context) []types.RecoveryRecord {
	store := ctx.KVStore(k.storeKey)

	records := []types.RecoveryRecord{}
	it := sdk.KVStorePrefixIterator(store, types.RecoveryRecordKeyPrefix)
	defer it.Close()

	for ; it.Valid(); it.Next() {
		record := types.RecoveryRecord{}
		k.cdc.MustUnmarshal(it.Value(), &record)
		records = append(records, record)
	}
	return records
}

func (k Keeper) GetRecoveryToken(ctx sdk.Context, address string) (types.RecoveryToken, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.RecoveryTokenKey(address))
	if bz == nil {
		return types.RecoveryToken{}, errorsmod.Wrapf(types.ErrRecoveryTokenDoesNotExist, "RecoveryToken: %s does not exist", address)
	}
	recovery := types.RecoveryToken{}
	k.cdc.MustUnmarshal(bz, &recovery)
	return recovery, nil
}

func (k Keeper) GetRecoveryTokenByDenom(ctx sdk.Context, denom string) (types.RecoveryToken, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.RecoveryTokenByDenomKey(denom))
	if bz == nil {
		return types.RecoveryToken{}, errorsmod.Wrapf(types.ErrRecoveryTokenDoesNotExist, "RecoveryTokenByDenom: %s does not exist", denom)
	}
	address := string(bz)
	return k.GetRecoveryToken(ctx, address)
}

func (k Keeper) SetRecoveryToken(ctx sdk.Context, recovery types.RecoveryToken) {
	bz := k.cdc.MustMarshal(&recovery)
	store := ctx.KVStore(k.storeKey)
	store.Set(types.RecoveryTokenKey(recovery.Address), bz)
	store.Set(types.RecoveryTokenByDenomKey(recovery.Token), []byte(recovery.Address))
}

func (k Keeper) DeleteRecoveryToken(ctx sdk.Context, recovery types.RecoveryToken) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.RecoveryTokenKey(recovery.Address))
	store.Delete(types.RecoveryTokenByDenomKey(recovery.Token))
}

func calcPortion(coins sdk.Coins, portion sdk.Int, supply sdk.Int) sdk.Coins {
	portionCoins := sdk.Coins{}
	for _, coin := range coins {
		portionCoin := sdk.NewCoin(coin.Denom, coin.Amount.Mul(portion).Quo(supply))
		portionCoins = portionCoins.Add(portionCoin)
	}
	return portionCoins
}

func (k Keeper) IncreaseRecoveryTokenUnderlying(ctx sdk.Context, addr sdk.AccAddress, amount sdk.Coins) error {
	recoveryToken, err := k.GetRecoveryToken(ctx, addr.String())
	if err != nil {
		return err
	}

	k.UnregisterNotEnoughAmountHolder(ctx, recoveryToken.Token)

	supply := k.bk.GetSupply(ctx, recoveryToken.Token).Amount

	holders := k.GetRRTokenHolders(ctx, recoveryToken.Token)
	totalAllocation := sdk.NewCoins()
	for _, holder := range holders {
		balances := k.bk.GetAllBalances(ctx, holder)
		balance := balances.AmountOf(recoveryToken.Token)
		allocation := calcPortion(amount, balance, supply)
		totalAllocation = totalAllocation.Add(allocation...)
		k.IncreaseRRTokenHolderRewards(ctx, holder, allocation)
	}

	unallocated := amount.Sub(totalAllocation...)
	recoveryToken.UnderlyingTokens = sdk.Coins(recoveryToken.UnderlyingTokens).Add(unallocated...)
	k.SetRecoveryToken(ctx, recoveryToken)
	return nil
}

func (k Keeper) GetAllRecoveryTokens(ctx sdk.Context) []types.RecoveryToken {
	store := ctx.KVStore(k.storeKey)

	recoveries := []types.RecoveryToken{}
	it := sdk.KVStorePrefixIterator(store, types.RecoveryTokenKeyPrefix)
	defer it.Close()

	for ; it.Valid(); it.Next() {
		recovery := types.RecoveryToken{}
		k.cdc.MustUnmarshal(it.Value(), &recovery)
		recoveries = append(recoveries, recovery)
	}
	return recoveries
}
