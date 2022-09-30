package keeper

import (
	"context"

	kiratypes "github.com/KiraCore/sekai/types"
	govtypes "github.com/KiraCore/sekai/x/gov/types"
	"github.com/KiraCore/sekai/x/slashing/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
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

	properties := k.gk.GetNetworkProperties(ctx)
	validators := k.sk.GetValidatorSet(ctx)
	if len(validators) <= int(properties.MinValidators) || len(validators) <= 1 {
		return nil, types.ErrPauseNotAllowedOnPoorNetwork
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

func (k msgServer) RefuteSlashingProposal(goCtx context.Context, msg *types.MsgRefuteSlashingProposal) (*types.MsgRefuteSlashingProposalResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	proposals, _ := k.gk.GetProposals(ctx)
	for _, proposal := range proposals {
		if proposal.Result == govtypes.Pending && proposal.GetContent().ProposalType() == kiratypes.ProposalTypeSlashValidator {
			content := proposal.GetContent().(*types.ProposalSlashValidator)
			if content.Offender == msg.Validator {
				content.Refutation = msg.Refutation
				any, err := codectypes.NewAnyWithValue(content)
				if err != nil {
					return nil, err
				}

				proposal.Content = any
				k.gk.SaveProposal(ctx, proposal)
				return &types.MsgRefuteSlashingProposalResponse{}, nil
			}
		}
	}
	return nil, types.ErrSlashProposalDoesNotExists
}
