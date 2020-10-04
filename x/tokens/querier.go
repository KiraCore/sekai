package tokens

import (
	"github.com/KiraCore/sekai/x/tokens/keeper"
	"github.com/KiraCore/sekai/x/tokens/types"
)

type Querier struct {
	keeper keeper.Keeper
}

func NewQuerier(keeper keeper.Keeper) types.QueryServer {
	return &Querier{keeper: keeper}
}
