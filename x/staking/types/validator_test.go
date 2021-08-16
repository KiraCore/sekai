package types_test

import (
	"strings"
	"testing"

	stakingtypes "github.com/KiraCore/sekai/x/staking/types"
	"github.com/cosmos/cosmos-sdk/simapp"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestNewValidator_Errors(t *testing.T) {
	valAddr, err := types.ValAddressFromBech32("kiravaloper1q24436yrnettd6v4eu6r4t9gycnnddac9nwqv0")
	require.NoError(t, err)

	pubkeys := simapp.CreateTestPubKeys(1)
	pubKey := pubkeys[0]

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
				_, err := stakingtypes.NewValidator(
					strings.Repeat("A", 65),
					types.NewDec(1234),
					valAddr,
					pubKey,
				)

				return err
			},
			err: stakingtypes.ErrInvalidMonikerLength,
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

	pubkeys := simapp.CreateTestPubKeys(1)
	pubKey := pubkeys[0]

	validator, err := stakingtypes.NewValidator(
		"the moniker",
		types.NewDec(1234),
		valAddr,
		pubKey,
	)
	require.NoError(t, err)
	require.True(t, validator.IsActive())
}
