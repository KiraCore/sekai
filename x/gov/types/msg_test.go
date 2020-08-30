package types

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
)

func TestMsgWhitelistPermissions_ValidateBasic(t *testing.T) {
	tests := []struct {
		name        string
		msg         *MsgWhitelistPermissions
		expectedErr *errors.Error
	}{
		{
			name: "empty proposer addr",
			msg: NewMsgWhitelistPermissions(
				types.AccAddress{},
				types.AccAddress("some addr"),
				0,
			),
			expectedErr: ErrEmptyProposerAccAddress,
		},
		{
			name: "empty addr",
			msg: NewMsgWhitelistPermissions(
				types.AccAddress("some addr"),
				types.AccAddress{},
				0,
			),
			expectedErr: ErrEmptyPermissionsAccAddress,
		},
	}
	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			require.Equal(t, test.expectedErr, test.msg.ValidateBasic())
		})
	}
}
