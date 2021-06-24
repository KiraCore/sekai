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

// Querier describes grpc querier
type Querier struct {
	keeper Keeper
}

// NewQuerier returns Querier instance
func NewQuerier(keeper Keeper) types.QueryServer {
	return &Querier{keeper: keeper}
}

var _ types.QueryServer = Querier{}

// RolesByAddress return roles associated to an address
func (q Querier) RolesByAddress(ctx context.Context, request *types.RolesByAddressRequest) (*types.RolesByAddressResponse, error) {
	actor, found := q.keeper.GetNetworkActorByAddress(sdk.UnwrapSDKContext(ctx), request.ValAddr)
	if !found {
		return nil, customstakingtypes.ErrNetworkActorNotFound
	}

	return &types.RolesByAddressResponse{
		Roles: actor.Roles,
	}, nil
}

// CouncilorByAddress return councilor object associated to an address
func (q Querier) CouncilorByAddress(ctx context.Context, request *types.CouncilorByAddressRequest) (*types.CouncilorResponse, error) {
	councilor, found := q.keeper.GetCouncilor(sdk.UnwrapSDKContext(ctx), request.ValAddr)
	if !found {
		return nil, types.ErrCouncilorNotFound
	}

	return &types.CouncilorResponse{Councilor: councilor}, nil
}

// CouncilorByMoniker return councilor object named moniker
func (q Querier) CouncilorByMoniker(ctx context.Context, request *types.CouncilorByMonikerRequest) (*types.CouncilorResponse, error) {
	councilor, found := q.keeper.GetCouncilorByMoniker(sdk.UnwrapSDKContext(ctx), request.Moniker)
	if !found {
		return nil, types.ErrCouncilorNotFound
	}

	return &types.CouncilorResponse{Councilor: councilor}, nil
}

// PermissionsByAddress returns permissions associated to an address
func (q Querier) PermissionsByAddress(ctx context.Context, request *types.PermissionsByAddressRequest) (*types.PermissionsResponse, error) {
	sdkContext := sdk.UnwrapSDKContext(ctx)

	networkActor, found := q.keeper.GetNetworkActorByAddress(sdkContext, request.ValAddr)
	if !found {
		return nil, customstakingtypes.ErrNetworkActorNotFound
	}

	return &types.PermissionsResponse{Permissions: networkActor.Permissions}, nil
}

// GetNetworkProperties return global network properties
func (q Querier) GetNetworkProperties(ctx context.Context, request *types.NetworkPropertiesRequest) (*types.NetworkPropertiesResponse, error) {
	sdkContext := sdk.UnwrapSDKContext(ctx)

	networkProperties := q.keeper.GetNetworkProperties(sdkContext)
	return &types.NetworkPropertiesResponse{Properties: networkProperties}, nil
}

// RolePermissions returns permissions associated to a role
func (q Querier) RolePermissions(ctx context.Context, request *types.RolePermissionsRequest) (*types.RolePermissionsResponse, error) {
	sdkContext := sdk.UnwrapSDKContext(ctx)

	perms, found := q.keeper.GetPermissionsForRole(sdkContext, types.Role(request.Role))
	if !found {
		return nil, types.ErrRoleDoesNotExist
	}

	return &types.RolePermissionsResponse{Permissions: &perms}, nil
}

// GetExecutionFee returns execution fee associated to a specific message type
func (q Querier) GetExecutionFee(ctx context.Context, request *types.ExecutionFeeRequest) (*types.ExecutionFeeResponse, error) {
	sdkContext := sdk.UnwrapSDKContext(ctx)
	fee := q.keeper.GetExecutionFee(sdkContext, request.TransactionType)
	if fee == nil {
		return nil, sdkerrors.Wrap(types.ErrFeeNotExist, fmt.Sprintf("fee does not exist for %s", request.TransactionType))
	}
	return &types.ExecutionFeeResponse{Fee: fee}, nil
}

// GetPoorNetworkMessages queries poor network messages
func (q Querier) GetPoorNetworkMessages(ctx context.Context, request *types.PoorNetworkMessagesRequest) (*types.PoorNetworkMessagesResponse, error) {
	sdkContext := sdk.UnwrapSDKContext(ctx)
	msg, ok := q.keeper.GetPoorNetworkMessages(sdkContext)
	if !ok {
		return nil, types.ErrPoorNetworkMessagesNotSet
	}
	return &types.PoorNetworkMessagesResponse{Messages: msg.Messages}, nil
}

// Proposal returns a proposal by id
func (q Querier) Proposal(ctx context.Context, request *types.QueryProposalRequest) (*types.QueryProposalResponse, error) {
	sdkContext := sdk.UnwrapSDKContext(ctx)
	proposal, found := q.keeper.GetProposal(sdkContext, request.ProposalId)
	if found == false {
		return nil, sdkerrors.Wrap(types.ErrGettingProposals, fmt.Sprintf("proposal does not exist for %d", request.ProposalId))
	}
	votes := q.keeper.GetProposalVotes(sdkContext, request.ProposalId)
	return &types.QueryProposalResponse{
		Proposal: proposal,
		Votes:    votes,
	}, nil
}

// Proposals query proposals by querying params with pagination
func (q Querier) Proposals(ctx context.Context, request *types.QueryProposalsRequest) (*types.QueryProposalsResponse, error) {
	c := sdk.UnwrapSDKContext(ctx)
	if request == nil {
		err := status.Error(codes.InvalidArgument, "empty request")
		return nil, sdkerrors.Wrap(types.ErrGettingProposals, fmt.Sprintf("error getting proposals: %s", err.Error()))
	}

	store := c.KVStore(q.keeper.storeKey)

	var proposals []types.Proposal
	var pageRes *query.PageResponse
	var err error

	proposalsStore := prefix.NewStore(store, ProposalsPrefix)

	onResult := func(key []byte, value []byte, accumulate bool) (bool, error) {
		var proposal types.Proposal
		err := q.keeper.cdc.UnmarshalBinaryBare(value, &proposal)
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

// GetWhitelistedProposalVoters returns whitelisted voters for a proposal for tracking
func (q Querier) GetWhitelistedProposalVoters(ctx context.Context, request *types.QueryWhitelistedProposalVotersRequest) (*types.QueryWhitelistedProposalVotersResponse, error) {
	sdkContext := sdk.UnwrapSDKContext(ctx)
	proposal, found := q.keeper.GetProposal(sdkContext, request.ProposalId)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrGettingProposals, fmt.Sprintf("proposal does not exist for %d", request.ProposalId))
	}
	actors := q.keeper.GetNetworkActorsByAbsoluteWhitelistPermission(sdkContext, proposal.GetContent().VotePermission())
	return &types.QueryWhitelistedProposalVotersResponse{Voters: actors}, nil
}

// Vote queries voted information based on proposalID, voterAddr.
func (q Querier) Vote(ctx context.Context, request *types.QueryVoteRequest) (*types.QueryVoteResponse, error) {
	sdkContext := sdk.UnwrapSDKContext(ctx)
	vote, found := q.keeper.GetVote(sdkContext, request.ProposalId, request.Voter)
	if !found {
		return &types.QueryVoteResponse{Vote: vote}, sdkerrors.Wrap(types.ErrGettingProposalVotes, fmt.Sprintf("error getting votes for proposal %d, voter %s", request.ProposalId, request.Voter.String()))
	}
	return &types.QueryVoteResponse{Vote: vote}, nil
}

// Votes queries votes of a given proposal.
func (q Querier) Votes(ctx context.Context, request *types.QueryVotesRequest) (*types.QueryVotesResponse, error) {
	sdkContext := sdk.UnwrapSDKContext(ctx)
	votes := q.keeper.GetProposalVotes(sdkContext, request.ProposalId)
	return &types.QueryVotesResponse{Votes: votes}, nil
}

// GetAllDataReferenceKeys queries all data reference keys with pagination
func (q Querier) GetAllDataReferenceKeys(ctx context.Context, request *types.QueryDataReferenceKeysRequest) (*types.QueryDataReferenceKeysResponse, error) {
	sdkContext := sdk.UnwrapSDKContext(ctx)
	return q.keeper.GetAllDataReferenceKeys(sdkContext, request)
}

// GetDataReferenceByKey queries data reference by key
func (q Querier) GetDataReferenceByKey(ctx context.Context, request *types.QueryDataReferenceRequest) (*types.QueryDataReferenceResponse, error) {
	sdkContext := sdk.UnwrapSDKContext(ctx)
	return q.keeper.GetDataReferenceByKey(sdkContext, request)
}

// GetIdentityRecord query identity record by id
func (q Querier) GetIdentityRecord(ctx context.Context, request *types.QueryIdentityRecordRequest) (*types.QueryIdentityRecordResponse, error) {
	sdkContext := sdk.UnwrapSDKContext(ctx)
	return q.keeper.GetIdentityRecord(sdkContext, request)
}

// GetIdentityRecords query identity records by creator
func (q Querier) GetIdentityRecords(ctx context.Context, request *types.QueryIdentityRecordsRequest) (*types.QueryIdentityRecordsResponse, error) {
	sdkContext := sdk.UnwrapSDKContext(ctx)
	return q.keeper.GetIdentityRecords(sdkContext, request)
}

// GetAllIdentityRecords query all identity records
func (q Querier) GetAllIdentityRecords(ctx context.Context, request *types.QueryAllIdentityRecordsRequest) (*types.QueryAllIdentityRecordsResponse, error) {
	sdkContext := sdk.UnwrapSDKContext(ctx)
	return q.keeper.GetAllIdentityRecords(sdkContext, request)
}

// GetIdentityRecordVerifyRequest query identity record verify request by id
func (q Querier) GetIdentityRecordVerifyRequest(ctx context.Context, request *types.QueryIdentityVerifyRecordRequest) (*types.QueryIdentityVerifyRecordResponse, error) {
	sdkContext := sdk.UnwrapSDKContext(ctx)
	return q.keeper.GetIdentityRecordVerifyRequest(sdkContext, request)
}

// GetIdentityRecordVerifyRequests query identity record verify request by id
func (q Querier) GetIdentityRecordVerifyRequests(ctx context.Context, request *types.QueryIdentityRecordVerifyRequests) (*types.QueryIdentityRecordVerifyRequestsResponse, error) {
	sdkContext := sdk.UnwrapSDKContext(ctx)
	return q.keeper.GetIdentityRecordVerifyRequests(sdkContext, request)
}

// GetIdentityRecordVerifyRequestsByApprover query identity records verify requests by approver
func (q Querier) GetIdentityRecordVerifyRequestsByApprover(ctx context.Context, request *types.QueryIdentityRecordVerifyRequestsByApprover) (*types.QueryIdentityRecordVerifyRequestsByApproverResponse, error) {
	sdkContext := sdk.UnwrapSDKContext(ctx)
	return q.keeper.GetIdentityRecordVerifyRequestsByApprover(sdkContext, request)
}

// GetAllIdentityRecordVerifyRequests query all identity records verify requests
func (q Querier) GetAllIdentityRecordVerifyRequests(ctx context.Context, request *types.QueryAllIdentityRecordVerifyRequests) (*types.QueryAllIdentityRecordVerifyRequestsResponse, error) {
	sdkContext := sdk.UnwrapSDKContext(ctx)
	return q.keeper.GetAllIdentityRecordVerifyRequests(sdkContext, request)
}
