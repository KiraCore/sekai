package queries

import (
	"fmt"

	sdk "github.com/KiraCore/cosmos-sdk/types"
	"github.com/KiraCore/sekai/x/kiraHub/keeper"
	"github.com/KiraCore/sekai/x/kiraHub/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func QueryListSignerKeys(ctx sdk.Context, path []string, req abci.RequestQuery, keeper keeper.Keeper) ([]byte, error) {

	var queryOutput []types.SignerKey
	curator, err := sdk.AccAddressFromBech32(path[0])
	if err != nil {
		return []byte{}, fmt.Errorf("Invalid curator address %s: %+v", path[0], err)
	}

	queryOutput = keeper.GetSignerKeys(ctx, curator)

	res, err := types.ModuleCdc.MarshalJSON(queryOutput)
	if err != nil {
		panic(err)
	}

	return res, nil
}
