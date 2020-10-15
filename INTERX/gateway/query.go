package gateway

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

const (
	queryStatus     = "/api/cosmos/status"
	queryRPCMethods = "/api/rpc_methods"
)

// RegisterQueryRoutes registers query routers.
func RegisterQueryRoutes(r *mux.Router, gwCosmosmux *runtime.ServeMux, rpcAddr string) {
	r.HandleFunc(queryStatus, QueryStatusRequest(rpcAddr)).Methods(GET)
	r.HandleFunc(queryRPCMethods, QueryRPCMethods(rpcAddr)).Methods(GET)

	AddRPCMethod(GET, queryStatus, "This is an API to query status.")
}

// QueryStatusRequest is a function to query status.
func QueryStatusRequest(rpcAddr string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := GetInterxRequest(r)

		if !rpcMethods[GET][queryStatus].Enabled {
			ServeError(w, request, rpcAddr, 0, "", "", http.StatusForbidden)
			return
		}

		if rpcMethods[GET][queryStatus].CachingEnabled {
			// Add Caching Here
		}

		r.Host = rpcAddr
		r.URL.Path = "/status"

		response, statusCode, err := MakeGetRequest(rpcAddr, "/status", "")
		if err != nil {
			ServeError(w, request, rpcAddr, 0, "", err.Error(), http.StatusInternalServerError)
		} else {
			ServeRPC(w, request, response, rpcAddr, statusCode)
		}
	}
}

// QueryRPCMethods is a function to query RPC methods.
func QueryRPCMethods(rpcAddr string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := GetInterxRequest(r)

		response := GetResponseFormat(request, rpcAddr)
		response.Response = rpcMethods

		WrapResponse(w, *response, 200)
	}
}
