package types

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestMarshalUnmarshalIdentityRecord(t *testing.T) {
	record := IdentityRecord{
		Id:        1,
		Address:   sdk.AccAddress{},
		Infos:     make(map[string]string),
		Date:      time.Now(),
		Verifiers: []sdk.AccAddress{},
	}

	bz, err := record.Marshal()
	require.NoError(t, err)

	parsed := IdentityRecord{}
	err = parsed.Unmarshal(bz)
	require.NoError(t, err)
}

func TestMarshalUnmarshalMsgCreateIdentityRecord(t *testing.T) {
	record := MsgCreateIdentityRecord{
		Address: sdk.AccAddress{},
		Infos:   WrapInfos(make(map[string]string)),
		Date:    time.Now(),
	}

	bz, err := record.Marshal()
	require.NoError(t, err)

	parsed := MsgCreateIdentityRecord{}
	err = parsed.Unmarshal(bz)
	require.NoError(t, err)
}
