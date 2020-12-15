package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/slashing module sentinel errors
var (
	ErrNoValidatorForAddress = sdkerrors.Register(ModuleName, 2, "address is not associated with any known validator")
	ErrBadValidatorAddr      = sdkerrors.Register(ModuleName, 3, "validator does not exist for that address")
	ErrValidatorJailed       = sdkerrors.Register(ModuleName, 4, "validator still jailed; cannot be activated")
	ErrValidatorNotJailed    = sdkerrors.Register(ModuleName, 5, "validator not jailed; cannot be activated")
	ErrNoSigningInfoFound    = sdkerrors.Register(ModuleName, 8, "no validator signing info found")
)
