package types_test

import (
	"testing"

	"github.com/KiraCore/sekai/app"

	types2 "github.com/KiraCore/sekai/x/staking/types"

	"github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestMsgClaimValidator_ValidateBasic(t *testing.T) {
	app.SetConfig()
	valAddr1, err := types.ValAddressFromBech32("kiravaloper15ky9du8a2wlstz6fpx3p4mqpjyrm5cgq38f2fp")
	require.NoError(t, err)

	pubKey, err := types.GetPubKeyFromBech32(types.Bech32PubKeyTypeConsPub, "kiravalconspub1zcjduepqylc5k8r40azmw0xt7hjugr4mr5w2am7jw77ux5w6s8hpjxyrjjsq4xg7em")
	require.NoError(t, err)

	tests := []struct {
		name        string
		constructor func() (*types2.MsgClaimValidator, error)
	}{
		{
			name: "nil val key",
			constructor: func() (*types2.MsgClaimValidator, error) {
				return types2.NewMsgClaimValidator("me", "web", "social", "id", types.NewDec(10), nil, pubKey)
			},
		},
		{
			name: "nil pub key",
			constructor: func() (*types2.MsgClaimValidator, error) {
				return types2.NewMsgClaimValidator("me", "web", "social", "id", types.NewDec(10), valAddr1, nil)
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
