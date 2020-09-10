package simapp

import (
	"fmt"

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
		NewValidateFeeRangeDecorator(ak),
		ante.NewSetPubKeyDecorator(ak), // SetPubKeyDecorator must be called before all signature verification decorators
		ante.NewValidateSigCountDecorator(ak),
		ante.NewDeductFeeDecorator(ak, bankKeeper),
		ante.NewSigGasConsumeDecorator(ak, sigGasConsumer),
		ante.NewSigVerificationDecorator(ak, signModeHandler),
		ante.NewIncrementSequenceDecorator(ak),
	)
}

// ValidateFeeRangeDecorator check if fee is within range defined as network properties
type ValidateFeeRangeDecorator struct {
	ak keeper.AccountKeeper
}

// NewValidateFeeRangeDecorator check if fee is within range defined as network properties
func NewValidateFeeRangeDecorator(ak keeper.AccountKeeper) ValidateFeeRangeDecorator {
	return ValidateFeeRangeDecorator{
		ak: ak,
	}
}

// AnteHandle is a handler for ValidateFeeRangeDecorator
func (svd ValidateFeeRangeDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	sigTx, ok := tx.(types.StdTx)
	if !ok {
		return ctx, sdkerrors.Wrap(sdkerrors.ErrTxDecode, "invalid transaction type")
	}

	properties := CustomGovKeeper.GetNetworkProperties(ctx)

	// TODO should check if use fee.Amount or Fee.Gas
	bondDenom := StakingKeeper.BondDenom(ctx)
	feeAmount := sigTx.Fee.Amount.AmountOf(bondDenom).Uint64() // || sigTx.Fee.Gas
	if feeAmount < properties.MinTxFee || feeAmount > properties.MaxTxFee {
		return ctx, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("fee out of range [%d, %d]", properties.MinTxFee, properties.MaxTxFee))
	}

	return next(ctx, tx, simulate)
}
