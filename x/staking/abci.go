package staking

import (
	"github.com/KiraCore/sekai/x/staking/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/encoding"
)

func EndBlocker(ctx sdk.Context, k keeper.Keeper) []abci.ValidatorUpdate {
	valSet := k.GetPendingValidatorSet(ctx)

	valUpdate := make([]abci.ValidatorUpdate, len(valSet))

	for i, val := range valSet {
		k.AddValidator(ctx, val)

		consPk, err := val.TmConsPubKey()
		if err != nil {
			panic(err)
		}

		pk, err := encoding.PubKeyToProto(consPk)
		if err != nil {
			panic(err)
		}

		valUpdate[i] = abci.ValidatorUpdate{
			Power:  1,
			PubKey: pk,
		}

		k.RemovePendingValidator(ctx, val)
	}

	return valUpdate
}
