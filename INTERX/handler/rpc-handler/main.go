package interx

import (
	"fmt"
	"net/http"
	"strings"

	handler "github.com/KiraCore/sekai/INTERX/handler"
)

// ServeRPC is a function to server RPC
func ServeRPC(w http.ResponseWriter, r *http.Request, rpcAddr string) bool {
	serve := false

	if strings.HasPrefix(r.URL.Path, handler.QueryStatus) && r.Method == http.MethodGet {
		serve = true

		r.URL.Path = strings.TrimPrefix(r.URL.Path, "/api/cosmos")
	} else if strings.HasPrefix(r.URL.Path, handler.QueryTransactionHash) && r.Method == http.MethodGet {
		serve = true

		hash := strings.TrimPrefix(r.URL.Path, handler.QueryTransactionHash)
		r.URL.RawQuery = fmt.Sprintf("hash=%s", hash)
		r.URL.Path = "/tx"
	} else if strings.HasPrefix(r.URL.Path, handler.PostTransaction) && r.Method == http.MethodPost {
		serve = true

		r.URL.Path = "/broadcast_tx_async"
		r.Method = http.MethodGet
	}

	if serve {
		response := handler.GetResponseFormat(rpcAddr)

		r.Host = rpcAddr
		if r.Method == http.MethodGet {
			result, err := handler.MakeGetRequest(w, r)

			if err != nil {
				result.Error = err
			}

			response.Response = result.Result
			response.Error = result.Error
		}

		handler.WrapResponse(w, *response)
	}

	return serve
}
