package keeper_test

import (
	"testing"

	types2 "github.com/KiraCore/cosmos-sdk/types"
	"github.com/magiconair/properties/assert"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/KiraCore/sekai/simapp"
	"github.com/KiraCore/sekai/x/staking/types"
)

func TestKeeper_AddValidator(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.NewContext(false, abci.Header{})

	addrs := simapp.AddTestAddrsIncremental(app, ctx, 1, types2.TokensFromConsensusPower(10))
	addr1 := addrs[0]
	valAddr := types2.ValAddress(addr1)
	pubKey, err := types2.GetPubKeyFromBech32(types2.Bech32PubKeyTypeConsPub, "kiravalconspub1zcjduepqylc5k8r40azmw0xt7hjugr4mr5w2am7jw77ux5w6s8hpjxyrjjsq4xg7em")
	require.NoError(t, err)

	validator, err := types.NewValidator(
		"aMoniker",
		"some-web.com",
		"A Social",
		"My Identity",
		types2.NewDec(1234),
		valAddr,
		pubKey,
	)
	require.NoError(t, err)

	app.CustomStakingKeeper.AddValidator(ctx, validator)
	getValidator := app.CustomStakingKeeper.GetValidator(ctx, validator.ValKey)

	assert.Equal(t, validator, getValidator)
}

func TestKeeper_GetValidatorSet(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.NewContext(false, abci.Header{})

	addrs := simapp.AddTestAddrsIncremental(app, ctx, 2, types2.TokensFromConsensusPower(10))
	addr1 := addrs[0]
	valAddr1 := types2.ValAddress(addr1)

	addr2 := addrs[1]
	valAddr2 := types2.ValAddress(addr2)

	pubKey, err := types2.GetPubKeyFromBech32(types2.Bech32PubKeyTypeConsPub, "kiravalconspub1zcjduepqylc5k8r40azmw0xt7hjugr4mr5w2am7jw77ux5w6s8hpjxyrjjsq4xg7em")
	require.NoError(t, err)

	validator1, err := types.NewValidator(
		"validator 1",
		"some-web.com",
		"A Social",
		"My Identity",
		types2.NewDec(1234),
		valAddr1,
		pubKey,
	)
	require.NoError(t, err)

	validator2, err := types.NewValidator(
		"validator 2",
		"some-web.com",
		"A Social",
		"My Identity",
		types2.NewDec(1234),
		valAddr2,
		pubKey,
	)
	require.NoError(t, err)

	app.CustomStakingKeeper.AddValidator(ctx, validator1)
	app.CustomStakingKeeper.AddValidator(ctx, validator2)

	validatorSet := app.CustomStakingKeeper.GetValidatorSet(ctx)
	require.Equal(t, 2, len(validatorSet))
}
