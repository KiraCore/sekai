package kiraHub

import (
	sdk "github.com/KiraCore/cosmos-sdk/types"
	"github.com/KiraCore/cosmos-sdk/types/errors"
	constants "github.com/KiraCore/sekai/x/kiraHub/constants"
	"github.com/KiraCore/sekai/x/kiraHub/queries/listOrderBooks"
	"github.com/KiraCore/sekai/x/kiraHub/queries/listOrders"
	signerkeys "github.com/KiraCore/sekai/x/kiraHub/queries/listSignerKeys"
	abciTypes "github.com/tendermint/tendermint/abci/types"
)

// define query keys
const (
	QueryListOrderBooks = "listOrderBooks"
	QueryListOrders     = "listOrders"
	QueryListSignerKeys = "listSignerKeys"
)

// NewQuerier is a router for different types of querying
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(context sdk.Context, path []string, requestQuery abciTypes.RequestQuery) ([]byte, error) {
		switch path[0] {

		case QueryListOrderBooks:
			return listOrderBooks.QueryGetOrderBooks(context, path[1:], requestQuery, keeper.getCreateOrderBookKeeper())

		case QueryListOrders:
			return listOrders.QueryGetOrders(context, path[1:], requestQuery, keeper.getCreateOrderKeeper())

		case QueryListSignerKeys:
			return signerkeys.QueryListSignerKeys(context, path[1:], requestQuery, keeper.getUpsertSignerKeyKeeper())

		default:
			return nil, errors.Wrapf(constants.UnknownQueryCode, "%v", path[0])
		}
	}
}
