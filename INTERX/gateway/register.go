package gateway

import (
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
)

// RegisterRequest is a function to register requests.
func RegisterRequest(router *mux.Router, gwCosmosmux *runtime.ServeMux, rpcAddr string) {
	RegisterQueryRoutes(router, gwCosmosmux, rpcAddr)
	RegisterFaucetRoutes(router, gwCosmosmux, rpcAddr)
	RegisterTxRoutes(router, gwCosmosmux, rpcAddr)
	RegisterBankRoutes(router, gwCosmosmux, rpcAddr)
	RegisterAuthRoutes(router, gwCosmosmux, rpcAddr)
	RegisterTxQueryRoutes(router, gwCosmosmux, rpcAddr)

	RegisterKiraGovRoutes(router, gwCosmosmux, rpcAddr)
}
