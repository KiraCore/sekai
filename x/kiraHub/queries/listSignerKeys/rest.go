package signerkeys

import (
	"fmt"
	"net/http"

	"github.com/KiraCore/cosmos-sdk/client/context"
	"github.com/KiraCore/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
)

func ListSignerKeys(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		curator := vars["curator"]

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/kiraHub/listsignerkeys/%s", curator), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}
