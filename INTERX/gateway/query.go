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
	r.HandleFunc(queryStatus, QueryStatusRequest(rpcAddr)).Methods("GET")
	r.HandleFunc(queryRPCMethods, QueryRPCMethods(rpcAddr)).Methods("GET")

	AddRPCMethod("Query Status", queryStatus, "GET")
}

// QueryStatusRequest is a function to query status.
func QueryStatusRequest(rpcAddr string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conf, err := getAPIConfig(queryStatus, "GET")
		if err == nil && conf.Disable {
			ServeError(w, rpcAddr, 0, "", "", http.StatusForbidden)
			return
		}

		r.Host = rpcAddr
		r.URL.Path = "/status"

		response, err := MakeGetRequest(w, rpcAddr, "/status", "")
		if err != nil {
			ServeError(w, rpcAddr, 0, "", err.Error(), http.StatusInternalServerError)
		} else {
			ServeRPC(w, response, rpcAddr)
		}
	}
}

// QueryRPCMethods is a function to query RPC methods.
func QueryRPCMethods(rpcAddr string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response := GetResponseFormat(rpcAddr)
		response.Response = rpcMethods

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(200)

		WrapResponse(w, *response)
	}
}
