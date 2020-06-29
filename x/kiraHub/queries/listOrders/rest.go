package listOrders

import (
	"fmt"
	"github.com/KiraCore/cosmos-sdk/client/context"
	"github.com/KiraCore/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
	"net/http"
)

func GetOrders(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]
		max := vars["max_orders"]
		min := vars["min_amount"]

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/kiraHub/listOrders/%s/%s/%s", id, max, min), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}
