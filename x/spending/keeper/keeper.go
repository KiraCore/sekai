package keeper

import (
	"github.com/KiraCore/sekai/x/spending/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Keeper is for managing token module
type Keeper struct {
	cdc      codec.BinaryCodec
	storeKey sdk.StoreKey

	bk types.BankKeeper
	gk types.CustomGovKeeper
}

// NewKeeper returns instance of a keeper
func NewKeeper(storeKey sdk.StoreKey, cdc codec.BinaryCodec, bk types.BankKeeper, gk types.CustomGovKeeper) Keeper {
	return Keeper{
		cdc:      cdc,
		storeKey: storeKey,
		bk:       bk,
		gk:       gk,
	}
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

func (k Keeper) DistributePoolRewards(ctx sdk.Context, pool types.SpendingPool) error {
	duplicateMap := map[string]bool{}
	var beneficiaries []string

	for _, acc := range pool.Beneficiaries.OwnerAccounts {
		if _, ok := duplicateMap[acc]; !ok {
			duplicateMap[acc] = true
			beneficiaries = append(beneficiaries, acc)
		}
	}
	for _, role := range pool.Beneficiaries.OwnerRoles {
		actorIter := k.gk.GetNetworkActorsByRole(ctx, role)

		for ; actorIter.Valid(); actorIter.Next() {
			if _, ok := duplicateMap[sdk.AccAddress(actorIter.Value()).String()]; !ok {
				duplicateMap[sdk.AccAddress(actorIter.Value()).String()] = true
				beneficiaries = append(beneficiaries, sdk.AccAddress(actorIter.Value()).String())
			}
		}
	}

	for _, beneficiary := range beneficiaries {
		beneficiaryAcc, err := sdk.AccAddressFromBech32(beneficiary)
		if err != nil {
			return err
		}

		err = k.ClaimSpendingPool(ctx, pool.Name, beneficiaryAcc)
		if err != nil {
			return err
		}
	}
	return nil
}
