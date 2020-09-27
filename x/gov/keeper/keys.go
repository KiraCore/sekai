package keeper

// Keys for Gov store.
// 0x00<uint64 in bytes> : The next proposalID.
// 0x01<proposalID_bytes> : The Proposal
//
// 0x10<role_uint64_Bytes> : The role permissions.
//
// 0x20<councilorAddress_Bytes> : NetworkActor.
//
// 0x30<actorAddress_Bytes> : Councilor.
var (
	NextProposalIDPrefix = []byte{0x00}
	ProposalsPrefix      = []byte{0x01}

	RolePermissionRegistry          = []byte{0x10}
	CouncilorIdentityRegistryPrefix = []byte{0x20}
	NetworkActorsPrefix             = []byte{0x30}
)
