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

	for _, dapp := range dapps {
		minDappBond := sdk.NewInt(int64(properties.MinDappBond)).Mul(sdk.NewInt(1000_000))

		if int64(dapp.CreationTime+properties.DappBondDuration) <= ctx.BlockTime().Unix() {
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

				dappBondLpToken := fmt.Sprintf("lp_%s", dapp.Name)
				err := sdk.ValidateDenom(dappBondLpToken)
				if err != nil {
					continue
				}

				spendingPoolDeposit := dapp.TotalBond.Amount.ToDec().Mul(dapp.Pool.Ratio).RoundInt()
				totalSupply := spendingPoolDeposit.Add(dapp.Issurance.Postmint).Add(dapp.Issurance.Premint)
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
				err = k.spk.DepositSpendingPoolFromModule(cacheCtx, types.ModuleName, spendingPoolName, sdk.NewCoin(dappBondLpToken, spendingPoolDeposit))
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
	}

}
