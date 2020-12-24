package keeper_test

import (
	"testing"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	stakingkeeper "github.com/KiraCore/sekai/x/staking/keeper"

	types2 "github.com/KiraCore/sekai/x/staking/types"

	"github.com/KiraCore/sekai/simapp"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestQuerier_ValidatorByAddress(t *testing.T) {
	valAddr1, err := types.ValAddressFromBech32("kiravaloper15ky9du8a2wlstz6fpx3p4mqpjyrm5cgq38f2fp")
	require.NoError(t, err)
	pubKey, err := types.GetPubKeyFromBech32(types.Bech32PubKeyTypeConsPub, "kiravalconspub1zcjduepqylc5k8r40azmw0xt7hjugr4mr5w2am7jw77ux5w6s8hpjxyrjjsq4xg7em")
	require.NoError(t, err)

	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	val, err := types2.NewValidator("Moniker", "Website", "Social", "identity", types.NewDec(123), valAddr1, pubKey)
	require.NoError(t, err)

	app.CustomStakingKeeper.AddValidator(ctx, val)

	querier := stakingkeeper.NewQuerier(app.CustomStakingKeeper)

	qValidatorResp, err := querier.ValidatorByAddress(types.WrapSDKContext(ctx), &types2.ValidatorByAddressRequest{ValAddr: valAddr1})
	require.NoError(t, err)

	require.True(t, val.Equal(qValidatorResp.Validator))
}
