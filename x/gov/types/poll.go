package types

import (
	"github.com/KiraCore/sekai/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"time"
)

// NewPoll creates a new poll proposal
func NewPoll(pollID uint64, creator sdk.AccAddress, title string, description string, reference string, checksum string, roles []uint64, options *PollOptions, expire time.Time) (Poll, error) {
	return Poll{
		PollId:        pollID,
		Creator:       creator,
		Title:         title,
		Description:   description,
		Reference:     reference,
		Checksum:      checksum,
		Roles:         roles,
		Options:       options,
		VotingEndTime: expire,
		Result:        PollPending,
	}, nil
}

func (m *Poll) Route() string {
	return ModuleName
}

func (m *Poll) Type() string {
	return types.MsgTypeCreatePoll
}

func (m *Poll) ValidateBasic() error {
	return nil
}

func (m *Poll) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m *Poll) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{
		m.Creator,
	}
}

func (m *AddressPolls) Route() string {
	return ModuleName
}

func (m *AddressPolls) Type() string {
	return types.MsgTypeAddressPoll
}

func (m *AddressPolls) ValidateBasic() error {
	return nil
}

func (m *AddressPolls) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m *AddressPolls) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{
		m.Address,
	}
}
