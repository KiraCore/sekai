package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// constants
const (
	ModuleName = "layer2"
	// RouterKey to be used for routing msgs
	RouterKey = ModuleName
	// QuerierRoute is the querier route for the layer2 module
	QuerierRoute = ModuleName
	// StoreKey is the store key string for layer2 module
	StoreKey = ModuleName
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
	BridgeRegistrarHelperKey        = []byte("bridge_registrar_helper")
	PrefixBridgeAccountKey          = []byte("bridge_account_key")
	PrefixBridgeTokenKey            = []byte("bridge_token_key")
	PrefixXAMKey                    = []byte("xam_key")
	PrefixTokenInfoKey              = []byte("token_info")
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

func DappSessionApprovalKey(dappName string, verifier string) []byte {
	return append(append(PrefixDappSessionApprovalKey, dappName...), verifier...)
}

func DappLeaderDenouncementKey(dappName string, leader string, sender string) []byte {
	return append(append(append(PrefixDappLeaderDenouncementKey, dappName...), leader...), sender...)
}

func BridgeAccountKey(index uint64) []byte {
	return append(PrefixBridgeAccountKey, sdk.Uint64ToBigEndian(index)...)
}

func BridgeTokenKey(index uint64) []byte {
	return append(PrefixBridgeTokenKey, sdk.Uint64ToBigEndian(index)...)
}

func XAMKey(xid uint64) []byte {
	return append(PrefixXAMKey, sdk.Uint64ToBigEndian(xid)...)
}

func TokenInfoKey(denom string) []byte {
	return append(PrefixTokenInfoKey, denom...)
}
