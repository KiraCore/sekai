package gov

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/KiraCore/sekai/x/gov/keeper"
	customgovtypes "github.com/KiraCore/sekai/x/gov/types"
)

func NewHandler(ck keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case *customgovtypes.MsgSetNetworkProperties:
			return handleSetNetworkProperties(ctx, ck, msg)
		case *customgovtypes.MsgSetExecutionFee:
			return handleSetExecutionFee(ctx, ck, msg)

		// Permission Related
		case *customgovtypes.MsgWhitelistPermissions:
			return handleWhitelistPermissions(ctx, ck, msg)
		case *customgovtypes.MsgBlacklistPermissions:
			return handleBlacklistPermissions(ctx, ck, msg)

		// Councilor Related
		case *customgovtypes.MsgClaimCouncilor:
			return handleClaimCouncilor(ctx, ck, msg)

		// Role Related
		case *customgovtypes.MsgWhitelistRolePermission:
			return handleWhitelistRolePermission(ctx, ck, msg)
		case *customgovtypes.MsgBlacklistRolePermission:
			return handleBlacklistRolePermission(ctx, ck, msg)
		case *customgovtypes.MsgRemoveWhitelistRolePermission:
			return handleRemoveWhitelistRolePermission(ctx, ck, msg)
		case *customgovtypes.MsgRemoveBlacklistRolePermission:
			return handleRemoveBlacklistRolePermission(ctx, ck, msg)
		case *customgovtypes.MsgCreateRole:
			return handleCreateRole(ctx, ck, msg)
		case *customgovtypes.MsgAssignRole:
			return handleAssignRole(ctx, ck, msg)
		case *customgovtypes.MsgRemoveRole:
			return handleMsgRemoveRole(ctx, ck, msg)

		// Proposal related
		case *customgovtypes.MsgProposalAssignPermission:
			return handleMsgProposalAssignPermission(ctx, ck, msg)
		case *customgovtypes.MsgVoteProposal:
			return handleMsgVoteProposal(ctx, ck, msg)

		default:
			return nil, errors.Wrapf(errors.ErrUnknownRequest, "unrecognized %s message type: %T", customgovtypes.ModuleName, msg)
		}
	}
}

func handleMsgVoteProposal(
	ctx sdk.Context,
	ck keeper.Keeper,
	msg *customgovtypes.MsgVoteProposal,
) (*sdk.Result, error) {
	isAllowed := keeper.CheckIfAllowedPermission(ctx, ck, msg.Voter, customgovtypes.PermVoteSetPermissionProposal)
	if !isAllowed {
		return nil, errors.Wrap(customgovtypes.ErrNotEnoughPermissions, "PermVoteSetPermissionProposal")
	}

	actor, found := ck.GetNetworkActorByAddress(ctx, msg.Voter)
	if !found || !actor.IsActive() {
		return nil, customgovtypes.ErrActorIsNotActive
	}

	_, found = ck.GetProposal(ctx, msg.ProposalId)
	if !found {
		return nil, customgovtypes.ErrProposalDoesNotExist
	}

	vote := customgovtypes.NewVote(msg.ProposalId, msg.Voter, msg.Option)
	ck.SaveVote(ctx, vote)

	return &sdk.Result{}, nil
}

func handleMsgProposalAssignPermission(
	ctx sdk.Context,
	ck keeper.Keeper,
	msg *customgovtypes.MsgProposalAssignPermission,
) (*sdk.Result, error) {
	isAllowed := keeper.CheckIfAllowedPermission(ctx, ck, msg.Proposer, customgovtypes.PermCreateSetPermissionsProposal)
	if !isAllowed {
		return nil, errors.Wrap(customgovtypes.ErrNotEnoughPermissions, "PermCreateSetPermissionsProposal")
	}

	actor, found := ck.GetNetworkActorByAddress(ctx, msg.Address)
	if found { // Actor exists
		if actor.Permissions.IsWhitelisted(customgovtypes.PermValue(msg.Permission)) {
			return nil, errors.Wrap(customgovtypes.ErrWhitelisting, "permission already whitelisted")
		}
	}

	blockTime := ctx.BlockTime()
	proposalID, err := ck.GetNextProposalID(ctx)
	if err != nil {
		return nil, err
	}

	properties := ck.GetNetworkProperties(ctx)

	proposal, err := customgovtypes.NewProposal(
		proposalID,
		customgovtypes.NewAssignPermissionProposal(
			msg.Address,
			customgovtypes.PermValue(msg.Permission),
		),
		blockTime,
		blockTime.Add(time.Minute*time.Duration(properties.ProposalEndTime)),
		blockTime.Add(time.Minute*time.Duration(properties.ProposalEnactmentTime)),
	)

	ck.SaveProposal(ctx, proposal)

	ck.AddToActiveProposals(ctx, proposal)

	return &sdk.Result{
		Data: keeper.ProposalIDToBytes(proposalID),
	}, nil
}

func handleMsgRemoveRole(ctx sdk.Context, ck keeper.Keeper, msg *customgovtypes.MsgRemoveRole) (*sdk.Result, error) {
	isAllowed := keeper.CheckIfAllowedPermission(ctx, ck, msg.Proposer, customgovtypes.PermSetPermissions)
	if !isAllowed {
		return nil, errors.Wrap(customgovtypes.ErrNotEnoughPermissions, "PermSetPermissions")
	}

	_, found := ck.GetPermissionsForRole(ctx, customgovtypes.Role(msg.Role))
	if !found {
		return nil, customgovtypes.ErrRoleDoesNotExist
	}

	actor, found := ck.GetNetworkActorByAddress(ctx, msg.Address)
	if !found {
		actor = customgovtypes.NewDefaultActor(msg.Address)
	}

	if !actor.HasRole(customgovtypes.Role(msg.Role)) {
		return nil, customgovtypes.ErrRoleNotAssigned
	}

	ck.RemoveRoleFromActor(ctx, actor, customgovtypes.Role(msg.Role))

	return &sdk.Result{}, nil
}

func handleAssignRole(ctx sdk.Context, ck keeper.Keeper, msg *customgovtypes.MsgAssignRole) (*sdk.Result, error) {
	isAllowed := keeper.CheckIfAllowedPermission(ctx, ck, msg.Proposer, customgovtypes.PermSetPermissions)
	if !isAllowed {
		return nil, errors.Wrap(customgovtypes.ErrNotEnoughPermissions, "PermSetPermissions")
	}

	_, found := ck.GetPermissionsForRole(ctx, customgovtypes.Role(msg.Role))
	if !found {
		return nil, customgovtypes.ErrRoleDoesNotExist
	}

	actor, found := ck.GetNetworkActorByAddress(ctx, msg.Address)
	if !found {
		actor = customgovtypes.NewDefaultActor(msg.Address)
	}

	if actor.HasRole(customgovtypes.Role(msg.Role)) {
		return nil, customgovtypes.ErrRoleAlreadyAssigned
	}

	ck.AssignRoleToActor(ctx, actor, customgovtypes.Role(msg.Role))

	return &sdk.Result{}, nil
}

func handleCreateRole(ctx sdk.Context, ck keeper.Keeper, msg *customgovtypes.MsgCreateRole) (*sdk.Result, error) {
	isAllowed := keeper.CheckIfAllowedPermission(ctx, ck, msg.Proposer, customgovtypes.PermUpsertRole)
	if !isAllowed {
		return nil, errors.Wrap(customgovtypes.ErrNotEnoughPermissions, "PermUpsertRole")
	}

	_, found := ck.GetPermissionsForRole(ctx, customgovtypes.Role(msg.Role))
	if found {
		return nil, customgovtypes.ErrRoleExist
	}

	ck.CreateRole(ctx, customgovtypes.Role(msg.Role))

	return &sdk.Result{}, nil
}

func handleRemoveBlacklistRolePermission(ctx sdk.Context, ck keeper.Keeper, msg *customgovtypes.MsgRemoveBlacklistRolePermission) (*sdk.Result, error) {
	_, err := validateAndGetPermissionsForRole(ctx, ck, msg.Proposer, customgovtypes.Role(msg.Role))
	if err != nil {
		return nil, err
	}

	err = ck.RemoveBlacklistRolePermission(ctx, customgovtypes.Role(msg.Role), customgovtypes.PermValue(msg.Permission))
	if err != nil {
		return nil, errors.Wrap(customgovtypes.ErrRemovingBlacklist, err.Error())
	}

	return &sdk.Result{}, nil
}

func handleRemoveWhitelistRolePermission(ctx sdk.Context, ck keeper.Keeper, msg *customgovtypes.MsgRemoveWhitelistRolePermission) (*sdk.Result, error) {
	_, err := validateAndGetPermissionsForRole(ctx, ck, msg.Proposer, customgovtypes.Role(msg.Role))
	if err != nil {
		return nil, err
	}

	err = ck.RemoveWhitelistRolePermission(ctx, customgovtypes.Role(msg.Role), customgovtypes.PermValue(msg.Permission))
	if err != nil {
		return nil, errors.Wrap(customgovtypes.ErrRemovingWhitelist, err.Error())
	}

	return &sdk.Result{}, nil
}

func handleBlacklistRolePermission(ctx sdk.Context, ck keeper.Keeper, msg *customgovtypes.MsgBlacklistRolePermission) (*sdk.Result, error) {
	_, err := validateAndGetPermissionsForRole(ctx, ck, msg.Proposer, customgovtypes.Role(msg.Role))
	if err != nil {
		return nil, err
	}

	err = ck.BlacklistRolePermission(ctx, customgovtypes.Role(msg.Role), customgovtypes.PermValue(msg.Permission))
	if err != nil {
		return nil, errors.Wrap(customgovtypes.ErrBlacklisting, err.Error())
	}

	return &sdk.Result{}, nil
}

func handleWhitelistRolePermission(ctx sdk.Context, ck keeper.Keeper, msg *customgovtypes.MsgWhitelistRolePermission) (*sdk.Result, error) {
	isAllowed := keeper.CheckIfAllowedPermission(ctx, ck, msg.Proposer, customgovtypes.PermSetPermissions)
	if !isAllowed {
		return nil, errors.Wrap(customgovtypes.ErrNotEnoughPermissions, "PermSetPermissions")
	}

	err := ck.WhitelistRolePermission(ctx, customgovtypes.Role(msg.Role), customgovtypes.PermValue(msg.Permission))
	if err != nil {
		return nil, err
	}

	return &sdk.Result{}, nil
}

func handleWhitelistPermissions(ctx sdk.Context, ck keeper.Keeper, msg *customgovtypes.MsgWhitelistPermissions) (*sdk.Result, error) {
	isAllowed := keeper.CheckIfAllowedPermission(ctx, ck, msg.Proposer, customgovtypes.PermSetPermissions)
	if !isAllowed {
		return nil, errors.Wrap(customgovtypes.ErrNotEnoughPermissions, "PermSetPermissions")
	}

	actor, found := ck.GetNetworkActorByAddress(ctx, msg.Address)
	if !found {
		actor = customgovtypes.NewDefaultActor(msg.Address)
	}

	err := ck.AddWhitelistPermission(ctx, actor, customgovtypes.PermValue(msg.Permission))
	if err != nil {
		return nil, errors.Wrapf(customgovtypes.ErrSetPermissions, "error setting %d to whitelist: %s", msg.Permission, err)
	}

	return &sdk.Result{}, nil
}

func handleBlacklistPermissions(ctx sdk.Context, ck keeper.Keeper, msg *customgovtypes.MsgBlacklistPermissions) (*sdk.Result, error) {
	isAllowed := keeper.CheckIfAllowedPermission(ctx, ck, msg.Proposer, customgovtypes.PermSetPermissions)
	if !isAllowed {
		return nil, errors.Wrap(customgovtypes.ErrNotEnoughPermissions, "PermSetPermissions")
	}

	actor, found := ck.GetNetworkActorByAddress(ctx, msg.Address)
	if !found {
		actor = customgovtypes.NewDefaultActor(msg.Address)
	}

	err := actor.Permissions.AddToBlacklist(customgovtypes.PermValue(msg.Permission))
	if err != nil {
		return nil, errors.Wrapf(customgovtypes.ErrSetPermissions, "error setting %d to whitelist", msg.Permission)
	}

	ck.SaveNetworkActor(ctx, actor)

	return &sdk.Result{}, nil
}

func handleSetNetworkProperties(ctx sdk.Context, ck keeper.Keeper, msg *customgovtypes.MsgSetNetworkProperties) (*sdk.Result, error) {
	isAllowed := keeper.CheckIfAllowedPermission(ctx, ck, msg.Proposer, customgovtypes.PermChangeTxFee)
	if !isAllowed {
		return nil, errors.Wrap(customgovtypes.ErrNotEnoughPermissions, "PermChangeTxFee")
	}
	ck.SetNetworkProperties(ctx, msg.NetworkProperties)
	return &sdk.Result{}, nil
}

func handleSetExecutionFee(ctx sdk.Context, ck keeper.Keeper, msg *customgovtypes.MsgSetExecutionFee) (*sdk.Result, error) {
	isAllowed := keeper.CheckIfAllowedPermission(ctx, ck, msg.Proposer, customgovtypes.PermChangeTxFee)
	if !isAllowed {
		return nil, errors.Wrap(customgovtypes.ErrNotEnoughPermissions, "PermChangeTxFee")
	}

	ck.SetExecutionFee(ctx, &customgovtypes.ExecutionFee{
		Name:              msg.Name,
		TransactionType:   msg.TransactionType,
		ExecutionFee:      msg.ExecutionFee,
		FailureFee:        msg.FailureFee,
		Timeout:           msg.Timeout,
		DefaultParameters: msg.DefaultParameters,
	})
	return &sdk.Result{}, nil
}

func handleClaimCouncilor(ctx sdk.Context, ck keeper.Keeper, msg *customgovtypes.MsgClaimCouncilor) (*sdk.Result, error) {
	isAllowed := keeper.CheckIfAllowedPermission(ctx, ck, msg.Address, customgovtypes.PermClaimCouncilor)
	if !isAllowed {
		return nil, errors.Wrap(customgovtypes.ErrNotEnoughPermissions, "PermClaimCouncilor")
	}

	councilor := customgovtypes.NewCouncilor(msg.Moniker, msg.Website, msg.Social, msg.Identity, msg.Address)

	ck.SaveCouncilor(ctx, councilor)

	return &sdk.Result{}, nil
}

// validateAndGetPermissionsForRole checks if:
// - Proposer has permissions to SetPermissions.
// - Role exists.
// And returns the permissions.
func validateAndGetPermissionsForRole(
	ctx sdk.Context,
	ck keeper.Keeper,
	proposer sdk.AccAddress,
	role customgovtypes.Role,
) (*customgovtypes.Permissions, error) {
	isAllowed := keeper.CheckIfAllowedPermission(ctx, ck, proposer, customgovtypes.PermSetPermissions)
	if !isAllowed {
		return nil, errors.Wrap(customgovtypes.ErrNotEnoughPermissions, "PermSetPermissions")
	}

	perms, found := ck.GetPermissionsForRole(ctx, role)
	if !found {
		return nil, customgovtypes.ErrRoleDoesNotExist
	}

	return &perms, nil
}
