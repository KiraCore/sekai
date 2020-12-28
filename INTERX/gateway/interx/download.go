package interx

import (
	"net/http"
	"strings"

	"github.com/KiraCore/sekai/INTERX/common"
	interx "github.com/KiraCore/sekai/INTERX/config"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
)

// RegisterInterxDownloadRoutes registers download routers.
func RegisterInterxDownloadRoutes(r *mux.Router, gwCosmosmux *runtime.ServeMux, rpcAddr string) {
	r.PathPrefix(common.Download).HandlerFunc(DownloadReference()).Methods("GET")

	common.AddRPCMethod("GET", common.Download, "This is an API to download files.", true)
}

// DownloadReference is a function to download reference.
func DownloadReference() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filename := strings.TrimPrefix(r.URL.Path, common.Download+"/")

		if len(filename) != 0 {
			http.ServeFile(w, r, interx.GetReferenceCacheDir()+filename)
		}
	}
}
