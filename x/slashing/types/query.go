package types

import (
	"github.com/cosmos/cosmos-sdk/types/query"
)

// SDKQueryPageReqFromCustomPageReq returns SDKQueryPageRequest from custom page request
func SDKQueryPageReqFromCustomPageReq(pageReq *PageRequest) *query.PageRequest {
	if pageReq == nil {
		return nil
	}
	return &query.PageRequest{
		Key:        pageReq.Key,
		Offset:     pageReq.Offset,
		Limit:      pageReq.Limit,
		CountTotal: pageReq.CountTotal,
	}
}
