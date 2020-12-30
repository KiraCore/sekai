package kira

import (
	"net/http"

	"github.com/KiraCore/sekai/INTERX/common"
	functions "github.com/KiraCore/sekai/INTERX/functions"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
)

// RegisterKiraQueryRoutes registers tx query routers.
func RegisterKiraQueryRoutes(r *mux.Router, gwCosmosmux *runtime.ServeMux, rpcAddr string) {
	r.HandleFunc(common.QueryKiraFunctions, QueryKiraFunctions(rpcAddr)).Methods("GET")

	common.AddRPCMethod("GET", common.QueryKiraFunctions, "This is an API to query kira functions and metadata.", true)
}

func queryKiraFunctionsHandle(rpcAddr string) (interface{}, interface{}, int) {
	functions := functions.GetKiraFunctions()

	return functions, nil, http.StatusOK
}

// QueryKiraFunctions is a function to list functions and metadata.
func QueryKiraFunctions(rpcAddr string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := common.GetInterxRequest(r)
		response := common.GetResponseFormat(request, rpcAddr)
		statusCode := http.StatusOK

		response.Response, response.Error, statusCode = queryKiraFunctionsHandle(rpcAddr)

		common.WrapResponse(w, request, *response, statusCode, false)
	}
}
