package gateway

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	interx "github.com/KiraCore/sekai/INTERX/config"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

const (
	postTransaction      = "/api/cosmos/txs"
	queryTransactionHash = "/api/cosmos/txs"
)

// RegisterTxRoutes registers query routers.
func RegisterTxRoutes(r *mux.Router, gwCosmosmux *runtime.ServeMux, rpcAddr string) {
	r.HandleFunc(postTransaction, PostTxRequest(rpcAddr)).Methods(POST)
	r.HandleFunc("/api/cosmos/txs/{hash}", QueryTxHashRequest(rpcAddr)).Methods(GET)

	AddRPCMethod(POST, postTransaction, "This is an API to post transaction.")
	AddRPCMethod(GET, queryTransactionHash, "This is an API to query transaction from transaction hash.")
}

// PostTxReq defines a tx broadcasting request.
type PostTxReq struct {
	Tx   json.RawMessage `json:"tx" yaml:"tx"`
	Mode string          `json:"mode" yaml:"mode"`
}

// PostTxRequest is a function to post transaction.
func PostTxRequest(rpcAddr string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !rpcMethods[POST][postTransaction].Enabled {
			ServeError(w, rpcAddr, 0, "", "", http.StatusForbidden)
			return
		}

		var req PostTxReq

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			ServeError(w, rpcAddr, 0, "", err.Error(), http.StatusBadRequest)
			return
		}

		err = json.Unmarshal(body, &req)
		if err != nil {
			ServeError(w, rpcAddr, 0, "failed to unmarshal", err.Error(), http.StatusBadRequest)
			return
		}

		signedTx, err := interx.EncodingCg.TxConfig.TxJSONDecoder()(req.Tx)
		if err != nil {
			ServeError(w, rpcAddr, 0, "failed to get signed TX", err.Error(), http.StatusBadRequest)
			return
		}

		txBuilder, err := interx.EncodingCg.TxConfig.WrapTxBuilder(signedTx)
		if err != nil {
			ServeError(w, rpcAddr, 0, "failed to get TX builder", err.Error(), http.StatusBadRequest)
			return
		}

		txBytes, err := interx.EncodingCg.TxConfig.TxEncoder()(txBuilder.GetTx())
		if err != nil {
			ServeError(w, rpcAddr, 0, "failed to get TX bytes", err.Error(), http.StatusBadRequest)
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

		response, statusCode, err := MakeGetRequest(w, rpcAddr, url, fmt.Sprintf("tx=0x%X", txBytes))
		if err != nil {
			ServeError(w, rpcAddr, 0, "", err.Error(), http.StatusInternalServerError)
		} else {
			ServeRPC(w, response, rpcAddr, statusCode)
		}
	}
}

// QueryTxHashRequest is a function to query transaction hash.
func QueryTxHashRequest(rpcAddr string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !rpcMethods[GET][queryTransactionHash].Enabled {
			ServeError(w, rpcAddr, 0, "", "", http.StatusForbidden)
			return
		}

		queries := mux.Vars(r)
		hash := queries["hash"]

		response, statusCode, err := MakeGetRequest(w, rpcAddr, "/tx", fmt.Sprintf("hash=%s", hash))
		if err != nil {
			ServeError(w, rpcAddr, 0, "", err.Error(), http.StatusInternalServerError)
		} else {
			ServeRPC(w, response, rpcAddr, statusCode)
		}
	}
}
