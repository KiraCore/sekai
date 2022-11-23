package types

// constants
var (
	ModuleName            = "collectives"
	DonationModuleAccount = "donation_module_account"

	// RouterKey to be used for routing msgs
	RouterKey    = ModuleName
	QuerierRoute = ModuleName

	PrefixCollectiveKey            = []byte("collective_by_name")
	PrefixCollectiveContributerKey = []byte("collective_contributer")
)

func CollectiveKey(collectiveName string) []byte {
	return append([]byte(PrefixCollectiveKey), collectiveName...)
}

func CollectiveContributerKey(collectiveName, contributer string) []byte {
	return append(append([]byte(PrefixCollectiveContributerKey), collectiveName...), contributer...)
}
