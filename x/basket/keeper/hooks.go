package keeper

import (
	"fmt"

	"github.com/KiraCore/sekai/x/basket/types"
	multistakingtypes "github.com/KiraCore/sekai/x/multistaking/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) AfterUpsertStakingPool(ctx sdk.Context, valAddr sdk.ValAddress, pool multistakingtypes.StakingPool) {
	rates := k.tk.GetAllTokenRates(ctx)
	for _, rate := range rates {
		if rate.StakeToken {
			basket, err := k.GetBasketByDenom(ctx, fmt.Sprintf("sdb/%s", rate.Denom))
			if err != nil {
				basket = types.Basket{
					Id:              1,
					Suffix:          fmt.Sprint("staking/%s", rate.Denom),
					Description:     fmt.Sprintf("Basket of staking derivatives for %s token", rate.Denom),
					Amount:          sdk.ZeroInt(),
					SwapFee:         sdk.ZeroDec(),
					SlipppageFeeMin: sdk.ZeroDec(),
					TokensCap:       sdk.ZeroDec(),
					LimitsPeriod:    86400,
					MintsMin:        sdk.OneInt(),
					MintsMax:        sdk.NewInt(1000_000_000_000), // 1M
					MintsDisabled:   false,
					BurnsMin:        sdk.OneInt(),
					BurnsMax:        sdk.NewInt(1000_000_000_000), // 1M
					BurnsDisabled:   false,
					SwapsMin:        sdk.OneInt(),
					SwapsMax:        sdk.NewInt(1000_000_000_000), // 1M
					SwapsDisabled:   false,
					Tokens:          []types.BasketToken{},
					Surplus:         []sdk.Coin{},
				}
				k.SetBasket(ctx, basket)
			}

			shareDenom := multistakingtypes.GetShareDenom(pool.Id, rate.Denom)
			tokenMap := make(map[string]bool)
			for _, token := range basket.Tokens {
				tokenMap[token.Denom] = true
			}
			if !tokenMap[shareDenom] {
				basket.Tokens = append(basket.Tokens, types.BasketToken{
					Denom:     shareDenom,
					Weight:    sdk.OneDec(),
					Amount:    sdk.ZeroInt(),
					Deposits:  true,
					Withdraws: true,
					Swaps:     true,
				})
			}
		}
	}
}

//_________________________________________________________________________________________

// Hooks wrapper struct for multistaking keeper
type Hooks struct {
	k Keeper
}

var _ types.MultistakingHooks = Hooks{}

// Return the wrapper struct
func (k Keeper) Hooks() Hooks {
	return Hooks{k}
}

// Implements sdk.ValidatorHooks
func (h Hooks) AfterUpsertStakingPool(ctx sdk.Context, valAddr sdk.ValAddress, pool multistakingtypes.StakingPool) {
	h.k.AfterUpsertStakingPool(ctx, valAddr, pool)
}
