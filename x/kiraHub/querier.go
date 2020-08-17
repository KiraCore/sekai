package kiraHub

import (
	sdk "github.com/KiraCore/cosmos-sdk/types"
	"github.com/KiraCore/cosmos-sdk/types/errors"
	"github.com/KiraCore/sekai/x/kiraHub/keeper"
	"github.com/KiraCore/sekai/x/kiraHub/queries"
	"github.com/KiraCore/sekai/x/kiraHub/types"
	abciTypes "github.com/tendermint/tendermint/abci/types"
)

// define query keys
const (
	QueryListOrderBooks = "listOrderBooks"
	QueryListOrders     = "listOrders"
	QueryListSignerKeys = "listSignerKeys"
)

// NewQuerier is a router for different types of querying
func NewQuerier(keeper keeper.Keeper) sdk.Querier {
	return func(context sdk.Context, path []string, requestQuery abciTypes.RequestQuery) ([]byte, error) {
		switch path[0] {

		case QueryListOrderBooks:
			return queries.QueryGetOrderBooks(context, path[1:], requestQuery, keeper)

		case QueryListOrders:
			return queries.QueryGetOrders(context, path[1:], requestQuery, keeper)

		case QueryListSignerKeys:
			return queries.QueryListSignerKeys(context, path[1:], requestQuery, keeper)

		default:
			return nil, errors.Wrapf(types.UnknownQueryCode, "%v", path[0])
		}
	}
}
