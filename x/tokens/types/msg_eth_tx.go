package types

import (
	"errors"
	"math/big"

	errorsmod "cosmossdk.io/errors"
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

func (msg *MsgEthereumTx) FromEthereumTx(tx *ethtypes.Transaction) error {
	data, err := rlp.EncodeToBytes(tx)
	if err != nil {
		return err
	}
	msg.Data = data
	msg.Hash = tx.Hash().Hex()
	return nil
}

func (msg MsgEthereumTx) Route() string { return RouterKey }

func (msg MsgEthereumTx) Type() string { return kiratypes.MsgTypeEthereumTx }

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
