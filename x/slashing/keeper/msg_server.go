package keeper

import (
	"context"
	"fmt"
	"time"

	govkeeper "github.com/KiraCore/sekai/x/gov/keeper"
	govtypes "github.com/KiraCore/sekai/x/gov/types"
	"github.com/KiraCore/sekai/x/slashing/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/pkg/errors"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the slashing MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

// Activate implements MsgServer.Activate method.
// Validators must submit a transaction to activate itself after
// having been inactivated (and thus unbonded) for downtime
func (k msgServer) Activate(goCtx context.Context, msg *types.MsgActivate) (*types.MsgActivateResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	valAddr, valErr := sdk.ValAddressFromBech32(msg.ValidatorAddr)
	if valErr != nil {
		return nil, valErr
	}
	err := k.Keeper.Activate(ctx, valAddr)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.ValidatorAddr),
		),
	)

	return &types.MsgActivateResponse{}, nil
}

// Pause implements MsgServer.Pause method.
// Validators must submit a transaction to pause itself after
// having been paused (and thus unbonded) for downtime
func (k msgServer) Pause(goCtx context.Context, msg *types.MsgPause) (*types.MsgPauseResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	valAddr, valErr := sdk.ValAddressFromBech32(msg.ValidatorAddr)
	if valErr != nil {
		return nil, valErr
	}
	err := k.Keeper.Pause(ctx, valAddr)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.ValidatorAddr),
		),
	)

	return &types.MsgPauseResponse{}, nil
}

// Unpause implements MsgServer.Unpause method.
// Validators must submit a transaction to unpause itself after
// having been paused (and thus unbonded) for downtime
func (k msgServer) Unpause(goCtx context.Context, msg *types.MsgUnpause) (*types.MsgUnpauseResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	valAddr, valErr := sdk.ValAddressFromBech32(msg.ValidatorAddr)
	if valErr != nil {
		return nil, valErr
	}
	err := k.Keeper.Unpause(ctx, valAddr)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.ValidatorAddr),
		),
	)

	return &types.MsgUnpauseResponse{}, nil
}

func (k msgServer) ProposalResetWholeValidatorRank(goCtx context.Context, msg *types.MsgProposalResetWholeValidatorRank) (*types.MsgProposalResetWholeValidatorRankResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	govKeeper := k.gk.(govkeeper.Keeper)
	isAllowed := govkeeper.CheckIfAllowedPermission(ctx, govKeeper, msg.Proposer, govtypes.PermCreateResetWholeValidatorRankProposal)
	if !isAllowed {
		return nil, errors.Wrap(govtypes.ErrNotEnoughPermissions, govtypes.PermCreateResetWholeValidatorRankProposal.String())
	}

	proposalID, err := k.CreateAndSaveProposalWithContent(
		ctx,
		msg.Description,
		types.NewProposalResetWholeValidatorRank(
			msg.Proposer,
		),
	)
	if err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			govtypes.EventTypeSubmitProposal,
			sdk.NewAttribute(govtypes.AttributeKeyProposalId, fmt.Sprintf("%d", proposalID)),
			sdk.NewAttribute(govtypes.AttributeKeyProposalType, msg.Type()),
			sdk.NewAttribute(types.AttributeKeyDescription, msg.Description),
		),
	)
	return &types.MsgProposalResetWholeValidatorRankResponse{
		ProposalID: proposalID,
	}, nil
}

func (k msgServer) CreateAndSaveProposalWithContent(ctx sdk.Context, description string, content govtypes.Content) (uint64, error) {

	govKeeper := k.gk.(govkeeper.Keeper)
	blockTime := ctx.BlockTime()
	proposalID, err := govKeeper.GetNextProposalID(ctx)
	if err != nil {
		return 0, err
	}

	properties := govKeeper.GetNetworkProperties(ctx)

	proposal, err := govtypes.NewProposal(
		proposalID,
		content,
		blockTime,
		blockTime.Add(time.Second*time.Duration(properties.ProposalEndTime)),
		blockTime.Add(time.Second*time.Duration(properties.ProposalEndTime)+
			time.Second*time.Duration(properties.ProposalEnactmentTime),
		),
		ctx.BlockHeight()+int64(properties.MinProposalEndBlocks),
		ctx.BlockHeight()+int64(properties.MinProposalEndBlocks+properties.MinProposalEnactmentBlocks),
		description,
	)

	govKeeper.SaveProposal(ctx, proposal)
	govKeeper.AddToActiveProposals(ctx, proposal)

	return proposalID, nil
}
