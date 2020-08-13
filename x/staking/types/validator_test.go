package types

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/KiraCore/cosmos-sdk/types"
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
				_, err := NewValidator(
					strings.Repeat("A", 65),
					"some-web.com",
					"some-web.com",
					"some-web.com",
					types.NewDec(1234),
					valAddr,
					pubKey,
				)

				return err
			},
			err: ErrInvalidMonikerLength,
		},
		{
			name:        "website longer than 64",
			expectError: true,
			newVal: func() error {
				_, err := NewValidator(
					"the moniker",
					strings.Repeat("A", 65),
					"some-web.com",
					"some-web.com",
					types.NewDec(1234),
					valAddr,
					pubKey,
				)

				return err
			},
			err: ErrInvalidWebsiteLength,
		},
		{
			name:        "social longer than 64",
			expectError: true,
			newVal: func() error {
				_, err := NewValidator(
					"the moniker",
					"some-web.com",
					strings.Repeat("A", 65),
					"some-web.com",
					types.NewDec(1234),
					valAddr,
					pubKey,
				)

				return err
			},
			err: ErrInvalidSocialLength,
		},
		{
			name:        "identity longer than 64",
			expectError: true,
			newVal: func() error {
				_, err := NewValidator(
					"the moniker",
					"some-web.com",
					"some-web.com",
					strings.Repeat("A", 65),
					types.NewDec(1234),
					valAddr,
					pubKey,
				)

				return err
			},
			err: ErrInvalidIdentityLength,
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
