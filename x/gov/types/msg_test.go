package types

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func TestMsgWhitelistPermissions_ValidateBasic(t *testing.T) {
	tests := []struct {
		name        string
		msg         *MsgWhitelistPermissions
		expectedErr *sdkerrors.Error
	}{
		{
			name: "empty proposer addr",
			msg: NewMsgWhitelistPermissions(
				sdk.AccAddress{},
				sdk.AccAddress("some addr"),
				0,
			),
			expectedErr: ErrEmptyProposerAccAddress,
		},
		{
			name: "empty addr",
			msg: NewMsgWhitelistPermissions(
				sdk.AccAddress("some addr"),
				sdk.AccAddress{},
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

func TestMsgRequestIdentityRecordsVerify_ValidateBasic(t *testing.T) {
	addr1 := sdk.AccAddress("foo1________________")
	addr3 := sdk.AccAddress("foo3________________")
	empty := sdk.AccAddress("")
	msg := MsgRequestIdentityRecordsVerify{
		Address:   addr1,
		Verifier:  addr3,
		RecordIds: []uint64{1},
		Tip:       sdk.Coin{},
	}
	require.Error(t, msg.ValidateBasic())

	msg = MsgRequestIdentityRecordsVerify{
		Address:   addr1,
		Verifier:  addr3,
		RecordIds: []uint64{1},
		Tip:       sdk.NewInt64Coin(sdk.DefaultBondDenom, 0),
	}
	require.NoError(t, msg.ValidateBasic())

	msg = MsgRequestIdentityRecordsVerify{
		Address:   addr1,
		Verifier:  addr3,
		RecordIds: []uint64{},
		Tip:       sdk.NewInt64Coin(sdk.DefaultBondDenom, 10),
	}
	require.Error(t, msg.ValidateBasic())

	msg = MsgRequestIdentityRecordsVerify{
		Address:   addr1,
		Verifier:  empty,
		RecordIds: []uint64{},
		Tip:       sdk.NewInt64Coin(sdk.DefaultBondDenom, 10),
	}
	require.Error(t, msg.ValidateBasic())
	msg = MsgRequestIdentityRecordsVerify{
		Address:   empty,
		Verifier:  addr1,
		RecordIds: []uint64{},
		Tip:       sdk.NewInt64Coin(sdk.DefaultBondDenom, 10),
	}
	require.Error(t, msg.ValidateBasic())

	msg = MsgRequestIdentityRecordsVerify{
		Address:   addr1,
		Verifier:  addr3,
		RecordIds: []uint64{1},
		Tip:       sdk.NewInt64Coin(sdk.DefaultBondDenom, 10),
	}
	require.NoError(t, msg.ValidateBasic())
}
