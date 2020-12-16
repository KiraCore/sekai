package staking_test

import (
	"os"
	"testing"

	customgovtypes "github.com/KiraCore/sekai/x/gov/types"

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
	os.Exit(m.Run())
}

func TestNewHandler_MsgClaimValidator_HappyPath(t *testing.T) {
	valAddr1, err := types.ValAddressFromBech32("kiravaloper15ky9du8a2wlstz6fpx3p4mqpjyrm5cgq38f2fp")
	require.NoError(t, err)

	pubKey, err := types.GetPubKeyFromBech32(types.Bech32PubKeyTypeConsPub, "kiravalconspub1zcjduepqylc5k8r40azmw0xt7hjugr4mr5w2am7jw77ux5w6s8hpjxyrjjsq4xg7em")
	require.NoError(t, err)

	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	// First we give user permissions
	networkActor := customgovtypes.NewNetworkActor(
		types.AccAddress(valAddr1),
		nil,
		1,
		nil,
		customgovtypes.NewPermissions([]customgovtypes.PermValue{
			customgovtypes.PermClaimValidator,
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

	validatorSet := app.CustomStakingKeeper.GetPendingValidatorSet(ctx)
	require.Len(t, validatorSet, 0)

	_, err = handler(ctx, theMsg)
	require.NoError(t, err)

	validatorSet = app.CustomStakingKeeper.GetPendingValidatorSet(ctx)
	require.Len(t, validatorSet, 1)
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
	require.EqualError(t, err, "PermClaimValidator: not enough permissions")
}

func TestNewHandler_SetPermissions_ActorWithRole(t *testing.T) {
	valAddr1, err := types.ValAddressFromBech32("kiravaloper15ky9du8a2wlstz6fpx3p4mqpjyrm5cgq38f2fp")
	require.NoError(t, err)

	pubKey, err := types.GetPubKeyFromBech32(types.Bech32PubKeyTypeConsPub, "kiravalconspub1zcjduepqylc5k8r40azmw0xt7hjugr4mr5w2am7jw77ux5w6s8hpjxyrjjsq4xg7em")
	require.NoError(t, err)

	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	// Save network actor With Role Validator
	networkActor := customgovtypes.NewDefaultActor(types.AccAddress(valAddr1))
	networkActor.SetRole(customgovtypes.RoleValidator)
	app.CustomGovKeeper.SaveNetworkActor(ctx, networkActor)

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
	require.NoError(t, err)

	validatorSet := app.CustomStakingKeeper.GetPendingValidatorSet(ctx)
	require.Len(t, validatorSet, 1)
}
