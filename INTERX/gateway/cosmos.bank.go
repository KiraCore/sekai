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
	queryTotalSupply = "/api/cosmos/bank/supply"
	queryBalances    = "/api/cosmos/bank/balances"
)

// RegisterCosmosBankRoutes registers query routers.
func RegisterCosmosBankRoutes(r *mux.Router, gwCosmosmux *runtime.ServeMux, rpcAddr string) {
	r.HandleFunc(queryTotalSupply, QuerySupplyRequest(gwCosmosmux, rpcAddr)).Methods(GET)
	r.HandleFunc("/api/cosmos/bank/balances/{address}", QueryBalancesRequest(gwCosmosmux, rpcAddr)).Methods(GET)

	AddRPCMethod(GET, queryTotalSupply, "This is an API to query total supply.", true)
	AddRPCMethod(GET, queryBalances, "This is an API to query balances of an address.", true)
}

func querySupplyHandle(r *http.Request, gwCosmosmux *runtime.ServeMux) (interface{}, interface{}, int) {
	return ServeGRPC(r, gwCosmosmux)
}

// QuerySupplyRequest is a function to query total supply.
func QuerySupplyRequest(gwCosmosmux *runtime.ServeMux, rpcAddr string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := GetInterxRequest(r)
		response := GetResponseFormat(request, rpcAddr)
		statusCode := http.StatusOK

		if !rpcMethods[GET][queryTotalSupply].Enabled {
			response.Response, response.Error, statusCode = ServeError(0, "", "", http.StatusForbidden)
		} else {
			if rpcMethods[GET][queryTotalSupply].CachingEnabled {
				found, cacheResponse, cacheError, cacheStatus := SearchCache(request, response)
				if found {
					response.Response, response.Error, statusCode = cacheResponse, cacheError, cacheStatus
					WrapResponse(w, request, *response, statusCode, false)
					return
				}
			}

			response.Response, response.Error, statusCode = querySupplyHandle(r, gwCosmosmux)
		}

		WrapResponse(w, request, *response, statusCode, rpcMethods[GET][queryTotalSupply].CachingEnabled)
	}
}

func queryBalancesHandle(r *http.Request, gwCosmosmux *runtime.ServeMux) (interface{}, interface{}, int) {
	queries := mux.Vars(r)
	bech32addr := queries["address"]

	addr, err := sdk.AccAddressFromBech32(bech32addr)
	if err != nil {
		return ServeError(0, "", err.Error(), http.StatusBadRequest)
	}

	r.URL.Path = fmt.Sprintf("/api/cosmos/bank/balances/%s", base64.URLEncoding.EncodeToString([]byte(addr)))
	return ServeGRPC(r, gwCosmosmux)
}

// QueryBalancesRequest is a function to query balances.
func QueryBalancesRequest(gwCosmosmux *runtime.ServeMux, rpcAddr string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		queries := mux.Vars(r)
		bech32addr := queries["address"]
		request := InterxRequest{
			Method:   r.Method,
			Endpoint: queryBalances,
			Params:   []byte(bech32addr),
		}
		response := GetResponseFormat(request, rpcAddr)
		statusCode := http.StatusOK

		if !rpcMethods[GET][queryBalances].Enabled {
			response.Response, response.Error, statusCode = ServeError(0, "", "", http.StatusForbidden)
		} else {
			if rpcMethods[GET][queryBalances].CachingEnabled {
				found, cacheResponse, cacheError, cacheStatus := SearchCache(request, response)
				if found {
					response.Response, response.Error, statusCode = cacheResponse, cacheError, cacheStatus
					WrapResponse(w, request, *response, statusCode, false)
					return
				}
			}

			response.Response, response.Error, statusCode = queryBalancesHandle(r, gwCosmosmux)
		}

		WrapResponse(w, request, *response, statusCode, rpcMethods[GET][queryBalances].CachingEnabled)
	}
}
