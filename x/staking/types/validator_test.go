package types_test

import (
	"strings"
	"testing"

	customstakingtypes "github.com/KiraCore/sekai/x/staking/types"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/types"
)

func TestNewValidator_Errors(t *testing.T) {
	valAddr, err := types.ValAddressFromBech32("kiravaloper1q24436yrnettd6v4eu6r4t9gycnnddac9nwqv0")
	require.NoError(t, err)

	pubKey, err := types.GetPubKeyFromBech32(types.Bech32PubKeyTypeConsPub, "kiravalconspub1zcjduepqylc5k8r40azmw0xt7hjugr4mr5w2am7jw77ux5w6s8hpjxyrjjsq4xg7em")
	require.NoError(t, err)

	tests := []struct {
		name        string
		expectError bool
		newVal      func() error
		err         error
	}{
		{
			name:        "moniker longer than 64",
			expectError: true,
			newVal: func() error {
				_, err := customstakingtypes.NewValidator(
					strings.Repeat("A", 65),
					types.NewDec(1234),
					valAddr,
					pubKey,
				)

				return err
			},
			err: customstakingtypes.ErrInvalidMonikerLength,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			err := tt.newVal()

			if tt.expectError {
				require.EqualError(t, err, tt.err.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestNewValidator_IsActiveByDefault(t *testing.T) {
	valAddr, err := types.ValAddressFromBech32("kiravaloper1q24436yrnettd6v4eu6r4t9gycnnddac9nwqv0")
	require.NoError(t, err)

	pubKey, err := types.GetPubKeyFromBech32(types.Bech32PubKeyTypeConsPub, "kiravalconspub1zcjduepqylc5k8r40azmw0xt7hjugr4mr5w2am7jw77ux5w6s8hpjxyrjjsq4xg7em")
	require.NoError(t, err)

	validator, err := customstakingtypes.NewValidator(
		"the moniker",
		types.NewDec(1234),
		valAddr,
		pubKey,
	)
	require.NoError(t, err)
	require.True(t, validator.IsActive())
}
