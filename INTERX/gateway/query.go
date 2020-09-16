package gateway

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

// RegisterQueryRoutes registers query routers.
func RegisterQueryRoutes(r *mux.Router, gwCosmosmux *runtime.ServeMux, rpcAddr string) {
	r.HandleFunc("/api/cosmos/status", QueryStatusRequest(rpcAddr)).Methods("GET")
}

// QueryStatusRequest is a function to query status.
func QueryStatusRequest(rpcAddr string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
