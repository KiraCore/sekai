package staking_test

import (
	"testing"

	"github.com/KiraCore/sekai/x/staking"

	"github.com/KiraCore/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/KiraCore/sekai/simapp"
	types2 "github.com/KiraCore/sekai/x/staking/types"
)

func TestNewHandler_MsgClaimValidator_HappyPath(t *testing.T) {
	addr1, err := types.AccAddressFromBech32("kira15ky9du8a2wlstz6fpx3p4mqpjyrm5cgqzp4f3d")
	require.NoError(t, err)
	valAddr1, err := types.ValAddressFromBech32("kiravaloper15ky9du8a2wlstz6fpx3p4mqpjyrm5cgq38f2fp")
	require.NoError(t, err)

	app := simapp.Setup(false)
	ctx := app.NewContext(false, abci.Header{})

	handler := staking.NewHandler(app.StakingKeeper, app.CustomStakingKeeper)

	theMsg, err := types2.NewMsgClaimValidator(
		"aMoniker",
		"some-web.com",
		"A Social",
		"My Identity",
		types.NewDec(1234),
		valAddr1,
		addr1,
	)
	require.NoError(t, err)

	_, err = handler(ctx, theMsg)
	require.NoError(t, err)

	validatorSet := app.CustomStakingKeeper.GetValidatorSet(ctx)
	require.Len(t, validatorSet, 1)
	val := app.CustomStakingKeeper.GetValidator(ctx, valAddr1)

	validatorIsEqualThanClaimMsg(t, val, theMsg)
}

func validatorIsEqualThanClaimMsg(t *testing.T, val types2.Validator, msg *types2.MsgClaimValidator) {
	require.Equal(t, msg.Moniker, val.Moniker)
	require.Equal(t, msg.PubKey, val.PubKey)
	require.Equal(t, msg.ValKey, val.ValKey)
	require.Equal(t, msg.Comission, val.Comission)
	require.Equal(t, msg.Identity, val.Identity)
	require.Equal(t, msg.Social, val.Social)
	require.Equal(t, msg.Website, val.Website)
}
