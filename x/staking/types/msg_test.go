package types

import (
	"testing"

	"github.com/KiraCore/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestNewMsgClaimValidator_Validate(t *testing.T) {
	msg := NewMsgClaimValidator("me", "web", "social", "id", types.NewDec(10), nil, nil)
	err := msg.ValidateBasic()

	require.Error(t, err)
}

func TestMsgClaimValidator_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		val  MsgClaimValidator
	}{
		{
			name: "nil val key",
			val:  NewMsgClaimValidator("me", "web", "social", "id", types.NewDec(10), nil, nil),
		},
		{
			name: "nil val key",
			val:  NewMsgClaimValidator("me", "web", "social", "id", types.NewDec(10), nil, nil),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			err := tt.val.ValidateBasic()
			require.Error(t, err)
		})
	}
}
