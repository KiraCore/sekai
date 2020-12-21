package gateway

import (
	"net/http"

	functions "github.com/KiraCore/sekai/INTERX/functions"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
)

const (
	queryKiraFunctions = "/api/kira/metadata"
)

// RegisterKiraQueryRoutes registers tx query routers.
func RegisterKiraQueryRoutes(r *mux.Router, gwCosmosmux *runtime.ServeMux, rpcAddr string) {
	r.HandleFunc(queryKiraFunctions, QueryKiraFunctions(rpcAddr)).Methods(GET)

	AddRPCMethod(GET, queryKiraFunctions, "This is an API to query kira functions and metadata.", true)
}

func queryKiraFunctionsHandle(rpcAddr string) (interface{}, interface{}, int) {
	functions := functions.GetKiraFunctions()

	return functions, nil, http.StatusOK
}

// QueryKiraFunctions is a function to list functions and metadata.
func QueryKiraFunctions(rpcAddr string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := GetInterxRequest(r)
		response := GetResponseFormat(request, rpcAddr)
		statusCode := http.StatusOK

		response.Response, response.Error, statusCode = queryKiraFunctionsHandle(rpcAddr)

		WrapResponse(w, request, *response, statusCode, false)
	}
}
