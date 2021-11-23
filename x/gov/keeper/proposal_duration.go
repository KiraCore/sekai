package keeper

import sdk "github.com/cosmos/cosmos-sdk/types"

func (k Keeper) SetProposalDuration(ctx sdk.Context, proposalType string, duration uint64) {
	// TODO: implement
	// KeyPrefixProposalDuration
}

func (k Keeper) GetProposalDuration(ctx sdk.Context, proposalType string) uint64 {
	// TODO: implement
	// KeyPrefixProposalDuration
	return 0
}

func (k Keeper) GetProposalDurations(ctx sdk.Context) map[string]uint64 {
	return make(map[string]uint64)
}
