package types

import (
	"github.com/KiraCore/sekai/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewPollVote creates a new poll proposal
func NewPollVote(pollID uint64, voter sdk.AccAddress) (PollVote, error) {
	return PollVote{
		Voter: voter,
	}, nil
}

func (m *PollVote) Route() string {
	return ModuleName
}

func (m *PollVote) Type() string {
	return types.MsgTypeVotePoll
}

func (m *PollVote) ValidateBasic() error {
	return nil
}

func (m *PollVote) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m *PollVote) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{
		m.Voter,
	}
}
