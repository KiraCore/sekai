package interx

import (
	"errors"
	"net/http"

	"github.com/KiraCore/sekai/INTERX/common"
	"github.com/KiraCore/sekai/INTERX/config"
	"github.com/KiraCore/sekai/INTERX/tasks"
	"github.com/KiraCore/sekai/INTERX/types"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
)

// RegisterSnapShotQueryRoutes registers snapshot query routers.
func RegisterSnapShotQueryRoutes(r *mux.Router, gwCosmosmux *runtime.ServeMux, rpcAddr string) {
	r.HandleFunc(config.QuerySnapShot, QuerySnapShot(rpcAddr)).Methods("GET")
	r.HandleFunc(config.QuerySnapShotInfo, QuerySnapShotInfo(rpcAddr)).Methods("GET")

	common.AddRPCMethod("GET", config.QuerySnapShot, "This is an API to query snapshot.", true)
	common.AddRPCMethod("GET", config.QuerySnapShotInfo, "This is an API to get snapshot checksum.", true)
}

func getSnapShotInfo() (*types.SnapShotChecksumResponse, error) {
	if !tasks.SnapshotChecksumAvailable {
		return nil, errors.New("snapshot checksum not available")
	}

	return &types.SnapShotChecksumResponse{
		Checksum: tasks.SnapshotChecksum,
		Size:     tasks.SnapshotLength,
	}, nil
}

// QuerySnapShot is a function to query snapshot.
func QuerySnapShot(rpcAddr string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, config.SnapshotPath())
	}
}

func querySnapShotInfoHandler(rpcAddr string) (interface{}, interface{}, int) {
	info, err := getSnapShotInfo()
	if err != nil {
		return common.ServeError(0, "", "interx error", http.StatusInternalServerError)
	}

	return info, nil, http.StatusOK
}

// QuerySnapShotInfo is a function to get snapshot checksum.
func QuerySnapShotInfo(rpcAddr string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := common.GetInterxRequest(r)
		response := common.GetResponseFormat(request, rpcAddr)
		statusCode := http.StatusOK

		response.Response, response.Error, statusCode = querySnapShotInfoHandler(rpcAddr)

		common.WrapResponse(w, request, *response, statusCode, false)
	}
}
