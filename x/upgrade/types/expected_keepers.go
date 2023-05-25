package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type CustomStakingKeeper interface {
	PauseProposalNotApprovedValidators(ctx sdk.Context, proposalID uint64) error
}
