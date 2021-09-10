package keeper_test

import (
	"testing"
	"time"

	simapp "github.com/KiraCore/sekai/app"
	govtypes "github.com/KiraCore/sekai/x/gov/types"
	"github.com/KiraCore/sekai/x/slashing/types"
	stakingtypes "github.com/KiraCore/sekai/x/staking/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

func TestQuerier_SigningInfo(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	consAddresses := []sdk.ConsAddress{}
	valCount := 1000
	for i := 0; i < valCount; i++ {
		valAddr := sdk.ValAddress(secp256k1.GenPrivKey().PubKey().Address())
		consPubKey := ed25519.GenPrivKey().PubKey()
		consAddress := sdk.ConsAddress(consPubKey.Address())
		consAddresses = append(consAddresses, consAddress)

		val, err := stakingtypes.NewValidator(valAddr, consPubKey)
		require.NoError(t, err)
		actor := govtypes.NewDefaultActor(sdk.AccAddress(valAddr))
		app.CustomGovKeeper.AddWhitelistPermission(ctx, actor, govtypes.PermClaimValidator)
		app.CustomStakingKeeper.AddValidator(ctx, val)

		newInfo := types.NewValidatorSigningInfo(
			consAddress,
			int64(4),
			time.Unix(2, 0),
			int64(10),
			int64(10),
			int64(10),
		)
		app.CustomSlashingKeeper.SetValidatorSigningInfo(ctx, consAddress, newInfo)
	}

	resp, err := app.CustomSlashingKeeper.SigningInfo(sdk.WrapSDKContext(ctx), &types.QuerySigningInfoRequest{
		ConsAddress:      consAddresses[0].String(),
		IncludeValidator: true,
	})
	require.NoError(t, err)
	require.NotEqual(t, resp.ValSigningInfo.Address, "")
	require.NotEqual(t, resp.Validator.Address, "")

	resp, err = app.CustomSlashingKeeper.SigningInfo(sdk.WrapSDKContext(ctx), &types.QuerySigningInfoRequest{
		ConsAddress:      consAddresses[0].String(),
		IncludeValidator: false,
	})
	require.NoError(t, err)
	require.NotEqual(t, resp.ValSigningInfo.Address, "")
	require.Equal(t, resp.Validator.Address, "")

	resp, err = app.CustomSlashingKeeper.SigningInfo(sdk.WrapSDKContext(ctx), &types.QuerySigningInfoRequest{
		ConsAddress:      "",
		IncludeValidator: false,
	})
	require.Error(t, err)
}

func TestQuerier_SigningInfos(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	consAddresses := []sdk.ConsAddress{}
	valCount := 1000
	for i := 0; i < valCount; i++ {
		valAddr := sdk.ValAddress(secp256k1.GenPrivKey().PubKey().Address())
		consPubKey := ed25519.GenPrivKey().PubKey()
		consAddress := sdk.ConsAddress(consPubKey.Address())
		consAddresses = append(consAddresses, consAddress)

		val, err := stakingtypes.NewValidator(valAddr, consPubKey)
		require.NoError(t, err)
		actor := govtypes.NewDefaultActor(sdk.AccAddress(valAddr))
		app.CustomGovKeeper.AddWhitelistPermission(ctx, actor, govtypes.PermClaimValidator)
		app.CustomStakingKeeper.AddValidator(ctx, val)

		newInfo := types.NewValidatorSigningInfo(
			consAddress,
			int64(4),
			time.Unix(2, 0),
			int64(10),
			int64(10),
			int64(10),
		)
		app.CustomSlashingKeeper.SetValidatorSigningInfo(ctx, consAddress, newInfo)
	}

	resp, err := app.CustomSlashingKeeper.SigningInfos(sdk.WrapSDKContext(ctx), &types.QuerySigningInfosRequest{
		IncludeValidator: true,
	})
	require.NoError(t, err)
	require.Greater(t, len(resp.Validators), 0)
	require.Greater(t, len(resp.Info), 0)

	resp, err = app.CustomSlashingKeeper.SigningInfos(sdk.WrapSDKContext(ctx), &types.QuerySigningInfosRequest{
		IncludeValidator: false,
	})
	require.NoError(t, err)
	require.Equal(t, len(resp.Validators), 0)
	require.Greater(t, len(resp.Info), 0)
}
