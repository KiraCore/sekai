package staking

import (
	"github.com/KiraCore/sekai/x/staking/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/encoding"
)

func EndBlocker(ctx sdk.Context, k keeper.Keeper) []abci.ValidatorUpdate {
	var valUpdate []abci.ValidatorUpdate

	valSet := k.GetPendingValidatorSet(ctx)
	for _, val := range valSet {
		k.AddValidator(ctx, val)

		consPk, err := val.TmConsPubKey()
		if err != nil {
			panic(err)
		}

		pk, err := encoding.PubKeyToProto(consPk)
		if err != nil {
			panic(err)
		}

		valUpdate = append(valUpdate, abci.ValidatorUpdate{
			Power:  1,
			PubKey: pk,
		})

		k.RemovePendingValidator(ctx, val)
	}

	// Remove validators from the set, paused or inactivated.
	removeVals := k.GetRemovingValidatorSet(ctx)
	for _, val := range removeVals {
		validator, err := k.GetValidator(ctx, val)
		if err != nil {
			panic("validator not found")
		}

		consPk, err := validator.TmConsPubKey()
		if err != nil {
			panic(err)
		}

		pk, err := encoding.PubKeyToProto(consPk)
		if err != nil {
			panic(err)
		}

		valUpdate = append(valUpdate, abci.ValidatorUpdate{
			Power:  0,
			PubKey: pk,
		})
		k.RemoveRemovingValidator(ctx, validator)
	}

	// Remove validators from the set, paused or inactivated.
	reactivateVals := k.GetReactivatingValidatorSet(ctx)
	for _, val := range reactivateVals {
		validator, err := k.GetValidator(ctx, val)
		if err != nil {
			panic("validator not found")
		}

		consPk, err := validator.TmConsPubKey()
		if err != nil {
			panic(err)
		}

		pk, err := encoding.PubKeyToProto(consPk)
		if err != nil {
			panic(err)
		}

		valUpdate = append(valUpdate, abci.ValidatorUpdate{
			Power:  1,
			PubKey: pk,
		})
		k.RemoveRemovingValidator(ctx, validator)
	}

	return valUpdate
}
