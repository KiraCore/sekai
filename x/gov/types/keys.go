package types

const ModuleName = "customgov"

var (
	KeyPrefixPermissionsRegistry = []byte("permissions_registry")
	KeyPrefixActors              = []byte("permissions_registry")

	// Roles
	RoleCouncilor Role = 0x1
	RoleValidator Role = 0x2
	RoleGovLeader Role = 0x3

	PermClaimValidator      PermValue = 1
	PermClaimGovernanceSeat PermValue = 2
)

// Role represents a Role in the registry.
type Role uint64

// PermValue represents a single permission value, like claim-role-validator.
type PermValue uint32
