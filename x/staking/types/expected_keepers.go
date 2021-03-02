package types

import (
	customgovtypes "github.com/KiraCore/sekai/x/gov/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// StakingHooks event hooks for staking validator object (noalias)
type StakingHooks interface {
	AfterValidatorJoined(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress)  // Must be called when a validator is joined
	AfterValidatorCreated(ctx sdk.Context, valAddr sdk.ValAddress)                           // Must be called when a validator is created
	BeforeValidatorModified(ctx sdk.Context, valAddr sdk.ValAddress)                         // Must be called when a validator's state changes
	AfterValidatorRemoved(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) // Must be called when a validator is deleted
}

// GovKeeper expected governance keeper
type GovKeeper interface {
	// returns network properties
	GetNetworkProperties(sdk.Context) *customgovtypes.NetworkProperties
	// GetNetworkActorsByAbsoluteWhitelistPermission returns all actors that have a specific whitelist permission,
	// it does not matter if it is by role or by individual permission.
	GetNetworkActorsByAbsoluteWhitelistPermission(ctx sdk.Context, perm customgovtypes.PermValue) []customgovtypes.NetworkActor
}
