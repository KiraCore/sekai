package keeper

import (
	"github.com/KiraCore/sekai/x/collectives/types"
	govkeeper "github.com/KiraCore/sekai/x/gov/keeper"
	govtypes "github.com/KiraCore/sekai/x/gov/types"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Keeper is for managing token module
type Keeper struct {
	cdc      codec.BinaryCodec
	storeKey storetypes.StoreKey
	bk       types.BankKeeper
	gk       govkeeper.Keeper
	mk       types.MultiStakingKeeper
	tk       types.TokensKeeper
	spk      types.SpendingKeeper
}

// NewKeeper returns instance of a keeper
func NewKeeper(storeKey storetypes.StoreKey, cdc codec.BinaryCodec, bk types.BankKeeper, gk govkeeper.Keeper, mk types.MultiStakingKeeper, tk types.TokensKeeper, spk types.SpendingKeeper) Keeper {
	return Keeper{
		cdc:      cdc,
		storeKey: storeKey,
		bk:       bk,
		gk:       gk,
		mk:       mk,
		tk:       tk,
		spk:      spk,
	}
}

func (k Keeper) CheckIfAllowedPermission(ctx sdk.Context, addr sdk.AccAddress, permValue govtypes.PermValue) bool {
	return govkeeper.CheckIfAllowedPermission(ctx, k.gk, addr, govtypes.PermHandleBasketEmergency)
}

func calcPortion(coins sdk.Coins, portion sdk.Dec) sdk.Coins {
	portionCoins := sdk.Coins{}
	for _, coin := range coins {
		portionCoin := sdk.NewCoin(coin.Denom, sdk.NewDecFromInt(coin.Amount).Mul(portion).RoundInt())
		portionCoins = portionCoins.Add(portionCoin)
	}
	return portionCoins
}

func (k Keeper) IsAllowedAddress(ctx sdk.Context, address sdk.AccAddress, permInfo types.OwnersWhitelist) bool {
	for _, owner := range permInfo.Accounts {
		if owner == address.String() {
			return true
		}
	}

	actor, found := k.gk.GetNetworkActorByAddress(ctx, address)
	if !found {
		return false
	}

	flags := make(map[uint64]bool)
	for _, role := range permInfo.Roles {
		flags[role] = true
	}

	for _, role := range actor.Roles {
		if flags[role] {
			return true
		}
	}
	return false
}

func (k Keeper) AllowedAddresses(ctx sdk.Context, permInfo types.OwnersWhitelist) []string {
	addrs := []string{}
	flags := make(map[string]bool)

	for _, owner := range permInfo.Accounts {
		if flags[owner] == false {
			flags[owner] = true
			addrs = append(addrs, owner)
		}
	}

	for _, role := range permInfo.Roles {
		actorIter := k.gk.GetNetworkActorsByRole(ctx, role)

		for ; actorIter.Valid(); actorIter.Next() {
			addr := sdk.AccAddress(actorIter.Value()).String()
			if flags[addr] == false {
				flags[addr] = true
				addrs = append(addrs, addr)
			}
		}
	}

	return addrs
}
