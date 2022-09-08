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
	ErrSwapsDisabledForInToken        = errors.Register(ModuleName, 9, "swap disabled for in token")
	ErrSwapsDisabledForOutToken       = errors.Register(ModuleName, 10, "swap disabled for out token")
	ErrInsufficientBasketDepositToken = errors.Register(ModuleName, 11, "insufficient token deposits on basket")
)
