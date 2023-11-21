package types

import (
	multistakingtypes "github.com/KiraCore/sekai/x/multistaking/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// combine multiple slashing hooks, all hook functions are run in array sequence
type MultiSlashingHooks []SlashingHooks

func NewMultiSlashingHooks(hooks ...SlashingHooks) MultiSlashingHooks {
	return hooks
}

func (h MultiSlashingHooks) AfterSlashProposalRaise(ctx sdk.Context, valAddr sdk.ValAddress, pool multistakingtypes.StakingPool) {
	for i := range h {
		h[i].AfterSlashProposalRaise(ctx, valAddr, pool)
	}
}
