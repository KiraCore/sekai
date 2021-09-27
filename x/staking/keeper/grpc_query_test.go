package keeper_test

import (
	"fmt"
	"testing"

	simapp "github.com/KiraCore/sekai/app"
	govtypes "github.com/KiraCore/sekai/x/gov/types"
	stakingkeeper "github.com/KiraCore/sekai/x/staking/keeper"
	stakingtypes "github.com/KiraCore/sekai/x/staking/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/cosmos/cosmos-sdk/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

func TestQuerier_ValidatorByAddress(t *testing.T) {
	valAddr1, err := types.ValAddressFromBech32("kiravaloper15ky9du8a2wlstz6fpx3p4mqpjyrm5cgq38f2fp")
	require.NoError(t, err)
	pubkeys := simapp.CreateTestPubKeys(1)
	pubKey := pubkeys[0]

	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	val, err := stakingtypes.NewValidator(valAddr1, pubKey)
	require.NoError(t, err)

	app.CustomStakingKeeper.AddValidator(ctx, val)

	querier := stakingkeeper.NewQuerier(app.CustomStakingKeeper)

	qValidatorResp, err := querier.ValidatorByAddress(types.WrapSDKContext(ctx), &stakingtypes.ValidatorByAddressRequest{ValAddr: valAddr1})
	require.NoError(t, err)

	require.True(t, val.Equal(qValidatorResp.Validator))
}

func TestQuerier_Validators(t *testing.T) {

	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	valCount := 1000
	for i := 0; i < valCount; i++ {
		valAddr := sdk.ValAddress(secp256k1.GenPrivKey().PubKey().Address())
		consPubKey := ed25519.GenPrivKey().PubKey()

		val, err := stakingtypes.NewValidator(valAddr, consPubKey)
		require.NoError(t, err)
		actor := govtypes.NewDefaultActor(sdk.AccAddress(valAddr))
		app.CustomGovKeeper.AddWhitelistPermission(ctx, actor, govtypes.PermClaimValidator)
		moniker := fmt.Sprintf("Moniker_%d", i+1)
		app.CustomGovKeeper.RegisterIdentityRecords(ctx, sdk.AccAddress(val.ValKey), []govtypes.IdentityInfoEntry{{
			Key:  "moniker",
			Info: moniker,
		}})
		app.CustomStakingKeeper.AddValidator(ctx, val)
	}

	querier := stakingkeeper.NewQuerier(app.CustomStakingKeeper)

	resp, err := querier.Validators(types.WrapSDKContext(ctx), &stakingtypes.ValidatorsRequest{
		// no restriction to query all
	})
	require.NoError(t, err)

	require.Len(t, resp.Validators, 100)
	require.Len(t, resp.Actors, 1000)
	require.Equal(t, resp.Pagination.Total, uint64(1000))
	require.NotNil(t, resp.Pagination.NextKey)
}
