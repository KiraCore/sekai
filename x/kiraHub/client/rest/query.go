package rest

import (
	"fmt"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
)

func GetOrderBooks(cliContext client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		by := vars["by"]
		value := vars["value"]

		// "/kira.kiraHub.Query/GetOrderBooks"
		res, _, err := cliContext.QueryWithData(fmt.Sprintf("custom/kiraHub/listOrderBooks/%s/%s", by, value), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliContext, res)
	}
}

func GetOrderBooksByTradingPair(cliContext client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		base := vars["base"]
		quote := vars["quote"]

		// "/kira.kiraHub.Query/GetOrderBooksByTradingPair"
		res, _, err := cliContext.QueryWithData(fmt.Sprintf("custom/kiraHub/listOrderBooks/tp/%s/%s", base, quote), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliContext, res)
	}
}

func GetOrders(cliContext client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]
		max := vars["max_orders"]
		min := vars["min_amount"]

		// "/kira.kiraHub.Query/GetOrders"
		res, _, err := cliContext.QueryWithData(fmt.Sprintf("custom/kiraHub/listOrders/%s/%s/%s", id, max, min), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliContext, res)
	}
}

func GetSignerKeys(cliContext client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		curator := vars["curator"]

		// "/kira.kiraHub.Query/GetSignerKeys"
		res, _, err := cliContext.QueryWithData(fmt.Sprintf("custom/kiraHub/listsignerkeys/%s", curator), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliContext, res)
	}
}
