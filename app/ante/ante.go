package ante

import (
	"fmt"

	kiratypes "github.com/KiraCore/sekai/types"
	feeprocessingkeeper "github.com/KiraCore/sekai/x/feeprocessing/keeper"
	feeprocessingtypes "github.com/KiraCore/sekai/x/feeprocessing/types"
	customgovkeeper "github.com/KiraCore/sekai/x/gov/keeper"
	customstakingkeeper "github.com/KiraCore/sekai/x/staking/keeper"
	tokenskeeper "github.com/KiraCore/sekai/x/tokens/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	"github.com/cosmos/cosmos-sdk/x/auth/keeper"
	"github.com/cosmos/cosmos-sdk/x/auth/signing"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	bank "github.com/cosmos/cosmos-sdk/x/bank/types"
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
	bk types.BankKeeper,
	sigGasConsumer ante.SignatureVerificationGasConsumer,
	signModeHandler signing.SignModeHandler,
) sdk.AnteHandler {
	return sdk.ChainAnteDecorators(
		ante.NewSetUpContextDecorator(), // outermost AnteDecorator. SetUpContext must be called first
		NewZeroGasMeterDecorator(),
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
		ante.NewDeductFeeDecorator(ak, bk, nil),
		// poor network management decorator
		NewPoorNetworkManagementDecorator(ak, cgk, sk),
		NewBlackWhiteTokensCheckDecorator(cgk, sk, tk),
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
	bondDenom := svd.sk.BondDenom(ctx)

	feeAmount := sdk.NewDec(0)
	feeCoins := feeTx.GetFee()
	tokensBlackWhite := svd.tk.GetTokenBlackWhites(ctx)
	for _, feeCoin := range feeCoins {
		rate := svd.tk.GetTokenRate(ctx, feeCoin.Denom)
		if !properties.EnableForeignFeePayments && feeCoin.Denom != bondDenom {
			return ctx, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("foreign fee payments is disabled by governance"))
		}
		if rate == nil || !rate.FeePayments {
			return ctx, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("currency you are trying to use was not whitelisted as fee payment"))
		}
		if tokensBlackWhite.IsFrozen(feeCoin.Denom, bondDenom, properties.EnableTokenBlacklist, properties.EnableTokenWhitelist) {
			return ctx, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("currency you are trying to use as fee is frozen"))
		}
		feeAmount = feeAmount.Add(feeCoin.Amount.ToDec().Mul(rate.Rate))
	}

	// execution fee should be prepaid
	executionMaxFee := uint64(0)
	for _, msg := range feeTx.GetMsgs() {
		fee := svd.cgk.GetExecutionFee(ctx, kiratypes.MsgType(msg))
		if fee != nil { // execution fee exist
			maxFee := fee.FailureFee
			if fee.ExecutionFee > maxFee {
				maxFee = fee.ExecutionFee
			}
			executionMaxFee += maxFee
		}
	}

	if feeAmount.LT(sdk.NewDec(int64(properties.MinTxFee))) || feeAmount.GT(sdk.NewDec(int64(properties.MaxTxFee))) {
		return ctx, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("fee %+v(%d) is out of range [%d, %d]%s", feeTx.GetFee(), feeAmount.RoundInt().Int64(), properties.MinTxFee, properties.MaxTxFee, bondDenom))
	}

	if feeAmount.LT(sdk.NewDec(int64(executionMaxFee))) {
		return ctx, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("fee %+v(%d) is less than max execution fee %d%s", feeTx.GetFee(), feeAmount.RoundInt().Int64(), executionMaxFee, bondDenom))
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
		ak,
		cgk,
		fk,
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
		fee := sgcd.cgk.GetExecutionFee(ctx, kiratypes.MsgType(msg))
		if fee != nil { // execution fee exist
			sgcd.fk.AddExecutionStart(ctx, msg)
		}
	}

	return next(ctx, tx, simulate)
}

// PoorNetworkManagementDecorator register poor network manager
type PoorNetworkManagementDecorator struct {
	ak  keeper.AccountKeeper
	cgk customgovkeeper.Keeper
	csk customstakingkeeper.Keeper
}

// NewPoorNetworkManagementDecorator returns instance of PoorNetworkManagementDecorator
func NewPoorNetworkManagementDecorator(ak keeper.AccountKeeper, cgk customgovkeeper.Keeper, csk customstakingkeeper.Keeper) PoorNetworkManagementDecorator {
	return PoorNetworkManagementDecorator{
		ak,
		cgk,
		csk,
	}
}

func findString(a []string, x string) int {
	for i, n := range a {
		if x == n {
			return i
		}
	}
	return -1
}

// AnteHandle handle PoorNetworkManagementDecorator
func (pnmd PoorNetworkManagementDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	sigTx, ok := tx.(sdk.FeeTx)
	if !ok {
		return ctx, sdkerrors.Wrap(sdkerrors.ErrTxDecode, "invalid transaction type")
	}

	// if not poor network, skip this process
	if pnmd.csk.IsNetworkActive(ctx) {
		return next(ctx, tx, simulate)
	}
	// handle messages on poor network
	pnmsgs := pnmd.cgk.GetPoorNetworkMessages(ctx)
	for _, msg := range sigTx.GetMsgs() {
		if kiratypes.MsgType(msg) == bank.TypeMsgSend {
			// on poor network, we introduce POOR_NETWORK_MAX_BANK_TX_SEND network property to limit transaction send amount
			msg := msg.(*bank.MsgSend)
			if len(msg.Amount) > 1 || msg.Amount[0].Denom != pnmd.csk.BondDenom(ctx) {
				return ctx, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "only bond denom is allowed on poor network")
			}
			if msg.Amount[0].Amount.Uint64() > pnmd.cgk.GetNetworkProperties(ctx).PoorNetworkMaxBankSend {
				return ctx, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "only restricted amount send is allowed on poor network")
			}
			// TODO: we could do restriction to send only when target account does not exist on chain yet for more restriction
			return next(ctx, tx, simulate)
		}
		if findString(pnmsgs.Messages, kiratypes.MsgType(msg)) >= 0 {
			return next(ctx, tx, simulate)
		}
		return ctx, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid transaction type on poor network")
	}

	return next(ctx, tx, simulate)
}

// BlackWhiteTokensCheckDecorator register black white tokens check decorator
type BlackWhiteTokensCheckDecorator struct {
	cgk customgovkeeper.Keeper
	csk customstakingkeeper.Keeper
	tk  tokenskeeper.Keeper
}

// NewBlackWhiteTokensCheckDecorator returns instance of BlackWhiteTokensCheckDecorator
func NewBlackWhiteTokensCheckDecorator(cgk customgovkeeper.Keeper, csk customstakingkeeper.Keeper, tk tokenskeeper.Keeper) BlackWhiteTokensCheckDecorator {
	return BlackWhiteTokensCheckDecorator{
		cgk,
		csk,
		tk,
	}
}

// AnteHandle handle NewPoorNetworkManagementDecorator
func (pnmd BlackWhiteTokensCheckDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	sigTx, ok := tx.(sdk.FeeTx)
	if !ok {
		return ctx, sdkerrors.Wrap(sdkerrors.ErrTxDecode, "invalid transaction type")
	}

	bondDenom := pnmd.csk.BondDenom(ctx)
	tokensBlackWhite := pnmd.tk.GetTokenBlackWhites(ctx)
	properties := pnmd.cgk.GetNetworkProperties(ctx)
	for _, msg := range sigTx.GetMsgs() {
		if kiratypes.MsgType(msg) == bank.TypeMsgSend {
			msg := msg.(*bank.MsgSend)
			for _, amt := range msg.Amount {
				if tokensBlackWhite.IsFrozen(amt.Denom, bondDenom, properties.EnableTokenBlacklist, properties.EnableTokenWhitelist) {
					return ctx, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "token is frozen")
				}
			}
		}
	}

	return next(ctx, tx, simulate)
}

// ZeroGasMeterDecorator uses infinite gas decorator to avoid gas usage in the network
type ZeroGasMeterDecorator struct{}

// NewZeroGasMeterDecorator returns instance of ZeroGasMeterDecorator
func NewZeroGasMeterDecorator() ZeroGasMeterDecorator {
	return ZeroGasMeterDecorator{}
}

func (igm ZeroGasMeterDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	newCtx = ctx.WithGasMeter(feeprocessingtypes.NewZeroGasMeter())
	return next(newCtx, tx, simulate)
}
