package keeper

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/types"
)

func TestValidatorKey(t *testing.T) {
	valAddr := types.ValAddress("valAddr")
	valAddr2 := types.ValAddress("valAddr2")

	require.Equal(t, append([]byte{0x00}, valAddr.Bytes()...), GetValidatorKey(valAddr))
	require.Equal(t, append([]byte{0x00}, valAddr2.Bytes()...), GetValidatorKey(valAddr2))
}
