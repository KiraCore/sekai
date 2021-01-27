package interx

import (
	"encoding/json"
	"net/http"

	"github.com/KiraCore/sekai/INTERX/common"
	"github.com/KiraCore/sekai/INTERX/config"
	functions "github.com/KiraCore/sekai/INTERX/functions"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
)

// RegisterInterxQueryRoutes registers query routers.
func RegisterInterxQueryRoutes(r *mux.Router, gwCosmosmux *runtime.ServeMux, rpcAddr string) {
	r.HandleFunc(common.QueryRPCMethods, QueryRPCMethods(rpcAddr)).Methods("GET")
	r.HandleFunc(common.QueryInterxFunctions, QueryInterxFunctions(rpcAddr)).Methods("GET")
	r.HandleFunc(common.QueryStatus, QueryStatusRequest(rpcAddr)).Methods("GET")

	common.AddRPCMethod("GET", common.QueryInterxFunctions, "This is an API to query interx functions.", true)
	common.AddRPCMethod("GET", common.QueryStatus, "This is an API to query interx status.", true)
}

// QueryRPCMethods is a function to query RPC methods.
func QueryRPCMethods(rpcAddr string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := common.GetInterxRequest(r)
		response := common.GetResponseFormat(request, rpcAddr)
		statusCode := http.StatusOK

		response.Response = common.RPCMethods

		common.WrapResponse(w, request, *response, statusCode, false)
	}
}

func queryInterxFunctionsHandle(rpcAddr string) (interface{}, interface{}, int) {
	functions := functions.GetInterxFunctions()

	return functions, nil, http.StatusOK
}

// QueryInterxFunctions is a function to list functions and metadata.
func QueryInterxFunctions(rpcAddr string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := common.GetInterxRequest(r)
		response := common.GetResponseFormat(request, rpcAddr)
		statusCode := http.StatusOK

		response.Response, response.Error, statusCode = queryInterxFunctionsHandle(rpcAddr)

		common.WrapResponse(w, request, *response, statusCode, false)
	}
}

func queryStatusHandle(rpcAddr string) (interface{}, interface{}, int) {
	type StatusTempResponse struct {
		InterxInfo struct {
			PubKey interface{} `json:"pub_key,omitempty"`
		} `json:"interx_info,omitempty"`
	}

	result := StatusTempResponse{}

	pubkeyBytes, err := config.EncodingCg.Amino.MarshalJSON(config.Config.PubKey)
	if err != nil {
		common.GetLogger().Error("[query-status] Failed to marshal interx pubkey", err)
		return common.ServeError(0, "", err.Error(), http.StatusInternalServerError)
	}

	err = json.Unmarshal(pubkeyBytes, &result.InterxInfo.PubKey)
	if err != nil {
		common.GetLogger().Error("[query-status] Failed to add interx pubkey to status response", err)
		return common.ServeError(0, "", err.Error(), http.StatusInternalServerError)
	}

	return result, nil, http.StatusOK
}

// QueryStatusRequest is a function to query status.
func QueryStatusRequest(rpcAddr string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := common.GetInterxRequest(r)
		response := common.GetResponseFormat(request, rpcAddr)
		statusCode := http.StatusOK

		common.GetLogger().Info("[query-status] Entering status query")

		if !common.RPCMethods["GET"][common.QueryStatus].Enabled {
			response.Response, response.Error, statusCode = common.ServeError(0, "", "API disabled", http.StatusForbidden)
		} else {
			if common.RPCMethods["GET"][common.QueryStatus].CachingEnabled {
				found, cacheResponse, cacheError, cacheStatus := common.SearchCache(request, response)
				if found {
					response.Response, response.Error, statusCode = cacheResponse, cacheError, cacheStatus
					common.WrapResponse(w, request, *response, statusCode, false)

					common.GetLogger().Info("[query-status] Returning from the cache")
					return
				}
			}

			response.Response, response.Error, statusCode = queryStatusHandle(rpcAddr)
		}

		common.WrapResponse(w, request, *response, statusCode, common.RPCMethods["GET"][common.QueryStatus].CachingEnabled)
	}
}
