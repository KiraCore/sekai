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
		}
	}

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
				Weight:  userBond.Bond.Amount.Uint64(),
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
		rate := spendingPoolDeposit.ToDec().Quo(sdk.NewDec(int64(drip)))
		err = k.bk.MintCoins(ctx, types.ModuleName, sdk.Coins{sdk.NewCoin(dappBondLpToken, totalSupply)})
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
		k.SetDapp(ctx, dapp)

		for _, userBond := range userBonds {
			k.DeleteUserDappBond(ctx, dapp.Name, userBond.User)
		}
	}
}
