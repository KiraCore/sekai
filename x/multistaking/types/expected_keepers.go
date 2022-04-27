package types

import (
	stakingtypes "github.com/KiraCore/sekai/x/staking/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// StakingKeeper expected staking keeper
type StakingKeeper interface {
	GetValidator(sdk.Context, sdk.ValAddress) (stakingtypes.Validator, error)
}
