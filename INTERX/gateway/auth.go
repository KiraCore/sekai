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
	queryAccounts = "/api/cosmos/auth/accounts"
)

// RegisterAuthRoutes registers query routers.
func RegisterAuthRoutes(r *mux.Router, gwCosmosmux *runtime.ServeMux, rpcAddr string) {
	r.HandleFunc("/api/cosmos/auth/accounts/{address}", QueryAccountsRequest(gwCosmosmux, rpcAddr)).Methods(GET)

	AddRPCMethod(GET, queryAccounts, "This is an API to query account address.")
}

// QueryAccountsRequest is a function to query balances.
func QueryAccountsRequest(gwCosmosmux *runtime.ServeMux, rpcAddr string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := GetInterxRequest(r)

		if !rpcMethods[GET][queryBalances].Enabled {
			ServeError(w, request, rpcAddr, 0, "", "", http.StatusForbidden)
			return
		}

		queries := mux.Vars(r)
		bech32addr := queries["address"]

		addr, err := sdk.AccAddressFromBech32(bech32addr)
		if err != nil {
			ServeError(w, request, rpcAddr, 0, "", err.Error(), http.StatusBadRequest)
		} else {
			r.URL.Path = fmt.Sprintf("/api/cosmos/auth/accounts/%s", base64.URLEncoding.EncodeToString([]byte(addr)))
			ServeGRPC(w, r, gwCosmosmux, request, rpcAddr)
		}
	}
}
