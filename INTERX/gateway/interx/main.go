package interx

import (
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
)

// RegisterRequest is a function to register requests.
func RegisterRequest(router *mux.Router, gwCosmosmux *runtime.ServeMux, rpcAddr string) {
	RegisterInterxDownloadRoutes(router, gwCosmosmux, rpcAddr)
	RegisterInterxFaucetRoutes(router, gwCosmosmux, rpcAddr)
	RegisterInterxQueryRoutes(router, gwCosmosmux, rpcAddr)
	RegisterInterxTxRoutes(router, gwCosmosmux, rpcAddr)
}
