package teststaking

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/KiraCore/sekai/x/staking/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewValidator is a testing helper method to create validators in tests
func NewValidator(t *testing.T, operator sdk.ValAddress, pubKey cryptotypes.PubKey) types.Validator {
	v, err := types.NewValidator(operator, pubKey)
	require.NoError(t, err)
	return v
}
