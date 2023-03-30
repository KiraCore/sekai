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
	KeyPrefixDapp                   = []byte("dapp_info")
	PrefixUserDappBondKey           = []byte("dapp_user_bond")
	PrefixDappOperatorCandidateKey  = []byte("dapp_operator_candidate_key")
	PrefixDappOperatorKey           = []byte("dapp_operator_key")
	PrefixDappSessionKey            = []byte("dapp_session_key")
	PrefixDappSessionApprovalKey    = []byte("dapp_session_approval_key")
	PrefixDappLeaderDenouncementKey = []byte("dapp_leader_denouncement_key")
)

func DappKey(name string) []byte {
	return append(KeyPrefixDapp, name...)
}

func UserDappBondKey(dappName string, user string) []byte {
	return append(append(PrefixUserDappBondKey, dappName...), user...)
}

func DappOperatorCandidateKey(dappName string, operator string) []byte {
	return append(append(PrefixDappOperatorCandidateKey, dappName...), operator...)
}

func DappOperatorKey(dappName string, operator string) []byte {
	return append(append(PrefixDappOperatorKey, dappName...), operator...)
}

func ExecutionRegistrarKey(dappName string) []byte {
	return append(PrefixDappSessionKey, dappName...)
}

func DappSessionApprovalKey(dappName string, leader string) []byte {
	return append(append(PrefixDappSessionApprovalKey, dappName...), leader...)
}

func DappLeaderDenouncementKey(dappName string, leader string, sender string) []byte {
	return append(append(append(PrefixDappLeaderDenouncementKey, dappName...), leader...), sender...)
}
