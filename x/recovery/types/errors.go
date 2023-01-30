package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/recovery module sentinel errors
var (
	ErrInvalidAccAddress           = sdkerrors.Register(ModuleName, 2, "invalid account address")
	ErrRecoveryRecordDoesNotExist  = sdkerrors.Register(ModuleName, 3, "recovery record does not exist")
	ErrRecoveryTokenDoesNotExist   = sdkerrors.Register(ModuleName, 4, "recovery token does not exist")
	ErrInvalidProof                = sdkerrors.Register(ModuleName, 5, "invalid proof")
	ErrRotatedAccountAlreadyExists = sdkerrors.Register(ModuleName, 6, "rotated account already exists")
	ErrAccountDoesNotExists        = sdkerrors.Register(ModuleName, 7, "account does not exist")
)
