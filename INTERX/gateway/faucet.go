package gateway

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	interx "github.com/KiraCore/sekai/INTERX/config"
	database "github.com/KiraCore/sekai/INTERX/database"
	sdk "github.com/cosmos/cosmos-sdk/types"
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
	response := GetResponseFormat(rpcAddr)

	faucetInfo := FaucetAccountInfo{}
	faucetInfo.Address = interx.FaucetCg.Address
	faucetInfo.Balances = GetAccountBalances(gwCosmosmux, r.Clone(r.Context()), interx.FaucetCg.Address)

	response.Response = faucetInfo

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(200)

	WrapResponse(w, *response)
}

func serveFaucet(w http.ResponseWriter, r *http.Request, gwCosmosmux *runtime.ServeMux, rpcAddr string, bech32addr string, token string) {
	response := GetResponseFormat(rpcAddr)
	// token := tokens[0]

	// check address
	_, err := sdk.AccAddressFromBech32(bech32addr)
	if err != nil {
		ServeError(w, rpcAddr, 0, "", fmt.Sprintf("invalid address: %s", err), http.StatusBadRequest)
		return
	}

	// check claim limit
	timeLeft := database.GetClaimTimeLeft(bech32addr)
	if timeLeft > 0 {
		ServeError(w, rpcAddr, 0, "", fmt.Sprintf("cliam limit: %d second(s) left", timeLeft), http.StatusBadRequest)
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
		ServeError(w, rpcAddr, 0, "", "invalid token", http.StatusBadRequest)
		return
	}

	if faucetAmount <= claimAmount {
		ServeError(w, rpcAddr, 0, "", "no need to send tokens", http.StatusBadRequest)
		return
	}

	if faucetAmount-claimAmount <= faucetMininumAmount {
		ServeError(w, rpcAddr, 0, "", "can't send tokens, less than minimum amount", http.StatusBadRequest)
		return
	}

	if faucetAmount-claimAmount > availableAmount-faucetMininumAmount {
		ServeError(w, rpcAddr, 0, "", "not enough tokens", http.StatusBadRequest)
		return
	}

	// send tokens
	response.Response = ""

	// add new claim
	database.AddNewClaim(bech32addr, time.Now())

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(200)

	WrapResponse(w, *response)
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
			ServeError(w, rpcAddr, 0, "", "invalid query parameters", http.StatusBadRequest)
		}
	}
}
