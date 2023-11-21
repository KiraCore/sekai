package types

import "github.com/cosmos/cosmos-sdk/types/errors"

// basket module errors
var (
	ErrBasketDoesNotExist             = errors.Register(ModuleName, 1, "basket not found")
	ErrMintsDisabledBasket            = errors.Register(ModuleName, 2, "mints disabled on the basket")
	ErrBurnsDisabledBasket            = errors.Register(ModuleName, 3, "burns disabled on the basket")
	ErrSwapsDisabledBasket            = errors.Register(ModuleName, 4, "swaps disabled on the basket")
	ErrInvalidBasketDepositDenom      = errors.Register(ModuleName, 5, "invalid basket deposit denom")
	ErrInvalidBasketWithdrawDenom     = errors.Register(ModuleName, 6, "invalid basket withdraw denom")
	ErrInvalidBasketDenom             = errors.Register(ModuleName, 7, "invalid basket denom")
	ErrNotAbleToWithdrawAnyTokens     = errors.Register(ModuleName, 8, "not able to withdraw any coins")
	ErrDepositsDisabledForToken       = errors.Register(ModuleName, 9, "deposits disabled for the token")
	ErrSwapsDisabledForInToken        = errors.Register(ModuleName, 10, "swap disabled for in token")
	ErrSwapsDisabledForOutToken       = errors.Register(ModuleName, 11, "swap disabled for out token")
	ErrInsufficientBasketDepositToken = errors.Register(ModuleName, 12, "insufficient token deposits on basket")
	ErrAmountBelowBaksetMintsMin      = errors.Register(ModuleName, 13, "mints amount is below basket minimum")
	ErrAmountAboveBaksetMintsMax      = errors.Register(ModuleName, 14, "mints amount is above basket maximum for the period")
	ErrAmountBelowBaksetBurnsMin      = errors.Register(ModuleName, 15, "burns amount is below basket minimum")
	ErrAmountAboveBaksetBurnsMax      = errors.Register(ModuleName, 16, "burns amount is above basket maximum for the period")
	ErrAmountBelowBaksetSwapsMin      = errors.Register(ModuleName, 17, "swaps amount is below basket minimum")
	ErrAmountAboveBaksetSwapsMax      = errors.Register(ModuleName, 18, "swaps amount is above basket maximum for the period")
	ErrTokenWeightShouldNotBeZero     = errors.Register(ModuleName, 19, "token weight should not be zero")
	ErrDuplicateDenomExistsOnTokens   = errors.Register(ModuleName, 20, "duplicated denom exists on tokens list")
	ErrTokenExceedingCap              = errors.Register(ModuleName, 21, "token exceeding cap")
	ErrEmptyUnderlyingTokens          = errors.Register(ModuleName, 22, "empty underlying tokens")
	ErrBasketDenomSupplyTooBig        = errors.Register(ModuleName, 23, "basket denom supply too big compared to underlying tokens")
)
