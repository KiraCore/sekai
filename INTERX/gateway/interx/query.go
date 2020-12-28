package interx

import (
	"encoding/json"
	"net/http"

	"github.com/KiraCore/sekai/INTERX/common"
	interx "github.com/KiraCore/sekai/INTERX/config"
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
	common.AddRPCMethod("GET", common.QueryStatus, "This is an API to query status.", true)
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
	success, failure, status := common.MakeGetRequest(rpcAddr, "/status", "")

	if success != nil {
		type StatusTempResponse struct {
			NodeInfo      interface{} `json:"node_info,omitempty"`
			SyncInfo      interface{} `json:"sync_info,omitempty"`
			ValidatorInfo interface{} `json:"validator_info,omitempty"`
			InterxInfo    struct {
				PubKey interface{} `json:"pub_key,omitempty"`
			} `json:"interx_info,omitempty"`
		}

		result := StatusTempResponse{}
		byteData, err := json.Marshal(success)
		err = json.Unmarshal(byteData, &result)
		if err != nil {
			panic(err)
		}

		pubkeyBytes, err := interx.EncodingCg.Amino.MarshalJSON(interx.Config.PubKey)
		if err != nil {
			panic(err)
		}

		err = json.Unmarshal(pubkeyBytes, &result.InterxInfo.PubKey)
		if err != nil {
			panic(err)
		}

		success = result
	}

	return success, failure, status
}

// QueryStatusRequest is a function to query status.
func QueryStatusRequest(rpcAddr string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := common.GetInterxRequest(r)
		response := common.GetResponseFormat(request, rpcAddr)
		statusCode := http.StatusOK

		if !common.RPCMethods["GET"][common.QueryStatus].Enabled {
			response.Response, response.Error, statusCode = common.ServeError(0, "", "API disabled", http.StatusForbidden)
		} else {
			if common.RPCMethods["GET"][common.QueryStatus].CachingEnabled {
				found, cacheResponse, cacheError, cacheStatus := common.SearchCache(request, response)
				if found {
					response.Response, response.Error, statusCode = cacheResponse, cacheError, cacheStatus
					common.WrapResponse(w, request, *response, statusCode, false)
					return
				}
			}

			response.Response, response.Error, statusCode = queryStatusHandle(rpcAddr)
		}

		common.WrapResponse(w, request, *response, statusCode, common.RPCMethods["GET"][common.QueryStatus].CachingEnabled)
	}
}
