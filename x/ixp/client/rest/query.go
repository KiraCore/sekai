package rest

import (
	"fmt"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
)

// GetOrderBooks is a rest endpoint utility for querying orderbooks
func GetOrderBooks(cliContext client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		by := vars["by"]
		value := vars["value"]

		// "/kira.ixp.Query/GetOrderBooks"
		res, _, err := cliContext.QueryWithData(fmt.Sprintf("custom/ixp/listOrderBooks/%s/%s", by, value), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliContext, res)
	}
}

// GetOrderBooksByTradingPair is a rest endpoint utility for querying orderbooks by trading pair
func GetOrderBooksByTradingPair(cliContext client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		base := vars["base"]
		quote := vars["quote"]

		// "/kira.ixp.Query/GetOrderBooksByTradingPair"
		res, _, err := cliContext.QueryWithData(fmt.Sprintf("custom/ixp/listOrderBooks/tp/%s/%s", base, quote), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliContext, res)
	}
}

// GetOrders is a rest endpoint utility for querying orders
func GetOrders(cliContext client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]
		max := vars["max_orders"]
		min := vars["min_amount"]

		// "/kira.ixp.Query/GetOrders"
		res, _, err := cliContext.QueryWithData(fmt.Sprintf("custom/ixp/listOrders/%s/%s/%s", id, max, min), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliContext, res)
	}
}

// GetSignerKeys is a rest endpoint utility for querying signer keys per curator
func GetSignerKeys(cliContext client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		curator := vars["curator"]

		// "/kira.ixp.Query/GetSignerKeys"
		res, _, err := cliContext.QueryWithData(fmt.Sprintf("custom/ixp/listsignerkeys/%s", curator), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliContext, res)
	}
}
