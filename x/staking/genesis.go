package staking

import (
	"github.com/KiraCore/sekai/x/staking/keeper"
	"github.com/KiraCore/sekai/x/staking/types"
	tmtypes "github.com/cometbft/cometbft/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func ExportGenesis(ctx sdk.Context, keeper keeper.Keeper) *types.GenesisState {
	return &types.GenesisState{
		Validators: keeper.GetValidatorSet(ctx),
	}
}

// WriteValidators returns a slice of bonded genesis validators.
func WriteValidators(ctx sdk.Context, k keeper.Keeper) (vals []tmtypes.GenesisValidator, err error) {
	k.IterateLastValidators(ctx, func(_ int64, validator types.Validator) (stop bool) {
		pk, err := validator.ConsPubKey()
		if err != nil {
			return true
		}
		tmPk, err := cryptocodec.ToTmPubKeyInterface(pk)
		if err != nil {
			return true
		}

		moniker := k.GetMonikerByAddress(ctx, sdk.AccAddress(validator.ValKey))
		vals = append(vals, tmtypes.GenesisValidator{
			Address: sdk.ConsAddress(tmPk.Address()).Bytes(),
			PubKey:  tmPk,
			Name:    moniker,
		})

		return false
	})

	return
}
