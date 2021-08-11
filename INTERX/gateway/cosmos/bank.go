package cosmos

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/KiraCore/sekai/INTERX/common"
	"github.com/KiraCore/sekai/INTERX/config"
	"github.com/KiraCore/sekai/INTERX/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
)

// RegisterCosmosBankRoutes registers query routers.
func RegisterCosmosBankRoutes(r *mux.Router, gwCosmosmux *runtime.ServeMux, rpcAddr string) {
	r.HandleFunc(config.QueryTotalSupply, QuerySupplyRequest(gwCosmosmux, rpcAddr)).Methods("GET")
	r.HandleFunc(config.QueryBalances, QueryBalancesRequest(gwCosmosmux, rpcAddr)).Methods("GET")

	common.AddRPCMethod("GET", config.QueryTotalSupply, "This is an API to query total supply.", true)
	common.AddRPCMethod("GET", config.QueryBalances, "This is an API to query balances of an address.", true)
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

		common.GetLogger().Info("[query-supply] Entering total supply query")

		if !common.RPCMethods["GET"][config.QueryTotalSupply].Enabled {
			response.Response, response.Error, statusCode = common.ServeError(0, "", "API disabled", http.StatusForbidden)
		} else {
			if common.RPCMethods["GET"][config.QueryTotalSupply].CachingEnabled {
				found, cacheResponse, cacheError, cacheStatus := common.SearchCache(request, response)
				if found {
					response.Response, response.Error, statusCode = cacheResponse, cacheError, cacheStatus
					common.WrapResponse(w, request, *response, statusCode, false)

					common.GetLogger().Info("[query-supply] Returning from the cache")
					return
				}
			}

			response.Response, response.Error, statusCode = querySupplyHandle(r, gwCosmosmux)
		}

		common.WrapResponse(w, request, *response, statusCode, common.RPCMethods["GET"][config.QueryTotalSupply].CachingEnabled)
	}
}

func queryBalancesHandle(r *http.Request, gwCosmosmux *runtime.ServeMux) (interface{}, interface{}, int) {
	params := mux.Vars(r)
	bech32addr := params["address"]
	queries := r.URL.Query()
	offset := queries["offset"]
	limit := queries["limit"]
	countTotal := queries["count_total"]

	var events = make([]string, 0, 3)
	if len(offset) == 1 {
		events = append(events, fmt.Sprintf("pagination.offset=%s", offset[0]))
	}
	if len(limit) == 1 {
		events = append(events, fmt.Sprintf("pagination.limit=%s", limit[0]))
	}
	if len(countTotal) == 1 {
		events = append(events, fmt.Sprintf("pagination.count_total=%s", countTotal[0]))
	}

	r.URL.RawQuery = strings.Join(events, "&")

	_, err := sdk.AccAddressFromBech32(bech32addr)
	if err != nil {
		common.GetLogger().Error("[query-balances] Invalid bech32addr: ", bech32addr)
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
			Endpoint: config.QueryBalances,
			Params:   []byte(bech32addr),
		}
		response := common.GetResponseFormat(request, rpcAddr)
		statusCode := http.StatusOK

		common.GetLogger().Info("[query-balances] Entering balances query: ", bech32addr)

		if !common.RPCMethods["GET"][config.QueryBalances].Enabled {
			response.Response, response.Error, statusCode = common.ServeError(0, "", "API disabled", http.StatusForbidden)
		} else {
			if common.RPCMethods["GET"][config.QueryBalances].CachingEnabled {
				found, cacheResponse, cacheError, cacheStatus := common.SearchCache(request, response)
				if found {
					response.Response, response.Error, statusCode = cacheResponse, cacheError, cacheStatus
					common.WrapResponse(w, request, *response, statusCode, false)

					common.GetLogger().Info("[query-balances] Returning from the cache: ", bech32addr)
					return
				}
			}

			response.Response, response.Error, statusCode = queryBalancesHandle(r, gwCosmosmux)
		}

		common.WrapResponse(w, request, *response, statusCode, common.RPCMethods["GET"][config.QueryBalances].CachingEnabled)
	}
}
