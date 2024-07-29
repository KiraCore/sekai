package keeper

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"time"

	errorsmod "cosmossdk.io/errors"
	"github.com/KiraCore/sekai/x/gov/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	"golang.org/x/exp/utf8string"
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

	// check special proposal with dynamic voter proposal handler
	if content.ProposalPermission() == types.PermZero {
		router := k.keeper.GetProposalRouter()
		isAllowed := router.IsAllowedAddressDynamicProposal(ctx, msg.Proposer, content)
		if !isAllowed {
			return nil, errors.Wrap(types.ErrNotEnoughPermissions, "not enough permission to create the proposal")
		}
	} else {
		isAllowed := CheckIfAllowedPermission(ctx, k.keeper, msg.Proposer, content.ProposalPermission())
		if !isAllowed {
			return nil, errors.Wrap(types.ErrNotEnoughPermissions, content.ProposalPermission().String())
		}
	}

	proposalID, err := k.keeper.CreateAndSaveProposalWithContent(ctx, msg.Title, msg.Description, content)
	if err != nil {
		return nil, err
	}

	// call councilor rank update function
	k.keeper.OnCouncilorAct(ctx, msg.Proposer)

	cacheCtx, _ := ctx.CacheContext()
	router := k.keeper.GetProposalRouter()
	proposal, found := k.keeper.GetProposal(cacheCtx, proposalID)
	if !found {
		return nil, types.ErrProposalDoesNotExist
	}

	err = router.ApplyProposal(cacheCtx, proposalID, proposal.GetContent(), sdk.ZeroDec())
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

	// check special proposal with dynamic voter proposal handler
	content := proposal.GetContent()
	if content.VotePermission() == types.PermZero {
		router := k.keeper.GetProposalRouter()
		isAllowed := router.IsAllowedAddressDynamicProposal(ctx, msg.Voter, content)
		if !isAllowed {
			return nil, errors.Wrap(types.ErrNotEnoughPermissions, "not enough permission to vote on the proposal")
		}
	} else {
		isAllowed := CheckIfAllowedPermission(ctx, k.keeper, msg.Voter, content.VotePermission())
		if !isAllowed {
			return nil, errors.Wrap(types.ErrNotEnoughPermissions, content.VotePermission().String())
		}
	}

	// call councilor rank update function when it's the first vote
	if _, found := k.keeper.GetVote(ctx, msg.ProposalId, msg.Voter); !found {
		k.keeper.OnCouncilorAct(ctx, msg.Voter)
	}

	vote := types.NewVote(msg.ProposalId, msg.Voter, msg.Option, msg.Slash)
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

func (k msgServer) PollCreate(goCtx context.Context, msg *types.MsgPollCreate) (*types.MsgPollCreateResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	blockTime := ctx.BlockTime()
	properties := k.keeper.GetNetworkProperties(ctx)
	allowedTypes := []string{"string", "uint", "int", "float", "bool"}

	isAllowed := CheckIfAllowedPermission(ctx, k.keeper, msg.Creator, types.PermCreatePollProposal)
	if !isAllowed {
		return nil, errors.Wrap(types.ErrNotEnoughPermissions, types.PermCreatePollProposal.String())
	}

	if len(msg.Title) > int(properties.MaxProposalTitleSize) {
		return nil, types.ErrProposalTitleSizeExceeds
	}

	if len(msg.Description) > int(properties.MaxProposalDescriptionSize) {
		return nil, types.ErrProposalDescriptionSizeExceeds
	}

	if len(msg.Reference) > int(properties.MaxProposalReferenceSize) {
		return nil, types.ErrProposalTitleSizeExceeds
	}

	if len(msg.Checksum) > int(properties.MaxProposalChecksumSize) {
		return nil, types.ErrProposalTitleSizeExceeds
	}

	if len(msg.PollValues) > int(properties.MaxProposalPollOptionCount) {
		return nil, types.ErrProposalOptionCountExceeds
	}

	duration, err := time.ParseDuration(msg.Duration)
	if err != nil || blockTime.Add(duration).Before(blockTime) {
		return nil, fmt.Errorf("invalid duration: %w", err)
	}

	for _, v := range msg.PollValues {
		if len(v) > int(properties.MaxProposalPollOptionSize) {
			return nil, types.ErrProposalOptionSizeExceeds
		}

		if !utf8string.NewString(v).IsASCII() {
			return nil, types.ErrProposalOptionOnlyAscii
		}
	}

	for _, v := range msg.Roles {
		_, err := k.keeper.GetRoleBySid(ctx, v)
		if err != nil {
			return nil, errors.Wrap(types.ErrRoleDoesNotExist, v)
		}
	}

	sort.Strings(allowedTypes)
	i := sort.SearchStrings(allowedTypes, msg.ValueType)

	if i == len(allowedTypes) && allowedTypes[i] != msg.ValueType {
		return nil, types.ErrProposalTypeNotAllowed
	}

	pollID, err := k.keeper.PollCreate(ctx, msg)

	if err != nil {
		return nil, err
	}

	return &types.MsgPollCreateResponse{
		PollID: pollID,
	}, nil
}

func (k msgServer) PollVote(goCtx context.Context, msg *types.MsgPollVote) (*types.MsgPollVoteResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	actor, found := k.keeper.GetNetworkActorByAddress(ctx, msg.Voter)
	if !found || !actor.IsActive() {
		return nil, types.ErrActorIsNotActive
	}

	poll, pErr := k.keeper.GetPoll(ctx, msg.PollId)
	if pErr != nil {
		return nil, pErr
	}

	if poll.VotingEndTime.Before(time.Now()) {
		return nil, types.ErrVotingTimeEnded
	}

	roles := intersection(poll.Roles, actor.Roles)

	if len(roles) == 0 {
		return nil, types.ErrNotEnoughPermissions
	}

	switch poll.Options.Type {
	case "uint":
		_, err := strconv.ParseUint(msg.Value, 10, 64)
		if err != nil {
			return nil, errors.Wrap(types.ErrPollWrongValue, "Can not be converted to the unsigned integer")
		}
	case "int":
		_, err := strconv.ParseInt(msg.Value, 10, 64)
		if err != nil {
			return nil, errors.Wrap(types.ErrPollWrongValue, "Can not be converted to the integer")
		}
	case "bool":
		_, err := strconv.ParseBool(msg.Value)
		if err != nil {
			return nil, errors.Wrap(types.ErrPollWrongValue, "Can not be converted to the boolean")
		}
	case "float":
		_, err := strconv.ParseFloat(msg.Value, 64)
		if err != nil {
			return nil, errors.Wrap(types.ErrPollWrongValue, "Can not be converted to the float")
		}
	}

	if msg.Option == types.PollOptionCustom && poll.Options.Count <= uint64(len(poll.Options.Values)) && !inSlice(poll.Options.Values, msg.Value) {
		return nil, errors.Wrap(types.ErrPollWrongValue, "Maximum custom values exceeded")
	}

	if msg.Option == types.PollOptionCustom && poll.Options.Count > uint64(len(poll.Options.Values)) && !inSlice(poll.Options.Values, msg.Value) {
		poll.Options.Values = append(poll.Options.Values, msg.Value)
		k.keeper.SavePoll(ctx, poll)
	}

	err := k.keeper.PollVote(ctx, msg)
	return &types.MsgPollVoteResponse{}, err
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

func (k msgServer) UnassignRole(
	goCtx context.Context,
	msg *types.MsgUnassignRole,
) (*types.MsgUnassignRoleResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	isAllowed := CheckIfAllowedPermission(ctx, k.keeper, msg.Proposer, types.PermUpsertRole)
	if !isAllowed {
		return nil, errors.Wrap(types.ErrNotEnoughPermissions, types.PermUpsertRole.String())
	}

	err := k.keeper.UnassignRoleFromAccount(ctx, msg.Address, uint64(msg.RoleId))
	if err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeUnassignRole,
			sdk.NewAttribute(types.AttributeKeyProposer, msg.Proposer.String()),
			sdk.NewAttribute(types.AttributeKeyAddress, msg.Address.String()),
			sdk.NewAttribute(types.AttributeKeyRoleId, fmt.Sprintf("%d", msg.RoleId)),
		),
	)
	return &types.MsgUnassignRoleResponse{}, nil
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

	err := k.keeper.AssignRoleToAccount(ctx, msg.Address, uint64(msg.RoleId))
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeAssignRole,
			sdk.NewAttribute(types.AttributeKeyProposer, msg.Proposer.String()),
			sdk.NewAttribute(types.AttributeKeyAddress, msg.Address.String()),
			sdk.NewAttribute(types.AttributeKeyRoleId, fmt.Sprintf("%d", msg.RoleId)),
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

	// check sid is good variable naming form
	if !ValidateRoleSidKey(msg.RoleSid) {
		return nil, errors.Wrap(types.ErrInvalidRoleSid, fmt.Sprintf("invalid role sid configuration: sid=%s", msg.RoleSid))
	}

	_, err := k.keeper.GetRoleBySid(ctx, msg.RoleSid)
	if err == nil {
		return nil, types.ErrRoleExist
	}

	roleId := k.keeper.CreateRole(ctx, msg.RoleSid, msg.RoleDescription)
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCreateRole,
			sdk.NewAttribute(types.AttributeKeyProposer, msg.Proposer.String()),
			sdk.NewAttribute(types.AttributeKeyRoleId, fmt.Sprintf("%d", roleId)),
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

	roleId, err := k.keeper.GetRoleIdFromIdentifierString(ctx, msg.RoleIdentifier)
	if err != nil {
		return nil, err
	}

	err = k.keeper.RemoveBlacklistRolePermission(ctx, roleId, types.PermValue(msg.Permission))
	if err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeRemoveBlacklistRolePermisison,
			sdk.NewAttribute(types.AttributeKeyProposer, msg.Proposer.String()),
			sdk.NewAttribute(types.AttributeKeyRoleId, msg.RoleIdentifier),
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

	roleId, err := k.keeper.GetRoleIdFromIdentifierString(ctx, msg.RoleIdentifier)
	if err != nil {
		return nil, err
	}

	err = k.keeper.RemoveWhitelistRolePermission(ctx, roleId, types.PermValue(msg.Permission))
	if err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeRemoveWhitelistRolePermisison,
			sdk.NewAttribute(types.AttributeKeyProposer, msg.Proposer.String()),
			sdk.NewAttribute(types.AttributeKeyRoleId, msg.RoleIdentifier),
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

	roleId, err := k.keeper.GetRoleIdFromIdentifierString(ctx, msg.RoleIdentifier)
	if err != nil {
		return nil, err
	}

	err = k.keeper.BlacklistRolePermission(ctx, roleId, types.PermValue(msg.Permission))
	if err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeBlacklistRolePermisison,
			sdk.NewAttribute(types.AttributeKeyProposer, msg.Proposer.String()),
			sdk.NewAttribute(types.AttributeKeyRoleId, msg.RoleIdentifier),
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

	roleId, err := k.keeper.GetRoleIdFromIdentifierString(ctx, msg.RoleIdentifier)
	if err != nil {
		return nil, err
	}

	err = k.keeper.WhitelistRolePermission(ctx, roleId, types.PermValue(msg.Permission))
	if err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeWhitelistRolePermisison,
			sdk.NewAttribute(types.AttributeKeyProposer, msg.Proposer.String()),
			sdk.NewAttribute(types.AttributeKeyRoleId, msg.RoleIdentifier),
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

func (k msgServer) RemoveWhitelistedPermissions(
	goCtx context.Context,
	msg *types.MsgRemoveWhitelistedPermissions,
) (*types.MsgRemoveWhitelistedPermissionsResponse, error) {
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

	err := k.keeper.RemoveWhitelistedPermission(ctx, actor, types.PermValue(msg.Permission))
	if err != nil {
		return nil, errors.Wrapf(types.ErrSetPermissions, "error setting %d to whitelist: %s", msg.Permission, err)
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeRemoveWhitelistedPermisison,
			sdk.NewAttribute(types.AttributeKeyProposer, msg.Proposer.String()),
			sdk.NewAttribute(types.AttributeKeyRoleId, msg.Address.String()),
			sdk.NewAttribute(types.AttributeKeyPermission, fmt.Sprintf("%d", msg.Permission)),
		),
	)
	return &types.MsgRemoveWhitelistedPermissionsResponse{}, nil
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
		return nil, errors.Wrapf(types.ErrSetPermissions, "error setting %d to whitelist: %s", msg.Permission, err)
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

func (k msgServer) RemoveBlacklistedPermissions(
	goCtx context.Context,
	msg *types.MsgRemoveBlacklistedPermissions,
) (*types.MsgRemoveBlacklistedPermissionsResponse, error) {
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

	err := actor.Permissions.RemoveFromBlacklist(types.PermValue(msg.Permission))
	if err != nil {
		return nil, errors.Wrapf(types.ErrSetPermissions, "error setting %d to whitelist: %s", msg.Permission, err)
	}

	k.keeper.SaveNetworkActor(ctx, actor)
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeRemoveBlacklistedPermisison,
			sdk.NewAttribute(types.AttributeKeyProposer, msg.Proposer.String()),
			sdk.NewAttribute(types.AttributeKeyRoleId, msg.Address.String()),
			sdk.NewAttribute(types.AttributeKeyPermission, fmt.Sprintf("%d", msg.Permission)),
		),
	)
	return &types.MsgRemoveBlacklistedPermissionsResponse{}, nil
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

	k.keeper.SetExecutionFee(ctx, types.ExecutionFee{
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

	actor, found := k.keeper.GetNetworkActorByAddress(ctx, msg.Address)
	if !found {
		return nil, errors.Wrap(types.ErrNotEnoughPermissions, "network actor not found")
	}
	err := actor.Permissions.AddToWhitelist(types.PermCreatePollProposal)
	if err == nil {
		k.keeper.SaveNetworkActor(ctx, actor)
	}

	councilor := types.NewCouncilor(msg.Address, types.CouncilorActive)
	k.keeper.SaveCouncilor(ctx, councilor)

	identityInfo := []types.IdentityInfoEntry{}
	if msg.Moniker != "" {
		identityInfo = append(identityInfo, types.IdentityInfoEntry{
			Key:  "moniker",
			Info: msg.Moniker,
		})
	}

	if msg.Username != "" {
		identityInfo = append(identityInfo, types.IdentityInfoEntry{
			Key:  "username",
			Info: msg.Username,
		})
	}

	if msg.Description != "" {
		identityInfo = append(identityInfo, types.IdentityInfoEntry{
			Key:  "description",
			Info: msg.Description,
		})
	}

	if msg.Social != "" {
		identityInfo = append(identityInfo, types.IdentityInfoEntry{
			Key:  "social",
			Info: msg.Social,
		})
	}

	if msg.Contact != "" {
		identityInfo = append(identityInfo, types.IdentityInfoEntry{
			Key:  "contact",
			Info: msg.Contact,
		})
	}

	if msg.Avatar != "" {
		identityInfo = append(identityInfo, types.IdentityInfoEntry{
			Key:  "avatar",
			Info: msg.Avatar,
		})
	}

	err = k.keeper.RegisterIdentityRecords(ctx, msg.Address, identityInfo)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeClaimCouncilor,
			sdk.NewAttribute(types.AttributeKeyAddress, msg.Address.String()),
		),
	)
	return &types.MsgClaimCouncilorResponse{}, nil
}

// CouncilorPause - signal to the network that Councilor will NOT be present for a prolonged period of time
func (k msgServer) CouncilorPause(
	goCtx context.Context,
	msg *types.MsgCouncilorPause,
) (*types.MsgCouncilorPauseResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	councilor, found := k.keeper.GetCouncilor(ctx, sender)
	if !found {
		return nil, types.ErrCouncilorNotFound
	}

	// cannot be paused if not jailed already
	if councilor.Status == types.CouncilorJailed {
		return nil, errorsmod.Wrap(types.ErrCouncilorJailed, "Can NOT pause jailed councilor")
	}

	// cannot be paused if not inactive already
	if councilor.Status == types.CouncilorInactive {
		return nil, errorsmod.Wrap(types.ErrCouncilorInactivated, "Can NOT pause inactivated councilor")
	}

	// cannot be paused if not paused already
	if councilor.Status == types.CouncilorPaused {
		return nil, errorsmod.Wrap(types.ErrCouncilorPaused, "Can NOT pause already paused councilor")
	}

	councilor.Status = types.CouncilorPaused
	k.keeper.SaveCouncilor(ctx, councilor)
	return &types.MsgCouncilorPauseResponse{}, nil
}

// CouncilorUnpause - signal to the network that Councilor wishes to regain voting ability after planned absence
func (k msgServer) CouncilorUnpause(
	goCtx context.Context,
	msg *types.MsgCouncilorUnpause,
) (*types.MsgCouncilorUnpauseResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	councilor, found := k.keeper.GetCouncilor(ctx, sender)
	if !found {
		return nil, types.ErrCouncilorNotFound
	}

	// cannot be paused if not paused already
	if councilor.Status != types.CouncilorPaused {
		return nil, errorsmod.Wrap(types.ErrCouncilorNotPaused, "Can NOT unpause not paused councilor")
	}

	councilor.Status = types.CouncilorActive
	k.keeper.SaveCouncilor(ctx, councilor)

	return &types.MsgCouncilorUnpauseResponse{}, nil
}

// CouncilorActivate - signal to the network that Councilor wishes to regain voting ability after planned absence
func (k msgServer) CouncilorActivate(
	goCtx context.Context,
	msg *types.MsgCouncilorActivate,
) (*types.MsgCouncilorActivateResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	councilor, found := k.keeper.GetCouncilor(ctx, sender)
	if !found {
		return nil, types.ErrCouncilorNotFound
	}

	// cannot be activated if not inactive already
	if councilor.Status != types.CouncilorInactive {
		return nil, errorsmod.Wrap(types.ErrCouncilorNotInactivated, "Can NOT activate NOT inactive councilor")
	}

	councilor.Status = types.CouncilorActive
	councilor.AbstentionCounter = 0
	k.keeper.SaveCouncilor(ctx, councilor)

	return &types.MsgCouncilorActivateResponse{}, nil
}

func intersection(first, second []uint64) []uint64 {
	out := []uint64{}
	bucket := map[uint64]bool{}

	for _, i := range first {
		for _, j := range second {
			if i == j && !bucket[i] {
				out = append(out, i)
				bucket[i] = true
			}
		}
	}

	return out
}

func inSlice(sl []string, name string) bool {
	for _, value := range sl {
		if value == name {
			return true
		}
	}
	return false
}
