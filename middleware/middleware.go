package middleware

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
  customgovkeeper "github.com/KiraCore/sekai/x/gov/keeper"
  authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
  bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
  customstakingkeeper "github.com/KiraCore/sekai/x/staking/keeper"
  "errors"
)

var (
  customGovKeeper customgovkeeper.Keeper
	customStakingKeeper customstakingkeeper.Keeper
  bankKeeper bankkeeper.Keeper
)

// SetKeepers set keepers to be used on middlewares
func SetKeepers(cgk customgovkeeper.Keeper, csk customstakingkeeper.Keeper, bk bankkeeper.Keeper) {
  customGovKeeper = cgk
  customStakingKeeper = csk
  bankKeeper = bk
}

func MergeTwoErrors(hErr, err error) error {
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

    // on failure case, should pay
    if hErr == nil {
      ctx.GasMeter().ConsumeGas(fee.ExecutionFee, "consume execution fee")
      return hResult, hErr
    }
  
    if fee.FailureFee < fee.ExecutionFee { // should refund for ExecutionFee - FailureFee
      amount := int64(fee.ExecutionFee - fee.FailureFee)
      fees := sdk.Coins{sdk.NewInt64Coin(bondDenom, amount)}
      
      err := bankKeeper.SendCoinsFromModuleToAccount(ctx, authtypes.FeeCollectorName, feePayer, fees)
      if err != nil {
        hErr = MergeTwoErrors(hErr, err)
      }
    } 
    
    if fee.FailureFee > fee.ExecutionFee { // should pay more fee if handler fails
      amount := int64(fee.FailureFee - fee.ExecutionFee)
      fees := sdk.Coins{sdk.NewInt64Coin(bondDenom, amount)}
      
      err := bankKeeper.SendCoinsFromAccountToModule(ctx, feePayer, authtypes.FeeCollectorName, fees)
      if err != nil {
        hErr = MergeTwoErrors(hErr, err)
      }
    }

    ctx.GasMeter().ConsumeGas(fee.FailureFee, "consume execution failure fee")
    return hResult, hErr
	}
	return sdk.NewRoute(p, newHandler)
}
