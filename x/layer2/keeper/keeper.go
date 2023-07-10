package keeper

import (
	appparams "github.com/KiraCore/sekai/app/params"
	govkeeper "github.com/KiraCore/sekai/x/gov/keeper"
	govtypes "github.com/KiraCore/sekai/x/gov/types"
	"github.com/KiraCore/sekai/x/layer2/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Keeper struct {
	cdc      codec.BinaryCodec
	storeKey sdk.StoreKey
	bk       types.BankKeeper
	sk       types.StakingKeeper
	gk       govkeeper.Keeper
	spk      types.SpendingKeeper
}

func NewKeeper(storeKey sdk.StoreKey, cdc codec.BinaryCodec, bk types.BankKeeper, sk types.StakingKeeper, gk govkeeper.Keeper, spk types.SpendingKeeper) Keeper {
	return Keeper{
		cdc:      cdc,
		storeKey: storeKey,
		bk:       bk,
		sk:       sk,
		gk:       gk,
		spk:      spk,
	}
}

// BondDenom returns the denom that is basically used for fee payment
func (k Keeper) BondDenom(ctx sdk.Context) string {
	return appparams.BondDenom
}

func (k Keeper) CheckIfAllowedPermission(ctx sdk.Context, addr sdk.AccAddress, permValue govtypes.PermValue) bool {
	return govkeeper.CheckIfAllowedPermission(ctx, k.gk, addr, govtypes.PermHandleBasketEmergency)
}

func (k Keeper) IsAllowedAddress(ctx sdk.Context, address sdk.AccAddress, permInfo types.Controllers) bool {
	for _, owner := range permInfo.Whitelist.Addresses {
		if owner == address.String() {
			return true
		}
	}

	actor, found := k.gk.GetNetworkActorByAddress(ctx, address)
	if !found {
		return false
	}

	flags := make(map[uint64]bool)
	for _, role := range permInfo.Whitelist.Roles {
		flags[role] = true
	}

	for _, role := range actor.Roles {
		if flags[role] {
			return true
		}
	}
	return false
}

func (k Keeper) AllowedAddresses(ctx sdk.Context, permInfo types.Controllers) []string {
	addrs := []string{}
	flags := make(map[string]bool)

	for _, owner := range permInfo.Whitelist.Addresses {
		if flags[owner] == false {
			flags[owner] = true
			addrs = append(addrs, owner)
		}
	}

	for _, role := range permInfo.Whitelist.Roles {
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
