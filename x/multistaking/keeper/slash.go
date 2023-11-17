package keeper

import (
	"github.com/KiraCore/sekai/x/multistaking/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

func (k Keeper) SlashStakingPool(ctx sdk.Context, validator string, slash sdk.Dec) {
	pool, found := k.GetStakingPoolByValidator(ctx, validator)
	if !found {
		return
	}
	pool.Slashed = slash
	pool.Enabled = false

	totalStakingTokens := sdk.Coins{}
	for _, stakingToken := range pool.TotalStakingTokens {
		totalStakingTokens = totalStakingTokens.Add(sdk.NewCoin(stakingToken.Denom, sdk.NewDecFromInt(stakingToken.Amount).Mul(sdk.OneDec().Sub(pool.Slashed)).RoundInt()))
	}
	totalSlashedTokens := sdk.Coins(pool.TotalStakingTokens).Sub(totalStakingTokens...)
	pool.TotalStakingTokens = totalStakingTokens

	defaultDenom := k.sk.DefaultDenom(ctx)
	defaultDenomAmount := totalSlashedTokens.AmountOf(defaultDenom)
	burnAmount := sdk.Coins{sdk.NewCoin(defaultDenom, defaultDenomAmount)}
	err := k.bankKeeper.BurnCoins(ctx, types.ModuleName, burnAmount)
	if err != nil {
		panic(err)
	}

	treasurySendAmount := totalSlashedTokens.Sub(burnAmount...)
	err = k.bankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, authtypes.FeeCollectorName, treasurySendAmount)
	if err != nil {
		panic(err)
	}
	feesTreasury := k.distrKeeper.GetFeesTreasury(ctx)
	feesTreasury = feesTreasury.Add(treasurySendAmount...)
	k.distrKeeper.SetFeesTreasury(ctx, feesTreasury)
	k.SetStakingPool(ctx, pool)

	// TODO: pause the basket when the pool's slashed
	// TODO: raise hook for basket module
	valAddr, err := sdk.ValAddressFromBech32(validator)
	if err != nil {
		panic(err)
	}
	if k.hooks != nil {
		k.hooks.AfterSlashStakingPool(ctx, valAddr, pool, slash)
	}
}
