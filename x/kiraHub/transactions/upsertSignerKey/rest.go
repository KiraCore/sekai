package signerkey

import (
	"net/http"

	"github.com/KiraCore/cosmos-sdk/client/context"
	sdk "github.com/KiraCore/cosmos-sdk/types"
	"github.com/KiraCore/cosmos-sdk/types/rest"
	"github.com/KiraCore/cosmos-sdk/x/auth/client"
	"github.com/KiraCore/sekai/types"
	"github.com/asaskevich/govalidator"
)

// Request describes the fields for rest endpoint
type Request struct {
	BaseReq     rest.BaseReq        `json:"base_req"       yaml:"base_req"       valid:"required~base_req"`
	PubKey      string              `json:"pubkey" yaml:"pubkey" valid:"required~PubKey is required"`
	KeyType     types.SignerKeyType `json:"type" yaml:"type" valid:"required~Type is required"`
	ExpiryTime  int64               `json:"expires" yaml:"expires" valid:"required~Expires is required"`
	Enabled     bool                `json:"enabled" yaml:"enabled" valid:"required~Enabled is required"`
	Permissions []int               `json:"permissions" yaml:"permissions" valid:"required~Permissions is required"`
	Curator     sdk.AccAddress      `json:"curator"  yaml:"curator" valid:"required~Curator is required"`
}

// RestRequestHandler handles rest endpoint
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
			PubKey:      request.PubKey,
			KeyType:     request.KeyType,
			ExpiryTime:  request.ExpiryTime,
			Enabled:     request.Enabled,
			Permissions: request.Permissions,
			Curator:     request.Curator,
		}

		client.WriteGenerateStdTxResponse(responseWriter, cliContext, request.BaseReq, []sdk.Msg{message})
	}
}
