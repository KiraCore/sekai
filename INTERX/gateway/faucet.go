package gateway

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	interx "github.com/KiraCore/sekai/INTERX/config"
	database "github.com/KiraCore/sekai/INTERX/database"
	sdk "github.com/cosmos/cosmos-sdk/types"
	auth "github.com/cosmos/cosmos-sdk/x/auth/types"
	bank "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

const (
	faucetRequestURL = "/api/faucet"
)

// RegisterFaucetRoutes registers faucet services.
func RegisterFaucetRoutes(r *mux.Router, gwCosmosmux *runtime.ServeMux, rpcAddr string) {
	r.HandleFunc(faucetRequestURL, FaucetRequest(gwCosmosmux, rpcAddr)).Methods(GET)

	AddRPCMethod(GET, faucetRequestURL, "This is an API for faucet service.", false)
}

func serveFaucetInfo(r *http.Request, gwCosmosmux *runtime.ServeMux) (interface{}, interface{}, int) {
	faucetInfo := FaucetAccountInfo{}
	faucetInfo.Address = interx.Config.Faucet.Address
	faucetInfo.Balances = GetAccountBalances(gwCosmosmux, r.Clone(r.Context()), interx.Config.Faucet.Address)

	return faucetInfo, nil, http.StatusOK
}

func serveFaucet(r *http.Request, gwCosmosmux *runtime.ServeMux, request InterxRequest, rpcAddr string, bech32addr string, token string) (interface{}, interface{}, int) {
	// check address
	faucetAccAddr, err := sdk.AccAddressFromBech32(interx.Config.Faucet.Address)
	if err != nil {
		return ServeError(0, "", fmt.Sprintf("internal server error: %s", err), http.StatusBadRequest)
	}

	// check address
	claimAccAddr, err := sdk.AccAddressFromBech32(bech32addr)
	if err != nil {
		return ServeError(0, "", fmt.Sprintf("invalid address: %s", err), http.StatusBadRequest)
	}

	// check claim limit
	timeLeft := database.GetClaimTimeLeft(bech32addr)
	if timeLeft > 0 {
		return ServeError(0, "", fmt.Sprintf("cliam limit: %d second(s) left", timeLeft), http.StatusBadRequest)
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

	faucetAmount, ok := interx.Config.Faucet.FaucetAmounts[token]               // X
	faucetMininumAmount, ok := interx.Config.Faucet.FaucetMinimumAmounts[token] // M
	if !ok {
		return ServeError(0, "", "invalid token", http.StatusBadRequest)
	}

	if faucetAmount <= claimAmount {
		return ServeError(0, "", "no need to send tokens", http.StatusBadRequest)
	}

	if faucetAmount-claimAmount <= faucetMininumAmount {
		return ServeError(0, "", "can't send tokens, less than minimum amount", http.StatusBadRequest)
	}

	if faucetAmount-claimAmount > availableAmount-faucetMininumAmount {
		return ServeError(0, "", "not enough tokens", http.StatusBadRequest)
	}

	// GET AccountNumber and Sequence
	accountNumber, sequence := GetAccountNumberSequence(gwCosmosmux, r.Clone(r.Context()), interx.Config.Faucet.Address)
	fmt.Println("accountNumber: ", accountNumber)
	fmt.Println("sequence: ", sequence)

	msgSend := &bank.MsgSend{
		FromAddress: faucetAccAddr,
		ToAddress:   claimAccAddr,
		Amount:      sdk.NewCoins(sdk.NewInt64Coin(token, faucetAmount-claimAmount)),
	}

	msgs := []sdk.Msg{msgSend}
	fee := auth.NewStdFee(200000, sdk.NewCoins())
	memo := "Faucet Transfer"

	sigs := make([]auth.StdSignature, 1)
	signBytes := auth.StdSignBytes(GetChainID(rpcAddr), accountNumber, sequence, 0, fee, msgs, memo)

	sig, err := interx.Config.Faucet.PrivKey.Sign(signBytes)
	if err != nil {
		panic(err)
	}

	sigs[0] = auth.StdSignature{PubKey: interx.Config.Faucet.PubKey, Signature: sig}

	stdTx := auth.NewStdTx(msgs, fee, sigs, memo)

	txBuilder := interx.EncodingCg.TxConfig.NewTxBuilder()
	err = txBuilder.SetMsgs(stdTx.GetMsgs()...)
	if err != nil {
		return ServeError(0, "failed to set TX Msgs", err.Error(), http.StatusBadRequest)
	}

	sigV2, err := stdTx.GetSignaturesV2()
	if err != nil {
		return ServeError(0, "failed to get SignaturesV2", err.Error(), http.StatusBadRequest)
	}

	sigV2[0].Sequence = sequence

	err = txBuilder.SetSignatures(sigV2...)
	if err != nil {
		return ServeError(0, "failed to set Signatures", err.Error(), http.StatusBadRequest)
	}

	txBuilder.SetMemo(stdTx.GetMemo())
	txBuilder.SetFeeAmount(stdTx.GetFee())
	txBuilder.SetGasLimit(stdTx.GetGas())

	txBytes, err := interx.EncodingCg.TxConfig.TxEncoder()(txBuilder.GetTx())
	if err != nil {
		return ServeError(0, "failed to get TX bytes", err.Error(), http.StatusBadRequest)
	}

	// send tokens
	txHash, err := BroadcastTransaction(rpcAddr, txBytes)
	if err != nil {
		return ServeError(0, "", err.Error(), http.StatusInternalServerError)
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
