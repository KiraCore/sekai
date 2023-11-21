package keeper_test

import (
	"time"

	"github.com/KiraCore/sekai/x/basket/types"
	"github.com/cometbft/cometbft/crypto/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
)

func (suite *KeeperTestSuite) TestMintBasketToken() {
	testCases := map[string]struct {
		basketId          uint64
		mintDisabled      bool
		userBalance       sdk.Coins
		depositBalance    sdk.Coins
		denomDisabled     bool
		minDepositAmount  sdk.Int
		maxDepositAmount  sdk.Int
		limitPeriod       uint64
		prevDepositAmount sdk.Int
		tokensCap         sdk.Dec
		expectErr         bool
		expectedOutAmount sdk.Int
	}{
		"case not available basket": {
			basketId:          0,
			mintDisabled:      false,
			userBalance:       sdk.NewCoins(sdk.NewInt64Coin("ukex", 1000_000)),
			depositBalance:    sdk.NewCoins(sdk.NewInt64Coin("ukex", 1000_000)),
			denomDisabled:     false,
			minDepositAmount:  sdk.NewInt(1000_000),
			maxDepositAmount:  sdk.NewInt(100_000_000),
			limitPeriod:       3600,
			prevDepositAmount: sdk.NewInt(0),
			tokensCap:         sdk.NewDec(1),
			expectErr:         true,
			expectedOutAmount: sdk.NewInt(0),
		},
		"case mint disabled basket": {
			basketId:          1,
			mintDisabled:      true,
			userBalance:       sdk.NewCoins(sdk.NewInt64Coin("ukex", 1000_000)),
			depositBalance:    sdk.NewCoins(sdk.NewInt64Coin("ukex", 1000_000)),
			denomDisabled:     false,
			minDepositAmount:  sdk.NewInt(1000_000),
			maxDepositAmount:  sdk.NewInt(100_000_000),
			limitPeriod:       3600,
			prevDepositAmount: sdk.NewInt(0),
			tokensCap:         sdk.NewDec(1),
			expectErr:         true,
			expectedOutAmount: sdk.NewInt(0),
		},
		"case not enough balance on user": {
			basketId:          1,
			mintDisabled:      true,
			userBalance:       sdk.NewCoins(sdk.NewInt64Coin("ukex", 100_000)),
			depositBalance:    sdk.NewCoins(sdk.NewInt64Coin("ukex", 1000_000)),
			denomDisabled:     false,
			minDepositAmount:  sdk.NewInt(1000_000),
			maxDepositAmount:  sdk.NewInt(100_000_000),
			limitPeriod:       3600,
			prevDepositAmount: sdk.NewInt(0),
			tokensCap:         sdk.NewDec(1),
			expectErr:         true,
			expectedOutAmount: sdk.NewInt(0),
		},
		"case not deposit denom": {
			basketId:          1,
			mintDisabled:      false,
			userBalance:       sdk.NewCoins(sdk.NewInt64Coin("xxx", 1000_000)),
			depositBalance:    sdk.NewCoins(sdk.NewInt64Coin("xxx", 1000_000)),
			denomDisabled:     false,
			minDepositAmount:  sdk.NewInt(1000_000),
			maxDepositAmount:  sdk.NewInt(100_000_000),
			limitPeriod:       3600,
			prevDepositAmount: sdk.NewInt(0),
			tokensCap:         sdk.NewDec(1),
			expectErr:         true,
			expectedOutAmount: sdk.NewInt(0),
		},
		"case deposit denom disabled": {
			basketId:          1,
			mintDisabled:      false,
			userBalance:       sdk.NewCoins(sdk.NewInt64Coin("ukex", 1000_000)),
			depositBalance:    sdk.NewCoins(sdk.NewInt64Coin("ukex", 1000_000)),
			denomDisabled:     true,
			minDepositAmount:  sdk.NewInt(1000_000),
			maxDepositAmount:  sdk.NewInt(100_000_000),
			limitPeriod:       3600,
			prevDepositAmount: sdk.NewInt(0),
			tokensCap:         sdk.NewDec(1),
			expectErr:         true,
			expectedOutAmount: sdk.NewInt(0),
		},
		"case lower than deposit min amount": {
			basketId:          1,
			mintDisabled:      false,
			userBalance:       sdk.NewCoins(sdk.NewInt64Coin("ukex", 100_000)),
			depositBalance:    sdk.NewCoins(sdk.NewInt64Coin("ukex", 100_000)),
			denomDisabled:     false,
			minDepositAmount:  sdk.NewInt(1000_000),
			maxDepositAmount:  sdk.NewInt(100_000_000),
			limitPeriod:       3600,
			prevDepositAmount: sdk.NewInt(0),
			tokensCap:         sdk.NewDec(1),
			expectErr:         true,
			expectedOutAmount: sdk.NewInt(0),
		},
		"case exceeding deposit max amount during limit period": {
			basketId:          1,
			mintDisabled:      false,
			userBalance:       sdk.NewCoins(sdk.NewInt64Coin("ukex", 2000_000)),
			depositBalance:    sdk.NewCoins(sdk.NewInt64Coin("ukex", 2000_000)),
			denomDisabled:     false,
			minDepositAmount:  sdk.NewInt(1000_000),
			maxDepositAmount:  sdk.NewInt(100_000_000),
			limitPeriod:       3600,
			prevDepositAmount: sdk.NewInt(99_000_000),
			tokensCap:         sdk.NewDec(1),
			expectErr:         true,
			expectedOutAmount: sdk.NewInt(0),
		},
		"case tokens cap is broken": {
			basketId:          1,
			mintDisabled:      false,
			userBalance:       sdk.NewCoins(sdk.NewInt64Coin("ukex", 2000_000)),
			depositBalance:    sdk.NewCoins(sdk.NewInt64Coin("ukex", 2000_000)),
			denomDisabled:     false,
			minDepositAmount:  sdk.NewInt(1000_000),
			maxDepositAmount:  sdk.NewInt(100_000_000),
			limitPeriod:       3600,
			prevDepositAmount: sdk.NewInt(0),
			tokensCap:         sdk.NewDecWithPrec(8, 1), // 80%
			expectErr:         true,
			expectedOutAmount: sdk.NewInt(0),
		},
		"case one token successful deposit": {
			basketId:          1,
			mintDisabled:      false,
			userBalance:       sdk.NewCoins(sdk.NewInt64Coin("ukex", 1000_000)),
			depositBalance:    sdk.NewCoins(sdk.NewInt64Coin("ukex", 1000_000)),
			denomDisabled:     false,
			minDepositAmount:  sdk.NewInt(1000_000),
			maxDepositAmount:  sdk.NewInt(100_000_000),
			limitPeriod:       3600,
			prevDepositAmount: sdk.NewInt(0),
			tokensCap:         sdk.NewDec(1),
			expectErr:         false,
			expectedOutAmount: sdk.NewInt(1000_000),
		},
		"case two tokens successful deposit": {
			basketId:          1,
			mintDisabled:      false,
			userBalance:       sdk.NewCoins(sdk.NewInt64Coin("ukex", 1000_000), sdk.NewInt64Coin("ueth", 100_000)),
			depositBalance:    sdk.NewCoins(sdk.NewInt64Coin("ukex", 1000_000), sdk.NewInt64Coin("ueth", 100_000)),
			denomDisabled:     false,
			minDepositAmount:  sdk.NewInt(1000_000),
			maxDepositAmount:  sdk.NewInt(100_000_000),
			limitPeriod:       3600,
			prevDepositAmount: sdk.NewInt(0),
			tokensCap:         sdk.NewDec(1),
			expectErr:         false,
			expectedOutAmount: sdk.NewInt(2000_000),
		},
	}

	for name, tc := range testCases {
		tc := tc

		suite.Run(name, func() {
			suite.SetupTest()
			now := time.Now().UTC()
			suite.ctx = suite.ctx.WithBlockTime(now)

			addr1 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
			basket := types.Basket{
				Id:              1,
				Suffix:          "usd",
				Description:     "usd basket",
				LimitsPeriod:    tc.limitPeriod,
				Amount:          sdk.NewInt(0),
				SwapFee:         sdk.NewDecWithPrec(1, 2), // 1%
				SlipppageFeeMin: sdk.NewDecWithPrec(1, 2), // 1%
				TokensCap:       tc.tokensCap,
				MintsMin:        tc.minDepositAmount,
				MintsMax:        tc.maxDepositAmount,
				MintsDisabled:   tc.mintDisabled,
				BurnsMin:        sdk.NewInt(1),
				BurnsMax:        sdk.NewInt(100000000),
				BurnsDisabled:   false,
				SwapsMin:        sdk.NewInt(1),
				SwapsMax:        sdk.NewInt(100000000),
				SwapsDisabled:   false,
				Tokens: []types.BasketToken{
					{
						Denom:     "ukex",
						Weight:    sdk.NewDec(1),
						Amount:    sdk.NewInt(0),
						Deposits:  !tc.denomDisabled,
						Withdraws: !tc.denomDisabled,
						Swaps:     !tc.denomDisabled,
					},
					{
						Denom:     "ueth",
						Weight:    sdk.NewDec(10),
						Amount:    sdk.NewInt(0),
						Deposits:  !tc.denomDisabled,
						Withdraws: !tc.denomDisabled,
						Swaps:     !tc.denomDisabled,
					},
				},
				Surplus: sdk.NewCoins(sdk.NewInt64Coin("usd2", 1)),
			}
			suite.app.BasketKeeper.SetBasket(suite.ctx, basket)

			suite.app.BasketKeeper.SetMintAmount(suite.ctx, now.Add(time.Second*60), tc.basketId, tc.prevDepositAmount)

			err := suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, tc.userBalance)
			suite.Require().NoError(err)
			err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, addr1, tc.userBalance)
			suite.Require().NoError(err)

			err = suite.app.BasketKeeper.MintBasketToken(suite.ctx, &types.MsgBasketTokenMint{
				Sender:   addr1.String(),
				BasketId: tc.basketId,
				Deposit:  tc.depositBalance,
			})
			if tc.expectErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)

				// check basket tokens balance increased
				balance := suite.app.BankKeeper.GetBalance(suite.ctx, addr1, basket.GetBasketDenom())
				suite.Require().Equal(balance.Amount, tc.expectedOutAmount)

				// check user's deposit balance decrease
				balances := suite.app.BankKeeper.GetAllBalances(suite.ctx, addr1)
				suite.Require().Equal(balances, sdk.Coins{balance})

				// check basket total amount increased
				savedBasket, err := suite.app.BasketKeeper.GetBasketById(suite.ctx, basket.Id)
				suite.Require().NoError(err)
				suite.Require().Equal(savedBasket.Amount, tc.expectedOutAmount)

				// check basket tokens balance increase
				basketUnderlyingCoins := sdk.Coins{}
				for _, token := range savedBasket.Tokens {
					basketUnderlyingCoins = basketUnderlyingCoins.Add(sdk.NewCoin(token.Denom, token.Amount))
				}
				suite.Require().Equal(basketUnderlyingCoins, tc.depositBalance)

				// check limit period amount increased
				historicalAmount := suite.app.BasketKeeper.GetLimitsPeriodMintAmount(suite.ctx, 1, tc.limitPeriod)
				suite.Require().Equal(historicalAmount, tc.prevDepositAmount.Add(tc.expectedOutAmount))
			}
		})
	}
}

func (suite *KeeperTestSuite) TestBurnBasketToken() {
	sampleBasket := types.Basket{
		Id:          1,
		Suffix:      "usd",
		Description: "usd basket",
	}
	basketDenom := sampleBasket.GetBasketDenom()

	testCases := map[string]struct {
		basketId          uint64
		burnDisabled      bool
		userBalance       sdk.Coins
		burnBalance       sdk.Coin
		denomDisabled     bool
		minBurnAmount     sdk.Int
		maxBurnAmount     sdk.Int
		limitPeriod       uint64
		prevBurnAmount    sdk.Int
		tokensCap         sdk.Dec
		expectErr         bool
		expectedOutAmount sdk.Coins
	}{
		"case not available basket": {
			basketId:          0,
			burnDisabled:      false,
			userBalance:       sdk.NewCoins(sdk.NewInt64Coin(basketDenom, 1000_000)),
			burnBalance:       sdk.NewInt64Coin(basketDenom, 1000_000),
			denomDisabled:     false,
			minBurnAmount:     sdk.NewInt(1000_000),
			maxBurnAmount:     sdk.NewInt(100_000_000),
			limitPeriod:       3600,
			prevBurnAmount:    sdk.NewInt(0),
			tokensCap:         sdk.NewDec(1),
			expectErr:         true,
			expectedOutAmount: sdk.NewCoins(),
		},
		"case burn disabled basket": {
			basketId:          1,
			burnDisabled:      true,
			userBalance:       sdk.NewCoins(sdk.NewInt64Coin(basketDenom, 1000_000)),
			burnBalance:       sdk.NewInt64Coin(basketDenom, 1000_000),
			denomDisabled:     false,
			minBurnAmount:     sdk.NewInt(1000_000),
			maxBurnAmount:     sdk.NewInt(100_000_000),
			limitPeriod:       3600,
			prevBurnAmount:    sdk.NewInt(0),
			tokensCap:         sdk.NewDec(1),
			expectErr:         true,
			expectedOutAmount: sdk.NewCoins(),
		},
		"case not enough balance on user": {
			basketId:          1,
			burnDisabled:      true,
			userBalance:       sdk.NewCoins(sdk.NewInt64Coin(basketDenom, 100_000)),
			burnBalance:       sdk.NewInt64Coin(basketDenom, 1000_000),
			denomDisabled:     false,
			minBurnAmount:     sdk.NewInt(1000_000),
			maxBurnAmount:     sdk.NewInt(100_000_000),
			limitPeriod:       3600,
			prevBurnAmount:    sdk.NewInt(0),
			tokensCap:         sdk.NewDec(1),
			expectErr:         true,
			expectedOutAmount: sdk.NewCoins(),
		},
		"case not basket denom": {
			basketId:          1,
			burnDisabled:      false,
			userBalance:       sdk.NewCoins(sdk.NewInt64Coin("xxx", 1000_000)),
			burnBalance:       sdk.NewInt64Coin("xxx", 1000_000),
			denomDisabled:     false,
			minBurnAmount:     sdk.NewInt(1000_000),
			maxBurnAmount:     sdk.NewInt(100_000_000),
			limitPeriod:       3600,
			prevBurnAmount:    sdk.NewInt(0),
			tokensCap:         sdk.NewDec(1),
			expectErr:         true,
			expectedOutAmount: sdk.NewCoins(),
		},
		"case withdraw denom disabled": {
			basketId:          1,
			burnDisabled:      false,
			userBalance:       sdk.NewCoins(sdk.NewInt64Coin(basketDenom, 1000_000)),
			burnBalance:       sdk.NewInt64Coin(basketDenom, 1000_000),
			denomDisabled:     true,
			minBurnAmount:     sdk.NewInt(1000_000),
			maxBurnAmount:     sdk.NewInt(100_000_000),
			limitPeriod:       3600,
			prevBurnAmount:    sdk.NewInt(0),
			tokensCap:         sdk.NewDec(1),
			expectErr:         true,
			expectedOutAmount: sdk.NewCoins(),
		},
		"case lower than withdraw min amount": {
			basketId:          1,
			burnDisabled:      false,
			userBalance:       sdk.NewCoins(sdk.NewInt64Coin(basketDenom, 100_000)),
			burnBalance:       sdk.NewInt64Coin(basketDenom, 100_000),
			denomDisabled:     false,
			minBurnAmount:     sdk.NewInt(1000_000),
			maxBurnAmount:     sdk.NewInt(100_000_000),
			limitPeriod:       3600,
			prevBurnAmount:    sdk.NewInt(0),
			tokensCap:         sdk.NewDec(1),
			expectErr:         true,
			expectedOutAmount: sdk.NewCoins(),
		},
		"case exceeding withdraw max amount during limit period": {
			basketId:          1,
			burnDisabled:      false,
			userBalance:       sdk.NewCoins(sdk.NewInt64Coin(basketDenom, 2000_000)),
			burnBalance:       sdk.NewInt64Coin(basketDenom, 2000_000),
			denomDisabled:     false,
			minBurnAmount:     sdk.NewInt(1000_000),
			maxBurnAmount:     sdk.NewInt(100_000_000),
			limitPeriod:       3600,
			prevBurnAmount:    sdk.NewInt(99_000_000),
			tokensCap:         sdk.NewDec(1),
			expectErr:         true,
			expectedOutAmount: sdk.NewCoins(),
		},
		"case successful withdraw": {
			basketId:          1,
			burnDisabled:      false,
			userBalance:       sdk.NewCoins(sdk.NewInt64Coin(basketDenom, 1000_000)),
			burnBalance:       sdk.NewInt64Coin(basketDenom, 1000_000),
			denomDisabled:     false,
			minBurnAmount:     sdk.NewInt(1000_000),
			maxBurnAmount:     sdk.NewInt(100_000_000),
			limitPeriod:       3600,
			prevBurnAmount:    sdk.NewInt(0),
			tokensCap:         sdk.NewDec(1),
			expectErr:         false,
			expectedOutAmount: sdk.NewCoins(sdk.NewInt64Coin("ukex", 500_000), sdk.NewInt64Coin("ueth", 50_000)),
		},
	}

	for name, tc := range testCases {
		tc := tc

		suite.Run(name, func() {
			suite.SetupTest()
			now := time.Now().UTC()
			suite.ctx = suite.ctx.WithBlockTime(now)

			addr1 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
			addr2 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
			basket := types.Basket{
				Id:              1,
				Suffix:          "usd",
				Description:     "usd basket",
				LimitsPeriod:    tc.limitPeriod,
				Amount:          sdk.NewInt(0),
				SwapFee:         sdk.NewDecWithPrec(1, 2), // 1%
				SlipppageFeeMin: sdk.NewDecWithPrec(1, 2), // 1%
				TokensCap:       tc.tokensCap,
				MintsMin:        sdk.NewInt(1),
				MintsMax:        sdk.NewInt(100000000),
				MintsDisabled:   false,
				BurnsMin:        tc.minBurnAmount,
				BurnsMax:        tc.maxBurnAmount,
				BurnsDisabled:   tc.burnDisabled,
				SwapsMin:        sdk.NewInt(1),
				SwapsMax:        sdk.NewInt(100000000),
				SwapsDisabled:   false,
				Tokens: []types.BasketToken{
					{
						Denom:     "ukex",
						Weight:    sdk.NewDec(1),
						Amount:    sdk.NewInt(0),
						Deposits:  true,
						Withdraws: !tc.denomDisabled,
						Swaps:     !tc.denomDisabled,
					},
					{
						Denom:     "ueth",
						Weight:    sdk.NewDec(10),
						Amount:    sdk.NewInt(0),
						Deposits:  true,
						Withdraws: !tc.denomDisabled,
						Swaps:     !tc.denomDisabled,
					},
				},
				Surplus: sdk.NewCoins(sdk.NewInt64Coin("usd2", 1)),
			}
			suite.app.BasketKeeper.SetBasket(suite.ctx, basket)

			depositCoins := sdk.NewCoins(sdk.NewInt64Coin("ukex", 1000_000), sdk.NewInt64Coin("ueth", 100_000))
			err := suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, depositCoins)
			suite.Require().NoError(err)
			err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, addr1, depositCoins)
			suite.Require().NoError(err)

			err = suite.app.BasketKeeper.MintBasketToken(suite.ctx, &types.MsgBasketTokenMint{
				Sender:   addr1.String(),
				BasketId: 1,
				Deposit:  depositCoins,
			})
			suite.Require().NoError(err)

			err = suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, tc.userBalance)
			suite.Require().NoError(err)
			err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, addr2, tc.userBalance)
			suite.Require().NoError(err)

			suite.app.BasketKeeper.SetBurnAmount(suite.ctx, now.Add(time.Second*60), tc.basketId, tc.prevBurnAmount)

			err = suite.app.BasketKeeper.BurnBasketToken(suite.ctx, &types.MsgBasketTokenBurn{
				Sender:     addr2.String(),
				BasketId:   tc.basketId,
				BurnAmount: tc.burnBalance,
			})

			if tc.expectErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)

				// check basket tokens balance decreased
				balance := suite.app.BankKeeper.GetBalance(suite.ctx, addr2, basket.GetBasketDenom())
				suite.Require().Equal(balance.Amount, sdk.ZeroInt())

				// check user's withdraw balance increase
				balances := suite.app.BankKeeper.GetAllBalances(suite.ctx, addr2)
				suite.Require().Equal(balances, tc.expectedOutAmount)

				// check basket total amount increased
				savedBasket, err := suite.app.BasketKeeper.GetBasketById(suite.ctx, basket.Id)
				suite.Require().NoError(err)
				suite.Require().Equal(savedBasket.Amount.Add(tc.burnBalance.Amount), sdk.NewInt(2000_000))

				// check basket tokens balance increase
				basketUnderlyingCoins := sdk.Coins{}
				for _, token := range savedBasket.Tokens {
					basketUnderlyingCoins = basketUnderlyingCoins.Add(sdk.NewCoin(token.Denom, token.Amount))
				}
				suite.Require().Equal(basketUnderlyingCoins.Add(tc.expectedOutAmount...).String(), depositCoins.String())

				// check limit period amount increased
				historicalAmount := suite.app.BasketKeeper.GetLimitsPeriodBurnAmount(suite.ctx, 1, tc.limitPeriod)
				suite.Require().Equal(historicalAmount, tc.prevBurnAmount.Add(tc.burnBalance.Amount))
			}
		})
	}
}

func (suite *KeeperTestSuite) TestBasketSwap() {
	testCases := map[string]struct {
		basketId          uint64
		swapDisabled      bool
		userBalance       sdk.Coins
		swapBalance       sdk.Coin
		denomDisabled     bool
		minSwapAmount     sdk.Int
		maxSwapAmount     sdk.Int
		limitPeriod       uint64
		prevSwapAmount    sdk.Int
		tokensCap         sdk.Dec
		expectErr         bool
		expectedOutAmount sdk.Coins
	}{
		"case not available basket": {
			basketId:          0,
			swapDisabled:      false,
			userBalance:       sdk.NewCoins(sdk.NewInt64Coin("ukex", 1000_000)),
			swapBalance:       sdk.NewInt64Coin("ukex", 100_000),
			denomDisabled:     false,
			minSwapAmount:     sdk.NewInt(100_000),
			maxSwapAmount:     sdk.NewInt(100_000_000),
			limitPeriod:       3600,
			prevSwapAmount:    sdk.NewInt(0),
			tokensCap:         sdk.NewDec(1),
			expectErr:         true,
			expectedOutAmount: sdk.NewCoins(),
		},
		"case swap disabled basket": {
			basketId:          1,
			swapDisabled:      true,
			userBalance:       sdk.NewCoins(sdk.NewInt64Coin("ukex", 1000_000)),
			swapBalance:       sdk.NewInt64Coin("ukex", 100_000),
			denomDisabled:     false,
			minSwapAmount:     sdk.NewInt(100_000),
			maxSwapAmount:     sdk.NewInt(100_000_000),
			limitPeriod:       3600,
			prevSwapAmount:    sdk.NewInt(0),
			tokensCap:         sdk.NewDec(1),
			expectErr:         true,
			expectedOutAmount: sdk.NewCoins(),
		},
		"case not enough balance on user": {
			basketId:          1,
			swapDisabled:      true,
			userBalance:       sdk.NewCoins(sdk.NewInt64Coin("ukex", 100)),
			swapBalance:       sdk.NewInt64Coin("ukex", 1000),
			denomDisabled:     false,
			minSwapAmount:     sdk.NewInt(100_000),
			maxSwapAmount:     sdk.NewInt(100_000_000),
			limitPeriod:       3600,
			prevSwapAmount:    sdk.NewInt(0),
			tokensCap:         sdk.NewDec(1),
			expectErr:         true,
			expectedOutAmount: sdk.NewCoins(),
		},
		"case not basket denom": {
			basketId:          1,
			swapDisabled:      false,
			userBalance:       sdk.NewCoins(sdk.NewInt64Coin("xxx", 100_000)),
			swapBalance:       sdk.NewInt64Coin("xxx", 100_000),
			denomDisabled:     false,
			minSwapAmount:     sdk.NewInt(100_000),
			maxSwapAmount:     sdk.NewInt(100_000_000),
			limitPeriod:       3600,
			prevSwapAmount:    sdk.NewInt(0),
			tokensCap:         sdk.NewDec(1),
			expectErr:         true,
			expectedOutAmount: sdk.NewCoins(),
		},
		"case swap denom disabled": {
			basketId:          1,
			swapDisabled:      false,
			userBalance:       sdk.NewCoins(sdk.NewInt64Coin("ukex", 100_000)),
			swapBalance:       sdk.NewInt64Coin("ukex", 100_000),
			denomDisabled:     true,
			minSwapAmount:     sdk.NewInt(100_000),
			maxSwapAmount:     sdk.NewInt(100_000_000),
			limitPeriod:       3600,
			prevSwapAmount:    sdk.NewInt(0),
			tokensCap:         sdk.NewDec(1),
			expectErr:         true,
			expectedOutAmount: sdk.NewCoins(),
		},
		"case lower than swap min amount": {
			basketId:          1,
			swapDisabled:      false,
			userBalance:       sdk.NewCoins(sdk.NewInt64Coin("ukex", 10_000)),
			swapBalance:       sdk.NewInt64Coin("ukex", 10_000),
			denomDisabled:     false,
			minSwapAmount:     sdk.NewInt(100_000),
			maxSwapAmount:     sdk.NewInt(100_000_000),
			limitPeriod:       3600,
			prevSwapAmount:    sdk.NewInt(0),
			tokensCap:         sdk.NewDec(1),
			expectErr:         true,
			expectedOutAmount: sdk.NewCoins(),
		},
		"case exceeding swap max amount during limit period": {
			basketId:          1,
			swapDisabled:      false,
			userBalance:       sdk.NewCoins(sdk.NewInt64Coin("ukex", 2000_000)),
			swapBalance:       sdk.NewInt64Coin("ukex", 200_000),
			denomDisabled:     false,
			minSwapAmount:     sdk.NewInt(100_000),
			maxSwapAmount:     sdk.NewInt(100_000_000),
			limitPeriod:       3600,
			prevSwapAmount:    sdk.NewInt(99_900_000),
			tokensCap:         sdk.NewDec(1),
			expectErr:         true,
			expectedOutAmount: sdk.NewCoins(),
		},
		"case tokens cap broken": {
			basketId:          1,
			swapDisabled:      false,
			userBalance:       sdk.NewCoins(sdk.NewInt64Coin("ukex", 500_000)),
			swapBalance:       sdk.NewInt64Coin("ukex", 500_000),
			denomDisabled:     false,
			minSwapAmount:     sdk.NewInt(100_000),
			maxSwapAmount:     sdk.NewInt(100_000_000),
			limitPeriod:       3600,
			prevSwapAmount:    sdk.NewInt(99_000_000),
			tokensCap:         sdk.NewDecWithPrec(6, 1), // 60%
			expectErr:         true,
			expectedOutAmount: sdk.NewCoins(),
		},
		"case successful swap": {
			basketId:          1,
			swapDisabled:      false,
			userBalance:       sdk.NewCoins(sdk.NewInt64Coin("ukex", 100_000)),
			swapBalance:       sdk.NewInt64Coin("ukex", 100_000),
			denomDisabled:     false,
			minSwapAmount:     sdk.NewInt(100_000),
			maxSwapAmount:     sdk.NewInt(100_000_000),
			limitPeriod:       3600,
			prevSwapAmount:    sdk.NewInt(0),
			tokensCap:         sdk.NewDec(1),
			expectErr:         false,
			expectedOutAmount: sdk.NewCoins(sdk.NewInt64Coin("ueth", 8_919)),
		},
	}

	for name, tc := range testCases {
		tc := tc

		suite.Run(name, func() {
			suite.SetupTest()
			now := time.Now().UTC()
			suite.ctx = suite.ctx.WithBlockTime(now)

			addr1 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
			addr2 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
			basket := types.Basket{
				Id:              1,
				Suffix:          "usd",
				Description:     "usd basket",
				LimitsPeriod:    tc.limitPeriod,
				Amount:          sdk.NewInt(0),
				SwapFee:         sdk.NewDecWithPrec(1, 2), // 1%
				SlipppageFeeMin: sdk.NewDecWithPrec(1, 2), // 1%
				TokensCap:       tc.tokensCap,
				MintsMin:        sdk.NewInt(1),
				MintsMax:        sdk.NewInt(100000000),
				MintsDisabled:   false,
				BurnsMin:        sdk.NewInt(1),
				BurnsMax:        sdk.NewInt(100000000),
				BurnsDisabled:   false,
				SwapsMin:        tc.minSwapAmount,
				SwapsMax:        tc.maxSwapAmount,
				SwapsDisabled:   tc.swapDisabled,
				Tokens: []types.BasketToken{
					{
						Denom:     "ukex",
						Weight:    sdk.NewDec(1),
						Amount:    sdk.NewInt(0),
						Deposits:  true,
						Withdraws: !tc.denomDisabled,
						Swaps:     !tc.denomDisabled,
					},
					{
						Denom:     "ueth",
						Weight:    sdk.NewDec(10),
						Amount:    sdk.NewInt(0),
						Deposits:  true,
						Withdraws: !tc.denomDisabled,
						Swaps:     !tc.denomDisabled,
					},
				},
				Surplus: sdk.NewCoins(sdk.NewInt64Coin("usd2", 1)),
			}
			suite.app.BasketKeeper.SetBasket(suite.ctx, basket)

			depositCoins := sdk.NewCoins(sdk.NewInt64Coin("ukex", 1000_000), sdk.NewInt64Coin("ueth", 100_000))
			err := suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, depositCoins)
			suite.Require().NoError(err)
			err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, addr1, depositCoins)
			suite.Require().NoError(err)

			err = suite.app.BasketKeeper.MintBasketToken(suite.ctx, &types.MsgBasketTokenMint{
				Sender:   addr1.String(),
				BasketId: 1,
				Deposit:  depositCoins,
			})
			suite.Require().NoError(err)

			err = suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, tc.userBalance)
			suite.Require().NoError(err)
			err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, addr2, tc.userBalance)
			suite.Require().NoError(err)

			suite.app.BasketKeeper.SetSwapAmount(suite.ctx, now.Add(time.Second*60), tc.basketId, tc.prevSwapAmount)

			err = suite.app.BasketKeeper.BasketSwap(suite.ctx, &types.MsgBasketTokenSwap{
				Sender:   addr2.String(),
				BasketId: tc.basketId,
				Pairs: []types.SwapPair{
					{
						InAmount: tc.swapBalance,
						OutToken: "ueth",
					},
				},
			})

			if tc.expectErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)

				// check in token balance decreased
				balance := suite.app.BankKeeper.GetBalance(suite.ctx, addr2, "ukex")
				suite.Require().Equal(balance.Amount, sdk.ZeroInt())

				// check user's withdraw balance increase
				balances := suite.app.BankKeeper.GetAllBalances(suite.ctx, addr2)
				suite.Require().Equal(balances.String(), tc.expectedOutAmount.String())

				// check basket total amount kept as it is
				savedBasket, err := suite.app.BasketKeeper.GetBasketById(suite.ctx, basket.Id)
				suite.Require().NoError(err)
				suite.Require().Equal(savedBasket.Amount, sdk.NewInt(2000_000))

				// check basket tokens balance changes
				basketUnderlyingCoins := sdk.Coins{}
				for _, token := range savedBasket.Tokens {
					basketUnderlyingCoins = basketUnderlyingCoins.Add(sdk.NewCoin(token.Denom, token.Amount))
				}
				suite.Require().True(basketUnderlyingCoins.Add(tc.expectedOutAmount...).Sub(tc.swapBalance).IsAllLTE(depositCoins))

				// check limit period amount increased
				historicalAmount := suite.app.BasketKeeper.GetLimitsPeriodSwapAmount(suite.ctx, 1, tc.limitPeriod)
				suite.Require().Equal(historicalAmount, tc.prevSwapAmount.Add(tc.swapBalance.Amount))

				// check correct slippage amount + surplus
				suite.Require().True(sdk.Coins(savedBasket.Surplus).Sub(basket.Surplus...).IsAllPositive())
			}
		})
	}
}

func (suite *KeeperTestSuite) TestBasketWithdrawSurplus() {
	suite.SetupTest()
	now := time.Now().UTC()
	suite.ctx = suite.ctx.WithBlockTime(now)

	addr1 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	addr2 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	basket := types.Basket{
		Id:              1,
		Suffix:          "usd",
		Description:     "usd basket",
		LimitsPeriod:    3600,
		Amount:          sdk.NewInt(0),
		SwapFee:         sdk.NewDecWithPrec(1, 2), // 1%
		SlipppageFeeMin: sdk.NewDecWithPrec(1, 2), // 1%
		TokensCap:       sdk.NewDecWithPrec(9, 1), // 90%
		MintsMin:        sdk.NewInt(1),
		MintsMax:        sdk.NewInt(100000000),
		MintsDisabled:   false,
		BurnsMin:        sdk.NewInt(1),
		BurnsMax:        sdk.NewInt(100000000),
		BurnsDisabled:   false,
		SwapsMin:        sdk.NewInt(1),
		SwapsMax:        sdk.NewInt(100000000),
		SwapsDisabled:   false,
		Tokens: []types.BasketToken{
			{
				Denom:     "ukex",
				Weight:    sdk.NewDec(1),
				Amount:    sdk.NewInt(0),
				Deposits:  true,
				Withdraws: true,
				Swaps:     true,
			},
			{
				Denom:     "ueth",
				Weight:    sdk.NewDec(10),
				Amount:    sdk.NewInt(0),
				Deposits:  true,
				Withdraws: true,
				Swaps:     true,
			},
		},
		Surplus: sdk.NewCoins(sdk.NewInt64Coin("ueth", 1)),
	}
	suite.app.BasketKeeper.SetBasket(suite.ctx, basket)

	depositCoins := sdk.NewCoins(sdk.NewInt64Coin("ukex", 1000_000), sdk.NewInt64Coin("ueth", 100_000))
	err := suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, depositCoins)
	suite.Require().NoError(err)
	err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, addr1, depositCoins)
	suite.Require().NoError(err)

	err = suite.app.BasketKeeper.MintBasketToken(suite.ctx, &types.MsgBasketTokenMint{
		Sender:   addr1.String(),
		BasketId: 1,
		Deposit:  depositCoins,
	})
	suite.Require().NoError(err)

	err = suite.app.BasketKeeper.BasketWithdrawSurplus(suite.ctx, types.ProposalBasketWithdrawSurplus{
		BasketIds:      []uint64{1},
		WithdrawTarget: addr2.String(),
	})
	suite.Require().NoError(err)

	// check account balance increased
	balance := suite.app.BankKeeper.GetAllBalances(suite.ctx, addr2)
	suite.Require().Equal(balance, sdk.Coins(basket.Surplus))

	// check surplus removal
	savedBasket, err := suite.app.BasketKeeper.GetBasketById(suite.ctx, basket.Id)
	suite.Require().NoError(err)
	suite.Require().Nil(savedBasket.Surplus)
}
