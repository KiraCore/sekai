package types_test

import (
	"os"
	"testing"

	"github.com/KiraCore/sekai/app"
	ixptypes "github.com/KiraCore/sekai/x/ixp/types"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	app.SetConfig()
	os.Exit(m.Run())
}

func TestMsgCreateOrderbook_ValidateBasic(t *testing.T) {

	kiraAddr1, err := types.AccAddressFromBech32("kira1da22wd7slpxpptasczs679mr5c8xtucqdzxc3n")
	require.NoError(t, err)

	tests := []struct {
		name          string
		constructor   func() (*ixptypes.MsgCreateOrderBook, error)
		validationErr string
	}{
		{
			name: "basic path test",
			constructor: func() (*ixptypes.MsgCreateOrderBook, error) {
				return ixptypes.NewMsgCreateOrderBook("base", "quote", "mnemonic", kiraAddr1)
			},
			validationErr: "",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			msg, err := tt.constructor()
			require.Equal(t, err, nil)

			err = msg.ValidateBasic()
			if tt.validationErr != "" {
				require.NotEqual(t, err, nil)
				require.Contains(t, err.Error(), tt.validationErr)
			} else {
				require.Equal(t, err, nil)
			}
		})
	}
}
