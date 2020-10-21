package types

import "github.com/cosmos/cosmos-sdk/types"

func NewNetworkActor(
	addr types.AccAddress,
	roles Roles,
	status ActorStatus,
	votes []uint32,
	perm *Permissions,
	skin uint64,
) NetworkActor {
	return NetworkActor{
		Address:     addr,
		Roles:       roles,
		Status:      status,
		Votes:       votes,
		Permissions: perm,
		Skin:        skin,
	}
}

func (m *NetworkActor) HasRole(role Role) bool {
	for _, r := range m.Roles {
		if r == uint64(role) {
			return true
		}
	}
	return false
}

func (m *NetworkActor) SetRole(role Role) {
	if !m.HasRole(role) {
		m.Roles = append(m.Roles, uint64(role))
	}
}

func (m *NetworkActor) RemoveRole(role Role) {
	for i, r := range m.Roles {
		if r == uint64(role) {
			m.Roles = append(m.Roles[:i], m.Roles[i+1:]...)
			return
		}
	}
}

func (m *NetworkActor) IsActive() bool {
	return m.Status == Active
}

// NewDefaultActor returns a default actor with:
// - The provided addr.
// - Roles set to nil
// - Status set to 0
// - Votes set to nil
// - Empty permissions
// - Skin set to 0
func NewDefaultActor(addr types.AccAddress) NetworkActor {
	return NewNetworkActor(
		addr,
		nil,
		Undefined,
		nil,
		NewPermissions(nil, nil),
		0,
	)
}
