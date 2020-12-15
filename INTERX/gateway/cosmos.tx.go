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
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
)

const (
	postTransaction      = "/api/cosmos/txs"
	queryTransactionHash = "/api/cosmos/txs"
	encodeTransaction    = "/api/cosmos/txs/encode"
)

// RegisterCosmosTxRoutes registers query routers.
func RegisterCosmosTxRoutes(r *mux.Router, gwCosmosmux *runtime.ServeMux, rpcAddr string) {
	r.HandleFunc(postTransaction, PostTxRequest(rpcAddr)).Methods(POST)
	r.HandleFunc(encodeTransaction, EncodeTransaction(rpcAddr)).Methods(POST)
	r.HandleFunc("/api/cosmos/txs/{hash}", QueryTxHashRequest(rpcAddr)).Methods(GET)

	AddRPCMethod(POST, postTransaction, "This is an API to post transaction.", false)
	AddRPCMethod(POST, encodeTransaction, "This is an API to encode transaction.", true)
	AddRPCMethod(GET, queryTransactionHash, "This is an API to query transaction from transaction hash.", true)
}

// PostTxReq defines a tx broadcasting request.
type PostTxReq struct {
	Tx   json.RawMessage `json:"tx" yaml:"tx"`
	Mode string          `json:"mode" yaml:"mode"`
}

func postTxHandle(r *http.Request, request InterxRequest, rpcAddr string) (interface{}, interface{}, int) {
	var req PostTxReq
	err := json.Unmarshal(request.Params, &req)
	if err != nil {
		return ServeError(0, "failed to unmarshal", err.Error(), http.StatusBadRequest)
	}

	signedTx, err := interx.EncodingCg.TxConfig.TxJSONDecoder()(req.Tx)
	if err != nil {
		return ServeError(0, "failed to get signed TX", err.Error(), http.StatusBadRequest)
	}

	txBuilder, err := interx.EncodingCg.TxConfig.WrapTxBuilder(signedTx)
	if err != nil {
		return ServeError(0, "failed to get TX builder", err.Error(), http.StatusBadRequest)
	}

	txBytes, err := interx.EncodingCg.TxConfig.TxEncoder()(txBuilder.GetTx())
	if err != nil {
		return ServeError(0, "failed to get TX bytes", err.Error(), http.StatusBadRequest)
	}

	url := ""
	if req.Mode == "block" {
		url = "/broadcast_tx_commit"
	} else if req.Mode == "sync" {
		url = "/broadcast_tx_sync"
	} else if req.Mode == "async" {
		url = "/broadcast_tx_async"
	} else {
		return ServeError(0, "", "invalid mode", http.StatusBadRequest)
	}

	return MakeGetRequest(rpcAddr, url, fmt.Sprintf("tx=0x%X", txBytes))
}

// PostTxRequest is a function to post transaction.
func PostTxRequest(rpcAddr string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := GetInterxRequest(r)
		response := GetResponseFormat(request, rpcAddr)
		statusCode := http.StatusOK

		if !rpcMethods[POST][postTransaction].Enabled {
			response.Response, response.Error, statusCode = ServeError(0, "", "", http.StatusForbidden)
		} else {
			response.Response, response.Error, statusCode = postTxHandle(r, request, rpcAddr)
		}

		WrapResponse(w, request, *response, statusCode, false)
	}
}

func queryTxHashHandle(hash string, rpcAddr string) (interface{}, interface{}, int) {
	return MakeGetRequest(rpcAddr, "/tx", fmt.Sprintf("hash=%s", hash))
}

// QueryTxHashRequest is a function to query transaction hash.
func QueryTxHashRequest(rpcAddr string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		queries := mux.Vars(r)
		hash := queries["hash"]
		request := InterxRequest{
			Method:   r.Method,
			Endpoint: queryTransactionHash,
			Params:   []byte(hash),
		}
		response := GetResponseFormat(request, rpcAddr)
		statusCode := http.StatusOK

		if !rpcMethods[GET][queryTransactionHash].Enabled {
			response.Response, response.Error, statusCode = ServeError(0, "", "", http.StatusForbidden)
		} else {
			if rpcMethods[GET][queryTransactionHash].CachingEnabled {
				found, cacheResponse, cacheError, cacheStatus := SearchCache(request, response)
				if found {
					response.Response, response.Error, statusCode = cacheResponse, cacheError, cacheStatus
					WrapResponse(w, request, *response, statusCode, false)
					return
				}
			}

			response.Response, response.Error, statusCode = queryTxHashHandle(hash, rpcAddr)
		}

		WrapResponse(w, request, *response, statusCode, rpcMethods[GET][queryTransactionHash].CachingEnabled)
	}
}

func encodeTransactionHandle(r *http.Request, request InterxRequest, rpcAddr string) (interface{}, interface{}, int) {
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
		return ServeError(0, "failed to unmarshal", err.Error(), http.StatusBadRequest)
	}

	signBytes := auth.StdSignBytes(req.ChainID, req.AccountNumber, req.Sequence, 0, req.Tx.Fee, req.Tx.Msgs, req.Tx.Memo)

	// TxEncodeResponse defines base64 encoded transaction.
	type TxEncodeResponse struct {
		Tx string `json:"tx" yaml:"tx"`
	}

	return TxEncodeResponse{
		Tx: base64.StdEncoding.EncodeToString(signBytes),
	}, nil, http.StatusOK
}

// EncodeTransaction is a function to encode unsigned transaction.
func EncodeTransaction(rpcAddr string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := GetInterxRequest(r)
		response := GetResponseFormat(request, rpcAddr)
		statusCode := http.StatusOK

		if !rpcMethods[POST][encodeTransaction].Enabled {
			response.Response, response.Error, statusCode = ServeError(0, "", "", http.StatusForbidden)
		} else {
			if rpcMethods[POST][encodeTransaction].CachingEnabled {
				found, cacheResponse, cacheError, cacheStatus := SearchCache(request, response)
				if found {
					response.Response, response.Error, statusCode = cacheResponse, cacheError, cacheStatus
					WrapResponse(w, request, *response, statusCode, false)
					return
				}
			}

			response.Response, response.Error, statusCode = encodeTransactionHandle(r, request, rpcAddr)
		}

		WrapResponse(w, request, *response, statusCode, rpcMethods[POST][encodeTransaction].CachingEnabled)
	}
}
