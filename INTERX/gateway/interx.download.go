package gateway

import (
	"net/http"
	"strings"

	interx "github.com/KiraCore/sekai/INTERX/config"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
)

const (
	download = "/download"
)

// RegisterInterxDownloadRoutes registers download routers.
func RegisterInterxDownloadRoutes(r *mux.Router, gwCosmosmux *runtime.ServeMux, rpcAddr string) {
	r.PathPrefix(download).HandlerFunc(DownloadReference()).Methods(GET)

	AddRPCMethod(GET, download, "This is an API to download files.", true)
}

// DownloadReference is a function to download reference.
func DownloadReference() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filename := strings.TrimPrefix(r.URL.Path, download+"/")

		if len(filename) != 0 {
			http.ServeFile(w, r, interx.GetReferenceCacheDir()+filename)
		}
	}
}
