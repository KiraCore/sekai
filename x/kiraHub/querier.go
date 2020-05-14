package kiraHub

import (
	sdk "github.com/KiraCore/cosmos-sdk/types"
	constants "github.com/KiraCore/sekai/x/kiraHub/constants"
	"github.com/KiraCore/cosmos-sdk/types/errors"
	abciTypes "github.com/tendermint/tendermint/abci/types"
)

func NewQuerier(keeper Keeper) sdk.Querier {
	return func(context sdk.Context, path []string, requestQuery abciTypes.RequestQuery) ([]byte, error) {
		switch path[0] {


		default:
			return nil, errors.Wrapf(constants.UnknownQueryCode, "%v", path[0])
		}
	}
}
