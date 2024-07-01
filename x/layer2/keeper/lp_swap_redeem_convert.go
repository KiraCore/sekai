package keeper

import (
	"github.com/KiraCore/sekai/x/layer2/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) LpTokenPrice(ctx sdk.Context, dapp types.Dapp) sdk.Dec {
	lpToken := dapp.LpToken()
	lpSupply := k.bk.GetSupply(ctx, lpToken).Amount
	totalBond := dapp.TotalBond.Amount
	if lpSupply.IsZero() {
		return sdk.ZeroDec()
	}
	return sdk.NewDecFromInt(totalBond).Quo(sdk.NewDecFromInt(lpSupply))
}

func (k Keeper) OnCollectFee(ctx sdk.Context, fee sdk.Coins) error {
	// TODO: The fixed fee will be applied after the swap from where
	// - `50%` of the corresponding tokens must be **burned** (deminted)
	// - `25%` given as a reward to liquidity providers
	// - `25%` will be split between **ACTIVE** dApp executors, and verifiers (fisherman).
	// Additionally, the premint and postmint tokens can be used to incentivize operators before
	// dApp starts to generate revenue.
	err := k.tk.BurnCoins(ctx, types.ModuleName, fee)
	if err != nil {
		return err
	}
	return nil
}

func (k Keeper) RedeemDappPoolTx(ctx sdk.Context, addr sdk.AccAddress, dapp types.Dapp, poolFee sdk.Dec, lpTokenAmount sdk.Coin) (sdk.Coin, error) {
	// 	totalBond * lpSupply = (totalBond - swapBond) * (lpSupply + swapLpAmount)
	lpToken := dapp.LpToken()
	if lpToken != lpTokenAmount.Denom {
		return sdk.Coin{}, types.ErrInvalidLpToken
	}
	lpSupply := k.bk.GetSupply(ctx, lpToken).Amount
	totalBond := dapp.TotalBond.Amount
	swapLpAmount := lpTokenAmount.Amount
	totalBondAfterSwap := totalBond.Mul(lpSupply).Quo(lpSupply.Add(swapLpAmount))
	swapBond := totalBond.Sub(totalBondAfterSwap)

	dapp.TotalBond.Amount = totalBondAfterSwap
	k.SetDapp(ctx, dapp)

	fee := sdk.NewDecFromInt(swapBond).Mul(poolFee).RoundInt()
	if fee.IsPositive() {
		feeCoin := sdk.NewCoin(dapp.TotalBond.Denom, fee)
		err := k.OnCollectFee(ctx, sdk.Coins{feeCoin})
		if err != nil {
			return sdk.Coin{}, err
		}
	}

	// send lp tokens to the module account
	err := k.bk.SendCoinsFromAccountToModule(ctx, addr, types.ModuleName, sdk.Coins{lpTokenAmount})
	if err != nil {
		return sdk.Coin{}, err
	}

	// send tokens to user
	userReceiveCoin := sdk.NewCoin(dapp.TotalBond.Denom, swapBond.Sub(fee))
	err = k.bk.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr, sdk.Coins{userReceiveCoin})
	if err != nil {
		return sdk.Coin{}, err
	}

	properties := k.gk.GetNetworkProperties(ctx)
	threshold := sdk.NewInt(int64(properties.DappLiquidationThreshold)).Mul(sdk.NewInt(1000_000))
	if dapp.LiquidationStart == 0 && dapp.TotalBond.Amount.LT(threshold) {
		dapp.LiquidationStart = uint64(ctx.BlockTime().Unix())
		k.SetDapp(ctx, dapp)
	}

	return userReceiveCoin, nil
}

func (k Keeper) SwapDappPoolTx(ctx sdk.Context, addr sdk.AccAddress, dapp types.Dapp, poolFee sdk.Dec, swapBond sdk.Coin) (sdk.Coin, error) {
	// 	totalBond * lpSupply = (totalBond + swapBond) * (lpSupply - swapLpAmount)
	lpToken := dapp.LpToken()
	if swapBond.Denom != k.DefaultDenom(ctx) {
		return sdk.Coin{}, types.ErrInvalidLpToken
	}
	lpSupply := k.bk.GetSupply(ctx, lpToken).Amount
	totalBond := dapp.TotalBond.Amount
	swapBondAmount := swapBond.Amount
	totalLpAfterSwap := totalBond.Mul(lpSupply).Quo(totalBond.Add(swapBondAmount))
	swapLpAmount := lpSupply.Sub(totalLpAfterSwap)

	dapp.TotalBond.Amount = dapp.TotalBond.Amount.Add(swapBond.Amount)
	k.SetDapp(ctx, dapp)

	fee := sdk.NewDecFromInt(swapLpAmount).Mul(poolFee).RoundInt()
	if fee.IsPositive() {
		feeCoin := sdk.NewCoin(lpToken, fee)
		err := k.OnCollectFee(ctx, sdk.Coins{feeCoin})
		if err != nil {
			return sdk.Coin{}, err
		}
	}

	// send lp tokens to the module account
	err := k.bk.SendCoinsFromAccountToModule(ctx, addr, types.ModuleName, sdk.Coins{swapBond})
	if err != nil {
		return sdk.Coin{}, err
	}

	// send tokens to user
	userReceiveCoin := sdk.NewCoin(lpToken, swapLpAmount.Sub(fee))
	err = k.bk.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr, sdk.Coins{userReceiveCoin})
	if err != nil {
		return sdk.Coin{}, err
	}

	properties := k.gk.GetNetworkProperties(ctx)
	threshold := sdk.NewInt(int64(properties.DappLiquidationThreshold)).Mul(sdk.NewInt(1000_000))
	if dapp.LiquidationStart != 0 && dapp.TotalBond.Amount.GTE(threshold) {
		dapp.LiquidationStart = 0
		k.SetDapp(ctx, dapp)
	}

	return userReceiveCoin, nil
}

func (k Keeper) ConvertDappPoolTx(
	ctx sdk.Context,
	addr sdk.AccAddress,
	dapp1 types.Dapp,
	dapp2 types.Dapp,
	lpToken sdk.Coin,
) (sdk.Coin, error) {
	swapBond, err := k.RedeemDappPoolTx(ctx, addr, dapp1, dapp1.PoolFee.Quo(sdk.NewDec(2)), lpToken)
	if err != nil {
		return sdk.Coin{}, err
	}
	return k.SwapDappPoolTx(ctx, addr, dapp2, dapp2.PoolFee.Quo(sdk.NewDec(2)), swapBond)
}
