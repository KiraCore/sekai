package keeper

import (
	"github.com/KiraCore/sekai/x/gov/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
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

	store.Set(GetPollKey(poll.PollId), bz)
}

func (k Keeper) PollVote(ctx sdk.Context, msg *types.MsgPollVote) error {
	return nil
}

func (k Keeper) GetPoll(ctx sdk.Context, pollID uint64) (types.Poll, bool) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(GetPollKey(pollID))
	if bz == nil {
		return types.Poll{}, false
	}

	var prop types.Poll
	k.cdc.MustUnmarshal(bz, &prop)

	return prop, true
}

func GetPollKey(pollID uint64) []byte {
	return append(PollPrefix, ProposalIDToBytes(pollID)...)
}
