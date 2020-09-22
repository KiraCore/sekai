package gateway

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

const (
	postTransaction      = "/api/cosmos/tx"
	queryTransactionHash = "/api/cosmos/tx"
)

// RegisterTxRoutes registers query routers.
func RegisterTxRoutes(r *mux.Router, gwCosmosmux *runtime.ServeMux, rpcAddr string) {
	r.HandleFunc(postTransaction, PostTxRequest(rpcAddr)).Methods("POST")
	r.HandleFunc("/api/cosmos/tx/{hash}", QueryTxHashRequest(rpcAddr)).Methods("GET")

	AddRPCMethod("Post Transaction", postTransaction, "POST")
	AddRPCMethod("Query Transaction Hash", queryTransactionHash, "GET")
}

// PostTxReq defines a tx broadcasting request.
type PostTxReq struct {
	Tx   types.StdTx `json:"tx" yaml:"tx"`
	Mode string      `json:"mode" yaml:"mode"`
}

// PostTxRequest is a function to post transaction.
func PostTxRequest(rpcAddr string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conf, err := getAPIConfig(postTransaction, "POST")
		if err == nil && conf.Disable {
			ServeError(w, rpcAddr, 0, "", "", http.StatusForbidden)
			return
		}

		var req PostTxReq

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			ServeError(w, rpcAddr, 0, "", err.Error(), http.StatusBadRequest)
			return
		}

		err = encodingConfig.Amino.UnmarshalJSON(body, &req)
		if err != nil {
			ServeError(w, rpcAddr, 0, "", err.Error(), http.StatusBadRequest)
			return
		}

		txBytes, err := encodingConfig.Amino.MarshalBinaryLengthPrefixed(req.Tx)
		if err != nil {
			ServeError(w, rpcAddr, 0, "", err.Error(), http.StatusInternalServerError)
			return
		}

		url := ""
		if req.Mode == "block" {
			url = "/broadcast_tx_commit"
		} else if req.Mode == "sync" {
			url = "/broadcast_tx_sync"
		} else if req.Mode == "async" {
			url = "/broadcast_tx_async"
		} else {
			ServeError(w, rpcAddr, 0, "", "invalid mode", http.StatusBadRequest)
			return
		}

		response, err := MakeGetRequest(w, rpcAddr, url, fmt.Sprintf("tx=0x%X", txBytes))
		if err != nil {
			ServeError(w, rpcAddr, 0, "", err.Error(), http.StatusInternalServerError)
		} else {
			ServeRPC(w, response, rpcAddr)
		}
	}
}

// QueryTxHashRequest is a function to query transaction hash.
func QueryTxHashRequest(rpcAddr string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conf, err := getAPIConfig(queryTransactionHash, "GET")
		if err == nil && conf.Disable {
			ServeError(w, rpcAddr, 0, "", "", http.StatusForbidden)
			return
		}

		queries := mux.Vars(r)
		hash := queries["hash"]

		response, err := MakeGetRequest(w, rpcAddr, "/tx", fmt.Sprintf("hash=%s", hash))
		if err != nil {
			ServeError(w, rpcAddr, 0, "", err.Error(), http.StatusInternalServerError)
		} else {
			ServeRPC(w, response, rpcAddr)
		}
	}
}
