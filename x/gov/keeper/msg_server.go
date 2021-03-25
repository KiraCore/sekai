package keeper

import (
	"context"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/KiraCore/sekai/x/gov/types"
	customgovtypes "github.com/KiraCore/sekai/x/gov/types"
)

type msgServer struct {
	keeper Keeper
}

// NewMsgServerImpl returns an implementation of the bank MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) customgovtypes.MsgServer {
	return &msgServer{
		keeper: keeper,
	}
}

var _ customgovtypes.MsgServer = msgServer{}

func (k msgServer) VoteProposal(
	goCtx context.Context,
	msg *customgovtypes.MsgVoteProposal,
) (*customgovtypes.MsgVoteProposalResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	actor, found := k.keeper.GetNetworkActorByAddress(ctx, msg.Voter)
	if !found || !actor.IsActive() {
		return nil, customgovtypes.ErrActorIsNotActive
	}

	proposal, found := k.keeper.GetProposal(ctx, msg.ProposalId)
	if !found {
		return nil, customgovtypes.ErrProposalDoesNotExist
	}

	if proposal.VotingEndTime.Before(ctx.BlockTime()) {
		return nil, customgovtypes.ErrVotingTimeEnded
	}

	isAllowed := CheckIfAllowedPermission(ctx, k.keeper, msg.Voter, proposal.GetContent().VotePermission())
	if !isAllowed {
		return nil, errors.Wrap(customgovtypes.ErrNotEnoughPermissions, proposal.GetContent().VotePermission().String())
	}

	vote := customgovtypes.NewVote(msg.ProposalId, msg.Voter, msg.Option)
	k.keeper.SaveVote(ctx, vote)

	return &customgovtypes.MsgVoteProposalResponse{}, nil
}

func (k msgServer) ProposalSetPoorNetworkMsgs(
	goCtx context.Context,
	msg *customgovtypes.MsgProposalSetPoorNetworkMessages,
) (*customgovtypes.MsgProposalSetPoorNetworkMessagesResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	isAllowed := CheckIfAllowedPermission(ctx, k.keeper, msg.Proposer, customgovtypes.PermCreateSetNetworkPropertyProposal)
	if !isAllowed {
		return nil, errors.Wrap(customgovtypes.ErrNotEnoughPermissions, customgovtypes.PermCreateSetNetworkPropertyProposal.String())
	}

	proposalID, err := k.CreateAndSaveProposalWithContent(ctx, msg.Description, customgovtypes.NewSetPoorNetworkMessagesProposal(msg.Messages))

	return &types.MsgProposalSetPoorNetworkMessagesResponse{
		ProposalID: proposalID,
	}, err
}

func (k msgServer) ProposalUpsertDataRegistry(
	goCtx context.Context,
	msg *customgovtypes.MsgProposalUpsertDataRegistry,
) (*customgovtypes.MsgProposalUpsertDataRegistryResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	isAllowed := CheckIfAllowedPermission(ctx, k.keeper, msg.Proposer, customgovtypes.PermUpsertDataRegistryProposal)
	if !isAllowed {
		return nil, errors.Wrap(customgovtypes.ErrNotEnoughPermissions, customgovtypes.PermUpsertDataRegistryProposal.String())
	}

	proposalID, err := k.CreateAndSaveProposalWithContent(ctx,
		msg.Description,
		customgovtypes.NewUpsertDataRegistryProposal(
			msg.Key,
			msg.Hash,
			msg.Reference,
			msg.Encoding,
			msg.Size_,
		),
	)

	return &types.MsgProposalUpsertDataRegistryResponse{
		ProposalID: proposalID,
	}, err
}

func (k msgServer) ProposalAssignPermission(
	goCtx context.Context,
	msg *customgovtypes.MsgProposalAssignPermission,
) (*customgovtypes.MsgProposalAssignPermissionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	isAllowed := CheckIfAllowedPermission(ctx, k.keeper, msg.Proposer, customgovtypes.PermCreateSetPermissionsProposal)
	if !isAllowed {
		return nil, errors.Wrap(customgovtypes.ErrNotEnoughPermissions, "PermCreateSetPermissionsProposal")
	}

	actor, found := k.keeper.GetNetworkActorByAddress(ctx, msg.Address)
	if found { // Actor exists
		if actor.Permissions.IsWhitelisted(customgovtypes.PermValue(msg.Permission)) {
			return nil, errors.Wrap(customgovtypes.ErrWhitelisting, "permission already whitelisted")
		}
	}

	proposalID, err := k.CreateAndSaveProposalWithContent(
		ctx,
		msg.Description,
		customgovtypes.NewAssignPermissionProposal(
			msg.Address,
			customgovtypes.PermValue(msg.Permission),
		))
	return &types.MsgProposalAssignPermissionResponse{
		ProposalID: proposalID,
	}, err
}

func (k msgServer) CreateAndSaveProposalWithContent(ctx sdk.Context, description string, content customgovtypes.Content) (uint64, error) {
	blockTime := ctx.BlockTime()
	proposalID, err := k.keeper.GetNextProposalID(ctx)
	if err != nil {
		return 0, err
	}

	properties := k.keeper.GetNetworkProperties(ctx)

	proposal, err := customgovtypes.NewProposal(
		proposalID,
		content,
		blockTime,
		blockTime.Add(time.Second*time.Duration(properties.ProposalEndTime)),
		blockTime.Add(time.Second*time.Duration(properties.ProposalEndTime)+
			time.Second*time.Duration(properties.ProposalEnactmentTime),
		),
		description,
	)

	k.keeper.SaveProposal(ctx, proposal)
	k.keeper.AddToActiveProposals(ctx, proposal)

	return proposalID, nil
}

func (k msgServer) ProposalSetNetworkProperty(
	goCtx context.Context,
	msg *customgovtypes.MsgProposalSetNetworkProperty,
) (*customgovtypes.MsgProposalSetNetworkPropertyResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	isAllowed := CheckIfAllowedPermission(ctx, k.keeper, msg.Proposer, customgovtypes.PermCreateSetNetworkPropertyProposal)
	if !isAllowed {
		return nil, errors.Wrap(customgovtypes.ErrNotEnoughPermissions, customgovtypes.PermCreateSetNetworkPropertyProposal.String())
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
		customgovtypes.NewSetNetworkPropertyProposal(
			msg.NetworkProperty,
			msg.Value,
		),
	)
	if err != nil {
		return nil, err
	}

	return &customgovtypes.MsgProposalSetNetworkPropertyResponse{
		ProposalID: proposalID,
	}, nil
}

func (k msgServer) RemoveRole(
	goCtx context.Context,
	msg *customgovtypes.MsgRemoveRole,
) (*customgovtypes.MsgRemoveRoleResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	isAllowed := CheckIfAllowedPermission(ctx, k.keeper, msg.Proposer, customgovtypes.PermUpsertRole)
	if !isAllowed {
		return nil, errors.Wrap(customgovtypes.ErrNotEnoughPermissions, customgovtypes.PermUpsertRole.String())
	}

	_, found := k.keeper.GetPermissionsForRole(ctx, customgovtypes.Role(msg.Role))
	if !found {
		return nil, customgovtypes.ErrRoleDoesNotExist
	}

	actor, found := k.keeper.GetNetworkActorByAddress(ctx, msg.Address)
	if !found {
		actor = customgovtypes.NewDefaultActor(msg.Address)
	}

	if !actor.HasRole(customgovtypes.Role(msg.Role)) {
		return nil, customgovtypes.ErrRoleNotAssigned
	}

	k.keeper.RemoveRoleFromActor(ctx, actor, customgovtypes.Role(msg.Role))

	return &customgovtypes.MsgRemoveRoleResponse{}, nil
}

func (k msgServer) AssignRole(
	goCtx context.Context,
	msg *customgovtypes.MsgAssignRole,
) (*customgovtypes.MsgAssignRoleResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	isAllowed := CheckIfAllowedPermission(ctx, k.keeper, msg.Proposer, customgovtypes.PermUpsertRole)
	if !isAllowed {
		return nil, errors.Wrap(customgovtypes.ErrNotEnoughPermissions, customgovtypes.PermUpsertRole.String())
	}

	_, found := k.keeper.GetPermissionsForRole(ctx, customgovtypes.Role(msg.Role))
	if !found {
		return nil, customgovtypes.ErrRoleDoesNotExist
	}

	actor, found := k.keeper.GetNetworkActorByAddress(ctx, msg.Address)
	if !found {
		actor = customgovtypes.NewDefaultActor(msg.Address)
	}

	if actor.HasRole(customgovtypes.Role(msg.Role)) {
		return nil, customgovtypes.ErrRoleAlreadyAssigned
	}

	k.keeper.AssignRoleToActor(ctx, actor, customgovtypes.Role(msg.Role))

	return &customgovtypes.MsgAssignRoleResponse{}, nil
}

func (k msgServer) CreateRole(
	goCtx context.Context,
	msg *customgovtypes.MsgCreateRole,
) (*customgovtypes.MsgCreateRoleResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	isAllowed := CheckIfAllowedPermission(ctx, k.keeper, msg.Proposer, customgovtypes.PermUpsertRole)
	if !isAllowed {
		return nil, errors.Wrap(customgovtypes.ErrNotEnoughPermissions, "PermUpsertRole")
	}

	_, found := k.keeper.GetPermissionsForRole(ctx, customgovtypes.Role(msg.Role))
	if found {
		return nil, customgovtypes.ErrRoleExist
	}

	k.keeper.CreateRole(ctx, customgovtypes.Role(msg.Role))

	return &customgovtypes.MsgCreateRoleResponse{}, nil
}

func (k msgServer) RemoveBlacklistRolePermission(
	goCtx context.Context,
	msg *customgovtypes.MsgRemoveBlacklistRolePermission,
) (*customgovtypes.MsgRemoveBlacklistRolePermissionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	isAllowed := CheckIfAllowedPermission(ctx, k.keeper, msg.Proposer, customgovtypes.PermUpsertRole)
	if !isAllowed {
		return nil, errors.Wrap(customgovtypes.ErrNotEnoughPermissions, customgovtypes.PermUpsertRole.String())
	}

	err := k.keeper.RemoveBlacklistRolePermission(ctx, customgovtypes.Role(msg.Role), customgovtypes.PermValue(msg.Permission))
	if err != nil {
		return nil, err
	}

	return &customgovtypes.MsgRemoveBlacklistRolePermissionResponse{}, nil
}

func (k msgServer) RemoveWhitelistRolePermission(
	goCtx context.Context,
	msg *customgovtypes.MsgRemoveWhitelistRolePermission,
) (*customgovtypes.MsgRemoveWhitelistRolePermissionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	isAllowed := CheckIfAllowedPermission(ctx, k.keeper, msg.Proposer, customgovtypes.PermUpsertRole)
	if !isAllowed {
		return nil, errors.Wrap(customgovtypes.ErrNotEnoughPermissions, customgovtypes.PermUpsertRole.String())
	}

	err := k.keeper.RemoveWhitelistRolePermission(ctx, customgovtypes.Role(msg.Role), customgovtypes.PermValue(msg.Permission))
	if err != nil {
		return nil, err
	}

	return &customgovtypes.MsgRemoveWhitelistRolePermissionResponse{}, nil
}

func (k msgServer) BlacklistRolePermission(
	goCtx context.Context,
	msg *customgovtypes.MsgBlacklistRolePermission,
) (*customgovtypes.MsgBlacklistRolePermissionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	isAllowed := CheckIfAllowedPermission(ctx, k.keeper, msg.Proposer, customgovtypes.PermUpsertRole)
	if !isAllowed {
		return nil, errors.Wrap(customgovtypes.ErrNotEnoughPermissions, customgovtypes.PermUpsertRole.String())
	}

	err := k.keeper.BlacklistRolePermission(ctx, customgovtypes.Role(msg.Role), customgovtypes.PermValue(msg.Permission))
	if err != nil {
		return nil, err
	}

	return &customgovtypes.MsgBlacklistRolePermissionResponse{}, nil
}

func (k msgServer) WhitelistRolePermission(
	goCtx context.Context,
	msg *customgovtypes.MsgWhitelistRolePermission,
) (*customgovtypes.MsgWhitelistRolePermissionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	isAllowed := CheckIfAllowedPermission(ctx, k.keeper, msg.Proposer, customgovtypes.PermUpsertRole)
	if !isAllowed {
		return nil, errors.Wrap(customgovtypes.ErrNotEnoughPermissions, customgovtypes.PermUpsertRole.String())
	}

	err := k.keeper.WhitelistRolePermission(ctx, customgovtypes.Role(msg.Role), customgovtypes.PermValue(msg.Permission))
	if err != nil {
		return nil, err
	}

	return &customgovtypes.MsgWhitelistRolePermissionResponse{}, nil
}

func (k msgServer) WhitelistPermissions(
	goCtx context.Context,
	msg *customgovtypes.MsgWhitelistPermissions,
) (*customgovtypes.MsgWhitelistPermissionsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	isAllowed := CheckIfAllowedPermission(ctx, k.keeper, msg.Proposer, customgovtypes.PermSetPermissions)
	if !isAllowed {
		return nil, errors.Wrap(customgovtypes.ErrNotEnoughPermissions, "PermSetPermissions")
	}

	actor, found := k.keeper.GetNetworkActorByAddress(ctx, msg.Address)
	if !found {
		actor = customgovtypes.NewDefaultActor(msg.Address)
	}

	err := k.keeper.AddWhitelistPermission(ctx, actor, customgovtypes.PermValue(msg.Permission))
	if err != nil {
		return nil, errors.Wrapf(customgovtypes.ErrSetPermissions, "error setting %d to whitelist: %s", msg.Permission, err)
	}

	return &customgovtypes.MsgWhitelistPermissionsResponse{}, nil
}

func (k msgServer) BlacklistPermissions(
	goCtx context.Context,
	msg *customgovtypes.MsgBlacklistPermissions,
) (*customgovtypes.MsgBlacklistPermissionsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	isAllowed := CheckIfAllowedPermission(ctx, k.keeper, msg.Proposer, customgovtypes.PermSetPermissions)
	if !isAllowed {
		return nil, errors.Wrap(customgovtypes.ErrNotEnoughPermissions, "PermSetPermissions")
	}

	actor, found := k.keeper.GetNetworkActorByAddress(ctx, msg.Address)
	if !found {
		actor = customgovtypes.NewDefaultActor(msg.Address)
	}

	err := actor.Permissions.AddToBlacklist(customgovtypes.PermValue(msg.Permission))
	if err != nil {
		return nil, errors.Wrapf(customgovtypes.ErrSetPermissions, "error setting %d to whitelist", msg.Permission)
	}

	k.keeper.SaveNetworkActor(ctx, actor)

	return &customgovtypes.MsgBlacklistPermissionsResponse{}, nil
}

func (k msgServer) SetNetworkProperties(
	goCtx context.Context,
	msg *customgovtypes.MsgSetNetworkProperties,
) (*customgovtypes.MsgSetNetworkPropertiesResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	isAllowed := CheckIfAllowedPermission(ctx, k.keeper, msg.Proposer, customgovtypes.PermChangeTxFee)
	if !isAllowed {
		return nil, errors.Wrap(customgovtypes.ErrNotEnoughPermissions, "PermChangeTxFee")
	}
	k.keeper.SetNetworkProperties(ctx, msg.NetworkProperties)
	return &customgovtypes.MsgSetNetworkPropertiesResponse{}, nil
}

func (k msgServer) SetExecutionFee(
	goCtx context.Context,
	msg *customgovtypes.MsgSetExecutionFee,
) (*customgovtypes.MsgSetExecutionFeeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	isAllowed := CheckIfAllowedPermission(ctx, k.keeper, msg.Proposer, customgovtypes.PermChangeTxFee)
	if !isAllowed {
		return nil, errors.Wrap(customgovtypes.ErrNotEnoughPermissions, "PermChangeTxFee")
	}

	k.keeper.SetExecutionFee(ctx, &customgovtypes.ExecutionFee{
		Name:              msg.Name,
		TransactionType:   msg.TransactionType,
		ExecutionFee:      msg.ExecutionFee,
		FailureFee:        msg.FailureFee,
		Timeout:           msg.Timeout,
		DefaultParameters: msg.DefaultParameters,
	})
	return &customgovtypes.MsgSetExecutionFeeResponse{}, nil
}

func (k msgServer) ClaimCouncilor(
	goCtx context.Context,
	msg *customgovtypes.MsgClaimCouncilor,
) (*customgovtypes.MsgClaimCouncilorResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	isAllowed := CheckIfAllowedPermission(ctx, k.keeper, msg.Address, customgovtypes.PermClaimCouncilor)
	if !isAllowed {
		return nil, errors.Wrap(customgovtypes.ErrNotEnoughPermissions, "PermClaimCouncilor")
	}

	councilor := customgovtypes.NewCouncilor(msg.Moniker, msg.Website, msg.Social, msg.Identity, msg.Address)

	k.keeper.SaveCouncilor(ctx, councilor)

	return &customgovtypes.MsgClaimCouncilorResponse{}, nil
}

func (k msgServer) ProposalCreateRole(goCtx context.Context, msg *customgovtypes.MsgProposalCreateRole) (*customgovtypes.MsgProposalCreateRoleResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if len(msg.WhitelistedPermissions) == 0 && len(msg.BlacklistedPermissions) == 0 {
		return nil, customgovtypes.ErrEmptyPermissions
	}

	isAllowed := CheckIfAllowedPermission(ctx, k.keeper, msg.Proposer, customgovtypes.PermCreateRoleProposal)
	if !isAllowed {
		return nil, errors.Wrap(customgovtypes.ErrNotEnoughPermissions, customgovtypes.PermCreateRoleProposal.String())
	}

	_, exists := k.keeper.GetPermissionsForRole(ctx, customgovtypes.Role(msg.Role))
	if exists {
		return nil, customgovtypes.ErrRoleExist
	}

	proposalID, err := k.CreateAndSaveProposalWithContent(
		ctx,
		msg.Description,
		customgovtypes.NewCreateRoleProposal(
			customgovtypes.Role(msg.Role),
			msg.WhitelistedPermissions,
			msg.BlacklistedPermissions,
		),
	)
	if err != nil {
		return nil, err
	}

	return &customgovtypes.MsgProposalCreateRoleResponse{
		ProposalID: proposalID,
	}, nil
}
