package gateway

import (
	"fmt"
    "io"
	"net/http"
	"strings"
)

const (
	QueryStatus string = "/api/cosmos/status"
	QueryTransactionHash string = "/api/cosmos/tx/"
)

var Endpoints = []string {
	QueryStatus,
	QueryTransactionHash,
}

func copyHeader(dst, src http.Header) {
    for k, vv := range src {
        for _, v := range vv {
            dst.Add(k, v)
        }
    }
}

func makeGetRequest(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get(fmt.Sprintf("%s%s", r.Host, r.URL))
	if err != nil {
		fmt.Printf("RPC error: %s", err)
	}
	defer resp.Body.Close()
	
	copyHeader(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func ServeRPC(w http.ResponseWriter, r *http.Request, rpcAddr string) bool {
	serve := false

	if strings.HasPrefix(r.URL.Path, QueryStatus) {
		serve = true
		
		r.URL.Path = strings.TrimPrefix(r.URL.Path, "/api/cosmos")
	} else if strings.HasPrefix(r.URL.Path, QueryTransactionHash) {
		serve = true

		hash := strings.TrimPrefix(r.URL.Path, QueryTransactionHash)
		r.URL.RawQuery = fmt.Sprintf("hash=%s", hash)
		r.URL.Path = "/tx"
	}

	if serve {
		r.Host = rpcAddr
		if r.Method == http.MethodGet {
			makeGetRequest(w, r)
		}
	}

	return serve
}
