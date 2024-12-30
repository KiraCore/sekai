package keeper

import (
	"bytes"
	"context"
	"cosmossdk.io/errors"
	"encoding/hex"
	"github.com/KiraCore/sekai/x/ethereum/types"
	"github.com/armon/go-metrics"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	types2 "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/gogoproto/proto"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
	"strconv"
)

type msgServer struct {
	keeper Keeper
	cgk    types.CustomGovKeeper
	bk     types.BankKeeper
}

// NewMsgServerImpl returns an implementation of the bank MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper, cgk types.CustomGovKeeper, bk types.BankKeeper) types.MsgServer {
	return &msgServer{
		keeper: keeper,
		cgk:    cgk,
		bk:     bk,
	}
}

var _ types.MsgServer = msgServer{}

func (m msgServer) Relay(goCtx context.Context, relay *types.MsgRelay) (*types.MsgRelayResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	txData, err := hex.DecodeString(relay.Data)
	if err != nil {
		return &types.MsgRelayResponse{}, errors.Wrap(types.ErrEthTxNotValid, err.Error())
	}

	var tx = new(types.EVMTx)
	err = proto.Unmarshal(txData, tx)
	if err != nil {
		return &types.MsgRelayResponse{}, errors.Wrap(types.ErrEthTxNotValid, err.Error())
	}

	ethFrom := common.HexToAddress(tx.From)
	ethTo := common.HexToAddress(tx.To)

	value := new(big.Int)
	value.SetString(tx.Value, 10)

	gas, err := strconv.Atoi(tx.Gas)
	if err != nil {
		return &types.MsgRelayResponse{}, errors.Wrap(types.ErrEthTxNotValid, err.Error())
	}

	gasPrice := new(big.Int)
	gasPrice.SetString(tx.GasPrice, 10)

	nonce, err := strconv.Atoi(tx.Nonce)
	if err != nil {
		return &types.MsgRelayResponse{}, errors.Wrap(types.ErrEthTxNotValid, err.Error())
	}

	data, _ := hex.DecodeString(tx.Data[2:])
	chainID := big.NewInt(int64(tx.ChainId))

	rBytes, err := hex.DecodeString(tx.R[2:])
	if err != nil {
		return &types.MsgRelayResponse{}, errors.Wrap(types.ErrEthTxNotValid, err.Error())
	}
	sBytes, err := hex.DecodeString(tx.S[2:])
	if err != nil {
		return &types.MsgRelayResponse{}, errors.Wrap(types.ErrEthTxNotValid, err.Error())
	}
	vBytes, err := hex.DecodeString(tx.V[2:])
	if err != nil {
		return &types.MsgRelayResponse{}, errors.Wrap(types.ErrEthTxNotValid, err.Error())
	}

	ntx := ethtypes.NewTx(&ethtypes.LegacyTx{
		Nonce:    uint64(nonce),
		To:       &ethTo,
		Value:    value,
		Gas:      uint64(gas),
		GasPrice: gasPrice,
		Data:     data,
	})

	signer := ethtypes.NewEIP155Signer(chainID)
	hash := signer.Hash(ntx)

	sig := append(rBytes, append(sBytes, vBytes...)...)
	pubKey, err := crypto.SigToPub(hash.Bytes(), sig)
	if err != nil {
		return &types.MsgRelayResponse{}, errors.Wrap(types.ErrEthTxNotValid, err.Error())
	}

	recoveredAddress := crypto.PubkeyToAddress(*pubKey)

	if !bytes.Equal(recoveredAddress.Bytes(), ethFrom.Bytes()) {
		return &types.MsgRelayResponse{}, errors.Wrap(types.ErrEthTxNotValid, "Recovered address does not equal sender")
	}

	var msg = new(types2.MsgSend)
	err = proto.Unmarshal([]byte(tx.Data), msg)
	if err != nil {
		return &types.MsgRelayResponse{}, errors.Wrap(types.ErrEthTxNotValid, err.Error())
	}

	if err := m.bk.IsSendEnabledCoins(ctx, msg.Amount...); err != nil {
		return nil, err
	}

	from, err := sdk.AccAddressFromBech32(msg.FromAddress)
	if err != nil {
		return nil, err
	}

	to, err := sdk.AccAddressFromBech32(msg.ToAddress)
	if err != nil {
		return nil, err
	}

	if m.bk.BlockedAddr(to) {
		return nil, errors.Wrapf(sdkerrors.ErrUnauthorized, "%s is not allowed to receive funds", msg.ToAddress)
	}

	err = m.bk.SendCoins(ctx, from, to, msg.Amount)
	if err != nil {
		return nil, err
	}

	defer func() {
		for _, a := range msg.Amount {
			if a.Amount.IsInt64() {
				telemetry.SetGaugeWithLabels(
					[]string{"tx", "msg", "send"},
					float32(a.Amount.Int64()),
					[]metrics.Label{telemetry.NewLabel("denom", a.Denom)},
				)
			}
		}
	}()

	m.keeper.SetRelay(ctx, relay)

	return &types.MsgRelayResponse{}, nil
}
