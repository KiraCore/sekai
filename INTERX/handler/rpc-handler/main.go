package rpc_handler

import (
    "encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	QueryStatus 			string = "/api/cosmos/status"
	QueryTransactionHash 	string = "/api/cosmos/tx/"
	PostTransaction 		string = "/api/cosmos/tx"
)

var Endpoints = []string {
	QueryStatus,
	QueryTransactionHash,
	PostTransaction,
}

type RPCResponse struct { 
	Jsonrpc		string     		`json:"jsonrpc"` 
	Id			int  			`json:"id"` 
	Result		interface{}  	`json:"result"` 
	Error		interface{}  	`json:"error"` 
}

type ProxyResponse struct {
	Chainid    	string    		`json:"chain_id"` 
	Block    	int64     		`json:"block"` 
	Timestamp   int64     		`json:"timestamp"` 
	Response    interface{}     `json:"response"` 
	Error    	interface{}     `json:"error"` 
}

func CopyHeader(dst, src http.Header) {
    for k, vv := range src {
        for _, v := range vv {
			if k != "Content-Length" {
				dst.Add(k, v)
			}
        }
    }
}

func makeGetRequest(w http.ResponseWriter, r *http.Request) (*RPCResponse, error) {
	resp, err := http.Get(fmt.Sprintf("%s%s", r.Host, r.URL))
	if err != nil {
		return nil, fmt.Errorf("error: %s", err)
	}
	defer resp.Body.Close()

	CopyHeader(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)

	result := new(RPCResponse)
	err = json.NewDecoder(resp.Body).Decode(result)
	if err != nil {
		return nil, fmt.Errorf("error: %s", err)
	}

	return result, nil
}

func makePostRequest(w http.ResponseWriter, r *http.Request) (*RPCResponse, error) {
	resp, err := http.PostForm(fmt.Sprintf("%s%s", r.Host, r.URL), r.Form)
	if err != nil {
		fmt.Printf("error: %s", err)
	}
	defer resp.Body.Close()

	CopyHeader(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)

	result := new(RPCResponse)
	err = json.NewDecoder(resp.Body).Decode(result)
	if err != nil {
		return nil, fmt.Errorf("error: %s", err)
	}

	return result, nil
}

func GetResponseFormat(rpcAddr string) *ProxyResponse  {
	response := new(ProxyResponse)
	response.Chainid = ""
	response.Block = 0
	response.Timestamp = time.Now().Unix()
	response.Response = nil
	response.Error = nil

    r, err := http.Get(fmt.Sprintf("%s/block", rpcAddr))
    if err != nil {
        return response
    }
    defer r.Body.Close()

	type RPCTempResponse struct { 
		Jsonrpc		string     		`json:"jsonrpc"` 
		Id			int  			`json:"id"` 
		Result		struct {
			Block	struct {
				Header struct {
					Chainid	string	`json:"chain_id"`
					Height	string	`json:"height"`
				}					`json:"header"`
			}						`json:"block"`
		}  							`json:"result"` 
		Error		interface{}  	`json:"error"` 
	}

	result := new(RPCTempResponse)
	if json.NewDecoder(r.Body).Decode(result) != nil {
		return response
	}

	response.Chainid = result.Result.Block.Header.Chainid
	response.Block, _ = strconv.ParseInt(result.Result.Block.Header.Height, 10, 64)

	return response;
}

func ServeRPC(w http.ResponseWriter, r *http.Request, rpcAddr string) bool {
	serve := false

	if strings.HasPrefix(r.URL.Path, QueryStatus) && r.Method == http.MethodGet {
		serve = true
		
		r.URL.Path = strings.TrimPrefix(r.URL.Path, "/api/cosmos")
	} else if strings.HasPrefix(r.URL.Path, QueryTransactionHash) && r.Method == http.MethodGet {
		serve = true

		hash := strings.TrimPrefix(r.URL.Path, QueryTransactionHash)
		r.URL.RawQuery = fmt.Sprintf("hash=%s", hash)
		r.URL.Path = "/tx"
	} else if strings.HasPrefix(r.URL.Path, PostTransaction) && r.Method == http.MethodPost {
		serve = true
		
		r.URL.Path = "/broadcast_tx_async"
		r.Method = http.MethodGet
	}

	if serve {
		response := GetResponseFormat(rpcAddr)

		r.Host = rpcAddr
		if r.Method == http.MethodGet {
			result, err := makeGetRequest(w, r)

			if err != nil {
				result.Error = err;
			}
			
			response.Response = result.Result
			response.Error = result.Error
		}

		json.NewEncoder(w).Encode(response)	
	}

	return serve
}
