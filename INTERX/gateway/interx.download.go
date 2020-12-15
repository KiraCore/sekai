package gateway

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
)

const (
	downloadReference = "/api/download/{key}"
)

// RegisterInterxDownloadRoutes registers interx download routes.
func RegisterInterxDownloadRoutes(r *mux.Router, gwCosmosmux *runtime.ServeMux, rpcAddr string) {
	r.HandleFunc(downloadReference, DownloadReference(rpcAddr)).Methods(GET)

	AddRPCMethod(GET, downloadReference, "This is an API to download references.", true)
}

func downloadReferenceHandler(key string) (interface{}, interface{}, int) {
	return nil, nil, http.StatusOK
}

// DownloadReference is a function to download reference.
func DownloadReference(rpcAddr string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		queries := mux.Vars(r)
		key := queries["key"]
		request := InterxRequest{
			Method:   r.Method,
			Endpoint: queryAccounts,
			Params:   []byte(key),
		}
		response := GetResponseFormat(request, rpcAddr)
		statusCode := http.StatusOK

		if !rpcMethods[GET][queryAccounts].Enabled {
			response.Response, response.Error, statusCode = ServeError(0, "", "", http.StatusForbidden)
		} else {
			response.Response, response.Error, statusCode = downloadReferenceHandler(key)
		}

		WrapResponse(w, request, *response, statusCode, rpcMethods[GET][queryBalances].CachingEnabled)
	}
}
