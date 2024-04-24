package types

import (
	"errors"
	"math/big"

	errorsmod "cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"

	kiratypes "github.com/KiraCore/sekai/types"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

var (
	_ sdk.Msg = &MsgEthereumTx{}
	_ sdk.Tx  = &MsgEthereumTx{}
)

// FromEthereumTx populates the message fields from the given ethereum transaction
func (msg *MsgEthereumTx) FromEthereumTx(tx *ethtypes.Transaction) error {
	data, err := rlp.EncodeToBytes(tx)
	if err != nil {
		return err
	}
	msg.Data = data
	msg.Hash = tx.Hash().Hex()
	return nil
}

// Route returns the route value of an MsgEthereumTx.
func (msg MsgEthereumTx) Route() string { return RouterKey }

// Type returns the type value of an MsgEthereumTx.
func (msg MsgEthereumTx) Type() string { return kiratypes.MsgTypeEthereumTx }

// ValidateBasic implements the sdk.Msg interface. It performs basic validation
// checks of a Transaction. If returns an error if validation fails.
func (msg MsgEthereumTx) ValidateBasic() error {
	if msg.Sender != "" {
		if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
			return errorsmod.Wrap(err, "invalid sender address")
		}
	}

	ethMsg, err := msg.AsMessage()
	if err != nil {
		return err
	}

	if !validateTx(msg.Data, ethMsg.From()) {
		return errors.New("validation failed")
	}

	return nil
}

// GetMsgs returns a single MsgEthereumTx as an sdk.Msg.
func (msg *MsgEthereumTx) GetMsgs() []sdk.Msg {
	return []sdk.Msg{msg}
}

// GetSigners returns the expected signers for an Ethereum transaction message.
// For such a message, there should exist only a single 'signer'.
//
// NOTE: This method panics if 'Sign' hasn't been called first.
func (msg *MsgEthereumTx) GetSigners() []sdk.AccAddress {
	signer := sdk.MustAccAddressFromBech32(msg.Sender)
	return []sdk.AccAddress{signer}
}

// GetSignBytes returns the Amino bytes of an Ethereum transaction message used
// for signing.
//
// NOTE: This method cannot be used as a chain ID is needed to create valid bytes
// to sign over. Use 'RLPSignBytes' instead.
func (msg MsgEthereumTx) GetSignBytes() []byte {
	panic("must use 'RLPSignBytes' with a chain ID to get the valid bytes to sign")
}

// Sign calculates a secp256k1 ECDSA signature and signs the transaction. It
// takes a keyring signer and the chainID to sign an Ethereum transaction according to
// EIP155 standard.
// This method mutates the transaction as it populates the V, R, S
// fields of the Transaction's Signature.
// The function will fail if the sender address is not defined for the msg or if
// the sender is not registered on the keyring
func (msg *MsgEthereumTx) Sign(ethSigner ethtypes.Signer, keyringSigner keyring.Signer) error {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return err
	}

	tx := msg.AsTransaction()
	txHash := ethSigner.Hash(tx)

	sig, _, err := keyringSigner.SignByAddress(sender, txHash.Bytes())
	if err != nil {
		return err
	}

	tx, err = tx.WithSignature(ethSigner, sig)
	if err != nil {
		return err
	}

	return msg.FromEthereumTx(tx)
}

func (msg *MsgEthereumTx) GetEthSender(chainID *big.Int) (common.Address, error) {
	signer := ethtypes.LatestSignerForChainID(chainID)
	from, err := signer.Sender(msg.AsTransaction())
	if err != nil {
		return common.Address{}, err
	}

	return from, nil
}

func (msg *MsgEthereumTx) AsMessage() (ethtypes.Message, error) {
	tx := msg.AsTransaction()
	return tx.AsMessage(ethtypes.NewEIP155Signer(tx.ChainId()), big.NewInt(1))
}

// AsTransaction creates an Ethereum Transaction type from the msg fields
func (msg MsgEthereumTx) AsTransaction() *ethtypes.Transaction {
	tx := new(ethtypes.Transaction)
	rlp.DecodeBytes(msg.Data, &tx)

	return tx
}

func GetSenderAddrFromRawTxBytes(rawTxBytes []byte) (common.Address, error) {
	var rawTx ethtypes.Transaction
	if err := rlp.DecodeBytes(rawTxBytes, &rawTx); err != nil {
		return common.Address{}, err
	}

	signer := ethtypes.NewEIP155Signer(rawTx.ChainId())
	sender, err := signer.Sender(&rawTx)
	if err != nil {
		return common.Address{}, err
	}
	return sender, nil
}

func validateTx(rawTxBytes []byte, sender common.Address) bool {
	senderFromTx, err := GetSenderAddrFromRawTxBytes(rawTxBytes)
	if err != nil {
		return false
	}

	if senderFromTx.Hex() == sender.Hex() {
		return true
	}
	return false
}
