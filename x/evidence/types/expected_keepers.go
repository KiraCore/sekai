package types

import (
	"time"

	stakingtypes "github.com/KiraCore/sekai/x/staking/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type (
	// StakingKeeper defines the staking module interface contract needed by the
	// evidence module.
	StakingKeeper interface {
		GetValidatorByConsAddr(sdk.Context, sdk.ConsAddress) (stakingtypes.Validator, error)
	}

	// SlashingKeeper defines the slashing module interface contract needed by the
	// evidence module.
	SlashingKeeper interface {
		GetPubkey(sdk.Context, cryptotypes.Address) (cryptotypes.PubKey, error)
		HasValidatorSigningInfo(sdk.Context, sdk.ConsAddress) bool
		Jail(sdk.Context, sdk.ConsAddress)
		JailUntil(sdk.Context, sdk.ConsAddress, time.Time)
	}
)
