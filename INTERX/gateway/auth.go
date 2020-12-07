package gateway

import (
	"encoding/base64"
	"fmt"
	"net/http"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
)

const (
	queryAccounts = "/api/cosmos/auth/accounts"
)

// RegisterAuthRoutes registers query routers.
func RegisterAuthRoutes(r *mux.Router, gwCosmosmux *runtime.ServeMux, rpcAddr string) {
	r.HandleFunc("/api/cosmos/auth/accounts/{address}", QueryAccountsRequest(gwCosmosmux, rpcAddr)).Methods(GET)

	AddRPCMethod(GET, queryAccounts, "This is an API to query account address.", true)
}

func queryAccountsHandle(r *http.Request, gwCosmosmux *runtime.ServeMux) (interface{}, interface{}, int) {
	queries := mux.Vars(r)
	bech32addr := queries["address"]

	addr, err := sdk.AccAddressFromBech32(bech32addr)
	if err != nil {
		return ServeError(0, "", err.Error(), http.StatusBadRequest)
	}

	r.URL.Path = fmt.Sprintf("/api/cosmos/auth/accounts/%s", base64.URLEncoding.EncodeToString([]byte(addr)))
	return ServeGRPC(r, gwCosmosmux)

}

// QueryAccountsRequest is a function to query balances.
func QueryAccountsRequest(gwCosmosmux *runtime.ServeMux, rpcAddr string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		queries := mux.Vars(r)
		bech32addr := queries["address"]
		request := InterxRequest{
			Method:   r.Method,
			Endpoint: queryAccounts,
			Params:   []byte(bech32addr),
		}
		response := GetResponseFormat(request, rpcAddr)
		statusCode := http.StatusOK

		if !rpcMethods[GET][queryAccounts].Enabled {
			response.Response, response.Error, statusCode = ServeError(0, "", "", http.StatusForbidden)
		} else {
			if rpcMethods[GET][queryAccounts].CachingEnabled {
				found, cacheResponse, cacheError, cacheStatus := SearchCache(request, response)
				if found {
					response.Response, response.Error, statusCode = cacheResponse, cacheError, cacheStatus
					WrapResponse(w, request, *response, statusCode, false)
					return
				}
			}

			response.Response, response.Error, statusCode = queryAccountsHandle(r, gwCosmosmux)
		}

		WrapResponse(w, request, *response, statusCode, rpcMethods[GET][queryBalances].CachingEnabled)
	}
}
