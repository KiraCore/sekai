package createOrderBook

//
//import (
//	"github.com/KiraCore/cosmos-sdk/client/tx"
//	"net/http"
//
//	"github.com/asaskevich/govalidator"
//
//	. "github.com/KiraCore/cosmos-sdk/client"
//	sdk "github.com/KiraCore/cosmos-sdk/types"
//	"github.com/KiraCore/cosmos-sdk/types/rest"
//)
//
//type Request struct {
//	BaseReq  rest.BaseReq `json:"base_req" yaml:"base_req" valid:"required~base_req"`
//	Base     string       `json:"base"     yaml:"base"     valid:"required~base"`
//	Quote    string       `json:"quote"    yaml:"quote"    valid:"required~quote"`
//	Mnemonic string       `json:"mnemonic" yaml:"mnemonic" valid:"required~mnemonic"`
//	Curator  string       `json:"curator"  yaml:"curator"  valid:"required~curator"`
//}
//
//func RestRequestHandler(cliContext Context) http.HandlerFunc {
//	return func(responseWriter http.ResponseWriter, httpRequest *http.Request) {
//		var request Request
//		if !rest.ReadRESTReq(responseWriter, httpRequest, cliContext.Codec, &request) {
//			return
//		}
//
//		request.BaseReq = request.BaseReq.Sanitize()
//		if !request.BaseReq.ValidateBasic(responseWriter) {
//			rest.WriteErrorResponse(responseWriter, http.StatusBadRequest, "")
//			return
//		}
//
//		_, Error := govalidator.ValidateStruct(request)
//		if Error != nil {
//			rest.WriteErrorResponse(responseWriter, http.StatusBadRequest, Error.Error())
//			return
//		}
//
//		curator, Error := sdk.AccAddressFromBech32(request.Curator)
//		if Error != nil {
//			rest.WriteErrorResponse(responseWriter, http.StatusBadRequest, Error.Error())
//			return
//		}
//
//		//var message = Message{
//		//	Base:     request.Base,
//		//	Quote:    request.Quote,
//		//	Mnemonic: request.Mnemonic,
//		//	Curator:  curator,
//		//}
//
//		//tx.WriteGeneratedTxResponse(responseWriter, cliContext, request.BaseReq, []sdk.Msg{message})
//	}
//}
