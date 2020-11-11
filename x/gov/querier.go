package gov

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/KiraCore/sekai/x/gov/keeper"
	"github.com/KiraCore/sekai/x/gov/types"
	cumstomtypes "github.com/KiraCore/sekai/x/staking/types"
)

// Querier describes grpc querier
type Querier struct {
	keeper keeper.Keeper
}

// NewQuerier returns Querier instance
func NewQuerier(keeper keeper.Keeper) types.QueryServer {
	return &Querier{keeper: keeper}
}

// RolesByAddress return roles associated to an address
func (q Querier) RolesByAddress(ctx context.Context, request *types.RolesByAddressRequest) (*types.RolesByAddressResponse, error) {
	actor, found := q.keeper.GetNetworkActorByAddress(sdk.UnwrapSDKContext(ctx), request.ValAddr)
	if !found {
		return nil, cumstomtypes.ErrNetworkActorNotFound
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
		return nil, cumstomtypes.ErrNetworkActorNotFound
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
		return nil, fmt.Errorf("fee does not exist for %s", request.TransactionType)
	}
	return &types.ExecutionFeeResponse{Fee: fee}, nil
}

// Proposal returns a proposal by id
func (q Querier) Proposal(ctx context.Context, request *types.QueryProposalRequest) (*types.QueryProposalResponse, error) {
	sdkContext := sdk.UnwrapSDKContext(ctx)
	proposal, found := q.keeper.GetProposal(sdkContext, request.ProposalId)
	if found == false {
		return nil, fmt.Errorf("proposal does not exist for %d", request.ProposalId)
	}
	return &types.QueryProposalResponse{Proposal: proposal}, nil
}

// Proposals query proposals by querying params
func (q Querier) Proposals(ctx context.Context, request *types.QueryProposalsRequest) (*types.QueryProposalsResponse, error) {
	sdkContext := sdk.UnwrapSDKContext(ctx)
	proposals, err := q.keeper.GetProposals(sdkContext)
	if err != nil {
		return nil, fmt.Errorf("error getting proposals: %s", err.Error())
	}
	return &types.QueryProposalsResponse{Proposals: proposals}, nil
}
