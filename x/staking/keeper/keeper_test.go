package keeper

import (
	"testing"

	"github.com/KiraCore/cosmos-sdk/codec"
	"github.com/KiraCore/cosmos-sdk/store"
	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tm-db"

	types2 "github.com/KiraCore/cosmos-sdk/types"
	"github.com/KiraCore/sekai/x/staking/types"
	"github.com/stretchr/testify/require"
)

func TestKeeper_AddValidator(t *testing.T) {
	valAddr, err := types2.ValAddressFromBech32("kiravaloper1q24436yrnettd6v4eu6r4t9gycnnddac9nwqv0")
	require.NoError(t, err)

	accAddr := types2.AccAddress(valAddr)

	validator, err := types.NewValidator(
		"aMoniker",
		"some-web.com",
		"A Social",
		"My Identity",
		types2.NewDec(1234),
		valAddr,
		accAddr,
	)
	require.NoError(t, err)

	key := types2.NewKVStoreKey(types.ModuleName)
	db := dbm.NewMemDB()
	cms := store.NewCommitMultiStore(db)
	cms.MountStoreWithDB(key, types2.StoreTypeIAVL, nil)
	ctx := types2.NewContext(cms, abci.Header{}, false, nil)

	keeper := NewKeeper(key, codec.New())
	keeper.AddValidator(ctx, validator)
}
