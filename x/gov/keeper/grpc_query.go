package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/KiraCore/sekai/x/gov/types"
	customstakingtypes "github.com/KiraCore/sekai/x/staking/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
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
	msg, ok := q.keeper.GetPoorNetworkMsgs(sdkContext)
	if !ok {
		return nil, types.ErrPoorNetworkMsgsNotSet
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

// Proposals query proposals by querying params
func (q Querier) Proposals(ctx context.Context, request *types.QueryProposalsRequest) (*types.QueryProposalsResponse, error) {
	fmt.Println("proposals request", request)
	sdkContext := sdk.UnwrapSDKContext(ctx)
	proposals, err := q.keeper.GetProposals(sdkContext)
	fmt.Println("proposals", proposals, err)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrGettingProposals, fmt.Sprintf("error getting proposals: %s", err.Error()))
	}
	return &types.QueryProposalsResponse{Proposals: proposals}, nil
}

// GetWhitelistedProposalVoters returns whitelisted voters for a proposal for tracking
func (q Querier) GetWhitelistedProposalVoters(ctx context.Context, request *types.QueryWhitelistedProposalVotersRequest) (*types.QueryWhitelistedProposalVotersResponse, error) {
	sdkContext := sdk.UnwrapSDKContext(ctx)
	// TODO: this should get availableVoters by proposal type
	actors := q.keeper.GetNetworkActorsByAbsoluteWhitelistPermission(sdkContext, types.PermVoteSetPermissionProposal)
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
