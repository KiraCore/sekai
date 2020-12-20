package gateway

import (
	"io"
	"net/http"
	"os"
	"strconv"

	common "github.com/KiraCore/sekai/INTERX/common"
	tasks "github.com/KiraCore/sekai/INTERX/tasks"
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

func downloadReferenceHandler(w http.ResponseWriter, response common.ProxyResponse, key string) error {
	cache, err := tasks.LoadRefCacheMeta(key)
	if err != nil || len(cache.Path) == 0 {
		return err
	}

	w.Header().Add("Interx_chain_id", response.Chainid)
	w.Header().Add("Interx_block", strconv.FormatInt(response.Block, 10))
	w.Header().Add("Interx_blocktime", response.Blocktime)
	w.Header().Add("Interx_timestamp", strconv.FormatInt(response.Timestamp, 10))
	w.Header().Add("Interx_request_hash", response.RequestHash)

	w.Header().Add("Content-Type", cache.Header.Get("Content-Type"))
	w.Header().Add("Content-Length", cache.Header.Get("Content-Length"))

	w.WriteHeader(http.StatusOK)

	common.Mutex.Lock()

	data, err := os.Open(cache.Path)

	if err != nil {
		return err
	}

	io.Copy(w, data)

	common.Mutex.Unlock()

	return nil
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
			response.Response, response.Error, statusCode = ServeError(0, "", "API disabled", http.StatusForbidden)
		} else {
			if err := downloadReferenceHandler(w, *response, key); err == nil {
				return
			}

			response.Response, response.Error, statusCode = ServeError(0, "", "ref not found", http.StatusBadRequest)
		}

		WrapResponse(w, request, *response, statusCode, false)
	}
}
