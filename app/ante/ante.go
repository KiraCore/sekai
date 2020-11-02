package ante

import (
	"fmt"

	feeprocessingkeeper "github.com/KiraCore/sekai/x/feeprocessing/keeper"
	customgovkeeper "github.com/KiraCore/sekai/x/gov/keeper"
	customstakingkeeper "github.com/KiraCore/sekai/x/staking/keeper"
	tokenskeeper "github.com/KiraCore/sekai/x/tokens/keeper"
	tokenstypes "github.com/KiraCore/sekai/x/tokens/types"
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
	tk tokenskeeper.Keeper,
	fk feeprocessingkeeper.Keeper,
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
		NewValidateFeeRangeDecorator(sk, cgk, tk, ak),
		ante.NewSetPubKeyDecorator(ak), // SetPubKeyDecorator must be called before all signature verification decorators
		ante.NewValidateSigCountDecorator(ak),
		ante.NewDeductFeeDecorator(ak, fk),
		// custom execution fee consume decorator
		NewExecutionFeeRegistrationDecorator(ak, cgk, fk),
		ante.NewSigGasConsumeDecorator(ak, sigGasConsumer),
		ante.NewSigVerificationDecorator(ak, signModeHandler),
		ante.NewIncrementSequenceDecorator(ak),
	)
}

// ValidateFeeRangeDecorator check if fee is within range defined as network properties
type ValidateFeeRangeDecorator struct {
	sk  customstakingkeeper.Keeper
	cgk customgovkeeper.Keeper
	tk  tokenskeeper.Keeper
	ak  keeper.AccountKeeper
}

// NewValidateFeeRangeDecorator check if fee is within range defined as network properties
func NewValidateFeeRangeDecorator(
	sk customstakingkeeper.Keeper,
	cgk customgovkeeper.Keeper,
	tk tokenskeeper.Keeper,
	ak keeper.AccountKeeper,
) ValidateFeeRangeDecorator {
	return ValidateFeeRangeDecorator{
		sk:  sk,
		cgk: cgk,
		ak:  ak,
		tk:  tk,
	}
}

// AnteHandle is a handler for ValidateFeeRangeDecorator
func (svd ValidateFeeRangeDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	feeTx, ok := tx.(sdk.FeeTx)
	if !ok {
		return ctx, sdkerrors.Wrap(sdkerrors.ErrTxDecode, "Tx must be a FeeTx")
	}

	properties := svd.cgk.GetNetworkProperties(ctx)

	feeAmount := uint64(0)
	feeCoins := feeTx.GetFee()
	for _, feeCoin := range feeCoins {
		rate := svd.tk.GetTokenRate(ctx, feeCoin.Denom)
		if rate == nil || !rate.FeePayments {
			return ctx, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("currency you are tying to use was not whitelisted as fee payment"))
		}
		// TODO it might be required to use safemath for this operation, in case of user set too much fee which can overflow 10^19
		feeAmount += uint64(feeCoin.Amount.Int64()) * rate.Rate / uint64(tokenstypes.RateDecimalDenominator)
	}

	// execution fee should be prepaid
	executionMaxFee := uint64(0)
	for _, msg := range feeTx.GetMsgs() {
		fee := svd.cgk.GetExecutionFee(ctx, msg.Type())
		if fee != nil { // execution fee exist
			maxFee := fee.FailureFee
			if fee.ExecutionFee > maxFee {
				maxFee = fee.ExecutionFee
			}
			executionMaxFee += maxFee
		}
	}

	bondDenom := svd.sk.BondDenom(ctx)
	if feeAmount < properties.MinTxFee || feeAmount > properties.MaxTxFee {
		return ctx, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("fee %+v(%d) is out of range [%d, %d]%s", feeTx.GetFee(), feeAmount, properties.MinTxFee, properties.MaxTxFee, bondDenom))
	}

	if feeAmount < executionMaxFee {
		return ctx, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("fee %+v(%d) is less than max execution fee %d%s", feeTx.GetFee(), feeAmount, executionMaxFee, bondDenom))
	}

	return next(ctx, tx, simulate)
}

// ExecutionFeeRegistrationDecorator register paid execution fee
type ExecutionFeeRegistrationDecorator struct {
	ak  keeper.AccountKeeper
	cgk customgovkeeper.Keeper
	fk  feeprocessingkeeper.Keeper
}

// NewExecutionFeeRegistrationDecorator returns instance of CustomExecutionFeeConsumeDecorator
func NewExecutionFeeRegistrationDecorator(ak keeper.AccountKeeper, cgk customgovkeeper.Keeper, fk feeprocessingkeeper.Keeper) ExecutionFeeRegistrationDecorator {
	return ExecutionFeeRegistrationDecorator{
		ak:  ak,
		cgk: cgk,
		fk:  fk,
	}
}

// AnteHandle handle ExecutionFeeRegistrationDecorator
func (sgcd ExecutionFeeRegistrationDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	sigTx, ok := tx.(sdk.FeeTx)
	if !ok {
		return ctx, sdkerrors.Wrap(sdkerrors.ErrTxDecode, "invalid transaction type")
	}

	// execution fee consume gas
	for _, msg := range sigTx.GetMsgs() {
		fee := sgcd.cgk.GetExecutionFee(ctx, msg.Type())
		if fee != nil { // execution fee exist
			sgcd.fk.AddExecutionStart(ctx, msg)
		}
	}

	return next(ctx, tx, simulate)
}
