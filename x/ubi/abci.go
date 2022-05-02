package ubi

import (
	"github.com/KiraCore/sekai/x/ubi/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func EndBlocker(ctx sdk.Context, k keeper.Keeper) {
	// Mint tokens and call spend pool deposit on Endblocker
	allRecords := k.GetUBIRecords(ctx)

	// TODO: should consider InflationRate and InflationPeriod on governance params
	for _, record := range allRecords {
		currUnixTimestamp := uint64(ctx.BlockTime().Unix())
		if currUnixTimestamp > record.DistributionLast+record.Period && (record.DistributionEnd == 0 || record.DistributionLast < record.DistributionEnd) {
			cacheCtx, write := ctx.CacheContext()
			err := k.ProcessUBIRecord(cacheCtx, record)
			if err == nil {
				write()
			}
		}
	}
}
