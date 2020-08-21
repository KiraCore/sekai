package types

const ModuleName = "customgov"

var (
	KeyPrefixPermissionsRegistry = []byte("permissions_registry")

	// Roles
	RoleCouncilor Role = []byte{0x1}
	RoleValidator Role = []byte{0x2}
	RoleGovLeader Role = []byte{0x3}

	PermClaimValidator PermValue = 1
)

// Role represents a Role in the registry.
type Role []byte

// PermValue represents a single permission value, like claim-role-validator.
type PermValue uint32
