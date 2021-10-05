package cosmos

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/KiraCore/sekai/INTERX/common"
	"github.com/KiraCore/sekai/INTERX/config"
	"github.com/KiraCore/sekai/INTERX/types"
	legacytx "github.com/cosmos/cosmos-sdk/x/auth/legacy/legacytx"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
)

// RegisterCosmosTxRoutes registers query routers.
func RegisterCosmosTxRoutes(r *mux.Router, gwCosmosmux *runtime.ServeMux, rpcAddr string) {
	r.HandleFunc(config.PostTransaction, PostTxRequest(rpcAddr)).Methods("POST")
	r.HandleFunc(config.EncodeTransaction, EncodeTransaction(rpcAddr)).Methods("POST")
	r.HandleFunc(config.QueryTransactionHash, QueryTxHashRequest(rpcAddr)).Methods("GET")

	common.AddRPCMethod("POST", config.PostTransaction, "This is an API to post transaction.", false)
	common.AddRPCMethod("POST", config.EncodeTransaction, "This is an API to encode transaction.", true)
	common.AddRPCMethod("GET", config.QueryTransactionHash, "This is an API to query transaction from transaction hash.", true)
}

// PostTxReq defines a tx broadcasting request.
type PostTxReq struct {
	Tx   json.RawMessage `json:"tx" yaml:"tx"`
	Mode string          `json:"mode" yaml:"mode"`
}

func postTxHandle(r *http.Request, request types.InterxRequest, rpcAddr string) (interface{}, interface{}, int) {
	var req PostTxReq
	err := json.Unmarshal(request.Params, &req)
	if err != nil {
		common.GetLogger().Error("[post-transaction] Failed to unmarshal request: ", err)
		return common.ServeError(0, "failed to unmarshal", err.Error(), http.StatusBadRequest)
	}

	txModeAllowed := false
	for _, txMode := range config.Config.TxModes {
		if txMode == req.Mode {
			txModeAllowed = true
			break
		}
	}

	if txModeAllowed == false {
		common.GetLogger().Error("[post-transaction] Invalid transaction mode")
		return common.ServeError(0, "invalid transaction mode: ", req.Mode, http.StatusBadRequest)
	}

	url := ""
	if req.Mode == "block" {
		url = "/broadcast_tx_commit"
	} else if req.Mode == "sync" {
		url = "/broadcast_tx_sync"
	} else if req.Mode == "async" {
		url = "/broadcast_tx_async"
	} else {
		common.GetLogger().Error("[post-transaction] Invalid mode: ", req.Mode)
		return common.ServeError(0, "", "invalid mode", http.StatusBadRequest)
	}

	signedTx, err := config.EncodingCg.TxConfig.TxJSONDecoder()(req.Tx)
	if err != nil {
		common.GetLogger().Error("[post-transaction] Failed to decode tx request: ", err)
		return common.ServeError(0, "failed to get signed TX", err.Error(), http.StatusBadRequest)
	}

	txBuilder, err := config.EncodingCg.TxConfig.WrapTxBuilder(signedTx)
	if err != nil {
		common.GetLogger().Error("[post-transaction] Failed to get tx builder: ", err)
		return common.ServeError(0, "failed to get TX builder", err.Error(), http.StatusBadRequest)
	}

	txBytes, err := config.EncodingCg.TxConfig.TxEncoder()(txBuilder.GetTx())
	if err != nil {
		common.GetLogger().Error("[post-transaction] Failed to get tx bytes: ", err)
		return common.ServeError(0, "failed to get TX bytes", err.Error(), http.StatusBadRequest)
	}

	return common.MakeTendermintRPCRequest(rpcAddr, url, fmt.Sprintf("tx=0x%X", txBytes))
}

// PostTxRequest is a function to post transaction.
func PostTxRequest(rpcAddr string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := common.GetInterxRequest(r)
		response := common.GetResponseFormat(request, rpcAddr)
		statusCode := http.StatusOK

		common.GetLogger().Info("[post-transaction] Entering transaction broadcast: ")

		if !common.RPCMethods["POST"][config.PostTransaction].Enabled {
			response.Response, response.Error, statusCode = common.ServeError(0, "", "API disabled", http.StatusForbidden)
		} else {
			response.Response, response.Error, statusCode = postTxHandle(r, request, rpcAddr)
		}

		common.WrapResponse(w, request, *response, statusCode, false)
	}
}

func queryTxHashHandle(hash string, rpcAddr string) (interface{}, interface{}, int) {
	return common.MakeTendermintRPCRequest(rpcAddr, "/tx", fmt.Sprintf("hash=%s", hash))
}

// QueryTxHashRequest is a function to query transaction hash.
func QueryTxHashRequest(rpcAddr string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		queries := mux.Vars(r)
		hash := queries["hash"]
		request := types.InterxRequest{
			Method:   r.Method,
			Endpoint: config.QueryTransactionHash,
			Params:   []byte(hash),
		}
		response := common.GetResponseFormat(request, rpcAddr)
		statusCode := http.StatusOK

		common.GetLogger().Info("[query-txhash] Entering transaction hash query: ", hash)

		if !common.RPCMethods["GET"][config.QueryTransactionHash].Enabled {
			response.Response, response.Error, statusCode = common.ServeError(0, "", "API disabled", http.StatusForbidden)
		} else {
			if common.RPCMethods["GET"][config.QueryTransactionHash].CachingEnabled {
				found, cacheResponse, cacheError, cacheStatus := common.SearchCache(request, response)
				if found {
					response.Response, response.Error, statusCode = cacheResponse, cacheError, cacheStatus
					common.WrapResponse(w, request, *response, statusCode, false)

					common.GetLogger().Info("[query-txhash] Returning from the cache: ", hash)
					return
				}
			}

			response.Response, response.Error, statusCode = queryTxHashHandle(hash, rpcAddr)
		}

		common.WrapResponse(w, request, *response, statusCode, common.RPCMethods["GET"][config.QueryTransactionHash].CachingEnabled)
	}
}

func encodeTransactionHandle(r *http.Request, request types.InterxRequest, rpcAddr string) (interface{}, interface{}, int) {
	// TxEncodeReq defines a tx to be encoded.
	type TxEncodeReq struct {
		ChainID       string         `json:"chain_id" yaml:"chain_id"`
		AccountNumber uint64         `json:"account_number" yaml:"account_number"`
		Sequence      uint64         `json:"sequence" yaml:"sequence"`
		Tx            legacytx.StdTx `json:"tx" yaml:"tx"`
	}
	var req TxEncodeReq

	err := config.EncodingCg.Amino.UnmarshalJSON(request.Params, &req)
	if err != nil {
		common.GetLogger().Error("[encode-transaction] Failed to decode tx request: ", err)
		return common.ServeError(0, "failed to unmarshal", err.Error(), http.StatusBadRequest)
	}

	signBytes := legacytx.StdSignBytes(req.ChainID, req.AccountNumber, req.Sequence, 0, req.Tx.Fee, req.Tx.Msgs, req.Tx.Memo)

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
		request := common.GetInterxRequest(r)
		response := common.GetResponseFormat(request, rpcAddr)
		statusCode := http.StatusOK

		common.GetLogger().Info("[encode-transaction] Entering transaction request encoding")

		if !common.RPCMethods["POST"][config.EncodeTransaction].Enabled {
			response.Response, response.Error, statusCode = common.ServeError(0, "", "API disabled", http.StatusForbidden)
		} else {
			if common.RPCMethods["POST"][config.EncodeTransaction].CachingEnabled {
				found, cacheResponse, cacheError, cacheStatus := common.SearchCache(request, response)
				if found {
					response.Response, response.Error, statusCode = cacheResponse, cacheError, cacheStatus
					common.WrapResponse(w, request, *response, statusCode, false)

					common.GetLogger().Info("[encode-transaction] Returning from the cache")
					return
				}
			}

			response.Response, response.Error, statusCode = encodeTransactionHandle(r, request, rpcAddr)
		}

		common.WrapResponse(w, request, *response, statusCode, common.RPCMethods["POST"][config.EncodeTransaction].CachingEnabled)
	}
}
