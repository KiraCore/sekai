package common

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"time"

	"github.com/KiraCore/sekai/INTERX/config"
	"github.com/KiraCore/sekai/INTERX/database"
	"github.com/KiraCore/sekai/INTERX/types"
	"github.com/KiraCore/sekai/INTERX/types/rosetta"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
)

// GetInterxRequest is a function to get Interx Request
func GetInterxRequest(r *http.Request) types.InterxRequest {
	request := types.InterxRequest{}

	request.Method = r.Method
	request.Endpoint = fmt.Sprintf("%s", r.URL)
	request.Params, _ = ioutil.ReadAll(r.Body)

	return request
}

// GetResponseFormat is a function to get response format
func GetResponseFormat(request types.InterxRequest, rpcAddr string) *types.ProxyResponse {
	response := new(types.ProxyResponse)
	response.Timestamp = time.Now().Unix()
	response.RequestHash = GetBlake2bHash(request)
	response.Chainid = NodeStatus.Chainid
	response.Block = NodeStatus.Block
	response.Blocktime = NodeStatus.Blocktime

	return response
}

// GetResponseSignature is a function to get response signature
func GetResponseSignature(response types.ProxyResponse) (string, string) {
	// Get Response Hash
	responseHash := GetBlake2bHash(response.Response)

	// Generate json to be signed
	sign := new(types.ResponseSign)
	sign.Chainid = response.Chainid
	sign.Block = response.Block
	sign.Blocktime = response.Blocktime
	sign.Timestamp = response.Timestamp
	sign.Response = responseHash
	signBytes, err := json.Marshal(sign)
	if err != nil {
		return "", responseHash
	}

	// Get Signature
	signature, err := config.Config.PrivKey.Sign(signBytes)
	if err != nil {
		return "", responseHash
	}

	return base64.StdEncoding.EncodeToString([]byte(signature)), responseHash
}

// SearchCache is a function to search response in cache
func SearchCache(request types.InterxRequest, response *types.ProxyResponse) (bool, interface{}, interface{}, int) {
	chainIDHash := GetBlake2bHash(response.Chainid)
	endpointHash := GetBlake2bHash(request.Endpoint)
	requestHash := GetBlake2bHash(request)

	// GetLogger().Info(chainIDHash, endpointHash, requestHash)
	result, err := GetCache(chainIDHash, endpointHash, requestHash)
	// GetLogger().Info(result)

	if err != nil {
		return false, nil, nil, -1
	}

	if result.ExpireAt.Before(time.Now()) && result.Response.Block != response.Block {
		return false, nil, nil, -1
	}

	return true, result.Response.Response, result.Response.Error, result.Status
}

// WrapResponse is a function to wrap response
func WrapResponse(w http.ResponseWriter, request types.InterxRequest, response types.ProxyResponse, statusCode int, saveToCache bool) {
	if statusCode == 0 {
		statusCode = 503 // Service Unavailable Error
	}
	if saveToCache {
		// GetLogger().Info("[gateway] Saving in the cache")

		chainIDHash := GetBlake2bHash(response.Chainid)
		endpointHash := GetBlake2bHash(request.Endpoint)
		requestHash := GetBlake2bHash(request)
		if conf, ok := RPCMethods[request.Method][request.Endpoint]; ok {
			err := PutCache(chainIDHash, endpointHash, requestHash, types.InterxResponse{
				Response: response,
				Status:   statusCode,
				ExpireAt: time.Now().Add(time.Duration(conf.CachingDuration) * time.Second),
			})
			if err != nil {
				// GetLogger().Error("[gateway] Failed to save in the cache: ", err.Error())
			}
			// GetLogger().Info("[gateway] Save finished")
		}
	}

	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("Interx_chain_id", response.Chainid)
	w.Header().Add("Interx_block", strconv.FormatInt(response.Block, 10))
	w.Header().Add("Interx_blocktime", response.Blocktime)
	w.Header().Add("Interx_timestamp", strconv.FormatInt(response.Timestamp, 10))
	w.Header().Add("Interx_request_hash", response.RequestHash)
	if request.Endpoint == config.QueryDataReference {
		reference, err := database.GetReference(string(request.Params))
		if err == nil {
			w.Header().Add("Interx_ref", "/download/"+reference.FilePath)
		}
	}

	if response.Response != nil {
		response.Signature, response.Hash = GetResponseSignature(response)

		w.Header().Add("Interx_signature", response.Signature)
		w.Header().Add("Interx_hash", response.Hash)
		w.WriteHeader(statusCode)

		json.NewEncoder(w).Encode(response.Response)
	} else {
		w.WriteHeader(statusCode)

		if response.Error == nil {
			response.Error = "service not available"
		}
		json.NewEncoder(w).Encode(response.Error)
	}
}

// ServeGRPC is a function to server GRPC
func ServeGRPC(r *http.Request, gwCosmosmux *runtime.ServeMux) (interface{}, interface{}, int) {
	recorder := httptest.NewRecorder()
	gwCosmosmux.ServeHTTP(recorder, r)
	resp := recorder.Result()

	result := new(interface{})
	if json.NewDecoder(resp.Body).Decode(result) == nil {
		if resp.StatusCode == http.StatusOK {
			return result, nil, resp.StatusCode
		}

		return nil, result, resp.StatusCode
	}

	return nil, nil, resp.StatusCode
}

// ServeError is a function to server GRPC
func ServeError(code int, data string, message string, statusCode int) (interface{}, interface{}, int) {
	return nil, types.ProxyResponseError{
		Code:    code,
		Data:    data,
		Message: message,
	}, statusCode
}

func RosettaBuildError(code int, message string, description string, retriable bool, details interface{}) rosetta.Error {
	return rosetta.Error{
		Code:        code,
		Message:     message,
		Description: description,
		Retriable:   retriable,
		Details:     details,
	}
}

func RosettaServeError(code int, data string, message string, statusCode int) (interface{}, interface{}, int) {
	return nil, RosettaBuildError(code, message, data, true, nil), statusCode
}
