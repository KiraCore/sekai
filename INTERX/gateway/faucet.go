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

	AddRPCMethod(GET, faucetRequestURL, "This is an API for faucet service.")
}

func serveFaucetInfo(w http.ResponseWriter, r *http.Request, gwCosmosmux *runtime.ServeMux, rpcAddr string) {
	request := GetInterxRequest(r)

	response := GetResponseFormat(request, rpcAddr)

	faucetInfo := FaucetAccountInfo{}
	faucetInfo.Address = interx.FaucetCg.Address
	faucetInfo.Balances = GetAccountBalances(gwCosmosmux, r.Clone(r.Context()), interx.FaucetCg.Address)

	response.Response = faucetInfo

	WrapResponse(w, *response, 200)
}

func serveFaucet(w http.ResponseWriter, r *http.Request, gwCosmosmux *runtime.ServeMux, rpcAddr string, bech32addr string, token string) {
	request := GetInterxRequest(r)

	// check address
	faucetAccAddr, err := sdk.AccAddressFromBech32(interx.FaucetCg.Address)
	if err != nil {
		ServeError(w, request, rpcAddr, 0, "", fmt.Sprintf("internal server error: %s", err), http.StatusBadRequest)
		return
	}

	// check address
	claimAccAddr, err := sdk.AccAddressFromBech32(bech32addr)
	if err != nil {
		ServeError(w, request, rpcAddr, 0, "", fmt.Sprintf("invalid address: %s", err), http.StatusBadRequest)
		return
	}

	// check claim limit
	timeLeft := database.GetClaimTimeLeft(bech32addr)
	if timeLeft > 0 {
		ServeError(w, request, rpcAddr, 0, "", fmt.Sprintf("cliam limit: %d second(s) left", timeLeft), http.StatusBadRequest)
		return
	}

	availableBalances := GetAccountBalances(gwCosmosmux, r.Clone(r.Context()), interx.FaucetCg.Address)
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

	faucetAmount, ok := interx.FaucetCg.FaucetAmounts[token]               // X
	faucetMininumAmount, ok := interx.FaucetCg.FaucetMinimumAmounts[token] // M
	if !ok {
		ServeError(w, request, rpcAddr, 0, "", "invalid token", http.StatusBadRequest)
		return
	}

	if faucetAmount <= claimAmount {
		ServeError(w, request, rpcAddr, 0, "", "no need to send tokens", http.StatusBadRequest)
		return
	}

	if faucetAmount-claimAmount <= faucetMininumAmount {
		ServeError(w, request, rpcAddr, 0, "", "can't send tokens, less than minimum amount", http.StatusBadRequest)
		return
	}

	if faucetAmount-claimAmount > availableAmount-faucetMininumAmount {
		ServeError(w, request, rpcAddr, 0, "", "not enough tokens", http.StatusBadRequest)
		return
	}

	// GET AccountNumber and Sequence
	accountNumber, sequence := GetAccountNumberSequence(gwCosmosmux, r.Clone(r.Context()), interx.FaucetCg.Address)
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

	sig, err := interx.FaucetCg.PrivKey.Sign(signBytes)
	if err != nil {
		panic(err)
	}

	sigs[0] = auth.StdSignature{PubKey: interx.FaucetCg.PubKey, Signature: sig}

	stdTx := auth.NewStdTx(msgs, fee, sigs, memo)

	txBuilder := interx.EncodingCg.TxConfig.NewTxBuilder()
	err = txBuilder.SetMsgs(stdTx.GetMsgs()...)
	if err != nil {
		ServeError(w, request, rpcAddr, 0, "failed to set TX Msgs", err.Error(), http.StatusBadRequest)
		return
	}

	sigV2, err := stdTx.GetSignaturesV2()
	if err != nil {
		ServeError(w, request, rpcAddr, 0, "failed to get SignaturesV2", err.Error(), http.StatusBadRequest)
		return
	}

	sigV2[0].Sequence = sequence

	err = txBuilder.SetSignatures(sigV2...)
	if err != nil {
		ServeError(w, request, rpcAddr, 0, "failed to set Signatures", err.Error(), http.StatusBadRequest)
		return
	}

	txBuilder.SetMemo(stdTx.GetMemo())
	txBuilder.SetFeeAmount(stdTx.GetFee())
	txBuilder.SetGasLimit(stdTx.GetGas())

	txBytes, err := interx.EncodingCg.TxConfig.TxEncoder()(txBuilder.GetTx())
	if err != nil {
		ServeError(w, request, rpcAddr, 0, "failed to get TX bytes", err.Error(), http.StatusBadRequest)
		return
	}

	// send tokens
	txHash, err := BroadcastTransaction(rpcAddr, txBytes)
	if err != nil {
		ServeError(w, request, rpcAddr, 0, "", err.Error(), http.StatusInternalServerError)
		return
	}

	// add new claim
	database.AddNewClaim(bech32addr, time.Now())

	response := GetResponseFormat(request, rpcAddr)
	type FaucetResponse struct {
		Hash string `json:"hash"`
	}
	response.Response = FaucetResponse{Hash: txHash}
	WrapResponse(w, *response, 200)
}

// FaucetRequest is a function to handle faucet service.
func FaucetRequest(gwCosmosmux *runtime.ServeMux, rpcAddr string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		queries := r.URL.Query()
		claims := queries["claim"]
		tokens := queries["token"]

		if len(claims) == 0 && len(tokens) == 0 {
			serveFaucetInfo(w, r, gwCosmosmux, rpcAddr)
		} else if len(claims) == 1 && len(tokens) == 1 {
			serveFaucet(w, r, gwCosmosmux, rpcAddr, claims[0], tokens[0])
		} else {
			ServeError(w, GetInterxRequest(r), rpcAddr, 0, "", "invalid query parameters", http.StatusBadRequest)
		}
	}
}
