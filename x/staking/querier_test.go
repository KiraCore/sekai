package staking_test

import (
	"testing"

	"github.com/KiraCore/sekai/x/staking"

	types2 "github.com/KiraCore/sekai/x/staking/types"

	"github.com/KiraCore/cosmos-sdk/types"
	"github.com/KiraCore/sekai/simapp"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
)

func TestQuerier_ValidatorByAddress(t *testing.T) {
	addr1, err := types.AccAddressFromBech32("kira15ky9du8a2wlstz6fpx3p4mqpjyrm5cgqzp4f3d")
	require.NoError(t, err)
	valAddr1, err := types.ValAddressFromBech32("kiravaloper15ky9du8a2wlstz6fpx3p4mqpjyrm5cgq38f2fp")
	require.NoError(t, err)

	app := simapp.Setup(false)
	ctx := app.NewContext(false, abci.Header{})

	val, err := types2.NewValidator("Moniker", "Website", "Social", "identity", types.NewDec(123), valAddr1, addr1)
	require.NoError(t, err)

	app.CustomStakingKeeper.AddValidator(ctx, val)

	querier := staking.NewQuerier(app.CustomStakingKeeper)

	qValidatorResp, err := querier.ValidatorByAddress(types.WrapSDKContext(ctx), &types2.ValidatorByAddressRequest{ValAddr: valAddr1})
	require.NoError(t, err)

	require.Equal(t, val, qValidatorResp.Validator)
}
