package types

// constants
const (
	ModuleName = "layer2"
	// RouterKey to be used for routing msgs
	RouterKey = ModuleName
	// QuerierRoute is the querier route for the layer2 module
	QuerierRoute = ModuleName
)

// constants
var (
	KeyPrefixDapp         = []byte("dapp_info")
	PrefixUserDappBondKey = []byte("dapp_user_bond")
)

func DappKey(name string) []byte {
	return append(KeyPrefixDapp, name...)
}

func UserDappBondKey(dappName string, user string) []byte {
	return append(append(PrefixUserDappBondKey, dappName...), user...)
}
