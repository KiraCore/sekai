package keeper

import (
	"time"

	"github.com/KiraCore/sekai/x/slashing/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ResetWholeValidatorRank whole validator rank
func (k Keeper) ResetWholeValidatorRank(ctx sdk.Context) error {

	k.IterateValidatorSigningInfos(ctx, func(address sdk.ConsAddress, info types.ValidatorSigningInfo) (stop bool) {
		info.StartHeight = ctx.BlockHeight()
		info.InactiveUntil = time.Unix(0, 0)
		info.MischanceConfidence = 0
		info.Mischance = 0
		info.MissedBlocksCounter = 0
		info.ProducedBlocksCounter = 0

		k.SetValidatorSigningInfo(ctx, address, info)
		return false
	})

	k.sk.ResetWholeValidatorRank(ctx)

	return nil
}
