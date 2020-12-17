package keeper_test

import (
	"os"
	"testing"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	app2 "github.com/KiraCore/sekai/app"
	"github.com/KiraCore/sekai/simapp"
	"github.com/KiraCore/sekai/x/staking/types"
	types2 "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	app2.SetConfig()
	os.Exit(m.Run())
}

func TestKeeper_AddValidator(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

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

	// Get By Validator Address.
	getValidator, err := app.CustomStakingKeeper.GetValidator(ctx, validator.ValKey)
	require.NoError(t, err)
	require.True(t, validator.Equal(getValidator))

	// Non existing validator Addr.
	_, err = app.CustomStakingKeeper.GetValidator(ctx, types2.ValAddress("non existing"))
	require.EqualError(t, err, "validator not found")

	// Get by AccAddress.
	getValidator, err = app.CustomStakingKeeper.GetValidatorByAccAddress(ctx, addr1)
	require.NoError(t, err)
	require.True(t, validator.Equal(getValidator))

	// Non existing AccAddress.
	_, err = app.CustomStakingKeeper.GetValidatorByAccAddress(ctx, types2.AccAddress("non existing"))
	require.EqualError(t, err, "validator not found")

	// Get by Moniker.
	getValidator, err = app.CustomStakingKeeper.GetValidatorByMoniker(ctx, validator.Moniker)
	require.NoError(t, err)
	require.True(t, validator.Equal(getValidator))

	// Get by ConsAddress
	getByConsAddrValidator, err := app.CustomStakingKeeper.GetValidatorByConsAddr(ctx, validator.GetConsAddr())
	require.NoError(t, err)
	require.True(t, validator.Equal(getByConsAddrValidator))

	// Non existing moniker
	_, err = app.CustomStakingKeeper.GetValidatorByMoniker(ctx, "UnexistingMoniker")
	require.EqualError(t, err, "validator with moniker UnexistingMoniker not found")
}

func TestKeeper_GetValidatorSet(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

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

	// Iterate validators
	passed := 0
	app.CustomStakingKeeper.IterateValidators(ctx, func(index int64, validator *types.Validator) (stop bool) {
		passed++
		return false
	})

	require.Equal(t, 2, passed)
}

func TestKeeper_GetPendingValidators(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

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

	app.CustomStakingKeeper.AddPendingValidator(ctx, validator1)
	app.CustomStakingKeeper.AddPendingValidator(ctx, validator2)

	validatorSet := app.CustomStakingKeeper.GetPendingValidatorSet(ctx)
	require.Equal(t, 2, len(validatorSet))

	app.CustomStakingKeeper.RemovePendingValidator(ctx, validator1)
	validatorSet = app.CustomStakingKeeper.GetPendingValidatorSet(ctx)
	require.Equal(t, 1, len(validatorSet))
}
