package createOrder

import (
	"github.com/KiraCore/cosmos-sdk/client/context"
	sdk "github.com/KiraCore/cosmos-sdk/types"
	"github.com/KiraCore/cosmos-sdk/types/rest"
	"github.com/KiraCore/cosmos-sdk/x/auth/client"
	"github.com/asaskevich/govalidator"
	"net/http"
)

type Request struct {
	BaseReq       rest.BaseReq `json:"base_req"       yaml:"base_req"       valid:"required~base_req"`
	OrderBookID   string       `json:"order_book_id"  yaml:"order_book_id"  valid:"required~order_book_id"`
	OrderType     uint8        `json:"order_type"     yaml:"order_type"     valid:"required~order_type"`
	Amount        int64        `json:"amount"         yaml:"amount"         valid:"required~amount"`
	LimitPrice    int64        `json:"limit_price"    yaml:"limit_price"    valid:"required~limit_price"`
}

func RestRequestHandler(cliContext context.CLIContext) http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, httpRequest *http.Request) {
		var request Request
		if !rest.ReadRESTReq(responseWriter, httpRequest, cliContext.Codec, &request) {
			return
		}

		request.BaseReq = request.BaseReq.Sanitize()
		if !request.BaseReq.ValidateBasic(responseWriter) {
			rest.WriteErrorResponse(responseWriter, http.StatusBadRequest, "")
			return
		}

		_, Error := govalidator.ValidateStruct(request)
		if Error != nil {
			rest.WriteErrorResponse(responseWriter, http.StatusBadRequest, Error.Error())
			return
		}

		var message = Message{
			OrderBookID: request.OrderBookID,
			OrderType: request.OrderType,
			Amount: request.Amount,
			LimitPrice: request.LimitPrice,
		}

		client.WriteGenerateStdTxResponse(responseWriter, cliContext, request.BaseReq, []sdk.Msg{message})
	}
}