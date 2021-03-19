package gov

type ActorStatus int32

const (
	// Undefined status
	ActorStatus_UNDEFINED ActorStatus = 0
	// Unclaimed status
	ActorStatus_UNCLAIMED ActorStatus = 1
	// Active status
	ActorStatus_ACTIVE ActorStatus = 2
	// Paused status
	ActorStatus_PAUSED ActorStatus = 3
	// Inactive status
	ActorStatus_INACTIVE ActorStatus = 4
	// Jailed status
	ActorStatus_JAILED ActorStatus = 5
	// Removed status
	ActorStatus_REMOVED ActorStatus = 6
)

// Enum value maps for ActorStatus.
var (
	ActorStatus_name = map[int32]string{
		0: "UNDEFINED",
		1: "UNCLAIMED",
		2: "ACTIVE",
		3: "PAUSED",
		4: "INACTIVE",
		5: "JAILED",
		6: "REMOVED",
	}
	ActorStatus_value = map[string]int32{
		"UNDEFINED": 0,
		"UNCLAIMED": 1,
		"ACTIVE":    2,
		"PAUSED":    3,
		"INACTIVE":  4,
		"JAILED":    5,
		"REMOVED":   6,
	}
)

type Permissions struct {
	Blacklist []PermValue `json:"blacklist"`
	Whitelist []PermValue `json:"whitelist"`
}

type Voter struct {
	Address     string   `json:"address"`
	Roles       []string `json:"roles"`
	Status      string   `json:"status"`
	Votes       []string `json:"votes"`
	Permissions struct {
		Blacklist []string `json:"blacklist"`
		Whitelist []string `json:"whitelist"`
	} `json:"permissions"`
	Skin uint64 `json:"skin"`
}
