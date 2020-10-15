package gateway

import (
	"encoding/base64"
	"fmt"
	"net/http"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

const (
	queryTotalSupply = "/api/cosmos/bank/supply"
	queryBalances    = "/api/cosmos/bank/balances"
)

// RegisterBankRoutes registers query routers.
func RegisterBankRoutes(r *mux.Router, gwCosmosmux *runtime.ServeMux, rpcAddr string) {
	r.HandleFunc(queryTotalSupply, QuerySupplyRequest(gwCosmosmux, rpcAddr)).Methods(GET)
	r.HandleFunc("/api/cosmos/bank/balances/{address}", QueryBalancesRequest(gwCosmosmux, rpcAddr)).Methods(GET)

	AddRPCMethod(GET, queryTotalSupply, "This is an API to query total supply.")
	AddRPCMethod(GET, queryBalances, "This is an API to query balances of an address.")
}

// QuerySupplyRequest is a function to query total supply.
func QuerySupplyRequest(gwCosmosmux *runtime.ServeMux, rpcAddr string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := GetInterxRequest(r)

		if !rpcMethods[GET][queryTotalSupply].Enabled {
			ServeError(w, request, rpcAddr, 0, "", "", http.StatusForbidden)
			return
		}

		if rpcMethods[GET][queryTotalSupply].CachingEnabled {
			// Add Caching Here
		}

		ServeGRPC(w, r, gwCosmosmux, request, rpcAddr)
	}
}

// QueryBalancesRequest is a function to query balances.
func QueryBalancesRequest(gwCosmosmux *runtime.ServeMux, rpcAddr string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := GetInterxRequest(r)

		if !rpcMethods[GET][queryBalances].Enabled {
			ServeError(w, request, rpcAddr, 0, "", "", http.StatusForbidden)
			return
		}

		if rpcMethods[GET][queryBalances].CachingEnabled {
			// Add Caching Here
		}

		queries := mux.Vars(r)
		bech32addr := queries["address"]

		addr, err := sdk.AccAddressFromBech32(bech32addr)
		if err != nil {
			ServeError(w, request, rpcAddr, 0, "", err.Error(), http.StatusBadRequest)
		} else {
			r.URL.Path = fmt.Sprintf("/api/cosmos/bank/balances/%s", base64.URLEncoding.EncodeToString([]byte(addr)))
			ServeGRPC(w, r, gwCosmosmux, request, rpcAddr)
		}
	}
}
