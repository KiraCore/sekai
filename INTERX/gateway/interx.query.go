package gateway

import (
	"net/http"

	functions "github.com/KiraCore/sekai/INTERX/functions"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
)

const (
	queryRPCMethods      = "/api/rpc_methods"
	queryInterxFunctions = "/api/metadata"
)

// RegisterInterxQueryRoutes registers query routers.
func RegisterInterxQueryRoutes(r *mux.Router, gwCosmosmux *runtime.ServeMux, rpcAddr string) {
	r.HandleFunc(queryRPCMethods, QueryRPCMethods(rpcAddr)).Methods(GET)
	r.HandleFunc(queryInterxFunctions, QueryInterxFunctions(rpcAddr)).Methods(GET)

	AddRPCMethod(GET, queryInterxFunctions, "This is an API to query interx functions.", true)
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
