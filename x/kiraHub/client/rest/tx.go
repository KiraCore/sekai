package rest

import (
	"net/http"

	"github.com/KiraCore/cosmos-sdk/client"
	"github.com/KiraCore/cosmos-sdk/client/tx"
	sdk "github.com/KiraCore/cosmos-sdk/types"
	"github.com/KiraCore/cosmos-sdk/types/rest"
	"github.com/KiraCore/sekai/x/kiraHub/types"
	"github.com/asaskevich/govalidator"
)

type CreateOrderRequest struct {
	BaseReq     rest.BaseReq         `json:"base_req"       yaml:"base_req"       valid:"required~base_req"`
	OrderBookID string               `json:"order_book_id"  yaml:"order_book_id"  valid:"required~order_book_id"`
	OrderType   types.LimitOrderType `json:"order_type"     yaml:"order_type"     valid:"required~order_type"`
	Amount      int64                `json:"amount"         yaml:"amount"         valid:"required~amount"`
	LimitPrice  int64                `json:"limit_price"    yaml:"limit_price"    valid:"required~limit_price"`
	Curator     string               `json:"curator"  yaml:"curator"  valid:"required~curator"`
}

func RestCreateOrderRequestHandler(cliContext client.Context) http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, httpRequest *http.Request) {
		var request CreateOrderRequest
		if !rest.ReadRESTReq(responseWriter, httpRequest, cliContext.JSONMarshaler, &request) {
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

		curator, Error := sdk.AccAddressFromBech32(request.Curator)
		if Error != nil {
			rest.WriteErrorResponse(responseWriter, http.StatusBadRequest, Error.Error())
			return
		}

		message, _ := types.NewMsgCreateOrder(
			request.OrderBookID,
			request.OrderType,
			request.Amount,
			request.LimitPrice,
			curator,
		)

		tx.WriteGeneratedTxResponse(cliContext, responseWriter, request.BaseReq, message)
	}
}

type CreateOrderBookRequest struct {
	BaseReq  rest.BaseReq `json:"base_req" yaml:"base_req" valid:"required~base_req"`
	Base     string       `json:"base"     yaml:"base"     valid:"required~base"`
	Quote    string       `json:"quote"    yaml:"quote"    valid:"required~quote"`
	Mnemonic string       `json:"mnemonic" yaml:"mnemonic" valid:"required~mnemonic"`
	Curator  string       `json:"curator"  yaml:"curator"  valid:"required~curator"`
}

func RestCreateOrderbookRequestHandler(cliContext client.Context) http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, httpRequest *http.Request) {
		var request CreateOrderBookRequest
		if !rest.ReadRESTReq(responseWriter, httpRequest, cliContext.JSONMarshaler, &request) {
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

		curator, Error := sdk.AccAddressFromBech32(request.Curator)
		if Error != nil {
			rest.WriteErrorResponse(responseWriter, http.StatusBadRequest, Error.Error())
			return
		}

		message, _ := types.NewMsgCreateOrderBook(
			request.Base,
			request.Quote,
			request.Mnemonic,
			curator,
		)

		tx.WriteGeneratedTxResponse(cliContext, responseWriter, request.BaseReq, message)
	}
}

// UpsertSignerKeyRequest describes the fields for rest endpoint
type UpsertSignerKeyRequest struct {
	BaseReq     rest.BaseReq        `json:"base_req"       yaml:"base_req"       valid:"required~base_req"`
	PubKey      string              `json:"pubkey" yaml:"pubkey" valid:"required~PubKey is required"`
	KeyType     types.SignerKeyType `json:"type" yaml:"type" valid:"required~Type is required"`
	ExpiryTime  int64               `json:"expires" yaml:"expires" valid:"required~Expires is required"`
	Enabled     bool                `json:"enabled" yaml:"enabled" valid:"required~Enabled is required"`
	Permissions []int64             `json:"permissions" yaml:"permissions" valid:"required~Permissions is required"`
	Curator     sdk.AccAddress      `json:"curator"  yaml:"curator" valid:"required~Curator is required"`
}

// RestRequestHandler handles rest endpoint
func RestUpsertSignerKeyRequestHandler(cliContext client.Context) http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, httpRequest *http.Request) {
		var request UpsertSignerKeyRequest
		if !rest.ReadRESTReq(responseWriter, httpRequest, cliContext.JSONMarshaler, &request) {
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

		// curator, Error := sdk.AccAddressFromBech32(request.Curator)
		// if Error != nil {
		// 	rest.WriteErrorResponse(responseWriter, http.StatusBadRequest, Error.Error())
		// 	return
		// }

		message, _ := types.NewMsgUpsertSignerKey(
			request.PubKey,
			request.KeyType,
			request.ExpiryTime,
			request.Enabled,
			request.Permissions,
			request.Curator,
		)

		tx.WriteGeneratedTxResponse(cliContext, responseWriter, request.BaseReq, message)
	}
}
