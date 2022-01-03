package cosmos

import (
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/KiraCore/sekai/INTERX/common"
	"github.com/KiraCore/sekai/INTERX/config"
	"github.com/KiraCore/sekai/INTERX/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
)

// RegisterCosmosAuthRoutes registers query routers.
func RegisterCosmosAuthRoutes(r *mux.Router, gwCosmosmux *runtime.ServeMux, rpcAddr string) {
	r.HandleFunc(config.QueryAccounts, QueryAccountsRequest(gwCosmosmux, rpcAddr)).Methods("GET")

	common.AddRPCMethod("GET", config.QueryAccounts, "This is an API to query account address.", true)
}

func queryAccountsHandle(r *http.Request, gwCosmosmux *runtime.ServeMux) (interface{}, interface{}, int) {
	queries := mux.Vars(r)
	bech32addr := queries["address"]

	_, err := sdk.AccAddressFromBech32(bech32addr)
	if err != nil {
		common.GetLogger().Error("[query-account] Invalid bech32addr: ", bech32addr)
		return common.ServeError(0, "", err.Error(), http.StatusBadRequest)
	}

	r.URL.Path = fmt.Sprintf("/api/cosmos/auth/accounts/%s", base64.URLEncoding.EncodeToString([]byte(bech32addr)))
	return common.ServeGRPC(r, gwCosmosmux)
}

// QueryAccountsRequest is a function to query balances.
func QueryAccountsRequest(gwCosmosmux *runtime.ServeMux, rpcAddr string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		queries := mux.Vars(r)
		bech32addr := queries["address"]
		request := types.InterxRequest{
			Method:   r.Method,
			Endpoint: config.QueryAccounts,
			Params:   []byte(bech32addr),
		}
		response := common.GetResponseFormat(request, rpcAddr)
		statusCode := http.StatusOK

		common.GetLogger().Info("[query-account] Entering account query: ", bech32addr)

		if !common.RPCMethods["GET"][config.QueryAccounts].Enabled {
			response.Response, response.Error, statusCode = common.ServeError(0, "", "API disabled", http.StatusForbidden)
		} else {
			if common.RPCMethods["GET"][config.QueryAccounts].CachingEnabled {
				found, cacheResponse, cacheError, cacheStatus := common.SearchCache(request, response)
				if found {
					response.Response, response.Error, statusCode = cacheResponse, cacheError, cacheStatus
					common.WrapResponse(w, request, *response, statusCode, false)

					common.GetLogger().Info("[query-account] Returning from the cache: ", bech32addr)
					return
				}
			}

			response.Response, response.Error, statusCode = queryAccountsHandle(r, gwCosmosmux)
		}

		common.WrapResponse(w, request, *response, statusCode, common.RPCMethods["GET"][config.QueryAccounts].CachingEnabled)
	}
}
