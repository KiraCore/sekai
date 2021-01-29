package interx

import (
	"net/http"
	"strings"

	"github.com/KiraCore/sekai/INTERX/common"
	"github.com/KiraCore/sekai/INTERX/config"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
)

// RegisterInterxDownloadRoutes registers download routers.
func RegisterInterxDownloadRoutes(r *mux.Router, gwCosmosmux *runtime.ServeMux, rpcAddr string) {
	r.PathPrefix(config.Download).HandlerFunc(DownloadReference()).Methods("GET")

	common.AddRPCMethod("GET", config.Download, "This is an API to download files.", true)
}

// DownloadReference is a function to download reference.
func DownloadReference() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filename := strings.TrimPrefix(r.URL.Path, config.Download+"/")

		common.GetLogger().Info("[download] Entering reference download: ", filename)

		if len(filename) != 0 {
			http.ServeFile(w, r, config.GetReferenceCacheDir()+"/"+filename)
		}
	}
}
