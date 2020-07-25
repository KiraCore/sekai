package staking

import (
	"testing"

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

	handler := NewHandler(app.StakingKeeper, app.CustomStakingKeeper)

	theMsg := types2.MsgClaimValidator{
		Moniker:   "aMoniker",
		Website:   "some-web.com",
		Social:    "A Social",
		Identity:  "My Identity",
		Comission: types.NewDec(1234),
		ValKey:    valAddr1,
		PubKey:    addr1,
	}

	_, err = handler(ctx, theMsg)
	require.NoError(t, err)

	validatorSet := app.CustomStakingKeeper.GetValidatorSet(ctx)
	require.Len(t, validatorSet, 1)
	val := app.CustomStakingKeeper.GetValidator(ctx, valAddr1)

	validatorIsEqualThanClaimMsg(t, val, theMsg)
}

func validatorIsEqualThanClaimMsg(t *testing.T, val types2.Validator, msg types2.MsgClaimValidator) {
	require.Equal(t, msg.Moniker, val.Moniker)
	require.Equal(t, msg.PubKey, val.PubKey)
	require.Equal(t, msg.ValKey, val.ValKey)
	require.Equal(t, msg.Comission, val.Comission)
	require.Equal(t, msg.Identity, val.Identity)
	require.Equal(t, msg.Social, val.Social)
	require.Equal(t, msg.Website, val.Website)
}
