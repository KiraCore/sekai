package keeper

import (
	"fmt"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"

	gogotypes "github.com/gogo/protobuf/types"

	"github.com/KiraCore/sekai/x/slashing/types"
	"github.com/cometbft/cometbft/crypto"
	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Keeper of the slashing store
type Keeper struct {
	storeKey storetypes.StoreKey
	cdc      codec.BinaryCodec
	sk       types.StakingKeeper
	gk       types.GovKeeper
	msk      types.MultiStakingKeeper
	hooks    types.SlashingHooks
}

// NewKeeper creates a slashing keeper
func NewKeeper(cdc codec.BinaryCodec, key storetypes.StoreKey, sk types.StakingKeeper, msk types.MultiStakingKeeper, gk types.GovKeeper) Keeper {
	return Keeper{
		storeKey: key,
		cdc:      cdc,
		sk:       sk,
		msk:      msk,
		gk:       gk,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// AddPubkey sets a address-pubkey relation
func (k Keeper) AddPubkey(ctx sdk.Context, pubkey cryptotypes.PubKey) error {
	bz, err := k.cdc.MarshalInterface(pubkey)
	if err != nil {
		return err
	}
	store := ctx.KVStore(k.storeKey)
	key := types.AddrPubkeyRelationKey(pubkey.Address())
	store.Set(key, bz)
	return nil
}

// GetPubkey returns the pubkey from the adddress-pubkey relation
func (k Keeper) GetPubkey(ctx sdk.Context, a cryptotypes.Address) (cryptotypes.PubKey, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.AddrPubkeyRelationKey(a))
	if bz == nil {
		return nil, fmt.Errorf("address %s not found", sdk.ConsAddress(a))
	}
	var pk cryptotypes.PubKey
	return pk, k.cdc.UnmarshalInterface(bz, &pk)
}

// Inactivate attempts to set validator's status to Inactive from Active.
func (k Keeper) Inactivate(ctx sdk.Context, consAddr sdk.ConsAddress) {
	validator, err := k.sk.GetValidatorByConsAddr(ctx, consAddr)
	if err == nil && validator.IsActive() {
		// only when validator is active, it could move to Inactive status
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeInactivate,
				sdk.NewAttribute(types.AttributeKeyInactivated, consAddr.String()),
			),
		)

		k.sk.Inactivate(ctx, validator.ValKey)
	}
}

func (k Keeper) setAddrPubkeyRelation(ctx sdk.Context, addr crypto.Address, pubkey string) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshal(&gogotypes.StringValue{Value: pubkey})
	store.Set(types.AddrPubkeyRelationKey(addr), bz)
}

func (k Keeper) deleteAddrPubkeyRelation(ctx sdk.Context, addr crypto.Address) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.AddrPubkeyRelationKey(addr))
}

// Set the slashing hooks
func (k *Keeper) SetHooks(sh types.SlashingHooks) *Keeper {
	if k.hooks != nil {
		panic("cannot set slashing hooks twice")
	}

	k.hooks = sh

	return k
}
