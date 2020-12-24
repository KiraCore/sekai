package gateway

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	interx "github.com/KiraCore/sekai/INTERX/config"
	database "github.com/KiraCore/sekai/INTERX/database"
	tasks "github.com/KiraCore/sekai/INTERX/tasks"
	sdk "github.com/cosmos/cosmos-sdk/types"
	legacytx "github.com/cosmos/cosmos-sdk/x/auth/legacy/legacytx"
	bank "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
)

const (
	faucetRequestURL = "/api/faucet"
)

// RegisterInterxFaucetRoutes registers faucet services.
func RegisterInterxFaucetRoutes(r *mux.Router, gwCosmosmux *runtime.ServeMux, rpcAddr string) {
	r.HandleFunc(faucetRequestURL, FaucetRequest(gwCosmosmux, rpcAddr)).Methods(GET)

	AddRPCMethod(GET, faucetRequestURL, "This is an API for faucet service.", false)
}

func serveFaucetInfo(r *http.Request, gwCosmosmux *runtime.ServeMux) (interface{}, interface{}, int) {
	faucetInfo := FaucetAccountInfo{}
	faucetInfo.Address = interx.Config.Faucet.Address
	faucetInfo.Balances = GetAccountBalances(gwCosmosmux, r.Clone(r.Context()), interx.Config.Faucet.Address)

	return faucetInfo, nil, http.StatusOK
}

/**
 * Error Codes
 * 0 : InternalServerError
 * 1 : Fail to send tokens
 * 100: Invalid address
 * 101: Claim time left
 * 102: Invalid token
 * 103: No need to send tokens
 * 104: Can't send tokens, less than minimum amount
 * 105: Not enough tokens
 */
func serveFaucet(r *http.Request, gwCosmosmux *runtime.ServeMux, request InterxRequest, rpcAddr string, bech32addr string, token string) (interface{}, interface{}, int) {
	// check address
	faucetAccAddr, err := sdk.AccAddressFromBech32(interx.Config.Faucet.Address)
	if err != nil {
		return ServeError(0, "", fmt.Sprintf("internal server error: %s", err), http.StatusInternalServerError)
	}

	// check address
	claimAccAddr, err := sdk.AccAddressFromBech32(bech32addr)
	if err != nil {
		return ServeError(100, "", fmt.Sprintf("invalid address: %s", err), http.StatusBadRequest)
	}

	// check claim limit
	timeLeft := database.GetClaimTimeLeft(bech32addr)
	if timeLeft > 0 {
		return ServeError(101, "", fmt.Sprintf("cliam limit: %d second(s) left", timeLeft), http.StatusBadRequest)
	}

	availableBalances := GetAccountBalances(gwCosmosmux, r.Clone(r.Context()), interx.Config.Faucet.Address)
	claimBalances := GetAccountBalances(gwCosmosmux, r.Clone(r.Context()), bech32addr)

	availableAmount := int64(0)
	for _, balance := range availableBalances {
		if balance.Denom == token {
			amount, err := strconv.ParseInt(balance.Amount, 10, 64)
			if err == nil {
				availableAmount = amount
			}
		}
	}

	claimAmount := int64(0) // Y
	for _, balance := range claimBalances {
		if balance.Denom == token {
			amount, err := strconv.ParseInt(balance.Amount, 10, 64)
			if err == nil {
				claimAmount = amount
			}
		}
	}

	faucetAmount, ok := interx.Config.Faucet.FaucetAmounts[token] // X

	if !ok {
		return ServeError(102, "", "invalid token", http.StatusBadRequest)
	}

	faucetMininumAmount, ok := interx.Config.Faucet.FaucetMinimumAmounts[token] // M

	if !ok {
		return ServeError(102, "", "invalid token", http.StatusBadRequest)
	}

	coinStr, ok := interx.Config.Faucet.FeeAmounts[token]

	if !ok {
		return ServeError(102, "", "invalid token", http.StatusBadRequest)
	}

	feeAmount, err := sdk.ParseCoin(coinStr)

	if err != nil {
		return ServeError(102, "", "invalid token", http.StatusBadRequest)
	}

	if faucetAmount <= claimAmount {
		return ServeError(103, "", "no need to send tokens", http.StatusBadRequest)
	}

	if faucetAmount-claimAmount <= faucetMininumAmount {
		return ServeError(104, "", "can't send tokens, less than minimum amount", http.StatusBadRequest)
	}

	if faucetAmount-claimAmount > availableAmount-faucetMininumAmount {
		return ServeError(105, "", "not enough tokens", http.StatusBadRequest)
	}

	// GET AccountNumber and Sequence
	accountNumber, sequence := GetAccountNumberSequence(gwCosmosmux, r.Clone(r.Context()), interx.Config.Faucet.Address)
	fmt.Println("accountNumber: ", accountNumber)
	fmt.Println("sequence: ", sequence)

	msgSend := &bank.MsgSend{
		FromAddress: faucetAccAddr.String(),
		ToAddress:   claimAccAddr.String(),
		Amount:      sdk.NewCoins(sdk.NewInt64Coin(token, faucetAmount-claimAmount)),
	}

	msgs := []sdk.Msg{msgSend}
	fee := legacytx.NewStdFee(200000, sdk.NewCoins(feeAmount)) //Fee handling
	memo := "Faucet Transfer"

	sigs := make([]legacytx.StdSignature, 1)
	signBytes := legacytx.StdSignBytes(tasks.NodeStatus.Chainid, accountNumber, sequence, 0, fee, msgs, memo)

	sig, err := interx.Config.Faucet.PrivKey.Sign(signBytes)
	if err != nil {
		panic(err)
	}

	sigs[0] = legacytx.StdSignature{PubKey: interx.Config.Faucet.PubKey, Signature: sig}

	stdTx := legacytx.NewStdTx(msgs, fee, sigs, memo)

	txBuilder := interx.EncodingCg.TxConfig.NewTxBuilder()
	err = txBuilder.SetMsgs(stdTx.GetMsgs()...)
	if err != nil {
		return ServeError(1, "failed to set TX Msgs", err.Error(), http.StatusInternalServerError)
	}

	sigV2, err := stdTx.GetSignaturesV2()
	if err != nil {
		return ServeError(1, "failed to get SignaturesV2", err.Error(), http.StatusInternalServerError)
	}

	sigV2[0].Sequence = sequence

	err = txBuilder.SetSignatures(sigV2...)
	if err != nil {
		return ServeError(1, "failed to set Signatures", err.Error(), http.StatusInternalServerError)
	}

	txBuilder.SetMemo(stdTx.GetMemo())
	txBuilder.SetFeeAmount(stdTx.GetFee())
	txBuilder.SetGasLimit(stdTx.GetGas())

	txBytes, err := interx.EncodingCg.TxConfig.TxEncoder()(txBuilder.GetTx())
	if err != nil {
		return ServeError(1, "failed to get TX bytes", err.Error(), http.StatusBadRequest)
	}

	// send tokens
	txHash, err := BroadcastTransaction(rpcAddr, txBytes)
	if err != nil {
		return ServeError(1, "", err.Error(), http.StatusInternalServerError)
	}

	// add new claim
	database.AddNewClaim(bech32addr, time.Now())

	type FaucetResponse struct {
		Hash string `json:"hash"`
	}
	return FaucetResponse{Hash: txHash}, nil, http.StatusOK
}

// FaucetRequest is a function to handle faucet service.
func FaucetRequest(gwCosmosmux *runtime.ServeMux, rpcAddr string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := GetInterxRequest(r)
		response := GetResponseFormat(request, rpcAddr)
		statusCode := http.StatusOK

		queries := r.URL.Query()
		claims := queries["claim"]
		tokens := queries["token"]

		if len(claims) == 0 && len(tokens) == 0 {
			response.Response, response.Error, statusCode = serveFaucetInfo(r, gwCosmosmux)
		} else if len(claims) == 1 && len(tokens) == 1 {
			response.Response, response.Error, statusCode = serveFaucet(r, gwCosmosmux, request, rpcAddr, claims[0], tokens[0])
		} else {
			response.Response, response.Error, statusCode = ServeError(0, "", "invalid query parameters", http.StatusBadRequest)
		}

		WrapResponse(w, request, *response, statusCode, false)
	}
}
