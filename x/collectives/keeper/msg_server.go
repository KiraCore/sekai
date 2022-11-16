package keeper

import (
	"context"

	"github.com/KiraCore/sekai/x/collectives/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type msgServer struct {
	keeper Keeper
	cgk    types.CustomGovKeeper
}

// NewMsgServerImpl returns an implementation of the bank MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper, cgk types.CustomGovKeeper) types.MsgServer {
	return &msgServer{
		keeper: keeper,
		cgk:    cgk,
	}
}

var _ types.MsgServer = msgServer{}

// CreateCollective defines a method for creating collective.
// allow ANY user to create Staking Collective even if they have no roles or
// permissions enabling that
func (k msgServer) CreateCollective(
	goCtx context.Context,
	msg *types.MsgCreateCollective,
) (*types.MsgCreateCollectiveResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	collective := k.keeper.GetCollective(ctx, msg.Name)
	if collective.Name != "" {
		return nil, types.ErrCollectiveDoesNotExist
	}
	k.keeper.SetCollective(ctx, types.Collective{
		Name:             msg.Name,
		Description:      msg.Description,
		Status:           types.CollectiveActive,
		DepositWhitelist: msg.DepositWhitelist,
		OwnersWhitelist:  msg.OwnersWhitelist,
		SpendingPools:    msg.SpendingPools,
		ClaimStart:       msg.ClaimStart,
		ClaimPeriod:      msg.ClaimPeriod,
		ClaimEnd:         msg.ClaimEnd,
		VoteQuorum:       msg.VoteQuorum,
		VotePeriod:       msg.VotePeriod,
		VoteEnactment:    msg.VoteEnactment,
		Donations:        sdk.NewCoins(),
		Rewards:          sdk.NewCoins(),
		LastDistribution: uint64(ctx.BlockTime().Unix()),
	})

	return &types.MsgCreateCollectiveResponse{}, nil
}

// ContributeCollective defines a method for putting bonds on collective.
// can be sent by any whitelisted “contributor” account that wants to add
// tokens to the Staking Collective during or after creation process
func (k msgServer) ContributeCollective(
	goCtx context.Context,
	msg *types.MsgBondCollective,
) (*types.MsgBondCollectiveResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}
	err = k.keeper.bk.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, msg.Bonds)
	if err != nil {
		return nil, err
	}

	cc := k.keeper.GetCollectiveContributer(ctx, msg.Name, msg.Sender)
	if cc.Name != "" {
		cc.Bonds = sdk.NewCoins(cc.Bonds...).Add(msg.Bonds...)
	} else {
		cc = types.CollectiveContributor{
			Address:      msg.Sender,
			Name:         msg.Name,
			Bonds:        msg.Bonds,
			Locking:      0,
			Donation:     sdk.ZeroDec(),
			DonationLock: false,
		}
	}
	k.keeper.SetCollectiveContributer(ctx, cc)
	return &types.MsgBondCollectiveResponse{}, nil
}

// DonateCollective defines a method to set lock and donation for bonds on the
// collection - allows to lock staking derivatives for a specific time period
// and donating a defined percentage of staking rewards to the collective.
func (k msgServer) DonateCollective(
	goCtx context.Context,
	msg *types.MsgDonateCollective,
) (*types.MsgDonateCollectiveResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	cc := k.keeper.GetCollectiveContributer(ctx, msg.Name, msg.Sender)
	if cc.Name != "" {
		return nil, types.ErrNotCollectiveContributer
	}

	// todo: collective contributer should have donation lock start time
	if cc.DonationLock {
		return nil, types.ErrDonationLocked
	}

	cc.Locking = msg.Locking
	cc.Donation = msg.Donation

	k.keeper.SetCollectiveContributer(ctx, cc)
	return &types.MsgDonateCollectiveResponse{}, nil
}

// WithdrawCollective can be sent by any whitelisted “contributor” to withdraw
// their tokens (unless locking is enabled)
func (k msgServer) WithdrawCollective(
	goCtx context.Context,
	msg *types.MsgWithdrawCollective,
) (*types.MsgWithdrawCollectiveResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	cc := k.keeper.GetCollectiveContributer(ctx, msg.Name, msg.Sender)
	if cc.Name != "" {
		return nil, types.ErrNotCollectiveContributer
	}

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}
	err = k.keeper.bk.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sender, cc.Bonds)
	if err != nil {
		return nil, err
	}
	k.keeper.DeleteCollectiveContributer(ctx, msg.Name, msg.Sender)
	return &types.MsgWithdrawCollectiveResponse{}, nil
}
