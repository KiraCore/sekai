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
		Key:       "",
		Value:     "",
		Date:      time.Now(),
		Verifiers: []sdk.AccAddress{},
	}

	bz, err := record.Marshal()
	require.NoError(t, err)

	parsed := IdentityRecord{}
	err = parsed.Unmarshal(bz)
	require.NoError(t, err)
}

func TestMarshalUnmarshalMsgRegisterIdentityRecords(t *testing.T) {
	record := MsgRegisterIdentityRecords{
		Address: sdk.AccAddress{},
		Infos:   WrapInfos(make(map[string]string)),
	}

	bz, err := record.Marshal()
	require.NoError(t, err)

	parsed := MsgRegisterIdentityRecords{}
	err = parsed.Unmarshal(bz)
	require.NoError(t, err)
}
