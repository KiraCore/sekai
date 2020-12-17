package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Activate a validator
func (k Keeper) Activate(ctx sdk.Context, consAddr sdk.ConsAddress) {
	// TODO: don't do anything for now, implement this
	// validator := k.GetValidatorByConsAddr(ctx, consAddr)
	// if !validator.Inactivated {
	// 	panic(fmt.Sprintf("cannot unjail already inactivated validator, validator: %v\n", validator))
	// }

	// validator.Inactivated = false
	// k.SetValidator(ctx, validator)
}

// slash the validator and delegators of the validator, specifying offence height, offence power, and slash fraction
func (k Keeper) Inactivate(sdk.Context, sdk.ConsAddress) { // inactivate a validator
	// TODO: don't do anything for now, implement this
}
