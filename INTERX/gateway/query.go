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
		if !rpcMethods[GET][queryStatus].Enabled {
			ServeError(w, rpcAddr, 0, "", "", http.StatusForbidden)
			return
		}

		r.Host = rpcAddr
		r.URL.Path = "/status"

		response, statusCode, err := MakeGetRequest(w, rpcAddr, "/status", "")
		if err != nil {
			ServeError(w, rpcAddr, 0, "", err.Error(), http.StatusInternalServerError)
		} else {
			ServeRPC(w, response, rpcAddr, statusCode)
		}
	}
}

// QueryRPCMethods is a function to query RPC methods.
func QueryRPCMethods(rpcAddr string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response := GetResponseFormat(rpcAddr)
		response.Response = rpcMethods

		WrapResponse(w, *response, 200)
	}
}
