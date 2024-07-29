package keeper

import (
	"context"
	"fmt"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	errorsmod "cosmossdk.io/errors"
	appparams "github.com/KiraCore/sekai/app/params"
	kiratypes "github.com/KiraCore/sekai/types"
	"github.com/KiraCore/sekai/x/gov/types"
	stakingtypes "github.com/KiraCore/sekai/x/staking/types"
	"github.com/cosmos/cosmos-sdk/types/query"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) PollsVotesByPollId(goCtx context.Context, request *types.QueryPollsVotesByPollId) (*types.QueryPollsVotesByPollIdResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	votes := k.GetPollVotes(ctx, request.PollId)

	return &types.QueryPollsVotesByPollIdResponse{
		Votes: votes,
	}, nil
}

// AllRoles return all roles registered
func (k Keeper) AllRoles(goCtx context.Context, request *types.AllRolesRequest) (*types.AllRolesResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	roles := k.GetAllRoles(ctx)
	roleQueries := []types.RoleQuery{}
	for _, role := range roles {
		permissions, _ := k.GetPermissionsForRole(ctx, uint64(role.Id))
		roleQueries = append(roleQueries, types.RoleQuery{
			Id:          role.Id,
			Sid:         role.Sid,
			Description: role.Description,
			Permissions: &permissions,
		})
	}

	return &types.AllRolesResponse{
		Roles: roleQueries,
	}, nil
}

// PollsListByAddress return polls associated to an address
func (k Keeper) PollsListByAddress(goCtx context.Context, request *types.QueryPollsListByAddress) (*types.QueryPollsListByAddressResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	polls, err := k.GetPollsByAddress(ctx, request.Creator)

	if err != nil {
		return nil, err
	}

	return &types.QueryPollsListByAddressResponse{
		Polls: polls,
	}, nil
}

// RolesByAddress return roles associated to an address
func (k Keeper) RolesByAddress(goCtx context.Context, request *types.RolesByAddressRequest) (*types.RolesByAddressResponse, error) {
	addr, err := sdk.AccAddressFromBech32(request.Addr)
	if err != nil {
		return nil, err
	}
	actor, found := k.GetNetworkActorByAddress(sdk.UnwrapSDKContext(goCtx), addr)
	if !found {
		return nil, stakingtypes.ErrNetworkActorNotFound
	}

	return &types.RolesByAddressResponse{
		RoleIds: actor.Roles,
	}, nil
}

// CouncilorByAddress return councilor object associated to an address
func (k Keeper) CouncilorByAddress(goCtx context.Context, request *types.CouncilorByAddressRequest) (*types.CouncilorResponse, error) {
	addr, err := sdk.AccAddressFromBech32(request.Addr)
	if err != nil {
		return nil, err
	}
	councilor, found := k.GetCouncilor(sdk.UnwrapSDKContext(goCtx), addr)
	if !found {
		return nil, types.ErrCouncilorNotFound
	}

	return &types.CouncilorResponse{Councilor: councilor}, nil
}

// CouncilorByMoniker return councilor object named moniker
func (k Keeper) CouncilorByMoniker(goCtx context.Context, request *types.CouncilorByMonikerRequest) (*types.CouncilorResponse, error) {
	councilor, found := k.GetCouncilorByMoniker(sdk.UnwrapSDKContext(goCtx), request.Moniker)
	if !found {
		return nil, types.ErrCouncilorNotFound
	}

	return &types.CouncilorResponse{Councilor: councilor}, nil
}

// PermissionsByAddress returns permissions associated to an address
func (k Keeper) PermissionsByAddress(goCtx context.Context, request *types.PermissionsByAddressRequest) (*types.PermissionsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	addr, err := sdk.AccAddressFromBech32(request.Addr)
	if err != nil {
		return nil, err
	}
	networkActor, found := k.GetNetworkActorByAddress(ctx, addr)
	if !found {
		return nil, stakingtypes.ErrNetworkActorNotFound
	}

	return &types.PermissionsResponse{Permissions: networkActor.Permissions}, nil
}

// NetworkProperties return global network properties
func (k Keeper) NetworkProperties(goCtx context.Context, request *types.NetworkPropertiesRequest) (*types.NetworkPropertiesResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	networkProperties := k.GetNetworkProperties(ctx)
	return &types.NetworkPropertiesResponse{Properties: networkProperties}, nil
}

// CustomPrefixes return default denom and bech32 prefix
func (k Keeper) CustomPrefixes(goCtx context.Context, request *types.QueryCustomPrefixesRequest) (*types.QueryCustomPrefixesResponse, error) {
	return &types.QueryCustomPrefixesResponse{
		DefaultDenom: appparams.DefaultDenom,
		Bech32Prefix: appparams.AccountAddressPrefix,
	}, nil
}

func (k Keeper) Role(goCtx context.Context, request *types.RoleRequest) (*types.RoleResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	roleId, err := k.GetRoleIdFromIdentifierString(ctx, request.Identifier)
	if err != nil {
		return nil, err
	}
	role, err := k.GetRole(ctx, uint64(roleId))
	if err != nil {
		return nil, err
	}

	permissions, found := k.GetPermissionsForRole(ctx, uint64(role.Id))
	if !found {
		return nil, types.ErrRoleDoesNotExist
	}

	return &types.RoleResponse{Role: &types.RoleQuery{
		Id:          role.Id,
		Sid:         role.Sid,
		Description: role.Description,
		Permissions: &permissions,
	}}, nil
}

// ExecutionFee returns execution fee associated to a specific message type
func (k Keeper) ExecutionFee(goCtx context.Context, request *types.ExecutionFeeRequest) (*types.ExecutionFeeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	fee := k.GetExecutionFee(ctx, request.TransactionType)
	if fee == nil {
		return nil, errorsmod.Wrap(types.ErrFeeNotExist, fmt.Sprintf("fee does not exist for %s", request.TransactionType))
	}
	return &types.ExecutionFeeResponse{Fee: fee}, nil
}

// ExecutionFee returns execution fee associated to a specific message type
func (k Keeper) AllExecutionFees(goCtx context.Context, request *types.AllExecutionFeesRequest) (*types.AllExecutionFeesResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	fees := k.GetExecutionFees(ctx)

	txTypes := []string{}
	for txType := range kiratypes.MsgFuncIDMapping {
		txTypes = append(txTypes, txType)
	}

	return &types.AllExecutionFeesResponse{
		Fees:    fees,
		TxTypes: txTypes,
	}, nil
}

// PoorNetworkMessages queries poor network messages
func (k Keeper) PoorNetworkMessages(goCtx context.Context, request *types.PoorNetworkMessagesRequest) (*types.PoorNetworkMessagesResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	msg := k.GetPoorNetworkMessages(ctx)
	return &types.PoorNetworkMessagesResponse{Messages: msg.Messages}, nil
}

// Proposal returns a proposal by id
func (k Keeper) Proposal(goCtx context.Context, request *types.QueryProposalRequest) (*types.QueryProposalResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	proposal, found := k.GetProposal(ctx, request.ProposalId)
	if found == false {
		return nil, errorsmod.Wrap(types.ErrGettingProposals, fmt.Sprintf("proposal does not exist for %d", request.ProposalId))
	}
	votes := k.GetProposalVotes(ctx, request.ProposalId)
	return &types.QueryProposalResponse{
		Proposal: proposal,
		Votes:    votes,
	}, nil
}

// Proposals query proposals by querying params with pagination
func (k Keeper) Proposals(goCtx context.Context, request *types.QueryProposalsRequest) (*types.QueryProposalsResponse, error) {
	c := sdk.UnwrapSDKContext(goCtx)
	if request == nil {
		err := status.Error(codes.InvalidArgument, "empty request")
		return nil, errorsmod.Wrap(types.ErrGettingProposals, fmt.Sprintf("error getting proposals: %s", err.Error()))
	}

	store := c.KVStore(k.storeKey)

	var proposals []types.Proposal
	var pageRes *query.PageResponse
	var err error

	proposalsStore := prefix.NewStore(store, ProposalsPrefix)

	onResult := func(key []byte, value []byte, accumulate bool) (bool, error) {
		var proposal types.Proposal
		err := k.cdc.Unmarshal(value, &proposal)
		if err != nil {
			return false, err
		}
		if request.Voter != "" {
			voter, err := sdk.AccAddressFromBech32(request.Voter)
			if err != nil {
				return false, err
			}

			_, found := k.GetVote(c, proposal.ProposalId, voter)
			if !found {
				return false, nil
			}
		}
		if accumulate {
			proposals = append(proposals, proposal)
		}
		return true, nil
	}

	// we set maximum limit for safety of iteration
	if request.Pagination != nil && request.Pagination.Limit > kiratypes.PageIterationLimit {
		request.Pagination.Limit = kiratypes.PageIterationLimit
	}

	pageRes, err = query.FilteredPaginate(proposalsStore, request.Pagination, onResult)

	if err != nil {
		return nil, errorsmod.Wrap(types.ErrGettingProposals, fmt.Sprintf("error getting proposals: %s", err.Error()))
	}

	res := types.QueryProposalsResponse{
		Proposals:  proposals,
		Pagination: pageRes,
	}

	return &res, nil
}

// WhitelistedProposalVoters returns whitelisted voters for a proposal for tracking
func (k Keeper) WhitelistedProposalVoters(goCtx context.Context, request *types.QueryWhitelistedProposalVotersRequest) (*types.QueryWhitelistedProposalVotersResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	proposal, found := k.GetProposal(ctx, request.ProposalId)
	if !found {
		return nil, errorsmod.Wrap(types.ErrGettingProposals, fmt.Sprintf("proposal does not exist for %d", request.ProposalId))
	}

	// dynamic proposal users for spending pool proposals
	content := proposal.GetContent()
	if content.VotePermission() == types.PermZero {
		router := k.GetProposalRouter()
		addrs := router.AllowedAddressesDynamicProposal(ctx, content)
		actors := []types.NetworkActor{}
		for _, addr := range addrs {
			sdkAddr, err := sdk.AccAddressFromBech32(addr)
			if err != nil {
				return nil, err
			}
			actor, ok := k.GetNetworkActorByAddress(ctx, sdkAddr)
			if !ok {
				actor = types.NetworkActor{
					Address: sdkAddr,
				}
			}
			actors = append(actors, actor)
		}
		return &types.QueryWhitelistedProposalVotersResponse{Voters: actors}, nil
	} else {
		actors := k.GetNetworkActorsByAbsoluteWhitelistPermission(ctx, proposal.GetContent().VotePermission())
		return &types.QueryWhitelistedProposalVotersResponse{Voters: actors}, nil
	}
}

func (k Keeper) actorsCountByPermissions(ctx sdk.Context, perms []types.PermValue) uint64 {
	count := uint64(0)
	counted := make(map[string]bool)
	for _, perm := range perms {
		actors := k.GetNetworkActorsByAbsoluteWhitelistPermission(ctx, perm)
		for _, actor := range actors {
			if !counted[actor.Address.String()] {
				counted[actor.Address.String()] = true
				count++
			}
		}
	}
	return count
}

// ProposerVotersCount returns proposer and voters count that can create at least a type of proposal
func (k Keeper) ProposerVotersCount(goCtx context.Context, request *types.QueryProposerVotersCountRequest) (*types.QueryProposerVotersCountResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: should consider dynamic proposals later where permission is set to PermZero

	proposalPerms := []types.PermValue{
		types.PermWhitelistAccountPermissionProposal,
		types.PermCreateUpsertDataRegistryProposal,
		types.PermCreateSetNetworkPropertyProposal,
		types.PermCreateSetPoorNetworkMessagesProposal,
		types.PermCreateUpsertTokenInfoProposal,
		types.PermCreateUnjailValidatorProposal,
		types.PermCreateRoleProposal,
		types.PermCreateTokensWhiteBlackChangeProposal,
		types.PermCreateResetWholeValidatorRankProposal,
		types.PermCreateSoftwareUpgradeProposal,
		types.PermCreateSetProposalDurationProposal,
		types.PermBlacklistAccountPermissionProposal,
		types.PermRemoveWhitelistedAccountPermissionProposal,
		types.PermRemoveBlacklistedAccountPermissionProposal,
		types.PermWhitelistRolePermissionProposal,
		types.PermBlacklistRolePermissionProposal,
		types.PermRemoveWhitelistedRolePermissionProposal,
		types.PermRemoveBlacklistedRolePermissionProposal,
		types.PermAssignRoleToAccountProposal,
		types.PermUnassignRoleFromAccountProposal,
		types.PermRemoveRoleProposal,
		types.PermCreateUpsertUBIProposal,
		types.PermCreateRemoveUBIProposal,
	}
	votePerms := []types.PermValue{
		types.PermVoteWhitelistAccountPermissionProposal,
		types.PermVoteUpsertDataRegistryProposal,
		types.PermVoteSetNetworkPropertyProposal,
		types.PermVoteSetPoorNetworkMessagesProposal,
		types.PermVoteUpsertTokenInfoProposal,
		types.PermVoteUnjailValidatorProposal,
		types.PermVoteCreateRoleProposal,
		types.PermVoteTokensWhiteBlackChangeProposal,
		types.PermVoteResetWholeValidatorRankProposal,
		types.PermVoteSoftwareUpgradeProposal,
		types.PermVoteSetProposalDurationProposal,
		types.PermVoteBlacklistAccountPermissionProposal,
		types.PermVoteRemoveWhitelistedAccountPermissionProposal,
		types.PermVoteRemoveBlacklistedAccountPermissionProposal,
		types.PermVoteWhitelistRolePermissionProposal,
		types.PermVoteBlacklistRolePermissionProposal,
		types.PermVoteRemoveWhitelistedRolePermissionProposal,
		types.PermVoteRemoveBlacklistedRolePermissionProposal,
		types.PermVoteAssignRoleToAccountProposal,
		types.PermVoteUnassignRoleFromAccountProposal,
		types.PermVoteRemoveRoleProposal,
		types.PermVoteUpsertUBIProposal,
		types.PermVoteRemoveUBIProposal,
	}

	return &types.QueryProposerVotersCountResponse{
		Proposers: k.actorsCountByPermissions(ctx, proposalPerms),
		Voters:    k.actorsCountByPermissions(ctx, votePerms),
	}, nil
}

// Vote queries voted information based on proposalID, voterAddr.
func (k Keeper) Vote(goCtx context.Context, request *types.QueryVoteRequest) (*types.QueryVoteResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	voter, err := sdk.AccAddressFromBech32(request.Voter)
	if err != nil {
		return nil, err
	}
	vote, found := k.GetVote(ctx, request.ProposalId, voter)
	if !found {
		return &types.QueryVoteResponse{Vote: vote}, errorsmod.Wrap(types.ErrGettingProposalVotes, fmt.Sprintf("error getting votes for proposal %d, voter %s", request.ProposalId, request.Voter))
	}
	return &types.QueryVoteResponse{Vote: vote}, nil
}

// Votes queries votes of a given proposal.
func (k Keeper) Votes(goCtx context.Context, request *types.QueryVotesRequest) (*types.QueryVotesResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	votes := k.GetProposalVotes(ctx, request.ProposalId)
	return &types.QueryVotesResponse{Votes: votes}, nil
}

// AllDataReferenceKeys queries all data reference keys with pagination
func (k Keeper) AllDataReferenceKeys(goCtx context.Context, request *types.QueryDataReferenceKeysRequest) (*types.QueryDataReferenceKeysResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	var keys []string
	store := ctx.KVStore(k.storeKey)
	dataReferenceStore := prefix.NewStore(store, DataRegistryPrefix)

	pageRes, err := query.Paginate(dataReferenceStore, request.Pagination, func(key []byte, value []byte) error {
		keys = append(keys, string(key))
		return nil
	})

	if err != nil {
		return &types.QueryDataReferenceKeysResponse{}, err
	}

	res := types.QueryDataReferenceKeysResponse{
		Keys:       keys,
		Pagination: pageRes,
	}

	return &res, nil
}

// DataReferenceByKey queries data reference by key
func (k Keeper) DataReferenceByKey(goCtx context.Context, request *types.QueryDataReferenceRequest) (*types.QueryDataReferenceResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	dataReference, ok := k.GetDataRegistryEntry(ctx, request.GetKey())

	if !ok {
		return nil, status.Error(codes.NotFound, "not found")
	}

	res := types.QueryDataReferenceResponse{
		Data: &dataReference,
	}

	return &res, nil
}

// IdentityRecord query identity record by id
func (k Keeper) IdentityRecord(goCtx context.Context, request *types.QueryIdentityRecordRequest) (*types.QueryIdentityRecordResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	res := types.QueryIdentityRecordResponse{
		Record: k.GetIdentityRecordById(ctx, request.Id),
	}

	return &res, nil
}

// IdentityRecordsByAddress query identity records by creator and keys
func (k Keeper) IdentityRecordsByAddress(goCtx context.Context, request *types.QueryIdentityRecordsByAddressRequest) (*types.QueryIdentityRecordsByAddressResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	records, err := k.GetIdRecordsByAddressAndKeys(ctx, request.Creator, request.Keys)

	return &types.QueryIdentityRecordsByAddressResponse{
		Records: records,
	}, err
}

// AllIdentityRecords query all identity records
func (k Keeper) AllIdentityRecords(goCtx context.Context, request *types.QueryAllIdentityRecordsRequest) (*types.QueryAllIdentityRecordsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	store := ctx.KVStore(k.storeKey)
	recordStore := prefix.NewStore(store, types.KeyPrefixIdentityRecord)

	records := []types.IdentityRecord{}
	onResult := func(key []byte, value []byte, accumulate bool) (bool, error) {
		var record types.IdentityRecord
		err := k.cdc.Unmarshal(value, &record)
		if err != nil {
			return false, err
		}
		if accumulate {
			records = append(records, record)
		}
		return true, nil
	}

	// we set maximum limit for safety of iteration
	if request.Pagination != nil && request.Pagination.Limit > kiratypes.PageIterationLimit {
		request.Pagination.Limit = kiratypes.PageIterationLimit
	}

	pageRes, err := query.FilteredPaginate(recordStore, request.Pagination, onResult)
	if err != nil {
		return nil, err
	}

	res := types.QueryAllIdentityRecordsResponse{
		Records:    records,
		Pagination: pageRes,
	}

	return &res, nil
}

// IdentityRecordVerifyRequest query identity record verify request by id
func (k Keeper) IdentityRecordVerifyRequest(goCtx context.Context, request *types.QueryIdentityVerifyRecordRequest) (*types.QueryIdentityVerifyRecordResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	res := types.QueryIdentityVerifyRecordResponse{
		VerifyRecord: k.GetIdRecordsVerifyRequest(ctx, request.RequestId),
	}

	return &res, nil
}

// IdentityRecordVerifyRequestsByRequester query identity record verify requests by requester
func (k Keeper) IdentityRecordVerifyRequestsByRequester(goCtx context.Context, request *types.QueryIdentityRecordVerifyRequestsByRequester) (*types.QueryIdentityRecordVerifyRequestsByRequesterResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	requests := []types.IdentityRecordsVerify{}
	store := ctx.KVStore(k.storeKey)
	requestByRequesterStore := prefix.NewStore(store, types.IdRecordVerifyRequestByRequesterPrefix(request.Requester.String()))
	onResult := func(key []byte, value []byte, accumulate bool) (bool, error) {
		requestId := sdk.BigEndianToUint64(value)
		request := k.GetIdRecordsVerifyRequest(ctx, requestId)
		if request == nil {
			return false, fmt.Errorf("invalid id available on requests: %d", requestId)
		}
		if accumulate {
			requests = append(requests, *request)
		}
		return true, nil
	}

	// we set maximum limit for safety of iteration
	if request.Pagination != nil && request.Pagination.Limit > kiratypes.PageIterationLimit {
		request.Pagination.Limit = kiratypes.PageIterationLimit
	}

	pageRes, err := query.FilteredPaginate(requestByRequesterStore, request.Pagination, onResult)
	if err != nil {
		return nil, err
	}

	res := types.QueryIdentityRecordVerifyRequestsByRequesterResponse{
		VerifyRecords: requests,
		Pagination:    pageRes,
	}

	return &res, nil
}

// IdentityRecordVerifyRequestsByApprover query identity records verify requests by approver
func (k Keeper) IdentityRecordVerifyRequestsByApprover(goCtx context.Context, request *types.QueryIdentityRecordVerifyRequestsByApprover) (*types.QueryIdentityRecordVerifyRequestsByApproverResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	requests := []types.IdentityRecordsVerify{}
	store := ctx.KVStore(k.storeKey)
	requestByApproverStore := prefix.NewStore(store, types.IdRecordVerifyRequestByApproverPrefix(request.Approver.String()))
	onResult := func(key []byte, value []byte, accumulate bool) (bool, error) {
		requestId := sdk.BigEndianToUint64(value)
		request := k.GetIdRecordsVerifyRequest(ctx, requestId)
		if request == nil {
			return false, fmt.Errorf("invalid id available on requests: %d", requestId)
		}
		if accumulate {
			requests = append(requests, *request)
		}
		return true, nil
	}

	// we set maximum limit for safety of iteration
	if request.Pagination != nil && request.Pagination.Limit > kiratypes.PageIterationLimit {
		request.Pagination.Limit = kiratypes.PageIterationLimit
	}

	pageRes, err := query.FilteredPaginate(requestByApproverStore, request.Pagination, onResult)
	if err != nil {
		return nil, err
	}

	res := types.QueryIdentityRecordVerifyRequestsByApproverResponse{
		VerifyRecords: requests,
		Pagination:    pageRes,
	}

	return &res, nil
}

// AllIdentityRecordVerifyRequests query all identity records verify requests
func (k Keeper) AllIdentityRecordVerifyRequests(goCtx context.Context, request *types.QueryAllIdentityRecordVerifyRequests) (*types.QueryAllIdentityRecordVerifyRequestsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	requests := []types.IdentityRecordsVerify{}
	store := ctx.KVStore(k.storeKey)
	requestStore := prefix.NewStore(store, types.KeyPrefixIdRecordVerifyRequest)
	onResult := func(key []byte, value []byte, accumulate bool) (bool, error) {
		request := types.IdentityRecordsVerify{}
		k.cdc.MustUnmarshal(value, &request)
		if accumulate {
			requests = append(requests, request)
		}
		return true, nil
	}

	// we set maximum limit for safety of iteration
	if request.Pagination != nil && request.Pagination.Limit > kiratypes.PageIterationLimit {
		request.Pagination.Limit = kiratypes.PageIterationLimit
	}

	pageRes, err := query.FilteredPaginate(requestStore, request.Pagination, onResult)
	if err != nil {
		return nil, err
	}

	res := types.QueryAllIdentityRecordVerifyRequestsResponse{
		VerifyRecords: requests,
		Pagination:    pageRes,
	}

	return &res, nil
}

// GetAllDataReferenceKeys implements the Query all data reference keys gRPC method
func (k Keeper) GetAllDataReferenceKeys(sdkCtx sdk.Context, req *types.QueryDataReferenceKeysRequest) (*types.QueryDataReferenceKeysResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	var keys []string
	store := sdkCtx.KVStore(k.storeKey)
	dataReferenceStore := prefix.NewStore(store, DataRegistryPrefix)

	pageRes, err := query.Paginate(dataReferenceStore, req.Pagination, func(key []byte, value []byte) error {
		keys = append(keys, string(key))
		return nil
	})

	if err != nil {
		return &types.QueryDataReferenceKeysResponse{}, err
	}

	res := types.QueryDataReferenceKeysResponse{
		Keys:       keys,
		Pagination: pageRes,
	}

	return &res, nil
}

// GetDataReferenceByKey implements the Query data reference by key gRPC method
func (k Keeper) GetDataReferenceByKey(sdkCtx sdk.Context, req *types.QueryDataReferenceRequest) (*types.QueryDataReferenceResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	dataReference, ok := k.GetDataRegistryEntry(sdkCtx, req.GetKey())

	if !ok {
		return nil, status.Error(codes.NotFound, "not found")
	}

	res := types.QueryDataReferenceResponse{
		Data: &dataReference,
	}

	return &res, nil
}

// Query all proposal durations
func (k Keeper) AllProposalDurations(goCtx context.Context, req *types.QueryAllProposalDurations) (*types.QueryAllProposalDurationsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	properties := k.GetNetworkProperties(ctx)
	allDurations := k.GetAllProposalDurations(ctx)
	for _, pt := range kiratypes.AllProposalTypes {
		if allDurations[pt] < properties.MinimumProposalEndTime {
			allDurations[pt] = properties.MinimumProposalEndTime
		}
	}
	return &types.QueryAllProposalDurationsResponse{
		ProposalDurations: allDurations,
	}, nil
}

// Query single proposal duration
func (k Keeper) ProposalDuration(goCtx context.Context, req *types.QueryProposalDuration) (*types.QueryProposalDurationResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	properties := k.GetNetworkProperties(ctx)
	duration := k.GetProposalDuration(ctx, req.ProposalType)
	if duration < properties.MinimumProposalEndTime {
		duration = properties.MinimumProposalEndTime
	}
	return &types.QueryProposalDurationResponse{
		Duration: duration,
	}, nil
}

// Councilors - all councilors (waiting or not), including their corresponding statuses,
// ranks & abstenation counters - add sub-query to search by specific KIRA address
func (k Keeper) Councilors(goCtx context.Context, req *types.QueryCouncilors) (*types.QueryCouncilorsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if req.Address == "" {
		return &types.QueryCouncilorsResponse{
			Councilors: k.GetAllCouncilors(ctx),
		}, nil
	}

	addr, err := sdk.AccAddressFromBech32(req.Address)
	if err != nil {
		return nil, err
	}
	councilor, found := k.GetCouncilor(ctx, addr)
	if !found {
		return nil, types.ErrCouncilorNotFound
	}
	return &types.QueryCouncilorsResponse{
		Councilors: []types.Councilor{councilor},
	}, nil
}

// NonCouncilors - list all governance members that are NOT Councilors
func (k Keeper) NonCouncilors(goCtx context.Context, req *types.QueryNonCouncilors) (*types.QueryNonCouncilorsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	networkActorsIterator := k.GetNetworkActorsIterator(ctx)
	defer networkActorsIterator.Close()
	networkActors := []types.NetworkActor{}
	for ; networkActorsIterator.Valid(); networkActorsIterator.Next() {
		actor := k.GetNetworkActorFromIterator(networkActorsIterator)

		_, found := k.GetCouncilor(ctx, actor.Address)
		if !found {
			networkActors = append(networkActors, *actor)
		}
	}

	return &types.QueryNonCouncilorsResponse{
		NonCouncilors: networkActors,
	}, nil
}

// AddressesByWhitelistedPermission - list all KIRA addresses by a specific whitelisted permission (address does NOT have to be a Councilor)
func (k Keeper) AddressesByWhitelistedPermission(goCtx context.Context, req *types.QueryAddressesByWhitelistedPermission) (*types.QueryAddressesByWhitelistedPermissionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	actors := k.GetNetworkActorsByAbsoluteWhitelistPermission(ctx, types.PermValue(req.Permission))

	addrs := []string{}
	for _, actor := range actors {
		addrs = append(addrs, actor.Address.String())
	}
	return &types.QueryAddressesByWhitelistedPermissionResponse{
		Addresses: addrs,
	}, nil
}

// AddressesByBlacklistedPermission - list all KIRA addresses by a specific whitelisted permission (address does NOT have to be a Councilor)
func (k Keeper) AddressesByBlacklistedPermission(goCtx context.Context, req *types.QueryAddressesByBlacklistedPermission) (*types.QueryAddressesByBlacklistedPermissionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	networkActorsIterator := k.GetNetworkActorsIterator(ctx)
	defer networkActorsIterator.Close()
	addrs := []string{}
	for ; networkActorsIterator.Valid(); networkActorsIterator.Next() {
		actor := k.GetNetworkActorFromIterator(networkActorsIterator)
		if actor.Permissions.IsBlacklisted(types.PermValue(req.Permission)) {

			addrs = append(addrs, actor.Address.String())
		}
	}

	return &types.QueryAddressesByBlacklistedPermissionResponse{
		Addresses: addrs,
	}, nil
}

// AddressesByWhitelistedRole - list all kira addresses by a specific whitelisted role (address does NOT have to be a Councilor)
func (k Keeper) AddressesByWhitelistedRole(goCtx context.Context, req *types.QueryAddressesByWhitelistedRole) (*types.QueryAddressesByWhitelistedRoleResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	addrs := []string{}
	actorIter := k.GetNetworkActorsByRole(ctx, uint64(req.Role))
	for ; actorIter.Valid(); actorIter.Next() {
		actor := k.GetNetworkActorOrFail(ctx, actorIter.Value())
		addrs = append(addrs, actor.Address.String())
	}

	return &types.QueryAddressesByWhitelistedRoleResponse{
		Addresses: addrs,
	}, nil
}
