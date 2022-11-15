package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// constants
var (
	ModuleName = "collectives"

	// RouterKey to be used for routing msgs
	RouterKey    = ModuleName
	QuerierRoute = ModuleName

	PrefixCollectiveKey            = []byte("collective_by_name")
	PrefixCollectiveContributerKey = []byte("collective_contributer")
)

func CollectiveKey(collectiveName string) []byte {
	return append([]byte(PrefixCollectiveKey), collectiveName...)
}

func CollectiveContributerKey(collectiveName string, contributer sdk.AccAddress) []byte {
	return append(append([]byte(PrefixCollectiveContributerKey), collectiveName...), contributer...)
}
