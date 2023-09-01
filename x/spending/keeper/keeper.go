package keeper

import (
	"github.com/KiraCore/sekai/x/spending/types"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Keeper is for managing token module
type Keeper struct {
	cdc      codec.BinaryCodec
	storeKey storetypes.StoreKey

	bk types.BankKeeper
	gk types.CustomGovKeeper
}

// NewKeeper returns instance of a keeper
func NewKeeper(storeKey storetypes.StoreKey, cdc codec.BinaryCodec, bk types.BankKeeper, gk types.CustomGovKeeper) Keeper {
	return Keeper{
		cdc:      cdc,
		storeKey: storeKey,
		bk:       bk,
		gk:       gk,
	}
}

func (k Keeper) GetBeneficiaryWeight(ctx sdk.Context, address sdk.AccAddress, permInfo types.WeightedPermInfo) uint64 {
	for _, owner := range permInfo.Accounts {
		if owner.Account == address.String() {
			return owner.Weight
		}
	}

	actor, found := k.gk.GetNetworkActorByAddress(ctx, address)
	if !found {
		return 0
	}

	weights := make(map[uint64]uint64)
	for _, role := range permInfo.Roles {
		weights[role.Role] = role.Weight
	}

	for _, role := range actor.Roles {
		if weights[role] != 0 {
			return weights[role]
		}
	}
	return 0
}

func (k Keeper) IsAllowedBeneficiary(ctx sdk.Context, address sdk.AccAddress, permInfo types.WeightedPermInfo) bool {
	for _, owner := range permInfo.Accounts {
		if owner.Account == address.String() {
			return true
		}
	}

	actor, found := k.gk.GetNetworkActorByAddress(ctx, address)
	if !found {
		return false
	}

	flags := make(map[uint64]bool)
	for _, role := range permInfo.Roles {
		flags[role.Role] = true
	}

	for _, role := range actor.Roles {
		if flags[role] {
			return true
		}
	}
	return false
}

func (k Keeper) IsAllowedAddress(ctx sdk.Context, address sdk.AccAddress, permInfo types.PermInfo) bool {
	for _, owner := range permInfo.OwnerAccounts {
		if owner == address.String() {
			return true
		}
	}

	actor, found := k.gk.GetNetworkActorByAddress(ctx, address)
	if !found {
		return false
	}

	flags := make(map[uint64]bool)
	for _, role := range permInfo.OwnerRoles {
		flags[role] = true
	}

	for _, role := range actor.Roles {
		if flags[role] {
			return true
		}
	}
	return false
}

func (k Keeper) AllowedAddresses(ctx sdk.Context, permInfo types.PermInfo) []string {
	addrs := []string{}
	flags := make(map[string]bool)

	for _, owner := range permInfo.OwnerAccounts {
		if flags[owner] == false {
			flags[owner] = true
			addrs = append(addrs, owner)
		}
	}

	for _, role := range permInfo.OwnerRoles {
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
