package cosmos

import (
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
)

// RegisterRequest is a function to register requests.
func RegisterRequest(router *mux.Router, gwCosmosmux *runtime.ServeMux, rpcAddr string) {
	RegisterCosmosAuthRoutes(router, gwCosmosmux, rpcAddr)
	RegisterCosmosBankRoutes(router, gwCosmosmux, rpcAddr)
	RegisterCosmosTxRoutes(router, gwCosmosmux, rpcAddr)
}
