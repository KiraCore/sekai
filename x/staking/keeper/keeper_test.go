package keeper_test

import (
	"testing"

	"github.com/magiconair/properties/assert"

	types2 "github.com/KiraCore/cosmos-sdk/types"
	"github.com/KiraCore/sekai/simapp"
	"github.com/KiraCore/sekai/x/staking/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
)

func TestKeeper_AddValidator(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.NewContext(false, abci.Header{})

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

	keeper := app.CustomStakingKeeper
	keeper.AddValidator(ctx, validator)

	getValidator := keeper.GetValidator(ctx, validator.ValKey)

	assert.Equal(t, validator, getValidator)
}
