package middleware

import (
	"errors"

	customgovkeeper "github.com/KiraCore/sekai/x/gov/keeper"
	customstakingkeeper "github.com/KiraCore/sekai/x/staking/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
)

var (
	customGovKeeper     customgovkeeper.Keeper
	customStakingKeeper customstakingkeeper.Keeper
	bankKeeper          bankkeeper.Keeper
)

// SetKeepers set keepers to be used on middlewares
func SetKeepers(cgk customgovkeeper.Keeper, csk customstakingkeeper.Keeper, bk bankkeeper.Keeper) {
	customGovKeeper = cgk
	customStakingKeeper = csk
	bankKeeper = bk
}

// combineTwoErrors combines two error into one to show two errors in a single error text
func combineTwoErrors(hErr, err error) error {
	if hErr == nil {
		return err
	}
	if err == nil {
		return hErr
	}
	return errors.New(hErr.Error() + ";" + err.Error())
}

// NewRoute returns an instance of Route.
func NewRoute(p string, h sdk.Handler) sdk.Route {
	newHandler := func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		hResult, hErr := h(ctx, msg)
		// handle extra fee based on handler result
		executionName := msg.Type()
		feePayer := msg.GetSigners()[0] // signers listing should be at least 1 always
		bondDenom := customStakingKeeper.BondDenom(ctx)

		fee := customGovKeeper.GetExecutionFee(ctx, executionName)
		if fee == nil {
			return hResult, hErr
		}

		if hErr != nil {
			// TODO this can be something that's not needed as this modified ctx won't be used for tx failure
			ctx.GasMeter().ConsumeGas(fee.FailureFee, "consume execution failure fee")
			return hResult, hErr
		}

		if fee.FailureFee < fee.ExecutionFee { // should pay extra fee
			amount := int64(fee.ExecutionFee - fee.FailureFee)
			fees := sdk.Coins{sdk.NewInt64Coin(bondDenom, amount)}
			err := bankKeeper.SendCoinsFromAccountToModule(ctx, feePayer, authtypes.FeeCollectorName, fees)

			if err != nil {
				// TODO this can be something that's not needed as this modified ctx won't be used for tx failure
				ctx.GasMeter().ConsumeGas(fee.FailureFee, "consume execution failure fee")
				return hResult, err
			}
		}

		if fee.FailureFee > fee.ExecutionFee { // should return risk fee on success
			amount := int64(fee.FailureFee - fee.ExecutionFee)
			fees := sdk.Coins{sdk.NewInt64Coin(bondDenom, amount)}
			err := bankKeeper.SendCoinsFromModuleToAccount(ctx, authtypes.FeeCollectorName, feePayer, fees)

			if err != nil {
				// TODO this can be something that's not needed as this modified ctx won't be used for tx failure
				ctx.GasMeter().ConsumeGas(fee.FailureFee, "consume execution failure fee")
				return hResult, err
			}
		}

		ctx.GasMeter().ConsumeGas(fee.ExecutionFee, "consume execution failure fee")
		return hResult, hErr
	}
	return sdk.NewRoute(p, newHandler)
}
