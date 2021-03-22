package data

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/KiraCore/sekai/INTERX/common"
	"github.com/KiraCore/sekai/INTERX/config"
	"github.com/KiraCore/sekai/INTERX/types"
	"github.com/KiraCore/sekai/INTERX/types/rosetta"
	"github.com/KiraCore/sekai/INTERX/types/rosetta/dataapi"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
)

// RegisterAccountRoutes registers Account API routers.
func RegisterAccountRoutes(r *mux.Router, gwCosmosmux *runtime.ServeMux, rpcAddr string) {
	r.HandleFunc(config.QueryRosettaAccountBalance, QueryAccountBalanceRequest(gwCosmosmux, rpcAddr)).Methods("POST")

	common.AddRPCMethod("POST", config.QueryRosettaAccountBalance, "This is an API to query account balance.", true)
}

func queryAccountBalanceHandler(r *http.Request, request types.InterxRequest, rpcAddr string, gwCosmosmux *runtime.ServeMux) (interface{}, interface{}, int) {
	var req dataapi.AccountBalanceRequest

	err := json.Unmarshal(request.Params, &req)
	if err != nil {
		common.GetLogger().Error("[rosetta-query-accountbalance] Failed to decode the request: ", err)
		return common.RosettaServeError(0, "failed to unmarshal", err.Error(), http.StatusBadRequest)
	}

	var response dataapi.AccountBalanceResponse

	balances := common.GetAccountBalances(gwCosmosmux, r.Clone(r.Context()), req.AccountIdentifier.Address)
	fmt.Println(balances)
	tokens := common.GetTokenAliases(gwCosmosmux, r.Clone(r.Context()))
	fmt.Println(tokens)

	response.Balances = make([]rosetta.Amount, 0)
	for _, balance := range balances {
		for _, token := range tokens {
			for _, denom := range token.Denoms {
				if denom == balance.Denom {
					response.Balances = append(response.Balances, rosetta.Amount{
						Value: balance.Amount,
						Currency: rosetta.Currency{
							Symbol:   token.Symbol,
							Decimals: token.Decimals,
						},
					})
					break
				}
			}
		}
	}

	return response, nil, http.StatusOK
}

// QueryAccountBalanceRequest is a function to query account balance.
func QueryAccountBalanceRequest(gwCosmosmux *runtime.ServeMux, rpcAddr string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := common.GetInterxRequest(r)
		response := common.GetResponseFormat(request, rpcAddr)
		statusCode := http.StatusOK

		common.GetLogger().Info("[rosetta-query-accountbalance] Entering account balance query")

		if !common.RPCMethods["POST"][config.QueryRosettaAccountBalance].Enabled {
			response.Response, response.Error, statusCode = common.ServeError(0, "", "API disabled", http.StatusForbidden)
		} else {
			if common.RPCMethods["POST"][config.QueryRosettaAccountBalance].CachingEnabled {
				found, cacheResponse, cacheError, cacheStatus := common.SearchCache(request, response)
				if found {
					response.Response, response.Error, statusCode = cacheResponse, cacheError, cacheStatus
					common.WrapResponse(w, request, *response, statusCode, false)

					common.GetLogger().Info("[rosetta-query-accountbalance] Returning from the cache")
					return
				}
			}

			response.Response, response.Error, statusCode = queryAccountBalanceHandler(r, request, rpcAddr, gwCosmosmux)
		}

		common.WrapResponse(w, request, *response, statusCode, common.RPCMethods["POST"][config.QueryRosettaAccountBalance].CachingEnabled)
	}
}
