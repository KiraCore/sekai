package cosmos

import (
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/KiraCore/sekai/INTERX/common"
	"github.com/KiraCore/sekai/INTERX/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
)

// RegisterCosmosBankRoutes registers query routers.
func RegisterCosmosBankRoutes(r *mux.Router, gwCosmosmux *runtime.ServeMux, rpcAddr string) {
	r.HandleFunc(common.QueryTotalSupply, QuerySupplyRequest(gwCosmosmux, rpcAddr)).Methods("GET")
	r.HandleFunc("/api/cosmos/bank/balances/{address}", QueryBalancesRequest(gwCosmosmux, rpcAddr)).Methods("GET")

	common.AddRPCMethod("GET", common.QueryTotalSupply, "This is an API to query total supply.", true)
	common.AddRPCMethod("GET", common.QueryBalances, "This is an API to query balances of an address.", true)
}

func querySupplyHandle(r *http.Request, gwCosmosmux *runtime.ServeMux) (interface{}, interface{}, int) {
	return common.ServeGRPC(r, gwCosmosmux)
}

// QuerySupplyRequest is a function to query total supply.
func QuerySupplyRequest(gwCosmosmux *runtime.ServeMux, rpcAddr string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := common.GetInterxRequest(r)
		response := common.GetResponseFormat(request, rpcAddr)
		statusCode := http.StatusOK

		if !common.RPCMethods["GET"][common.QueryTotalSupply].Enabled {
			response.Response, response.Error, statusCode = common.ServeError(0, "", "API disabled", http.StatusForbidden)
		} else {
			if common.RPCMethods["GET"][common.QueryTotalSupply].CachingEnabled {
				found, cacheResponse, cacheError, cacheStatus := common.SearchCache(request, response)
				if found {
					response.Response, response.Error, statusCode = cacheResponse, cacheError, cacheStatus
					common.WrapResponse(w, request, *response, statusCode, false)
					return
				}
			}

			response.Response, response.Error, statusCode = querySupplyHandle(r, gwCosmosmux)
		}

		common.WrapResponse(w, request, *response, statusCode, common.RPCMethods["GET"][common.QueryTotalSupply].CachingEnabled)
	}
}

func queryBalancesHandle(r *http.Request, gwCosmosmux *runtime.ServeMux) (interface{}, interface{}, int) {
	queries := mux.Vars(r)
	bech32addr := queries["address"]

	_, err := sdk.AccAddressFromBech32(bech32addr)
	if err != nil {
		return common.ServeError(0, "", err.Error(), http.StatusBadRequest)
	}

	r.URL.Path = fmt.Sprintf("/api/cosmos/bank/balances/%s", base64.URLEncoding.EncodeToString([]byte(bech32addr)))
	return common.ServeGRPC(r, gwCosmosmux)
}

// QueryBalancesRequest is a function to query balances.
func QueryBalancesRequest(gwCosmosmux *runtime.ServeMux, rpcAddr string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		queries := mux.Vars(r)
		bech32addr := queries["address"]
		request := types.InterxRequest{
			Method:   r.Method,
			Endpoint: common.QueryBalances,
			Params:   []byte(bech32addr),
		}
		response := common.GetResponseFormat(request, rpcAddr)
		statusCode := http.StatusOK

		if !common.RPCMethods["GET"][common.QueryBalances].Enabled {
			response.Response, response.Error, statusCode = common.ServeError(0, "", "API disabled", http.StatusForbidden)
		} else {
			if common.RPCMethods["GET"][common.QueryBalances].CachingEnabled {
				found, cacheResponse, cacheError, cacheStatus := common.SearchCache(request, response)
				if found {
					response.Response, response.Error, statusCode = cacheResponse, cacheError, cacheStatus
					common.WrapResponse(w, request, *response, statusCode, false)
					return
				}
			}

			response.Response, response.Error, statusCode = queryBalancesHandle(r, gwCosmosmux)
		}

		common.WrapResponse(w, request, *response, statusCode, common.RPCMethods["GET"][common.QueryBalances].CachingEnabled)
	}
}
