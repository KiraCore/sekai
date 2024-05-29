package keeper

import (
	"fmt"

	"github.com/KiraCore/sekai/x/layer2/types"
	spendingtypes "github.com/KiraCore/sekai/x/spending/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) BeginBlocker(ctx sdk.Context) {

}

func (k Keeper) EndBlocker(ctx sdk.Context) {
	dapps := k.GetAllDapps(ctx)
	properties := k.gk.GetNetworkProperties(ctx)
	currTimestamp := uint64(ctx.BlockTime().Unix())
	for _, dapp := range dapps {
		if dapp.Status == types.Bootstrap && dapp.CreationTime+properties.DappBondDuration <= currTimestamp {
			k.FinishDappBootstrap(ctx, dapp)
		}

		if dapp.PremintTime+dapp.Pool.Drip < uint64(ctx.BlockTime().Unix()) &&
			dapp.Issuance.Postmint.IsPositive() &&
			dapp.Status == types.Active {
			teamReserve := sdk.MustAccAddressFromBech32(dapp.TeamReserve)
			dappBondLpToken := dapp.LpToken()
			premintCoin := sdk.NewCoin(dappBondLpToken, dapp.Issuance.Premint)
			err := k.bk.SendCoinsFromModuleToAccount(ctx, types.ModuleName, teamReserve, sdk.Coins{premintCoin})
			if err != nil {
				panic(err)
			}
			dapp.PostMintPaid = true
			k.SetDapp(ctx, dapp)
		}

		if dapp.Status == types.Active {
			session := k.GetDappSession(ctx, dapp.Name)
			// session not started yet during denounce time
			if session.NextSession != nil &&
				session.NextSession.Start+properties.DappAutoDenounceTime <= currTimestamp &&
				session.NextSession.Status == types.SessionScheduled {
				operator := k.GetDappOperator(ctx, dapp.Name, session.NextSession.Leader)
				operator.Status = types.OperatorJailed
				k.SetDappOperator(ctx, operator)

				// update next session with new info
				k.ResetNewSession(ctx, dapp.Name, session.CurrSession.Leader)
			}

			// If the KEX collateral in the pool falls below dapp_liquidation_threshold (by default set to 100’000 KEX)
			// then the dApp will enter a depreciation phase lasting dapp_liquidation_period (by default set to 2419200, that is ~28d)
			// after which the execution will be stopped.
			if dapp.LiquidationStart+properties.DappLiquidationPeriod < uint64(ctx.BlockTime().Unix()) {
				dapp.Status = types.Halted
				k.SetDapp(ctx, dapp)
			}
		}
	}

	// handle bridge xam time outs
	for _, xam := range k.GetXAMs(ctx) {
		if xam.Req.SourceDapp != 0 && xam.Res.Src == 0 {
			account := k.GetBridgeAccount(ctx, xam.Req.SourceDapp)
			dapp := k.GetDapp(ctx, account.DappName)
			if xam.ReqTime+2*dapp.UpdateTimeMax < uint64(ctx.BlockTime().Unix()) {
				xam.Res.Src = 498
				k.SetXAM(ctx, xam)
			}
		}

		if xam.Req.DestDapp != 0 && xam.Res.Drc == 0 {
			account := k.GetBridgeAccount(ctx, xam.Req.DestDapp)
			dapp := k.GetDapp(ctx, account.DappName)
			if xam.ReqTime+2*dapp.UpdateTimeMax < uint64(ctx.BlockTime().Unix()) {
				xam.Res.Drc = 498
				k.SetXAM(ctx, xam)
			}
		}
	}
	// - **Deposit** funds from personal account (source application `src == 0`) to destination `dst` application beneficiary `ben` address.
	//     - Deposit must be accepted by the `dst` application or funds returned back to the user kira account outside of ABR
	//     - Max time for accepting deposit is `max(2 blocks,  2 * update_time_max<DstApp>)` after which transaction must fail (internal response `irc` code `522`)
	// - **Transfer** funds from the source `src` application account `acc` address to another destination `dst` application beneficiary `ben` address
	//     - Withdrawal from `src` application must be permissionless, the source `src` app does not need to confirm it but can speed it up by setting source response `src` code to `200`
	//     - Deposit to destination `dst` app must be accepted by the `dst` app with destination app response code `drc` set to `200` or funds **returned through deposit like process**
	//     - If Deposit fails to `dst` fails then `src` must accept the re-deposit otherwise if re-deposit fails the funds must be returned to the kira account outside of ABR
	//     - Max time for accepting withdrawal by `src` app must be `max(2 blocks,  2 * update_time_max<SrcApp>)` otherwise the withdrawal must automatically succeed
	//     - Max time for accepting deposit by `dst` app must be `max(2 blocks,  2 * update_time_max<DstApp>)` otherwise the deposit must automatically fail with internal response `irc` code set to `522`
	//     - If redeposit must be executed the new transaction must be created with new `xid`. Rules for re-deposit are same as for the deposit (if re-deposit fails funds must be returned to kira address outside of ABR)
	// - **Withdrawal** funds from the source `src` application account `acc` address to kira address outside the ABR (`dst == 0`) to beneficiary `ben` address
	//     - Withdrawal from `src` application must be permissionless, the source `src` app does not need to confirm it but can speed it up by setting source response `src` code to `200`
	//     - Max time for accepting withdrawal by `src` app must be `max(2 blocks,  2 * update_time_max<SrcApp>)` otherwise the withdrawal must automatically succeed

}

func (k Keeper) FinishDappBootstrap(ctx sdk.Context, dapp types.Dapp) {
	properties := k.gk.GetNetworkProperties(ctx)
	minDappBond := sdk.NewInt(int64(properties.MinDappBond)).Mul(sdk.NewInt(1000_000))
	if dapp.TotalBond.Amount.LT(minDappBond) {
		cacheCtx, write := ctx.CacheContext()
		err := k.ExecuteDappRemove(cacheCtx, dapp)
		if err == nil {
			write()
		}
	} else {
		userBonds := k.GetUserDappBonds(ctx, dapp.Name)
		beneficiaries := []spendingtypes.WeightedAccount{}
		for _, userBond := range userBonds {
			beneficiaries = append(beneficiaries, spendingtypes.WeightedAccount{
				Account: userBond.User,
				Weight:  userBond.Bond.Amount.ToLegacyDec(),
			})
		}

		dappBondLpToken := dapp.LpToken()
		err := sdk.ValidateDenom(dappBondLpToken)
		if err != nil {
			return
		}

		spendingPoolDeposit := dapp.GetSpendingPoolLpDeposit()
		totalSupply := dapp.GetLpTokenSupply()
		drip := dapp.Pool.Drip
		if drip == 0 {
			drip = 1
		}
		rate := sdk.NewDecFromInt(spendingPoolDeposit).Quo(sdk.NewDec(int64(drip)))
		err = k.tk.MintCoins(ctx, types.ModuleName, sdk.Coins{sdk.NewCoin(dappBondLpToken, totalSupply)})
		if err != nil {
			panic(err)
		}

		cacheCtx, write := ctx.CacheContext()
		blockTime := uint64(ctx.BlockTime().Unix())
		spendingPoolName := fmt.Sprintf("dp_%s", dapp.Name)
		err = k.spk.CreateSpendingPool(cacheCtx, spendingtypes.SpendingPool{
			Name:          spendingPoolName,
			ClaimStart:    blockTime,
			ClaimEnd:      blockTime + drip,
			Rates:         sdk.NewDecCoins(sdk.NewDecCoinFromDec(dappBondLpToken, rate)),
			VoteQuorum:    dapp.VoteQuorum,
			VotePeriod:    dapp.VotePeriod,
			VoteEnactment: dapp.VoteEnactment,
			Owners: &spendingtypes.PermInfo{
				OwnerRoles:    dapp.Controllers.Whitelist.Roles,
				OwnerAccounts: dapp.Controllers.Whitelist.Addresses,
			},
			Beneficiaries: &spendingtypes.WeightedPermInfo{
				Accounts: beneficiaries,
			},
			Balances:                []sdk.Coin{},
			DynamicRate:             false,
			DynamicRatePeriod:       0,
			LastDynamicRateCalcTime: 0,
		})
		if err == nil {
			write()
		}

		cacheCtx, write = ctx.CacheContext()
		coin := sdk.NewCoin(dappBondLpToken, spendingPoolDeposit)
		err = k.spk.DepositSpendingPoolFromModule(cacheCtx, types.ModuleName, spendingPoolName, sdk.Coins{coin})
		if err == nil {
			write()
		}

		dapp.Status = types.Halted
		dapp.Pool.Deposit = spendingPoolName
		dapp.PremintTime = blockTime
		k.SetDapp(ctx, dapp)

		// register bridge account for dapp
		helper := k.GetBridgeRegistrarHelper(ctx)
		k.SetBridgeAccount(ctx, types.BridgeAccount{
			Index:    helper.NextUser,
			Address:  dapp.GetAccount().String(),
			DappName: dapp.Name,
			Balances: []types.BridgeBalance{},
		})
		helper.NextUser += 1
		k.SetBridgeRegistrarHelper(ctx, helper)

		for _, userBond := range userBonds {
			k.DeleteUserDappBond(ctx, dapp.Name, userBond.User)
		}

		// send premint amount to team reserve
		if dapp.Issuance.Premint.IsPositive() {
			teamReserve := sdk.MustAccAddressFromBech32(dapp.TeamReserve)
			premintCoin := sdk.NewCoin(dappBondLpToken, dapp.Issuance.Premint)
			err = k.bk.SendCoinsFromModuleToAccount(ctx, types.ModuleName, teamReserve, sdk.Coins{premintCoin})
			if err != nil {
				panic(err)
			}
		}
	}
}
