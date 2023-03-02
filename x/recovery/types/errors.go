package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/recovery module sentinel errors
var (
	ErrInvalidAccAddress                      = sdkerrors.Register(ModuleName, 2, "invalid account address")
	ErrRecoveryRecordDoesNotExist             = sdkerrors.Register(ModuleName, 3, "recovery record does not exist")
	ErrRecoveryTokenDoesNotExist              = sdkerrors.Register(ModuleName, 4, "recovery token does not exist")
	ErrInvalidProof                           = sdkerrors.Register(ModuleName, 5, "invalid proof")
	ErrRotatedAccountAlreadyExists            = sdkerrors.Register(ModuleName, 6, "rotated account already exists")
	ErrAccountDoesNotExists                   = sdkerrors.Register(ModuleName, 7, "account does not exist")
	ErrInvalidMoniker                         = sdkerrors.Register(ModuleName, 8, "invalid moniker")
	ErrRecoveryTokenAlreadyExists             = sdkerrors.Register(ModuleName, 9, "recovery token already exists")
	ErrAddressHasValidatorRecoveryToken       = sdkerrors.Register(ModuleName, 10, "address already has validator recovery token")
	ErrNotEnoughRRTokenAmountForRotation      = sdkerrors.Register(ModuleName, 11, "not enough rr token amount for rotation")
	ErrTargetAddressAlreadyHasRotationHistory = sdkerrors.Register(ModuleName, 12, "target address already has rotation history")
)
