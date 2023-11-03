package keeper

import (
	"fmt"
	"time"

	kiratypes "github.com/KiraCore/sekai/types"
	govtypes "github.com/KiraCore/sekai/x/gov/types"
	"github.com/KiraCore/sekai/x/slashing/types"
	stakingtypes "github.com/KiraCore/sekai/x/staking/types"
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
		k.SetSlashedValidator(ctx, validator.ValKey)

		pool, found := k.msk.GetStakingPoolByValidator(ctx, validator.ValKey.String())
		if found {
			properties := k.gk.GetNetworkProperties(ctx)
			timestamp := ctx.BlockTime().Add(-time.Second * time.Duration(properties.SlashingPeriod))
			colluderVals := k.GetSlashedValidatorsAfter(ctx, timestamp)

			validators := k.sk.GetValidatorSet(ctx)
			numActiveValidators := 0
			for _, val := range validators {
				if val.Status == stakingtypes.Active {
					numActiveValidators++
				}
			}

			if len(colluderVals) <= int(sdk.NewDec(int64(numActiveValidators)).Mul(properties.MaxJailedPercentage).RoundInt64()) {
				return
			}

			// if same validator slash proposal already exists, do not raise it
			proposals, _ := k.gk.GetProposals(ctx)
			for _, proposal := range proposals {
				if proposal.Result == govtypes.Pending && proposal.GetContent().ProposalType() == kiratypes.ProposalTypeSlashValidator {
					content := proposal.GetContent().(*types.ProposalSlashValidator)
					if content.Offender == validator.ValKey.String() {
						return
					}
				}
			}

			colluders := []string{}
			for _, colluderVal := range colluderVals {
				colluders = append(colluders, colluderVal.String())
			}

			// create a proposal automatically to jail the validator
			content := types.NewSlashValidatorProposal(
				validator.ValKey.String(),
				pool.Id,
				ctx.BlockTime(),
				"double-sign",
				1,
				colluders,
				"",
			)
			cacheCtx, write := ctx.CacheContext()
			proposalID, err := k.gk.CreateAndSaveProposalWithContent(cacheCtx, "Slash proposal", "Slash for double sign", content)
			if err == nil {
				write()
				fmt.Println("proposal created", proposalID)
				if k.hooks != nil {
					k.hooks.AfterSlashProposalRaise(ctx, validator.ValKey, pool)
				}
			} else {
				fmt.Println("proposal creation error", err)
			}
		}
	}
}

func (k Keeper) SlashStakingPool(ctx sdk.Context, proposal *types.ProposalSlashValidator, slash sdk.Dec) {
	k.msk.SlashStakingPool(ctx, proposal.Offender, slash)
}
