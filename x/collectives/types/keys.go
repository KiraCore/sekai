package types

// constants
var (
	ModuleName = "collectives"

	// RouterKey to be used for routing msgs
	RouterKey    = ModuleName
	QuerierRoute = ModuleName
	StoreKey     = ModuleName

	PrefixCollectiveKey            = []byte("collective_by_name")
	PrefixCollectiveContributerKey = []byte("collective_contributer")
)

func CollectiveKey(collectiveName string) []byte {
	return append([]byte(PrefixCollectiveKey), collectiveName...)
}

func CollectiveContributerKey(collectiveName, contributer string) []byte {
	return append(append([]byte(PrefixCollectiveContributerKey), collectiveName...), contributer...)
}
