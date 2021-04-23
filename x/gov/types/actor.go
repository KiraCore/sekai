package types

import "github.com/cosmos/cosmos-sdk/types"

func NewNetworkActor(
	addr types.AccAddress,
	roles Roles,
	status ActorStatus,
	votes []VoteOption,
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

func (m *NetworkActor) IsInactive() bool {
	return m.Status == Inactive
}

// Deactivate the actor
func (m *NetworkActor) Deactivate() {
	m.Status = Inactive
}

// CanVote returns if the actor can vote a specific vote option.
func (m *NetworkActor) CanVote(voteOption VoteOption) bool {
	for _, v := range m.Votes {
		if v == voteOption {
			return true
		}
	}

	return false
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
		Active,
		[]VoteOption{
			OptionYes,
			OptionNo,
			OptionAbstain,
			OptionNoWithVeto,
		},
		NewPermissions(nil, nil),
		0,
	)
}

// GetActorsWithVoteWithVeto returns the actors that have permission to vote with Veto from a list of network actors.
func GetActorsWithVoteWithVeto(actors []NetworkActor) []NetworkActor {
	var actorsWithVeto []NetworkActor

	for _, actor := range actors {
		if actor.CanVote(OptionNoWithVeto) {
			actorsWithVeto = append(actorsWithVeto, actor)
		}
	}

	return actorsWithVeto
}
