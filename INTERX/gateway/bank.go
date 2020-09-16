package gateway

import (
	"fmt"
	"net/http"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

// RegisterBankRoutes registers query routers.
func RegisterBankRoutes(r *mux.Router, gwCosmosmux *runtime.ServeMux, rpcAddr string) {
	r.HandleFunc("/api/cosmos/bank/supply", QuerySupplyRequest(gwCosmosmux, rpcAddr)).Methods("GET")
	r.HandleFunc("/api/cosmos/bank/balances/{address}", QueryBalancesRequest(gwCosmosmux, rpcAddr)).Methods("GET")
}

// QuerySupplyRequest is a function to query total supply.
func QuerySupplyRequest(gwCosmosmux *runtime.ServeMux, rpcAddr string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ServeGRPC(w, r, gwCosmosmux, rpcAddr)
	}
}

// QueryBalancesRequest is a function to query balances.
func QueryBalancesRequest(gwCosmosmux *runtime.ServeMux, rpcAddr string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		queries := mux.Vars(r)
		bech32addr := queries["address"]

		addr, err := sdk.AccAddressFromBech32(bech32addr)
		if err != nil {
			ServeError(w, rpcAddr, 0, "", err.Error(), http.StatusBadRequest)
		} else {
			r.URL.Path = fmt.Sprintf("/api/cosmos/bank/balances/%x", addr)
			ServeGRPC(w, r, gwCosmosmux, rpcAddr)
		}
	}
}
