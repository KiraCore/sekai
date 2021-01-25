package kira

import (
	"net/http"

	"github.com/KiraCore/sekai/INTERX/common"
	functions "github.com/KiraCore/sekai/INTERX/functions"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
)

// RegisterKiraQueryRoutes registers tx query routers.
func RegisterKiraQueryRoutes(r *mux.Router, gwCosmosmux *runtime.ServeMux, rpcAddr string) {
	r.HandleFunc(common.QueryKiraFunctions, QueryKiraFunctions(rpcAddr)).Methods("GET")
	r.HandleFunc(common.QueryKiraStatus, QueryKiraStatusRequest(rpcAddr)).Methods("GET")

	common.AddRPCMethod("GET", common.QueryKiraFunctions, "This is an API to query kira functions and metadata.", true)
	common.AddRPCMethod("GET", common.QueryKiraStatus, "This is an API to query kira status.", true)
}

func queryKiraFunctionsHandle(rpcAddr string) (interface{}, interface{}, int) {
	functions := functions.GetKiraFunctions()

	return functions, nil, http.StatusOK
}

// QueryKiraFunctions is a function to list functions and metadata.
func QueryKiraFunctions(rpcAddr string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := common.GetInterxRequest(r)
		response := common.GetResponseFormat(request, rpcAddr)
		statusCode := http.StatusOK

		response.Response, response.Error, statusCode = queryKiraFunctionsHandle(rpcAddr)

		common.WrapResponse(w, request, *response, statusCode, false)
	}
}

// QueryKiraStatusRequest is a function to query kira status.
func QueryKiraStatusRequest(rpcAddr string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := common.GetInterxRequest(r)
		response := common.GetResponseFormat(request, rpcAddr)
		statusCode := http.StatusOK

		common.GetLogger().Info("[query-kira-status] Entering status query")

		if !common.RPCMethods["GET"][common.QueryKiraStatus].Enabled {
			response.Response, response.Error, statusCode = common.ServeError(0, "", "API disabled", http.StatusForbidden)
		} else {
			if common.RPCMethods["GET"][common.QueryKiraStatus].CachingEnabled {
				found, cacheResponse, cacheError, cacheStatus := common.SearchCache(request, response)
				if found {
					response.Response, response.Error, statusCode = cacheResponse, cacheError, cacheStatus
					common.WrapResponse(w, request, *response, statusCode, false)

					common.GetLogger().Info("[query-kira-status] Returning from the cache")
					return
				}
			}

			response.Response, response.Error, statusCode = common.MakeGetRequest(rpcAddr, "/status", "")
		}

		common.WrapResponse(w, request, *response, statusCode, common.RPCMethods["GET"][common.QueryKiraStatus].CachingEnabled)
	}
}
