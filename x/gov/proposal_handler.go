package gov

import (
	"github.com/KiraCore/sekai/x/gov/keeper"
	"github.com/KiraCore/sekai/x/gov/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

type ApplyAssignPermissionProposalHandler struct {
	keeper keeper.Keeper
}

func NewApplyAssignPermissionProposalHandler(keeper keeper.Keeper) *ApplyAssignPermissionProposalHandler {
	return &ApplyAssignPermissionProposalHandler{keeper: keeper}
}

func (a ApplyAssignPermissionProposalHandler) ProposalType() string {
	return types.AssignPermissionProposalType
}

func (a ApplyAssignPermissionProposalHandler) Apply(ctx sdk.Context, proposal types.Content) error {
	p := proposal.(*types.AssignPermissionProposal)

	actor, found := a.keeper.GetNetworkActorByAddress(ctx, p.Address)
	if found {
		if actor.Permissions.IsWhitelisted(types.PermValue(p.Permission)) {
			return sdkerrors.Wrap(types.ErrWhitelisting, "permission already whitelisted")
		}
	} else {
		actor = types.NewDefaultActor(p.Address)
	}

	return a.keeper.AddWhitelistPermission(ctx, actor, types.PermValue(p.Permission))
}

type ApplySetNetworkPropertyProposalHandler struct {
	keeper keeper.Keeper
}

func NewApplySetNetworkPropertyProposalHandler(keeper keeper.Keeper) *ApplySetNetworkPropertyProposalHandler {
	return &ApplySetNetworkPropertyProposalHandler{keeper: keeper}
}

func (a ApplySetNetworkPropertyProposalHandler) ProposalType() string {
	return types.SetNetworkPropertyProposalType
}

func (a ApplySetNetworkPropertyProposalHandler) Apply(ctx sdk.Context, proposal types.Content) error {
	p := proposal.(*types.SetNetworkPropertyProposal)

	property, err := a.keeper.GetNetworkProperty(ctx, p.NetworkProperty)
	if err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}
	if property == p.Value {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "network property already set as proposed value")
	}

	return a.keeper.SetNetworkProperty(ctx, p.NetworkProperty, p.Value)
}

type ApplyUpsertDataRegistryProposalHandler struct {
	keeper keeper.Keeper
}

func NewApplyUpsertDataRegistryProposalHandler(keeper keeper.Keeper) *ApplyUpsertDataRegistryProposalHandler {
	return &ApplyUpsertDataRegistryProposalHandler{keeper: keeper}
}

func (a ApplyUpsertDataRegistryProposalHandler) ProposalType() string {
	return types.UpsertDataRegistryProposalType
}

func (a ApplyUpsertDataRegistryProposalHandler) Apply(ctx sdk.Context, proposal types.Content) error {
	p := proposal.(*types.UpsertDataRegistryProposal)
	entry := types.NewDataRegistryEntry(p.Hash, p.Reference, p.Encoding, p.Size_)
	a.keeper.UpsertDataRegistryEntry(ctx, p.Key, entry)
	return nil
}

type ApplySetPoorNetworkMessagesProposalHandler struct {
	keeper keeper.Keeper
}

func NewApplySetPoorNetworkMessagesProposalHandler(keeper keeper.Keeper) *ApplySetPoorNetworkMessagesProposalHandler {
	return &ApplySetPoorNetworkMessagesProposalHandler{keeper: keeper}
}

func (a ApplySetPoorNetworkMessagesProposalHandler) ProposalType() string {
	return types.SetPoorNetworkMessagesProposalType
}

func (a ApplySetPoorNetworkMessagesProposalHandler) Apply(ctx sdk.Context, proposal types.Content) error {
	p := proposal.(*types.SetPoorNetworkMessagesProposal)
	msgs := types.AllowedMessages{Messages: p.Messages}
	a.keeper.SavePoorNetworkMessages(ctx, &msgs)
	return nil
}

type CreateRoleProposalHandler struct {
	keeper keeper.Keeper
}

func (c CreateRoleProposalHandler) ProposalType() string {
	return types.CreateRoleProposalType
}

func (c CreateRoleProposalHandler) Apply(ctx sdk.Context, proposal types.Content) error {
	p := proposal.(*types.CreateRoleProposal)

	_, exists := c.keeper.GetPermissionsForRole(ctx, types.Role(p.Role))
	if exists {
		return types.ErrRoleExist
	}

	c.keeper.CreateRole(ctx, types.Role(p.Role))

	for _, w := range p.WhitelistedPermissions {
		err := c.keeper.WhitelistRolePermission(ctx, types.Role(p.Role), w)
		if err != nil {
			return err
		}
	}

	for _, b := range p.BlacklistedPermissions {
		err := c.keeper.BlacklistRolePermission(ctx, types.Role(p.Role), b)
		if err != nil {
			return err
		}
	}
	return nil
}

func NewApplyCreateRoleProposalHandler(keeper keeper.Keeper) *CreateRoleProposalHandler {
	return &CreateRoleProposalHandler{keeper: keeper}
}
