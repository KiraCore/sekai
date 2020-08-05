package types

import (
	"testing"

	"github.com/KiraCore/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestMsgClaimValidator_ValidateBasic(t *testing.T) {
	addr1, err := types.AccAddressFromBech32("kira15ky9du8a2wlstz6fpx3p4mqpjyrm5cgqzp4f3d")
	require.NoError(t, err)
	valAddr1, err := types.ValAddressFromBech32("kiravaloper15ky9du8a2wlstz6fpx3p4mqpjyrm5cgq38f2fp")
	require.NoError(t, err)

	tests := []struct {
		name        string
		constructor func() (*MsgClaimValidator, error)
	}{
		{
			name: "nil val key",
			constructor: func() (*MsgClaimValidator, error) {
				return NewMsgClaimValidator("me", "web", "social", "id", types.NewDec(10), nil, addr1)
			},
		},
		{
			name: "nil pub key",
			constructor: func() (*MsgClaimValidator, error) {
				return NewMsgClaimValidator("me", "web", "social", "id", types.NewDec(10), valAddr1, nil)
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.constructor()
			require.Error(t, err)
		})
	}
}
