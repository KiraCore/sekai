package keeper

import (
	"context"
	"fmt"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	kiratypes "github.com/KiraCore/sekai/types"
	kiraquery "github.com/KiraCore/sekai/types/query"
	"github.com/KiraCore/sekai/x/gov/types"
	customstakingtypes "github.com/KiraCore/sekai/x/staking/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
)

var _ types.QueryServer = Keeper{}

// RolesByAddress return roles associated to an address
func (k Keeper) RolesByAddress(goCtx context.Context, request *types.RolesByAddressRequest) (*types.RolesByAddressResponse, error) {
	actor, found := k.GetNetworkActorByAddress(sdk.UnwrapSDKContext(goCtx), request.ValAddr)
	if !found {
		return nil, customstakingtypes.ErrNetworkActorNotFound
	}

	return &types.RolesByAddressResponse{
		Roles: actor.Roles,
	}, nil
}

// CouncilorByAddress return councilor object associated to an address
func (k Keeper) CouncilorByAddress(goCtx context.Context, request *types.CouncilorByAddressRequest) (*types.CouncilorResponse, error) {
	councilor, found := k.GetCouncilor(sdk.UnwrapSDKContext(goCtx), request.ValAddr)
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

	networkActor, found := k.GetNetworkActorByAddress(ctx, request.ValAddr)
	if !found {
		return nil, customstakingtypes.ErrNetworkActorNotFound
	}

	return &types.PermissionsResponse{Permissions: networkActor.Permissions}, nil
}

// NetworkProperties return global network properties
func (k Keeper) NetworkProperties(goCtx context.Context, request *types.NetworkPropertiesRequest) (*types.NetworkPropertiesResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	networkProperties := k.GetNetworkProperties(ctx)
	return &types.NetworkPropertiesResponse{Properties: networkProperties}, nil
}

// RolePermissions returns permissions associated to a role
func (k Keeper) RolePermissions(goCtx context.Context, request *types.RolePermissionsRequest) (*types.RolePermissionsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	perms, found := k.GetPermissionsForRole(ctx, types.Role(request.Role))
	if !found {
		return nil, types.ErrRoleDoesNotExist
	}

	return &types.RolePermissionsResponse{Permissions: &perms}, nil
}

// ExecutionFee returns execution fee associated to a specific message type
func (k Keeper) ExecutionFee(goCtx context.Context, request *types.ExecutionFeeRequest) (*types.ExecutionFeeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	fee := k.GetExecutionFee(ctx, request.TransactionType)
	if fee == nil {
		return nil, sdkerrors.Wrap(types.ErrFeeNotExist, fmt.Sprintf("fee does not exist for %s", request.TransactionType))
	}
	return &types.ExecutionFeeResponse{Fee: fee}, nil
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
		return nil, sdkerrors.Wrap(types.ErrGettingProposals, fmt.Sprintf("proposal does not exist for %d", request.ProposalId))
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
		return nil, sdkerrors.Wrap(types.ErrGettingProposals, fmt.Sprintf("error getting proposals: %s", err.Error()))
	}

	store := c.KVStore(k.storeKey)

	var proposals []types.Proposal
	var pageRes *query.PageResponse
	var err error

	proposalsStore := prefix.NewStore(store, ProposalsPrefix)

	onResult := func(key []byte, value []byte, accumulate bool) (bool, error) {
		var proposal types.Proposal
		err := k.cdc.UnmarshalBinaryBare(value, &proposal)
		if err != nil {
			return false, err
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

	if request.All {
		pageRes, err = kiraquery.IterateAll(proposalsStore, request.Pagination, onResult)
	} else if request.Reverse {
		pageRes, err = kiraquery.FilteredReversePaginate(proposalsStore, request.Pagination, onResult)
	} else {
		pageRes, err = query.FilteredPaginate(proposalsStore, request.Pagination, onResult)
	}

	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrGettingProposals, fmt.Sprintf("error getting proposals: %s", err.Error()))
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
		return nil, sdkerrors.Wrap(types.ErrGettingProposals, fmt.Sprintf("proposal does not exist for %d", request.ProposalId))
	}
	actors := k.GetNetworkActorsByAbsoluteWhitelistPermission(ctx, proposal.GetContent().VotePermission())
	return &types.QueryWhitelistedProposalVotersResponse{Voters: actors}, nil
}

// Vote queries voted information based on proposalID, voterAddr.
func (k Keeper) Vote(goCtx context.Context, request *types.QueryVoteRequest) (*types.QueryVoteResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	vote, found := k.GetVote(ctx, request.ProposalId, request.Voter)
	if !found {
		return &types.QueryVoteResponse{Vote: vote}, sdkerrors.Wrap(types.ErrGettingProposalVotes, fmt.Sprintf("error getting votes for proposal %d, voter %s", request.ProposalId, request.Voter.String()))
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
		Record: k.GetIdentityRecord(ctx, request.Id),
	}

	return &res, nil
}

// IdentityRecords query identity records by creator
func (k Keeper) IdentityRecordsByAddress(goCtx context.Context, request *types.QueryIdentityRecordsByAddressRequest) (*types.QueryIdentityRecordsByAddressResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	// TODO: add pagination for QueryIdentityRecordsByAddressRequest
	res := types.QueryIdentityRecordsByAddressResponse{
		Records: k.GetIdRecordsByAddress(ctx, request.Creator),
	}

	return &res, nil
}

// AllIdentityRecords query all identity records
func (k Keeper) AllIdentityRecords(goCtx context.Context, request *types.QueryAllIdentityRecordsRequest) (*types.QueryAllIdentityRecordsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	// TODO: add pagination for QueryAllIdentityRecordsRequest
	res := types.QueryAllIdentityRecordsResponse{
		Records: k.GetAllIdentityRecords(ctx),
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

	// TODO: add pagination for QueryIdentityRecordVerifyRequestsByRequester
	res := types.QueryIdentityRecordVerifyRequestsByRequesterResponse{
		VerifyRecords: k.GetIdRecordsVerifyRequestsByRequester(ctx, request.Requester),
	}

	return &res, nil
}

// IdentityRecordVerifyRequestsByApprover query identity records verify requests by approver
func (k Keeper) IdentityRecordVerifyRequestsByApprover(goCtx context.Context, request *types.QueryIdentityRecordVerifyRequestsByApprover) (*types.QueryIdentityRecordVerifyRequestsByApproverResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	// TODO: add pagination for QueryIdentityRecordVerifyRequestsByApprover
	res := types.QueryIdentityRecordVerifyRequestsByApproverResponse{
		VerifyRecords: k.GetIdRecordsVerifyRequestsByApprover(ctx, request.Approver),
	}

	return &res, nil
}

// AllIdentityRecordVerifyRequests query all identity records verify requests
func (k Keeper) AllIdentityRecordVerifyRequests(goCtx context.Context, request *types.QueryAllIdentityRecordVerifyRequests) (*types.QueryAllIdentityRecordVerifyRequestsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	// TODO: add pagination for QueryAllIdentityRecordVerifyRequests
	res := types.QueryAllIdentityRecordVerifyRequestsResponse{
		VerifyRecords: k.GetAllIdRecordsVerifyRequests(ctx),
	}

	return &res, nil
}
