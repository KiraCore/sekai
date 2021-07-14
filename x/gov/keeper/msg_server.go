package keeper

import (
	"context"
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/KiraCore/sekai/x/gov/types"
)

type msgServer struct {
	keeper Keeper
}

// NewMsgServerImpl returns an implementation of the bank MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{
		keeper: keeper,
	}
}

var _ types.MsgServer = msgServer{}

func (k msgServer) VoteProposal(
	goCtx context.Context,
	msg *types.MsgVoteProposal,
) (*types.MsgVoteProposalResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	actor, found := k.keeper.GetNetworkActorByAddress(ctx, msg.Voter)
	if !found || !actor.IsActive() {
		return nil, types.ErrActorIsNotActive
	}

	proposal, found := k.keeper.GetProposal(ctx, msg.ProposalId)
	if !found {
		return nil, types.ErrProposalDoesNotExist
	}

	if proposal.VotingEndTime.Before(ctx.BlockTime()) {
		return nil, types.ErrVotingTimeEnded
	}

	isAllowed := CheckIfAllowedPermission(ctx, k.keeper, msg.Voter, proposal.GetContent().VotePermission())
	if !isAllowed {
		return nil, errors.Wrap(types.ErrNotEnoughPermissions, proposal.GetContent().VotePermission().String())
	}

	vote := types.NewVote(msg.ProposalId, msg.Voter, msg.Option)
	k.keeper.SaveVote(ctx, vote)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeProposalVote,
			sdk.NewAttribute(types.AttributeKeyProposalId, fmt.Sprintf("%d", msg.ProposalId)),
			sdk.NewAttribute(types.AttributeKeyVoter, msg.Voter.String()),
			sdk.NewAttribute(types.AttributeKeyOption, msg.Option.String()),
		),
	)
	return &types.MsgVoteProposalResponse{}, nil
}

func (k msgServer) ProposalSetPoorNetworkMessages(
	goCtx context.Context,
	msg *types.MsgProposalSetPoorNetworkMessages,
) (*types.MsgProposalSetPoorNetworkMessagesResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	isAllowed := CheckIfAllowedPermission(ctx, k.keeper, msg.Proposer, types.PermCreateSetPoorNetworkMessagesProposal)
	if !isAllowed {
		return nil, errors.Wrap(types.ErrNotEnoughPermissions, types.PermCreateSetPoorNetworkMessagesProposal.String())
	}

	proposalID, err := k.CreateAndSaveProposalWithContent(ctx, msg.Description, types.NewSetPoorNetworkMessagesProposal(msg.Messages))
	if err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeSubmitProposal,
			sdk.NewAttribute(types.AttributeKeyProposalId, fmt.Sprintf("%d", proposalID)),
			sdk.NewAttribute(types.AttributeKeyProposalType, msg.Type()),
			sdk.NewAttribute(types.AttributeKeyProposalContent, msg.String()),
		),
	)
	return &types.MsgProposalSetPoorNetworkMessagesResponse{
		ProposalID: proposalID,
	}, nil
}

func (k msgServer) ProposalUpsertDataRegistry(
	goCtx context.Context,
	msg *types.MsgProposalUpsertDataRegistry,
) (*types.MsgProposalUpsertDataRegistryResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	isAllowed := CheckIfAllowedPermission(ctx, k.keeper, msg.Proposer, types.PermCreateUpsertDataRegistryProposal)
	if !isAllowed {
		return nil, errors.Wrap(types.ErrNotEnoughPermissions, types.PermCreateUpsertDataRegistryProposal.String())
	}

	proposalID, err := k.CreateAndSaveProposalWithContent(ctx,
		msg.Description,
		types.NewUpsertDataRegistryProposal(
			msg.Key,
			msg.Hash,
			msg.Reference,
			msg.Encoding,
			msg.Size_,
		),
	)
	if err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeSubmitProposal,
			sdk.NewAttribute(types.AttributeKeyProposalId, fmt.Sprintf("%d", proposalID)),
			sdk.NewAttribute(types.AttributeKeyProposalType, msg.Type()),
			sdk.NewAttribute(types.AttributeKeyProposalContent, msg.String()),
		),
	)
	return &types.MsgProposalUpsertDataRegistryResponse{
		ProposalID: proposalID,
	}, nil
}

// CreateIdentityRecord defines a method to create identity record
func (k msgServer) CreateIdentityRecord(goCtx context.Context, msg *types.MsgCreateIdentityRecord) (*types.MsgCreateIdentityRecordResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	recordId := k.keeper.CreateIdentityRecord(ctx, msg.Address, msg.Infos, msg.Date)
	return &types.MsgCreateIdentityRecordResponse{
		RecordId: recordId,
	}, nil
}

func (k msgServer) EditIdentityRecord(goCtx context.Context, msg *types.MsgEditIdentityRecord) (*types.MsgEditIdentityRecordResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	err := k.keeper.EditIdentityRecord(ctx, msg.RecordId, msg.Address, msg.Infos, msg.Date)
	return &types.MsgEditIdentityRecordResponse{}, err
}

// RequestIdentityRecordsVerify defines a method to request verify request from specific verifier
func (k msgServer) RequestIdentityRecordsVerify(goCtx context.Context, msg *types.MsgRequestIdentityRecordsVerify) (*types.MsgRequestIdentityRecordsVerifyResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	requestId, err := k.keeper.RequestIdentityRecordsVerify(ctx, msg.Address, msg.Verifier, msg.RecordIds, msg.Tip)
	return &types.MsgRequestIdentityRecordsVerifyResponse{
		RequestId: requestId,
	}, err
}

// ApproveIdentityRecords defines a method to accept verification request
func (k msgServer) ApproveIdentityRecords(goCtx context.Context, msg *types.MsgApproveIdentityRecords) (*types.MsgApproveIdentityRecordsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	err := k.keeper.ApproveIdentityRecords(ctx, msg.Verifier, msg.VerifyRequestId)
	return &types.MsgApproveIdentityRecordsResponse{}, err
}

func (k msgServer) CancelIdentityRecordsVerifyRequest(goCtx context.Context, msg *types.MsgCancelIdentityRecordsVerifyRequest) (*types.MsgCancelIdentityRecordsVerifyRequestResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	err := k.keeper.CancelIdentityRecordsVerifyRequest(ctx, msg.Executor, msg.VerifyRequestId)
	return &types.MsgCancelIdentityRecordsVerifyRequestResponse{}, err
}

func (k msgServer) ProposalAssignPermission(
	goCtx context.Context,
	msg *types.MsgProposalAssignPermission,
) (*types.MsgProposalAssignPermissionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	isAllowed := CheckIfAllowedPermission(ctx, k.keeper, msg.Proposer, types.PermCreateSetPermissionsProposal)
	if !isAllowed {
		return nil, errors.Wrap(types.ErrNotEnoughPermissions, types.PermCreateSetPermissionsProposal.String())
	}

	actor, found := k.keeper.GetNetworkActorByAddress(ctx, msg.Address)
	if found { // Actor exists
		if actor.Permissions.IsWhitelisted(types.PermValue(msg.Permission)) {
			return nil, errors.Wrap(types.ErrWhitelisting, "permission already whitelisted")
		}
	}

	proposalID, err := k.CreateAndSaveProposalWithContent(
		ctx,
		msg.Description,
		types.NewAssignPermissionProposal(
			msg.Address,
			types.PermValue(msg.Permission),
		))
	if err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeSubmitProposal,
			sdk.NewAttribute(types.AttributeKeyProposalId, fmt.Sprintf("%d", proposalID)),
			sdk.NewAttribute(types.AttributeKeyProposalType, msg.Type()),
			sdk.NewAttribute(types.AttributeKeyProposalContent, msg.String()),
		),
	)
	return &types.MsgProposalAssignPermissionResponse{
		ProposalID: proposalID,
	}, nil
}

func (k msgServer) CreateAndSaveProposalWithContent(ctx sdk.Context, description string, content types.Content) (uint64, error) {
	blockTime := ctx.BlockTime()
	proposalID := k.keeper.GetNextProposalIDAndIncrement(ctx)

	properties := k.keeper.GetNetworkProperties(ctx)

	proposal, err := types.NewProposal(
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

	if err != nil {
		return proposalID, err
	}

	k.keeper.SaveProposal(ctx, proposal)
	k.keeper.AddToActiveProposals(ctx, proposal)

	return proposalID, nil
}

func (k msgServer) ProposalSetNetworkProperty(
	goCtx context.Context,
	msg *types.MsgProposalSetNetworkProperty,
) (*types.MsgProposalSetNetworkPropertyResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	isAllowed := CheckIfAllowedPermission(ctx, k.keeper, msg.Proposer, types.PermCreateSetNetworkPropertyProposal)
	if !isAllowed {
		return nil, errors.Wrap(types.ErrNotEnoughPermissions, types.PermCreateSetNetworkPropertyProposal.String())
	}

	property, err := k.keeper.GetNetworkProperty(ctx, msg.NetworkProperty)
	if err != nil {
		return nil, errors.Wrap(errors.ErrInvalidRequest, err.Error())
	}
	if property == msg.Value {
		return nil, errors.Wrap(errors.ErrInvalidRequest, "network property already set as proposed value")
	}

	proposalID, err := k.CreateAndSaveProposalWithContent(
		ctx,
		msg.Description,
		types.NewSetNetworkPropertyProposal(
			msg.NetworkProperty,
			msg.Value,
		),
	)
	if err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeSubmitProposal,
			sdk.NewAttribute(types.AttributeKeyProposalId, fmt.Sprintf("%d", proposalID)),
			sdk.NewAttribute(types.AttributeKeyProposalType, msg.Type()),
			sdk.NewAttribute(types.AttributeKeyProposalContent, msg.String()),
		),
	)

	return &types.MsgProposalSetNetworkPropertyResponse{
		ProposalID: proposalID,
	}, nil
}

func (k msgServer) ProposalCreateRole(goCtx context.Context, msg *types.MsgProposalCreateRole) (*types.MsgProposalCreateRoleResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if len(msg.WhitelistedPermissions) == 0 && len(msg.BlacklistedPermissions) == 0 {
		return nil, types.ErrEmptyPermissions
	}

	isAllowed := CheckIfAllowedPermission(ctx, k.keeper, msg.Proposer, types.PermCreateRoleProposal)
	if !isAllowed {
		return nil, errors.Wrap(types.ErrNotEnoughPermissions, types.PermCreateRoleProposal.String())
	}

	_, exists := k.keeper.GetPermissionsForRole(ctx, types.Role(msg.Role))
	if exists {
		return nil, types.ErrRoleExist
	}

	proposalID, err := k.CreateAndSaveProposalWithContent(
		ctx,
		msg.Description,
		types.NewCreateRoleProposal(
			types.Role(msg.Role),
			msg.WhitelistedPermissions,
			msg.BlacklistedPermissions,
		),
	)
	if err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeSubmitProposal,
			sdk.NewAttribute(types.AttributeKeyProposalId, fmt.Sprintf("%d", proposalID)),
			sdk.NewAttribute(types.AttributeKeyProposalType, msg.Type()),
			sdk.NewAttribute(types.AttributeKeyProposalContent, msg.String()),
		),
	)

	return &types.MsgProposalCreateRoleResponse{
		ProposalID: proposalID,
	}, nil
}

func (k msgServer) RemoveRole(
	goCtx context.Context,
	msg *types.MsgRemoveRole,
) (*types.MsgRemoveRoleResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	isAllowed := CheckIfAllowedPermission(ctx, k.keeper, msg.Proposer, types.PermUpsertRole)
	if !isAllowed {
		return nil, errors.Wrap(types.ErrNotEnoughPermissions, types.PermUpsertRole.String())
	}

	_, found := k.keeper.GetPermissionsForRole(ctx, types.Role(msg.Role))
	if !found {
		return nil, types.ErrRoleDoesNotExist
	}

	actor, found := k.keeper.GetNetworkActorByAddress(ctx, msg.Address)
	if !found {
		actor = types.NewDefaultActor(msg.Address)
	}

	if !actor.HasRole(types.Role(msg.Role)) {
		return nil, types.ErrRoleNotAssigned
	}

	k.keeper.RemoveRoleFromActor(ctx, actor, types.Role(msg.Role))
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeRemoveRole,
			sdk.NewAttribute(types.AttributeKeyProposer, msg.Proposer.String()),
			sdk.NewAttribute(types.AttributeKeyAddress, msg.Address.String()),
			sdk.NewAttribute(types.AttributeKeyRoleId, fmt.Sprintf("%d", msg.Role)),
		),
	)
	return &types.MsgRemoveRoleResponse{}, nil
}

func (k msgServer) AssignRole(
	goCtx context.Context,
	msg *types.MsgAssignRole,
) (*types.MsgAssignRoleResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	isAllowed := CheckIfAllowedPermission(ctx, k.keeper, msg.Proposer, types.PermUpsertRole)
	if !isAllowed {
		return nil, errors.Wrap(types.ErrNotEnoughPermissions, types.PermUpsertRole.String())
	}

	_, found := k.keeper.GetPermissionsForRole(ctx, types.Role(msg.Role))
	if !found {
		return nil, types.ErrRoleDoesNotExist
	}

	actor, found := k.keeper.GetNetworkActorByAddress(ctx, msg.Address)
	if !found {
		actor = types.NewDefaultActor(msg.Address)
	}

	if actor.HasRole(types.Role(msg.Role)) {
		return nil, types.ErrRoleAlreadyAssigned
	}

	k.keeper.AssignRoleToActor(ctx, actor, types.Role(msg.Role))
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeAssignRole,
			sdk.NewAttribute(types.AttributeKeyProposer, msg.Proposer.String()),
			sdk.NewAttribute(types.AttributeKeyAddress, msg.Address.String()),
			sdk.NewAttribute(types.AttributeKeyRoleId, fmt.Sprintf("%d", msg.Role)),
		),
	)
	return &types.MsgAssignRoleResponse{}, nil
}

func (k msgServer) CreateRole(
	goCtx context.Context,
	msg *types.MsgCreateRole,
) (*types.MsgCreateRoleResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	isAllowed := CheckIfAllowedPermission(ctx, k.keeper, msg.Proposer, types.PermUpsertRole)
	if !isAllowed {
		return nil, errors.Wrap(types.ErrNotEnoughPermissions, "PermUpsertRole")
	}

	_, found := k.keeper.GetPermissionsForRole(ctx, types.Role(msg.Role))
	if found {
		return nil, types.ErrRoleExist
	}

	k.keeper.CreateRole(ctx, types.Role(msg.Role))
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCreateRole,
			sdk.NewAttribute(types.AttributeKeyProposer, msg.Proposer.String()),
			sdk.NewAttribute(types.AttributeKeyRoleId, fmt.Sprintf("%d", msg.Role)),
		),
	)
	return &types.MsgCreateRoleResponse{}, nil
}

func (k msgServer) RemoveBlacklistRolePermission(
	goCtx context.Context,
	msg *types.MsgRemoveBlacklistRolePermission,
) (*types.MsgRemoveBlacklistRolePermissionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	isAllowed := CheckIfAllowedPermission(ctx, k.keeper, msg.Proposer, types.PermUpsertRole)
	if !isAllowed {
		return nil, errors.Wrap(types.ErrNotEnoughPermissions, types.PermUpsertRole.String())
	}

	err := k.keeper.RemoveBlacklistRolePermission(ctx, types.Role(msg.Role), types.PermValue(msg.Permission))
	if err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeRemoveBlacklistRolePermisison,
			sdk.NewAttribute(types.AttributeKeyProposer, msg.Proposer.String()),
			sdk.NewAttribute(types.AttributeKeyRoleId, fmt.Sprintf("%d", msg.Role)),
			sdk.NewAttribute(types.AttributeKeyPermission, fmt.Sprintf("%d", msg.Permission)),
		),
	)
	return &types.MsgRemoveBlacklistRolePermissionResponse{}, nil
}

func (k msgServer) RemoveWhitelistRolePermission(
	goCtx context.Context,
	msg *types.MsgRemoveWhitelistRolePermission,
) (*types.MsgRemoveWhitelistRolePermissionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	isAllowed := CheckIfAllowedPermission(ctx, k.keeper, msg.Proposer, types.PermUpsertRole)
	if !isAllowed {
		return nil, errors.Wrap(types.ErrNotEnoughPermissions, types.PermUpsertRole.String())
	}

	err := k.keeper.RemoveWhitelistRolePermission(ctx, types.Role(msg.Role), types.PermValue(msg.Permission))
	if err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeRemoveWhitelistRolePermisison,
			sdk.NewAttribute(types.AttributeKeyProposer, msg.Proposer.String()),
			sdk.NewAttribute(types.AttributeKeyRoleId, fmt.Sprintf("%d", msg.Role)),
			sdk.NewAttribute(types.AttributeKeyPermission, fmt.Sprintf("%d", msg.Permission)),
		),
	)
	return &types.MsgRemoveWhitelistRolePermissionResponse{}, nil
}

func (k msgServer) BlacklistRolePermission(
	goCtx context.Context,
	msg *types.MsgBlacklistRolePermission,
) (*types.MsgBlacklistRolePermissionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	isAllowed := CheckIfAllowedPermission(ctx, k.keeper, msg.Proposer, types.PermUpsertRole)
	if !isAllowed {
		return nil, errors.Wrap(types.ErrNotEnoughPermissions, types.PermUpsertRole.String())
	}

	err := k.keeper.BlacklistRolePermission(ctx, types.Role(msg.Role), types.PermValue(msg.Permission))
	if err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeBlacklistRolePermisison,
			sdk.NewAttribute(types.AttributeKeyProposer, msg.Proposer.String()),
			sdk.NewAttribute(types.AttributeKeyRoleId, fmt.Sprintf("%d", msg.Role)),
			sdk.NewAttribute(types.AttributeKeyPermission, fmt.Sprintf("%d", msg.Permission)),
		),
	)
	return &types.MsgBlacklistRolePermissionResponse{}, nil
}

func (k msgServer) WhitelistRolePermission(
	goCtx context.Context,
	msg *types.MsgWhitelistRolePermission,
) (*types.MsgWhitelistRolePermissionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	isAllowed := CheckIfAllowedPermission(ctx, k.keeper, msg.Proposer, types.PermUpsertRole)
	if !isAllowed {
		return nil, errors.Wrap(types.ErrNotEnoughPermissions, types.PermUpsertRole.String())
	}

	err := k.keeper.WhitelistRolePermission(ctx, types.Role(msg.Role), types.PermValue(msg.Permission))
	if err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeWhitelistRolePermisison,
			sdk.NewAttribute(types.AttributeKeyProposer, msg.Proposer.String()),
			sdk.NewAttribute(types.AttributeKeyRoleId, fmt.Sprintf("%d", msg.Role)),
			sdk.NewAttribute(types.AttributeKeyPermission, fmt.Sprintf("%d", msg.Permission)),
		),
	)
	return &types.MsgWhitelistRolePermissionResponse{}, nil
}

func (k msgServer) WhitelistPermissions(
	goCtx context.Context,
	msg *types.MsgWhitelistPermissions,
) (*types.MsgWhitelistPermissionsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	isAllowed := CheckIfAllowedPermission(ctx, k.keeper, msg.Proposer, types.PermSetPermissions)
	if !isAllowed {
		return nil, errors.Wrap(types.ErrNotEnoughPermissions, "PermSetPermissions")
	}

	actor, found := k.keeper.GetNetworkActorByAddress(ctx, msg.Address)
	if !found {
		actor = types.NewDefaultActor(msg.Address)
	}

	err := k.keeper.AddWhitelistPermission(ctx, actor, types.PermValue(msg.Permission))
	if err != nil {
		return nil, errors.Wrapf(types.ErrSetPermissions, "error setting %d to whitelist: %s", msg.Permission, err)
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeWhitelistPermisison,
			sdk.NewAttribute(types.AttributeKeyProposer, msg.Proposer.String()),
			sdk.NewAttribute(types.AttributeKeyRoleId, msg.Address.String()),
			sdk.NewAttribute(types.AttributeKeyPermission, fmt.Sprintf("%d", msg.Permission)),
		),
	)
	return &types.MsgWhitelistPermissionsResponse{}, nil
}

func (k msgServer) BlacklistPermissions(
	goCtx context.Context,
	msg *types.MsgBlacklistPermissions,
) (*types.MsgBlacklistPermissionsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	isAllowed := CheckIfAllowedPermission(ctx, k.keeper, msg.Proposer, types.PermSetPermissions)
	if !isAllowed {
		return nil, errors.Wrap(types.ErrNotEnoughPermissions, "PermSetPermissions")
	}

	actor, found := k.keeper.GetNetworkActorByAddress(ctx, msg.Address)
	if !found {
		actor = types.NewDefaultActor(msg.Address)
	}

	err := actor.Permissions.AddToBlacklist(types.PermValue(msg.Permission))
	if err != nil {
		return nil, errors.Wrapf(types.ErrSetPermissions, "error setting %d to whitelist", msg.Permission)
	}

	k.keeper.SaveNetworkActor(ctx, actor)
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeBlacklistPermisison,
			sdk.NewAttribute(types.AttributeKeyProposer, msg.Proposer.String()),
			sdk.NewAttribute(types.AttributeKeyRoleId, msg.Address.String()),
			sdk.NewAttribute(types.AttributeKeyPermission, fmt.Sprintf("%d", msg.Permission)),
		),
	)
	return &types.MsgBlacklistPermissionsResponse{}, nil
}

func (k msgServer) SetNetworkProperties(
	goCtx context.Context,
	msg *types.MsgSetNetworkProperties,
) (*types.MsgSetNetworkPropertiesResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	isAllowed := CheckIfAllowedPermission(ctx, k.keeper, msg.Proposer, types.PermChangeTxFee)
	if !isAllowed {
		return nil, errors.Wrap(types.ErrNotEnoughPermissions, "PermChangeTxFee")
	}
	k.keeper.SetNetworkProperties(ctx, msg.NetworkProperties)
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeSetNetworkProperties,
			sdk.NewAttribute(types.AttributeKeyProposer, msg.Proposer.String()),
			sdk.NewAttribute(types.AttributeKeyProperties, msg.NetworkProperties.String()),
		),
	)
	return &types.MsgSetNetworkPropertiesResponse{}, nil
}

func (k msgServer) SetExecutionFee(
	goCtx context.Context,
	msg *types.MsgSetExecutionFee,
) (*types.MsgSetExecutionFeeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	isAllowed := CheckIfAllowedPermission(ctx, k.keeper, msg.Proposer, types.PermChangeTxFee)
	if !isAllowed {
		return nil, errors.Wrap(types.ErrNotEnoughPermissions, "PermChangeTxFee")
	}

	k.keeper.SetExecutionFee(ctx, &types.ExecutionFee{
		Name:              msg.Name,
		TransactionType:   msg.TransactionType,
		ExecutionFee:      msg.ExecutionFee,
		FailureFee:        msg.FailureFee,
		Timeout:           msg.Timeout,
		DefaultParameters: msg.DefaultParameters,
	})
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeSetExecutionFee,
			sdk.NewAttribute(types.AttributeKeyProposer, msg.Proposer.String()),
			sdk.NewAttribute(types.AttributeKeyTransactionType, msg.TransactionType),
			sdk.NewAttribute(types.AttributeKeyExecutionFee, fmt.Sprintf("%d", msg.ExecutionFee)),
			sdk.NewAttribute(types.AttributeKeyFailureFee, fmt.Sprintf("%d", msg.FailureFee)),
			sdk.NewAttribute(types.AttributeKeyTimeout, fmt.Sprintf("%d", msg.FailureFee)),
			sdk.NewAttribute(types.AttributeKeyDefaultParameters, fmt.Sprintf("%d", msg.DefaultParameters)),
		),
	)
	return &types.MsgSetExecutionFeeResponse{}, nil
}

func (k msgServer) ClaimCouncilor(
	goCtx context.Context,
	msg *types.MsgClaimCouncilor,
) (*types.MsgClaimCouncilorResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	isAllowed := CheckIfAllowedPermission(ctx, k.keeper, msg.Address, types.PermClaimCouncilor)
	if !isAllowed {
		return nil, errors.Wrap(types.ErrNotEnoughPermissions, "PermClaimCouncilor")
	}

	councilor := types.NewCouncilor(msg.Moniker, msg.Address)

	k.keeper.SaveCouncilor(ctx, councilor)
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeClaimCouncilor,
			sdk.NewAttribute(types.AttributeKeyAddress, msg.Address.String()),
		),
	)
	return &types.MsgClaimCouncilorResponse{}, nil
}
