package keeper

import (
	"encoding/binary"
	"fmt"
	"github.com/KiraCore/sekai/x/gov/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
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

	pollID := sdk.BigEndianToUint64(bz)
	return pollID
}

func (k Keeper) SetNextPollID(ctx sdk.Context, pollID uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(NextPollIDPrefix, sdk.Uint64ToBigEndian(pollID))
}

func (k Keeper) PollCreate(ctx sdk.Context, msg *types.MsgPollCreate) (uint64, error) {
	var roles []uint64
	pollID := k.GetNextPollIDAndIncrement(ctx)
	options := new(types.PollOptions)

	options.Type = msg.ValueType
	options.Count = msg.ValueCount
	options.Choices = msg.PossibleChoices

	duration, err := time.ParseDuration(msg.Duration)
	if err != nil {
		return pollID, fmt.Errorf("invalid duration: %w", err)
	}

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
		time.Now().Add(duration),
	)

	if err != nil {
		return pollID, err
	}

	k.SavePoll(ctx, poll)
	k.AddAddressPoll(ctx, poll)
	k.AddToActivePolls(ctx, poll)

	return pollID, nil
}

func (k Keeper) SavePoll(ctx sdk.Context, poll types.Poll) {
	store := ctx.KVStore(k.storeKey)
	key := append(PollPrefix, sdk.Uint64ToBigEndian(poll.PollId)...)
	store.Set(key, k.cdc.MustMarshal(&poll))
}

func (k Keeper) AddAddressPoll(ctx sdk.Context, poll types.Poll) {
	store := ctx.KVStore(k.storeKey)
	addressKey := append(PollPrefix, poll.Creator.Bytes()...)
	key := append(addressKey, sdk.Uint64ToBigEndian(poll.PollId)...)
	store.Set(key, sdk.Uint64ToBigEndian(poll.PollId))
}

func (k Keeper) AddToActivePolls(ctx sdk.Context, poll types.Poll) {
	store := ctx.KVStore(k.storeKey)
	key := append(PollsByTimeKey(poll.VotingEndTime), sdk.Uint64ToBigEndian(poll.PollId)...)
	store.Set(key, sdk.Uint64ToBigEndian(poll.PollId))
}

func (k Keeper) RemoveActivePoll(ctx sdk.Context, poll types.Poll) {
	store := ctx.KVStore(k.storeKey)
	key := append(PollsByTimeKey(poll.VotingEndTime), sdk.Uint64ToBigEndian(poll.PollId)...)
	store.Delete(key)
}

func (k Keeper) PollVote(ctx sdk.Context, msg *types.MsgPollVote) error {
	vote := types.PollVote{
		Voter:       msg.Voter,
		PollId:      msg.PollId,
		Option:      msg.Option,
		CustomValue: msg.Value,
	}

	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&vote)
	store.Set(PollVoteKey(vote.PollId, vote.Voter), bz)

	return nil
}

func (k Keeper) GetPoll(ctx sdk.Context, pollID uint64) (types.Poll, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(append(PollPrefix, sdk.Uint64ToBigEndian(pollID)...))

	if bz == nil {
		return types.Poll{}, types.ErrPollsNotFount
	}

	var poll types.Poll
	k.cdc.MustUnmarshal(bz, &poll)

	return poll, nil
}

func (k Keeper) GetPollsIdsByAddress(ctx sdk.Context, address sdk.AccAddress) []uint64 {
	var ids []uint64
	store := ctx.KVStore(k.storeKey)
	addressKey := append(PollPrefix, address.Bytes()...)
	iterator := sdk.KVStorePrefixIterator(store, addressKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		id := binary.BigEndian.Uint64(iterator.Value())
		ids = append(ids, id)
	}

	return ids
}

func (k Keeper) GetPollsByAddress(ctx sdk.Context, address sdk.AccAddress) ([]types.Poll, error) {
	ids := k.GetPollsIdsByAddress(ctx, address)
	var polls []types.Poll

	for _, v := range ids {
		poll, err := k.GetPoll(ctx, v)
		if err != nil {
			return nil, err
		}
		polls = append(polls, poll)
	}

	return polls, nil
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

func (k Keeper) GetPollsWithFinishedVotingEndTimeIterator(ctx sdk.Context, endTime time.Time) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return store.Iterator(ActivePollPrefix, sdk.PrefixEndBytes(PollsByTimeKey(endTime)))
}

func PollsByTimeKey(endTime time.Time) []byte {
	return append(ActivePollPrefix, sdk.FormatTimeBytes(endTime)...)
}

func PollVotesKey(pollId uint64) []byte {
	return append(PollVotesPrefix, sdk.Uint64ToBigEndian(pollId)...)
}

func PollVoteKey(pollId uint64, address sdk.AccAddress) []byte {
	return append(PollVotesKey(pollId), address.Bytes()...)
}
