package keeper

import (
	"context"
	"fmt"

	"github.com/KiraCore/sekai/x/gov/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
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

func (k msgServer) SubmitProposal(goCtx context.Context, msg *types.MsgSubmitProposal) (*types.MsgSubmitProposalResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	content := msg.GetContent()

	err := content.ValidateBasic()
	if err != nil {
		return nil, err
	}

	isAllowed := CheckIfAllowedPermission(ctx, k.keeper, msg.Proposer, content.ProposalPermission())
	if !isAllowed {
		return nil, errors.Wrap(types.ErrNotEnoughPermissions, content.ProposalPermission().String())
	}

	proposalID, err := k.keeper.CreateAndSaveProposalWithContent(ctx, msg.Title, msg.Description, content)
	if err != nil {
		return nil, err
	}

	cacheCtx, _ := ctx.CacheContext()
	router := k.keeper.GetProposalRouter()
	proposal, found := k.keeper.GetProposal(cacheCtx, proposalID)
	if !found {
		return nil, types.ErrProposalDoesNotExist
	}

	err = router.ApplyProposal(cacheCtx, proposal.GetContent())
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeSubmitProposal,
			sdk.NewAttribute(types.AttributeKeyProposalId, fmt.Sprintf("%d", proposalID)),
			sdk.NewAttribute(types.AttributeKeyProposalType, content.ProposalType()),
			sdk.NewAttribute(types.AttributeKeyProposalContent, msg.String()),
		),
	)

	return &types.MsgSubmitProposalResponse{
		ProposalID: proposalID,
	}, nil
}

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

// RegisterIdentityRecords defines a method to create identity record
func (k msgServer) RegisterIdentityRecords(goCtx context.Context, msg *types.MsgRegisterIdentityRecords) (*types.MsgRegisterIdentityRecordsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	err := k.keeper.RegisterIdentityRecords(ctx, msg.Address, msg.Infos)
	return &types.MsgRegisterIdentityRecordsResponse{}, err
}

func (k msgServer) DeleteIdentityRecords(goCtx context.Context, msg *types.MsgDeleteIdentityRecords) (*types.MsgDeleteIdentityRecordsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	err := k.keeper.DeleteIdentityRecords(ctx, msg.Address, msg.Keys)
	return &types.MsgDeleteIdentityRecordsResponse{}, err
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
func (k msgServer) HandleIdentityRecordsVerifyRequest(goCtx context.Context, msg *types.MsgHandleIdentityRecordsVerifyRequest) (*types.MsgHandleIdentityRecordsVerifyResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	err := k.keeper.HandleIdentityRecordsVerifyRequest(ctx, msg.Verifier, msg.VerifyRequestId, msg.Yes)
	return &types.MsgHandleIdentityRecordsVerifyResponse{}, err
}

func (k msgServer) CancelIdentityRecordsVerifyRequest(goCtx context.Context, msg *types.MsgCancelIdentityRecordsVerifyRequest) (*types.MsgCancelIdentityRecordsVerifyRequestResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	err := k.keeper.CancelIdentityRecordsVerifyRequest(ctx, msg.Executor, msg.VerifyRequestId)
	return &types.MsgCancelIdentityRecordsVerifyRequestResponse{}, err
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

	isSetClaimValidatorMsg := msg.Permission == uint32(types.PermClaimValidator)
	hasSetClaimValidatorPermission := CheckIfAllowedPermission(ctx, k.keeper, msg.Proposer, types.PermSetClaimValidatorPermission)
	hasSetPermissionsPermission := CheckIfAllowedPermission(ctx, k.keeper, msg.Proposer, types.PermSetPermissions)
	if !hasSetPermissionsPermission && !(isSetClaimValidatorMsg && hasSetClaimValidatorPermission) {
		return nil, errors.Wrap(types.ErrNotEnoughPermissions, "PermSetPermissions || (ClaimValidatorPermission && ClaimValidatorPermMsg)")
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

	isSetClaimValidatorMsg := msg.Permission == uint32(types.PermClaimValidator)
	hasSetClaimValidatorPermission := CheckIfAllowedPermission(ctx, k.keeper, msg.Proposer, types.PermSetClaimValidatorPermission)
	hasSetPermissionsPermission := CheckIfAllowedPermission(ctx, k.keeper, msg.Proposer, types.PermSetPermissions)
	if !hasSetPermissionsPermission && !(isSetClaimValidatorMsg && hasSetClaimValidatorPermission) {
		return nil, errors.Wrap(types.ErrNotEnoughPermissions, "PermSetPermissions || (ClaimValidatorPermission && ClaimValidatorPermMsg)")
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
	err := k.keeper.SetNetworkProperties(ctx, msg.NetworkProperties)
	if err != nil {
		return nil, err
	}
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
