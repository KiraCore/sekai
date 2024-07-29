package ante

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"

	kiratypes "github.com/KiraCore/sekai/types"
	tokenstypes "github.com/KiraCore/sekai/x/tokens/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	sdktx "github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	ante "github.com/cosmos/cosmos-sdk/x/auth/ante"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/gogo/protobuf/proto"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/ethereum/go-ethereum/common"
	ethmath "github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
	apitypes "github.com/ethereum/go-ethereum/signer/core/apitypes"
)

var (
	// simulation signature values used to estimate gas consumption
	key                = make([]byte, secp256k1.PubKeySize)
	simSecp256k1Pubkey = &secp256k1.PubKey{Key: key}
	EthChainID         = uint64(8789)
)

func init() {
	// This decodes a valid hex string into a sepc256k1Pubkey for use in transaction simulation
	bz, _ := hex.DecodeString("035AD6810A47F073553FF30D2FCC7E0D3B1C0B74B61A1AAA2582344037151E143A")
	copy(key, bz)
	simSecp256k1Pubkey.Key = key
}

// SetPubKeyDecorator sets PubKeys in context for any signer which does not already have pubkey set
// PubKeys must be set in context for all signers before any other sigverify decorators run
// CONTRACT: Tx must implement SigVerifiableTx interface
type SetPubKeyDecorator struct {
	ak ante.AccountKeeper
}

func NewSetPubKeyDecorator(ak ante.AccountKeeper) SetPubKeyDecorator {
	return SetPubKeyDecorator{
		ak: ak,
	}
}

func (spkd SetPubKeyDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	sigTx, ok := tx.(authsigning.SigVerifiableTx)
	if !ok {
		return ctx, sdkerrors.Wrap(sdkerrors.ErrTxDecode, "invalid tx type")
	}

	pubkeys, err := sigTx.GetPubKeys()
	if err != nil {
		return ctx, err
	}
	signers := sigTx.GetSigners()

	for i, pk := range pubkeys {
		// PublicKey was omitted from slice since it has already been set in context
		if pk == nil {
			if !simulate {
				continue
			}
			pk = simSecp256k1Pubkey
		}
		// Only make check if simulate=false
		// if !simulate && !bytes.Equal(pk.Address(), signers[i]) {
		// 	return ctx, sdkerrors.Wrapf(sdkerrors.ErrInvalidPubKey,
		// 		"pubKey does not match signer address %s with signer index: %d", signers[i], i)
		// }

		acc, err := GetSignerAcc(ctx, spkd.ak, signers[i])
		if err != nil {
			return ctx, err
		}
		// account already has pubkey set,no need to reset
		if acc.GetPubKey() != nil {
			continue
		}
		err = acc.SetPubKey(pk)
		if err != nil {
			return ctx, sdkerrors.Wrap(sdkerrors.ErrInvalidPubKey, err.Error())
		}
		spkd.ak.SetAccount(ctx, acc)
	}

	// Also emit the following events, so that txs can be indexed by these
	// indices:
	// - signature (via `tx.signature='<sig_as_base64>'`),
	// - concat(address,"/",sequence) (via `tx.acc_seq='cosmos1abc...def/42'`).
	sigs, err := sigTx.GetSignaturesV2()
	if err != nil {
		return ctx, err
	}

	var events sdk.Events
	for i, sig := range sigs {
		events = append(events, sdk.NewEvent(sdk.EventTypeTx,
			sdk.NewAttribute(sdk.AttributeKeyAccountSequence, fmt.Sprintf("%s/%d", signers[i], sig.Sequence)),
		))

		sigBzs, err := signatureDataToBz(sig.Data)
		if err != nil {
			return ctx, err
		}
		for _, sigBz := range sigBzs {
			events = append(events, sdk.NewEvent(sdk.EventTypeTx,
				sdk.NewAttribute(sdk.AttributeKeySignature, base64.StdEncoding.EncodeToString(sigBz)),
			))
		}
	}

	ctx.EventManager().EmitEvents(events)

	return next(ctx, tx, simulate)
}

// SignatureVerificationGasConsumer is the type of function that is used to both
// consume gas when verifying signatures and also to accept or reject different types of pubkeys
// This is where apps can define their own PubKey
type SignatureVerificationGasConsumer = func(meter sdk.GasMeter, sig signing.SignatureV2, params types.Params) error

// Verify all signatures for a tx and return an error if any are invalid. Note,
//
// CONTRACT: Pubkeys are set in context for all signers before this decorator runs
// CONTRACT: Tx must implement SigVerifiableTx interface
type SigVerificationDecorator struct {
	ak                ante.AccountKeeper
	signModeHandler   authsigning.SignModeHandler
	interfaceRegistry codectypes.InterfaceRegistry
}

func NewSigVerificationDecorator(ak ante.AccountKeeper, signModeHandler authsigning.SignModeHandler, interfaceRegistry codectypes.InterfaceRegistry) SigVerificationDecorator {
	return SigVerificationDecorator{
		ak:                ak,
		signModeHandler:   signModeHandler,
		interfaceRegistry: interfaceRegistry,
	}
}

// OnlyLegacyAminoSigners checks SignatureData to see if all
// signers are using SIGN_MODE_LEGACY_AMINO_JSON. If this is the case
// then the corresponding SignatureV2 struct will not have account sequence
// explicitly set, and we should skip the explicit verification of sig.Sequence
// in the SigVerificationDecorator's AnteHandler function.
func OnlyLegacyAminoSigners(sigData signing.SignatureData) bool {
	switch v := sigData.(type) {
	case *signing.SingleSignatureData:
		return v.SignMode == signing.SignMode_SIGN_MODE_LEGACY_AMINO_JSON
	case *signing.MultiSignatureData:
		for _, s := range v.Signatures {
			if !OnlyLegacyAminoSigners(s) {
				return false
			}
		}
		return true
	default:
		return false
	}
}

func (svd SigVerificationDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	sigTx, ok := tx.(authsigning.SigVerifiableTx)
	if !ok {
		return ctx, sdkerrors.Wrap(sdkerrors.ErrTxDecode, "invalid transaction type")
	}

	// stdSigs contains the sequence number, account number, and signatures.
	// When simulating, this would just be a 0-length slice.
	sigs, err := sigTx.GetSignaturesV2()
	if err != nil {
		return ctx, err
	}

	signerAddrs := sigTx.GetSigners()

	// check that signer length and signature length are the same
	if len(sigs) != len(signerAddrs) {
		return ctx, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "invalid number of signer;  expected: %d, got %d", len(signerAddrs), len(sigs))
	}

	for i, sig := range sigs {
		acc, err := GetSignerAcc(ctx, svd.ak, signerAddrs[i])
		if err != nil {
			return ctx, err
		}

		// retrieve pubkey
		pubKey := acc.GetPubKey()
		if !simulate && pubKey == nil {
			return ctx, sdkerrors.Wrap(sdkerrors.ErrInvalidPubKey, "pubkey on account is not set")
		}

		// Check account sequence number.
		if sig.Sequence != acc.GetSequence() {
			return ctx, sdkerrors.Wrapf(
				sdkerrors.ErrWrongSequence,
				"account sequence mismatch, expected %d, got %d", acc.GetSequence(), sig.Sequence,
			)
		}

		// retrieve signer data
		genesis := ctx.BlockHeight() == 0
		chainID := ctx.ChainID()
		var accNum uint64
		if !genesis {
			accNum = acc.GetAccountNumber()
		}

		signerData := authsigning.SignerData{
			Address:       acc.GetAddress().String(),
			ChainID:       chainID,
			AccountNumber: accNum,
			Sequence:      acc.GetSequence(),
			PubKey:        pubKey,
		}

		if !simulate {
			if !bytes.Equal(pubKey.Address(), acc.GetAddress()) {
				// try verifying ethereum signature
				if ethErr := VerifyEthereumSignature(pubKey, signerData, sig.Data, svd.signModeHandler, tx, svd.interfaceRegistry); ethErr == nil {
					return next(ctx, tx, simulate)
				} else {
					errMsg := fmt.Sprintf("ethereum signature verification failed; %s", ethErr.Error())
					return ctx, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, errMsg)
				}
			}

			err := authsigning.VerifySignature(pubKey, signerData, sig.Data, svd.signModeHandler, tx)
			if err != nil {
				var errMsg string
				if OnlyLegacyAminoSigners(sig.Data) {
					// If all signers are using SIGN_MODE_LEGACY_AMINO, we rely on VerifySignature to check account sequence number,
					// and therefore communicate sequence number as a potential cause of error.
					errMsg = fmt.Sprintf("signature verification failed; please verify account number (%d), sequence (%d) and chain-id (%s)", accNum, acc.GetSequence(), chainID)
				} else {
					errMsg = fmt.Sprintf("signature verification failed; please verify account number (%d) and chain-id (%s)", accNum, chainID)
				}
				return ctx, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, errMsg)
			}
		}
	}

	return next(ctx, tx, simulate)
}

func VerifyEthereumSignature(pubKey cryptotypes.PubKey, signerData authsigning.SignerData, sigData signing.SignatureData, handler authsigning.SignModeHandler, tx sdk.Tx, interfaceRegistry codectypes.InterfaceRegistry) error {
	switch data := sigData.(type) {
	case *signing.SingleSignatureData:
		signBytes, err := handler.GetSignBytes(data.SignMode, signerData, tx)
		if err != nil {
			return err
		}

		if data.SignMode == signing.SignMode_SIGN_MODE_DIRECT {
			signDoc := sdktx.SignDoc{}
			err = signDoc.Unmarshal(signBytes)
			if err != nil {
				return err
			}

			txBody := &sdktx.TxBody{}
			err = proto.Unmarshal(signDoc.BodyBytes, txBody)
			if err != nil {
				return err
			}

			if len(txBody.Messages) != 1 {
				return errors.New("only one message's enabled for EIP712 signature")
			}
			msg := txBody.Messages[0]
			err := txBody.UnpackInterfaces(interfaceRegistry)
			if err != nil {
				return err
			}

			switch msg.TypeUrl {
			case "/kira.tokens.MsgEthereumTx":
				msg := msg.GetCachedValue().(*tokenstypes.MsgEthereumTx)
				tx := msg.AsTransaction()
				if signerData.Sequence != tx.Nonce() {
					return sdkerrors.Wrapf(
						sdkerrors.ErrWrongSequence,
						"ethereum account sequence mismatch, expected %d, got %d", signerData.Sequence, tx.Nonce(),
					)
				}

				if EthChainID != tx.ChainId().Uint64() {
					return fmt.Errorf("invalid ethereum chain ID, expected %d, got %d", EthChainID, tx.ChainId().Uint64())
				}
				return msg.ValidateBasic()
			default:
				msg := msg.GetCachedValue().(sdk.Msg)
				signBytes, err = GenEIP712SignBytesFromMsg(msg, signerData.Sequence)
				if err != nil {
					return err
				}
			}
		}

		signatureData := data.Signature
		if len(signatureData) <= crypto.RecoveryIDOffset {
			return fmt.Errorf("not a correct ethereum signature")
		}
		signatureData[crypto.RecoveryIDOffset] -= 27 // Transform yellow paper V from 27/28 to 0/1
		recovered, err := crypto.SigToPub(signBytes, signatureData)
		if err != nil {
			return err
		}
		recoveredAddr := crypto.PubkeyToAddress(*recovered)

		sigTx, ok := tx.(authsigning.SigVerifiableTx)
		if !ok {
			return sdkerrors.Wrap(sdkerrors.ErrTxDecode, "invalid transaction type")
		}

		signerAddrs := sigTx.GetSigners()
		if len(signerAddrs) != 1 {
			return fmt.Errorf("only 1 signer transaction supported: got %d", len(signerAddrs))
		}

		address := signerAddrs[0]

		if recoveredAddr.String() != common.BytesToAddress(address).String() {
			return fmt.Errorf("mismatching recovered address and sender: %s != %s", recoveredAddr.String(), common.BytesToAddress(address).String())
		}
		return nil
	default:
		return fmt.Errorf("unexpected SignatureData %T", sigData)
	}
}

func GenEIP712SignBytesFromMsg(msg sdk.Msg, nonce uint64) ([]byte, error) {
	msgName := kiratypes.MsgType(msg)
	msgData, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	types := apitypes.Types{
		"EIP712Domain": {
			{Name: "name", Type: "string"},
			{Name: "version", Type: "string"},
			{Name: "chainId", Type: "uint256"},
		},
		msgName: {
			{Name: "param", Type: "string"},
			{Name: "nonce", Type: "uint256"},
		},
	}

	chainId := ethmath.NewHexOrDecimal256(int64(EthChainID))
	nonceEth := ethmath.NewHexOrDecimal256(int64(nonce))

	domain := apitypes.TypedDataDomain{
		Name:    "Kira",
		Version: "1",
		ChainId: chainId,
		// VerifyingContract: "0xdef1c0ded9bec7f1a1670819833240f027b25eff",
	}

	typedData := apitypes.TypedData{
		Types:       types,
		PrimaryType: msgName,
		Domain:      domain,
		Message: apitypes.TypedDataMessage{
			"param": string(msgData),
			"nonce": nonceEth,
		},
	}

	typedDataHash, err := typedData.HashStruct(typedData.PrimaryType, typedData.Message)
	if err != nil {
		return nil, err
	}

	domainSeparator, err := typedData.HashStruct("EIP712Domain", typedData.Domain.Map())
	if err != nil {
		return nil, err
	}

	rawData := []byte(fmt.Sprintf("\x19\x01%s%s", string(domainSeparator), string(typedDataHash)))
	hashBytes := crypto.Keccak256(rawData)

	return hashBytes, nil
}

// GetSignerAcc returns an account for a given address that is expected to sign
// a transaction.
func GetSignerAcc(ctx sdk.Context, ak ante.AccountKeeper, addr sdk.AccAddress) (types.AccountI, error) {
	if acc := ak.GetAccount(ctx, addr); acc != nil {
		return acc, nil
	}

	return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress, "account %s does not exist", addr)
}

// signatureDataToBz converts a SignatureData into raw bytes signature.
// For SingleSignatureData, it returns the signature raw bytes.
// For MultiSignatureData, it returns an array of all individual signatures,
// as well as the aggregated signature.
func signatureDataToBz(data signing.SignatureData) ([][]byte, error) {
	if data == nil {
		return nil, fmt.Errorf("got empty SignatureData")
	}

	switch data := data.(type) {
	case *signing.SingleSignatureData:
		return [][]byte{data.Signature}, nil
	case *signing.MultiSignatureData:
		sigs := [][]byte{}
		var err error

		for _, d := range data.Signatures {
			nestedSigs, err := signatureDataToBz(d)
			if err != nil {
				return nil, err
			}
			sigs = append(sigs, nestedSigs...)
		}

		multisig := cryptotypes.MultiSignature{
			Signatures: sigs,
		}
		aggregatedSig, err := multisig.Marshal()
		if err != nil {
			return nil, err
		}
		sigs = append(sigs, aggregatedSig)

		return sigs, nil
	default:
		return nil, sdkerrors.ErrInvalidType.Wrapf("unexpected signature data type %T", data)
	}
}
