package listOrderBooks

import (
	"fmt"
	"github.com/KiraCore/cosmos-sdk/client"
	"github.com/KiraCore/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
	"net/http"
)

func GetOrderBooks(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		by := vars["by"]
		value := vars["value"]

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/kiraHub/listOrderBooks/%s/%s", by, value), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func GetOrderBooksByTP(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		base := vars["base"]
		quote := vars["quote"]

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/kiraHub/listOrderBooks/tp/%s/%s", base, quote), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}