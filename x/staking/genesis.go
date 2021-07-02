package staking

import (
	"github.com/KiraCore/sekai/x/staking/keeper"
	"github.com/KiraCore/sekai/x/staking/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

func ExportGenesis(ctx sdk.Context, keeper keeper.Keeper) *types.GenesisState {
	return &types.GenesisState{
		Validators: keeper.GetValidatorSet(ctx),
	}
}

// WriteValidators returns a slice of bonded genesis validators.
func WriteValidators(ctx sdk.Context, keeper keeper.Keeper) (vals []tmtypes.GenesisValidator, err error) {
	keeper.IterateLastValidators(ctx, func(_ int64, validator types.Validator) (stop bool) {
		pk, err := validator.ConsPubKey()
		if err != nil {
			return true
		}
		tmPk, err := cryptocodec.ToTmPubKeyInterface(pk)
		if err != nil {
			return true
		}

		vals = append(vals, tmtypes.GenesisValidator{
			Address: sdk.ConsAddress(tmPk.Address()).Bytes(),
			PubKey:  tmPk,
			Name:    validator.GetMoniker(),
		})

		return false
	})

	return
}
