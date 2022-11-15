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
	_ = ctx
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
	_ = ctx
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
	_ = ctx
	return &types.MsgDonateCollectiveResponse{}, nil
}

// WithdrawCollective can be sent by any whitelisted “contributor” to withdraw
// their tokens (unless locking is enabled)
func (k msgServer) WithdrawCollective(
	goCtx context.Context,
	msg *types.MsgWithdrawCollective,
) (*types.MsgWithdrawCollectiveResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	_ = ctx
	return &types.MsgWithdrawCollectiveResponse{}, nil
}
