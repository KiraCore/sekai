package keeper

import (
	"github.com/KiraCore/sekai/x/upgrade/types"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
)

// Keeper is for managing upgrade module
type Keeper struct {
	cdc             codec.BinaryCodec
	storeKey        storetypes.StoreKey
	upgradeHandlers map[string]types.UpgradeHandler
	sk              types.CustomStakingKeeper
}

// NewKeeper constructs an upgrade Keeper
func NewKeeper(storeKey storetypes.StoreKey, cdc codec.BinaryCodec, sk types.CustomStakingKeeper) Keeper {
	return Keeper{
		storeKey:        storeKey,
		cdc:             cdc,
		upgradeHandlers: map[string]types.UpgradeHandler{},
		sk:              sk,
	}
}

// SetUpgradeHandler sets an UpgradeHandler for the upgrade specified by name. This handler will be called when the upgrade
// with this name is applied. In order for an upgrade with the given name to proceed, a handler for this upgrade
// must be set even if it is a no-op function.
func (k Keeper) SetUpgradeHandler(name string, upgradeHandler types.UpgradeHandler) {
	k.upgradeHandlers[name] = upgradeHandler
}
