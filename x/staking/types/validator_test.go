package types_test

import (
	"testing"

	stakingtypes "github.com/KiraCore/sekai/x/staking/types"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/types"
)

func TestNewValidator_IsActiveByDefault(t *testing.T) {
	valAddr, err := types.ValAddressFromBech32("kiravaloper1q24436yrnettd6v4eu6r4t9gycnnddac9nwqv0")
	require.NoError(t, err)

	pubKey, err := types.GetPubKeyFromBech32(types.Bech32PubKeyTypeConsPub, "kiravalconspub1zcjduepqylc5k8r40azmw0xt7hjugr4mr5w2am7jw77ux5w6s8hpjxyrjjsq4xg7em")
	require.NoError(t, err)

	validator, err := stakingtypes.NewValidator(
		valAddr,
		pubKey,
	)
	require.NoError(t, err)
	require.True(t, validator.IsActive())
}
