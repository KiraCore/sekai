package gateway

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	interx "github.com/KiraCore/sekai/INTERX/config"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	auth "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

const (
	postTransaction      = "/api/cosmos/txs"
	queryTransactionHash = "/api/cosmos/txs"
	encodeTransaction    = "/api/cosmos/txs/encode"
)

// RegisterTxRoutes registers query routers.
func RegisterTxRoutes(r *mux.Router, gwCosmosmux *runtime.ServeMux, rpcAddr string) {
	r.HandleFunc(postTransaction, PostTxRequest(rpcAddr)).Methods(POST)
	r.HandleFunc(encodeTransaction, EncodeTransaction(rpcAddr)).Methods(POST)
	r.HandleFunc("/api/cosmos/txs/{hash}", QueryTxHashRequest(rpcAddr)).Methods(GET)

	AddRPCMethod(POST, postTransaction, "This is an API to post transaction.")
	AddRPCMethod(POST, encodeTransaction, "This is an API to encode transaction.")
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
		request := GetInterxRequest(r)

		if !rpcMethods[POST][postTransaction].Enabled {
			ServeError(w, request, rpcAddr, 0, "", "", http.StatusForbidden)
			return
		}

		var req PostTxReq
		err := json.Unmarshal(request.Params, &req)
		if err != nil {
			ServeError(w, request, rpcAddr, 0, "failed to unmarshal", err.Error(), http.StatusBadRequest)
			return
		}

		signedTx, err := interx.EncodingCg.TxConfig.TxJSONDecoder()(req.Tx)
		if err != nil {
			ServeError(w, request, rpcAddr, 0, "failed to get signed TX", err.Error(), http.StatusBadRequest)
			return
		}

		txBuilder, err := interx.EncodingCg.TxConfig.WrapTxBuilder(signedTx)
		if err != nil {
			ServeError(w, request, rpcAddr, 0, "failed to get TX builder", err.Error(), http.StatusBadRequest)
			return
		}

		txBytes, err := interx.EncodingCg.TxConfig.TxEncoder()(txBuilder.GetTx())
		if err != nil {
			ServeError(w, request, rpcAddr, 0, "failed to get TX bytes", err.Error(), http.StatusBadRequest)
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
			ServeError(w, request, rpcAddr, 0, "", "invalid mode", http.StatusBadRequest)
			return
		}

		response, statusCode, err := MakeGetRequest(rpcAddr, url, fmt.Sprintf("tx=0x%X", txBytes))
		if err != nil {
			ServeError(w, request, rpcAddr, 0, "", err.Error(), http.StatusInternalServerError)
		} else {
			ServeRPC(w, request, response, rpcAddr, statusCode)
		}
	}
}

// QueryTxHashRequest is a function to query transaction hash.
func QueryTxHashRequest(rpcAddr string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := GetInterxRequest(r)

		if !rpcMethods[GET][queryTransactionHash].Enabled {
			ServeError(w, request, rpcAddr, 0, "", "", http.StatusForbidden)
			return
		}

		if rpcMethods[GET][queryTransactionHash].CachingEnabled {
			// Add Caching Here
		}

		queries := mux.Vars(r)
		hash := queries["hash"]

		response, statusCode, err := MakeGetRequest(rpcAddr, "/tx", fmt.Sprintf("hash=%s", hash))
		if err != nil {
			ServeError(w, request, rpcAddr, 0, "", err.Error(), http.StatusInternalServerError)
		} else {
			ServeRPC(w, request, response, rpcAddr, statusCode)
		}
	}
}

// EncodeTransaction is a function to encode unsigned transaction.
func EncodeTransaction(rpcAddr string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := GetInterxRequest(r)

		if !rpcMethods[POST][encodeTransaction].Enabled {
			ServeError(w, request, rpcAddr, 0, "", "", http.StatusForbidden)
			return
		}

		// TxEncodeReq defines a tx to be encoded.
		type TxEncodeReq struct {
			ChainID       string      `json:"chain_id" yaml:"chain_id"`
			AccountNumber uint64      `json:"account_number" yaml:"account_number"`
			Sequence      uint64      `json:"sequence" yaml:"sequence"`
			Tx            types.StdTx `json:"tx" yaml:"tx"`
		}
		var req TxEncodeReq

		err := interx.EncodingCg.Amino.UnmarshalJSON(request.Params, &req)
		if err != nil {
			ServeError(w, request, rpcAddr, 0, "failed to unmarshal", err.Error(), http.StatusBadRequest)
			return
		}

		signBytes := auth.StdSignBytes(req.ChainID, req.AccountNumber, req.Sequence, 0, req.Tx.Fee, req.Tx.Msgs, req.Tx.Memo)

		// TxEncodeResponse defines base64 encoded transaction.
		type TxEncodeResponse struct {
			Tx string `json:"tx" yaml:"tx"`
		}

		response := GetResponseFormat(request, rpcAddr)
		response.Response = TxEncodeResponse{
			Tx: base64.StdEncoding.EncodeToString(signBytes),
		}

		WrapResponse(w, *response, 200)
	}
}
