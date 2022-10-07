package keeper

import (
	"github.com/KiraCore/sekai/x/gov/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strconv"
	"time"
)

func (k Keeper) GetNextPollIDAndIncrement(ctx sdk.Context) uint64 {
	pollID := k.GetNextPollID(ctx)
	k.SetNextPollID(ctx, pollID+1)
	return pollID
}

func (k Keeper) GetNextPollID(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(NextPollIDPrefix)
	if bz == nil {
		return 1
	}

	pollID := BytesToProposalID(bz)
	return pollID
}

func (k Keeper) SetNextPollID(ctx sdk.Context, pollID uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(NextPollIDPrefix, ProposalIDToBytes(pollID))
}

func (k Keeper) PollCreate(ctx sdk.Context, msg *types.MsgPollCreate) (uint64, error) {
	var roles []uint64
	pollID := k.GetNextPollIDAndIncrement(ctx)
	options := new(types.PollOptions)

	options.Type = msg.ValueType
	options.Count = msg.ValueCount
	options.Choices = msg.PossibleChoices

	for _, v := range msg.PollValues {
		options.Values = append(options.Values, v)
	}

	for _, sid := range msg.Roles {
		role, _ := k.GetRoleBySid(ctx, sid)
		roles = append(roles, uint64(role.Id))
	}

	poll, err := types.NewPoll(
		pollID,
		msg.Creator,
		msg.Title,
		msg.Description,
		msg.Reference,
		msg.Checksum,
		roles,
		options,
		msg.Expiry,
	)

	if err != nil {
		return pollID, err
	}

	k.SavePoll(ctx, poll)

	return pollID, nil
}

func (k Keeper) SavePoll(ctx sdk.Context, poll types.Poll) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&poll)

	store.Set(append(PollPrefix, ProposalIDToBytes(poll.PollId)...), bz)
}

func (k Keeper) AddAddressPoll(ctx sdk.Context, pollID uint64, address sdk.AccAddress) {
	addressPolls := types.AddressPolls{
		Address: address,
		Ids:     k.GetPollsIdsByAddress(ctx, address),
	}

	addressPolls.Ids = append(addressPolls.Ids, pollID)

	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&addressPolls)

	store.Set(append(PollPrefix, address...), bz)
}

func (k Keeper) PollVote(ctx sdk.Context, msg *types.MsgPollVote) error {
	var option types.PollVoteOption
	value, err := strconv.Atoi(msg.Value)

	if err != nil {
		option = types.PollOptionCustom
	} else {
		option = types.PollVoteOption(value)
	}

	vote := types.PollVote{
		Voter:       msg.Voter,
		PollId:      msg.PollId,
		Option:      option,
		CustomValue: msg.Value,
	}

	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&vote)
	store.Set(PollVoteKey(vote.PollId, vote.Voter), bz)

	return nil
}

func (k Keeper) GetPoll(ctx sdk.Context, pollID uint64) (*types.Poll, bool) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(append(PollPrefix, ProposalIDToBytes(pollID)...))
	if bz == nil {
		return &types.Poll{}, false
	}

	var poll *types.Poll
	k.cdc.MustUnmarshal(bz, poll)

	return poll, true
}

func (k Keeper) GetPollsIdsByAddress(ctx sdk.Context, address sdk.AccAddress) []uint64 {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), PollPrefix)
	var ids []uint64

	bz := prefixStore.Get(address.Bytes())
	if bz == nil {
		return ids
	}

	var addressPolls types.AddressPolls
	k.cdc.MustUnmarshal(bz, &addressPolls)

	return addressPolls.Ids
}

func (k Keeper) GetPollsByAddress(ctx sdk.Context, address sdk.AccAddress) ([]*types.Poll, bool) {
	ids := k.GetPollsIdsByAddress(ctx, address)
	var polls []*types.Poll

	for _, id := range ids {
		poll, found := k.GetPoll(ctx, id)
		if !found {
			return polls, false
		}

		polls = append(polls, poll)
	}

	return polls, true
}

func (k Keeper) GetPollVotes(ctx sdk.Context, pollID uint64) types.PollVotes {
	var votes types.PollVotes

	iterator := k.GetPollVotesIterator(ctx, pollID)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var vote types.PollVote
		k.cdc.MustUnmarshal(iterator.Value(), &vote)
		votes = append(votes, vote)
	}

	return votes
}

func (k Keeper) GetPollVotesIterator(ctx sdk.Context, pollID uint64) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, PollVotesKey(pollID))
}

func PollsByTimeKey(endTime time.Time) []byte {
	return append(PollPrefix, sdk.FormatTimeBytes(endTime)...)
}

// GetPollsWithFinishedVotingEndTimeIterator returns the proposals that have endtime finished.
func (k Keeper) GetPollsWithFinishedVotingEndTimeIterator(ctx sdk.Context, endTime time.Time) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return store.Iterator(PollPrefix, sdk.PrefixEndBytes(PollsByTimeKey(endTime)))
}

func PollVotesKey(pollId uint64) []byte {
	return append(PollVotesPrefix, ProposalIDToBytes(pollId)...)
}

func PollVoteKey(pollId uint64, address sdk.AccAddress) []byte {
	return append(PollVotesKey(pollId), address.Bytes()...)
}
