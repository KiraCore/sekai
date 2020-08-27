package staking_test

import (
	"testing"

	types3 "github.com/KiraCore/sekai/x/gov/types"

	"github.com/KiraCore/sekai/app"

	"github.com/KiraCore/sekai/x/staking"

	"github.com/KiraCore/sekai/simapp"
	types2 "github.com/KiraCore/sekai/x/staking/types"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

func TestMain(m *testing.M) {
	app.SetConfig()
	m.Run()
}

func TestNewHandler_MsgClaimValidator_HappyPath(t *testing.T) {
	valAddr1, err := types.ValAddressFromBech32("kiravaloper15ky9du8a2wlstz6fpx3p4mqpjyrm5cgq38f2fp")
	require.NoError(t, err)

	pubKey, err := types.GetPubKeyFromBech32(types.Bech32PubKeyTypeConsPub, "kiravalconspub1zcjduepqylc5k8r40azmw0xt7hjugr4mr5w2am7jw77ux5w6s8hpjxyrjjsq4xg7em")
	require.NoError(t, err)

	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	// First we give user permissions
	networkActor := types3.NewNetworkActor(
		types.AccAddress(valAddr1),
		nil,
		1,
		nil,
		types3.NewPermissions([]types3.PermValue{
			types3.PermClaimValidator,
		}, nil),
		1,
	)
	app.CustomGovKeeper.SaveNetworkActor(ctx, networkActor)

	handler := staking.NewHandler(app.CustomStakingKeeper, app.CustomGovKeeper)

	theMsg, err := types2.NewMsgClaimValidator(
		"aMoniker",
		"some-web.com",
		"A Sociale",
		"My Identity",
		types.NewDec(1234),
		valAddr1,
		pubKey,
	)
	require.NoError(t, err)

	_, err = handler(ctx, theMsg)
	require.NoError(t, err)

	validatorSet := app.CustomStakingKeeper.GetValidatorSet(ctx)
	require.Len(t, validatorSet, 1)
	val, err := app.CustomStakingKeeper.GetValidator(ctx, valAddr1)
	require.NoError(t, err)

	validatorIsEqualThanClaimMsg(t, val, theMsg)
}

func TestNewHandler_MsgClaimValidator_ItFailsIfUserDoesNotHavePermissionsToClaimValidator(t *testing.T) {
	valAddr1, err := types.ValAddressFromBech32("kiravaloper15ky9du8a2wlstz6fpx3p4mqpjyrm5cgq38f2fp")
	require.NoError(t, err)

	pubKey, err := types.GetPubKeyFromBech32(types.Bech32PubKeyTypeConsPub, "kiravalconspub1zcjduepqylc5k8r40azmw0xt7hjugr4mr5w2am7jw77ux5w6s8hpjxyrjjsq4xg7em")
	require.NoError(t, err)

	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	handler := staking.NewHandler(app.CustomStakingKeeper, app.CustomGovKeeper)

	theMsg, err := types2.NewMsgClaimValidator(
		"aMoniker",
		"some-web.com",
		"A Social",
		"My Identity",
		types.NewDec(1234),
		valAddr1,
		pubKey,
	)
	require.NoError(t, err)

	_, err = handler(ctx, theMsg)
	require.EqualError(t, err, "network actor not found")
}

func validatorIsEqualThanClaimMsg(t *testing.T, val types2.Validator, msg *types2.MsgClaimValidator) {
	require.Equal(t, msg.Moniker, val.Moniker)
	require.Equal(t, msg.PubKey, val.PubKey)
	require.Equal(t, msg.ValKey, val.ValKey)
	require.Equal(t, msg.Commission, val.Commission)
	require.Equal(t, msg.Identity, val.Identity)
	require.Equal(t, msg.Social, val.Social)
	require.Equal(t, msg.Website, val.Website)
}
