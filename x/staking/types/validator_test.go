package types

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/KiraCore/cosmos-sdk/types"
)

func TestNewValidator_Errors(t *testing.T) {
	valAddr, err := types.ValAddressFromHex()
	require.NoError(t, err)

	accAddr, err := types.AccAddressFromHex()
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
					accAddr,
				)

				return err
			},
			err: nil,
		},
	}
}
