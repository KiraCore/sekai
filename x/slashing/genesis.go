package slashing

import (
	"github.com/KiraCore/sekai/x/slashing/keeper"
	"github.com/KiraCore/sekai/x/slashing/types"
	stakingtypes "github.com/KiraCore/sekai/x/staking/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis initialize default parameters
// and the keeper's address to pubkey map
func InitGenesis(ctx sdk.Context, keeper keeper.Keeper, stakingKeeper types.StakingKeeper, data *types.GenesisState) {
	stakingKeeper.IterateValidators(ctx,
		func(index int64, validator *stakingtypes.Validator) bool {
			consPk, err := validator.ConsPubKey()
			if err != nil {
				panic(err)
			}
			keeper.AddPubkey(ctx, consPk)
			return false
		},
	)

	for _, info := range data.SigningInfos {
		address, err := sdk.ConsAddressFromBech32(info.Address)
		if err != nil {
			panic(err)
		}
		keeper.SetValidatorSigningInfo(ctx, address, info.ValidatorSigningInfo)
	}

	keeper.SetParams(ctx, data.Params)
}

// ExportGenesis writes the current store values
// to a genesis file, which can be imported again
// with InitGenesis
func ExportGenesis(ctx sdk.Context, keeper keeper.Keeper) (data *types.GenesisState) {
	params := keeper.GetParams(ctx)
	signingInfos := make([]types.SigningInfo, 0)
	keeper.IterateValidatorSigningInfos(ctx, func(address sdk.ConsAddress, info types.ValidatorSigningInfo) (stop bool) {
		bechAddr := address.String()
		signingInfos = append(signingInfos, types.SigningInfo{
			Address:              bechAddr,
			ValidatorSigningInfo: info,
		})

		return false
	})

	return types.NewGenesisState(params, signingInfos)
}
