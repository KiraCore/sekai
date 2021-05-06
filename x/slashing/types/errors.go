package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/slashing module sentinel errors
var (
	ErrNoValidatorForAddress   = sdkerrors.Register(ModuleName, 2, "address is not associated with any known validator")
	ErrBadValidatorAddr        = sdkerrors.Register(ModuleName, 3, "validator does not exist for that address")
	ErrValidatorJailed         = sdkerrors.Register(ModuleName, 4, "validator jailed")
	ErrValidatorInactivated    = sdkerrors.Register(ModuleName, 5, "validator still inactivated; cannot be activated")
	ErrValidatorNotInactivated = sdkerrors.Register(ModuleName, 6, "validator not inactivated; cannot be activated")
	ErrValidatorNotPaused      = sdkerrors.Register(ModuleName, 7, "validator not paused")
	ErrValidatorPaused         = sdkerrors.Register(ModuleName, 8, "validator paused")
	ErrNoSigningInfoFound      = sdkerrors.Register(ModuleName, 9, "no validator signing info found")
)
