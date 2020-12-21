package gateway

import (
	"encoding/json"
	"net/http"

	interx "github.com/KiraCore/sekai/INTERX/config"
	functions "github.com/KiraCore/sekai/INTERX/functions"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
)

const (
	queryRPCMethods      = "/api/rpc_methods"
	queryInterxFunctions = "/api/metadata"
	queryStatus          = "/api/status"
)

// RegisterInterxQueryRoutes registers query routers.
func RegisterInterxQueryRoutes(r *mux.Router, gwCosmosmux *runtime.ServeMux, rpcAddr string) {
	r.HandleFunc(queryRPCMethods, QueryRPCMethods(rpcAddr)).Methods(GET)
	r.HandleFunc(queryInterxFunctions, QueryInterxFunctions(rpcAddr)).Methods(GET)
	r.HandleFunc(queryStatus, QueryStatusRequest(rpcAddr)).Methods(GET)

	AddRPCMethod(GET, queryInterxFunctions, "This is an API to query interx functions.", true)
	AddRPCMethod(GET, queryStatus, "This is an API to query status.", true)
}

// QueryRPCMethods is a function to query RPC methods.
func QueryRPCMethods(rpcAddr string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := GetInterxRequest(r)
		response := GetResponseFormat(request, rpcAddr)
		statusCode := http.StatusOK

		response.Response = rpcMethods

		WrapResponse(w, request, *response, statusCode, false)
	}
}

func queryInterxFunctionsHandle(rpcAddr string) (interface{}, interface{}, int) {
	functions := functions.GetInterxFunctions()

	return functions, nil, http.StatusOK
}

// QueryInterxFunctions is a function to list functions and metadata.
func QueryInterxFunctions(rpcAddr string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := GetInterxRequest(r)
		response := GetResponseFormat(request, rpcAddr)
		statusCode := http.StatusOK

		response.Response, response.Error, statusCode = queryInterxFunctionsHandle(rpcAddr)

		WrapResponse(w, request, *response, statusCode, false)
	}
}

func queryStatusHandle(rpcAddr string) (interface{}, interface{}, int) {
	success, failure, status := MakeGetRequest(rpcAddr, "/status", "")

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
		request := GetInterxRequest(r)
		response := GetResponseFormat(request, rpcAddr)
		statusCode := http.StatusOK

		if !rpcMethods[GET][queryStatus].Enabled {
			response.Response, response.Error, statusCode = ServeError(0, "", "API disabled", http.StatusForbidden)
		} else {
			if rpcMethods[GET][queryStatus].CachingEnabled {
				found, cacheResponse, cacheError, cacheStatus := SearchCache(request, response)
				if found {
					response.Response, response.Error, statusCode = cacheResponse, cacheError, cacheStatus
					WrapResponse(w, request, *response, statusCode, false)
					return
				}
			}

			response.Response, response.Error, statusCode = queryStatusHandle(rpcAddr)
		}

		WrapResponse(w, request, *response, statusCode, rpcMethods[GET][queryStatus].CachingEnabled)
	}
}
