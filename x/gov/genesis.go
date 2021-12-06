package gov

import (
	"bytes"

	"github.com/KiraCore/sekai/x/gov/keeper"
	"github.com/KiraCore/sekai/x/gov/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func InitGenesis(
	ctx sdk.Context,
	k keeper.Keeper,
	genesisState types.GenesisState,
) error {
	k.SetNextRoleId(ctx, genesisState.NextRoleId)
	for _, actor := range genesisState.NetworkActors {
		k.SaveNetworkActor(ctx, *actor)
		for _, role := range actor.Roles {
			k.AssignRoleToActor(ctx, *actor, role)
		}
		for _, perm := range actor.Permissions.Whitelist {
			k.SetWhitelistAddressPermKey(ctx, *actor, types.PermValue(perm))
		}
		// TODO when we add keeper function for managing blacklist mapping, we can just enable this
		// for _, perm := range actor.Permissions.Blacklist {
		// 	k.RemoveWhitelistPermission(ctx, *actor, types.PermValue(perm))
		// }
	}

	for _, role := range genesisState.Roles {
		k.SetRole(ctx, role)
	}

	for roleId, perm := range genesisState.RolePermissions {
		for _, white := range perm.Whitelist {
			k.WhitelistRolePermission(ctx, roleId, types.PermValue(white))
		}
		// TODO when we add keeper function for managing blacklist mapping, we can just enable this
		// for _, black := range perm.Blacklist {
		// 	err := k.BlacklistRolePermission(ctx, roleId, types.PermValue(black))
		// 	if err != nil {
		// 		// TODO: this is fine with current upgrade but from next time, it should panic
		// 		fmt.Println("There was an error blacklisting role permission", err)
		// 		// panic(err)
		// 	}
		// }
	}

	k.SetNextProposalID(ctx, genesisState.StartingProposalId)
	err := k.SetNetworkProperties(ctx, genesisState.NetworkProperties)
	if err != nil {
		panic(err)
	}

	for _, fee := range genesisState.ExecutionFees {
		k.SetExecutionFee(ctx, fee)
	}

	k.SavePoorNetworkMessages(ctx, genesisState.PoorNetworkMessages)

	for _, proposal := range genesisState.Proposals {
		k.SaveProposal(ctx, proposal)
	}

	for _, vote := range genesisState.Votes {
		k.SaveVote(ctx, vote)
	}

	for key, entry := range genesisState.DataRegistry {
		if entry == nil {
			continue
		}
		k.UpsertDataRegistryEntry(ctx, key, *entry)
	}

	for _, record := range genesisState.IdentityRecords {
		k.SetIdentityRecord(ctx, record)
	}
	for _, request := range genesisState.IdRecordsVerifyRequests {
		k.SetIdentityRecordsVerifyRequest(ctx, request)
	}

	k.SetLastIdentityRecordId(ctx, genesisState.LastIdentityRecordId)
	k.SetLastIdRecordVerifyRequestId(ctx, genesisState.LastIdRecordVerifyRequestId)

	for typeofProposal, duration := range genesisState.ProposalDurations {
		err := k.SetProposalDuration(ctx, typeofProposal, duration)
		if err != nil {
			return nil
		}
	}

	return nil
}

// ExportGenesis writes the current store values
// to a genesis file, which can be imported again
// with InitGenesis
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) (data *types.GenesisState) {
	rolesIterator := k.IterateRoles(ctx)
	defer rolesIterator.Close()
	rolePermissions := make(map[uint64]*types.Permissions)
	for ; rolesIterator.Valid(); rolesIterator.Next() {
		role := sdk.BigEndianToUint64(bytes.TrimPrefix(rolesIterator.Key(), keeper.RolePermissionRegistry))
		perms := k.GetPermissionsFromIterator(rolesIterator)
		rolePermissions[role] = &perms
	}

	networkActorsIterator := k.GetNetworkActorsIterator(ctx)
	defer networkActorsIterator.Close()
	networkActors := []*types.NetworkActor{}
	for ; networkActorsIterator.Valid(); networkActorsIterator.Next() {
		networkActors = append(networkActors, k.GetNetworkActorFromIterator(networkActorsIterator))
	}

	proposals, _ := k.GetProposals(ctx)

	return &types.GenesisState{
		StartingProposalId:          k.GetNextProposalID(ctx),
		NextRoleId:                  k.GetNextRoleId(ctx),
		Roles:                       k.GetAllRoles(ctx),
		RolePermissions:             rolePermissions,
		NetworkActors:               networkActors,
		NetworkProperties:           k.GetNetworkProperties(ctx),
		ExecutionFees:               k.GetExecutionFees(ctx),
		PoorNetworkMessages:         k.GetPoorNetworkMessages(ctx),
		Proposals:                   proposals,
		Votes:                       k.GetVotes(ctx),
		DataRegistry:                k.AllDataRegistry(ctx),
		IdentityRecords:             k.GetAllIdentityRecords(ctx),
		LastIdentityRecordId:        k.GetLastIdentityRecordId(ctx),
		IdRecordsVerifyRequests:     k.GetAllIdRecordsVerifyRequests(ctx),
		LastIdRecordVerifyRequestId: k.GetLastIdRecordVerifyRequestId(ctx),
		ProposalDurations:           k.GetAllProposalDurations(ctx),
	}
}
