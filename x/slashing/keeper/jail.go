package keeper

import (
	"fmt"

	"github.com/KiraCore/sekai/x/slashing/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Jail attempts to jail a validator. The slash is delegated to the staking module
// to make the necessary validator changes.
func (k Keeper) Jail(ctx sdk.Context, consAddr sdk.ConsAddress) {
	validator, err := k.sk.GetValidatorByConsAddr(ctx, consAddr)
	if err == nil && !validator.IsJailed() {
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeJail,
				sdk.NewAttribute(types.AttributeKeyJailed, consAddr.String()),
			),
		)

		k.sk.Jail(ctx, validator.ValKey)

		pool, found := k.msk.GetStakingPoolByValidator(ctx, validator.ValKey.String())
		if found {
			// create a proposal automatically to jail the validator
			content := types.NewSlashValidatorProposal(
				validator.ValKey.String(),
				pool.Id,
				ctx.BlockTime(),
				"double-sign",
				1,
				[]string{},
				"",
			)
			cacheCtx, write := ctx.CacheContext()
			proposalID, err := k.gk.CreateAndSaveProposalWithContent(cacheCtx, "Slash proposal", "Slash for double sign", content)
			if err == nil {
				write()
				fmt.Println("proposal created", proposalID)
			} else {
				fmt.Println("proposal creation error", err)
			}
		}
	}
}

func (k Keeper) SlashStakingPool(ctx sdk.Context, proposal *types.ProposalSlashValidator, slash uint64) {
	k.msk.SlashStakingPool(ctx, proposal.Offender, slash)
}
