package interx

import (
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
)

// RegisterRequest is a function to register requests.
func RegisterRequest(router *mux.Router, gwCosmosmux *runtime.ServeMux, rpcAddr string) {
	RegisterBlockRoutes(router, gwCosmosmux, rpcAddr)
	RegisterInterxDownloadRoutes(router, gwCosmosmux, rpcAddr)
	RegisterInterxFaucetRoutes(router, gwCosmosmux, rpcAddr)
	RegisterInterxQueryRoutes(router, gwCosmosmux, rpcAddr)
	RegisterInterxTxRoutes(router, gwCosmosmux, rpcAddr)
	RegisterValidatorsQueryRoutes(router, gwCosmosmux, rpcAddr)
	RegisterGenesisQueryRoutes(router, gwCosmosmux, rpcAddr)
	RegisterSnapShotQueryRoutes(router, gwCosmosmux, rpcAddr)
	RegisterNodeListQueryRoutes(router, gwCosmosmux, rpcAddr)
}
