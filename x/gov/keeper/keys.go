package keeper

// Keys for Gov store.
// 0x00<uint64 in bytes> : The next proposalID.
// 0x01<proposalID_bytes> : The Proposal
// 0x02<proposalID_Bytes + voterAddress_Bytes> : The vote
// 0x03<endTime_Bytes + proposalID_Bytes> : ActiveProposalID
// 0x04<EnactmentEndTime_Bytes + proposalID_Bytes> : ProposalID Holds all the proposals that are in process of enactment.
//
// 0x10<role_uint64_Bytes> : The role permissions.
//
// 0x20<councilorAddress_Bytes> : Councilor.
//
// 0x30<actorAddress_Bytes> : NetworkActor.
// 0x31<permissionID_Bytes + actor_address_bytes> : NetworkActorAddress. This is used to get all the actors that have a specific whitelist.
// 0x32<roleID_Bytes + actor_address_bytes> : NetworkActorAddress. This is used to get all the actors that have a specific role.
// 0x33<permissionID_Bytes + RoleID_Bytes> : RoleID_Bytes. This is used to get all the actors that have a specific role.
//
// 0x40<key_Bytes> : DataRegistryEntry
var (
	NextProposalIDPrefix     = []byte{0x00}
	ProposalsPrefix          = []byte{0x01}
	VotesPrefix              = []byte{0x02}
	ActiveProposalsPrefix    = []byte{0x03}
	EnactmentProposalsPrefix = []byte{0x04}

	RolePermissionRegistry          = []byte{0x10}
	CouncilorIdentityRegistryPrefix = []byte{0x20}

	NetworkActorsPrefix  = []byte{0x30}
	WhitelistActorPrefix = []byte{0x31}
	RoleActorPrefix      = []byte{0x32}
	WhitelistRolePrefix  = []byte{0x33}

	DataRegistryPrefix = []byte{0x40}

	PoorNetworkMessagesPrefix = []byte{0x41}
)
