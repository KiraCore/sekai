package types

const ModuleName = "customgov"

var (
	KeyPrefixPermissionsRegistry = []byte("permissions_registry")

	// Roles
	RoleCouncilor Role = []byte{0x1}
	RoleValidator Role = []byte{0x2}
	RoleGovLeader Role = []byte{0x3}

	PermClaimValidator uint32 = 1
)

type Role []byte
