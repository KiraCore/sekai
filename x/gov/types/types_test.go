package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPermissions_IsBlacklisted(t *testing.T) {
	perms := Permissions{
		Blacklist: []uint32{
			uint32(PermClaimValidator),
		},
		Whitelist: nil,
	}

	require.True(t, perms.IsBlacklisted(PermClaimValidator))
}
