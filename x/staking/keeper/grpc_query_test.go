package keeper_test

import (
	"fmt"
	"testing"

	"github.com/KiraCore/sekai/simapp"
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
	pubKey, err := types.GetPubKeyFromBech32(types.Bech32PubKeyTypeConsPub, "kiravalconspub1zcjduepqylc5k8r40azmw0xt7hjugr4mr5w2am7jw77ux5w6s8hpjxyrjjsq4xg7em")
	require.NoError(t, err)

	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	val, err := stakingtypes.NewValidator("Moniker", types.NewDec(123), valAddr1, pubKey)
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

		moniker := fmt.Sprintf("Moniker_%d", i+1)
		val, err := stakingtypes.NewValidator(moniker, types.NewDec(123), valAddr, consPubKey)
		require.NoError(t, err)
		actor := govtypes.NewDefaultActor(sdk.AccAddress(valAddr))
		app.CustomGovKeeper.AddWhitelistPermission(ctx, actor, govtypes.PermClaimValidator)
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
