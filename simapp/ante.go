package simapp

import (
	"fmt"

	customgovkeeper "github.com/KiraCore/sekai/x/gov/keeper"
	customstakingkeeper "github.com/KiraCore/sekai/x/staking/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	"github.com/cosmos/cosmos-sdk/x/auth/keeper"
	"github.com/cosmos/cosmos-sdk/x/auth/signing"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
)

// NewAnteHandler returns an AnteHandler that checks and increments sequence
// numbers, checks signatures & account numbers, and deducts fees from the first
// signer.
func NewAnteHandler(
	sk customstakingkeeper.Keeper,
	cgk customgovkeeper.Keeper,
	ak keeper.AccountKeeper,
	bankKeeper types.BankKeeper,
	sigGasConsumer ante.SignatureVerificationGasConsumer,
	signModeHandler signing.SignModeHandler,
) sdk.AnteHandler {
	return sdk.ChainAnteDecorators(
		ante.NewSetUpContextDecorator(), // outermost AnteDecorator. SetUpContext must be called first
		ante.NewRejectExtensionOptionsDecorator(),
		ante.NewMempoolFeeDecorator(),
		ante.NewValidateBasicDecorator(),
		ante.TxTimeoutHeightDecorator{},
		ante.NewValidateMemoDecorator(ak),
		ante.NewConsumeGasForTxSizeDecorator(ak),
		// custom fee range validator
		NewValidateFeeRangeDecorator(sk, cgk, ak),
		ante.NewSetPubKeyDecorator(ak), // SetPubKeyDecorator must be called before all signature verification decorators
		ante.NewValidateSigCountDecorator(ak),
		ante.NewDeductFeeDecorator(ak, bankKeeper),
		ante.NewSigGasConsumeDecorator(ak, sigGasConsumer),
		// custom execution fee consume decorator
		// NewCustomExecutionFeeConsumeDecorator(ak, cgk),
		ante.NewSigVerificationDecorator(ak, signModeHandler),
		ante.NewIncrementSequenceDecorator(ak),
	)
}

// ValidateFeeRangeDecorator check if fee is within range defined as network properties
type ValidateFeeRangeDecorator struct {
	sk  customstakingkeeper.Keeper
	cgk customgovkeeper.Keeper
	ak  keeper.AccountKeeper
}

// NewValidateFeeRangeDecorator check if fee is within range defined as network properties
func NewValidateFeeRangeDecorator(
	sk customstakingkeeper.Keeper,
	cgk customgovkeeper.Keeper,
	ak keeper.AccountKeeper,
) ValidateFeeRangeDecorator {
	return ValidateFeeRangeDecorator{
		sk:  sk,
		cgk: cgk,
		ak:  ak,
	}
}

// AnteHandle is a handler for ValidateFeeRangeDecorator
func (svd ValidateFeeRangeDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	feeTx, ok := tx.(sdk.FeeTx)
	if !ok {
		return ctx, sdkerrors.Wrap(sdkerrors.ErrTxDecode, "Tx must be a FeeTx")
	}

	properties := svd.cgk.GetNetworkProperties(ctx)

	bondDenom := svd.sk.BondDenom(ctx)
	feeAmount := feeTx.GetFee().AmountOf(bondDenom).Uint64()

	// execution failure fee should be prepaid
	executionFailureFee := uint64(0)
	for _, msg := range feeTx.GetMsgs() {
		executionName := msg.Type()
		fee := svd.cgk.GetExecutionFee(ctx, executionName)
		if fee != nil { // execution fee exist
			executionFailureFee += fee.FailureFee
		}
	}

	if feeAmount < properties.MinTxFee || feeAmount > properties.MaxTxFee {
		return ctx, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("fee out of range [%d, %d]", properties.MinTxFee, properties.MaxTxFee))
	}

	if feeAmount < executionFailureFee {
		return ctx, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("fee is less than execution failure fee %d", executionFailureFee))
	}

	return next(ctx, tx, simulate)
}

// // CustomExecutionFeeConsumeDecorator calculate custom gas consume
// type CustomExecutionFeeConsumeDecorator struct {
// 	ak  keeper.AccountKeeper
// 	cgk customgovkeeper.Keeper
// }

// // NewCustomExecutionFeeConsumeDecorator returns instance of CustomExecutionFeeConsumeDecorator
// func NewCustomExecutionFeeConsumeDecorator(ak keeper.AccountKeeper, cgk customgovkeeper.Keeper) CustomExecutionFeeConsumeDecorator {
// 	return CustomExecutionFeeConsumeDecorator{
// 		ak:  ak,
// 		cgk: cgk,
// 	}
// }

// // AnteHandle handle CustomExecutionFeeConsumeDecorator
// func (sgcd CustomExecutionFeeConsumeDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
// 	sigTx, ok := tx.(sdk.FeeTx)
// 	if !ok {
// 		return ctx, sdkerrors.Wrap(sdkerrors.ErrTxDecode, "invalid transaction type")
// 	}

// 	// execution fee consume gas
// 	for _, msg := range sigTx.GetMsgs() {
// 		executionName := msg.Type()
// 		fee := sgcd.cgk.GetExecutionFee(ctx, executionName)
// 		if fee != nil { // execution fee exist
// 			// TODO should check failure case and in that case should consume failure fee
// 			ctx.GasMeter().ConsumeGas(fee.ExecutionFee, "consume execution fee")
// 		}
// 	}

// 	return next(ctx, tx, simulate)
// }
