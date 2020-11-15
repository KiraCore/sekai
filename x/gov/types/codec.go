package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterInterface((*Content)(nil), nil)

	registerPermissionsCodec(cdc)
	registerRolesCodec(cdc)
	registerCouncilorCodec(cdc)
	registerProposalCodec(cdc)

	cdc.RegisterConcrete(&MsgSetNetworkProperties{}, "kiraHub/MsgSetNetworkProperties", nil)
	cdc.RegisterConcrete(&MsgSetExecutionFee{}, "kiraHub/MsgSetExecutionFee", nil)
}

func registerProposalCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgProposalAssignPermission{}, "kiraHub/MsgProposalAssignPermission", nil)
	cdc.RegisterConcrete(&MsgProposalUpsertDataRegistry{}, "kiraHub/MsgProposalUpsertDataRegistry", nil)
	cdc.RegisterConcrete(&MsgVoteProposal{}, "kiraHub/MsgVoteProposal", nil)
}

func registerCouncilorCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgClaimCouncilor{}, "kiraHub/MsgClaimCouncilor", nil)
}

func registerPermissionsCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgWhitelistPermissions{}, "kiraHub/MsgWhitelistPermissions", nil)
	cdc.RegisterConcrete(&MsgBlacklistPermissions{}, "kiraHub/MsgBlacklistPermissions", nil)
}

func registerRolesCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgCreateRole{}, "kiraHub/MsgCreateRole", nil)
	cdc.RegisterConcrete(&MsgAssignRole{}, "kiraHub/MsgAssignRole", nil)
	cdc.RegisterConcrete(&MsgRemoveRole{}, "kiraHub/MsgRemoveRole", nil)

	cdc.RegisterConcrete(&MsgWhitelistRolePermission{}, "kiraHub/MsgWhitelistRolePermission", nil)
	cdc.RegisterConcrete(&MsgBlacklistRolePermission{}, "kiraHub/MsgBlacklistRolePermission", nil)
	cdc.RegisterConcrete(&MsgRemoveWhitelistRolePermission{}, "kiraHub/MsgRemoveWhitelistRolePermission", nil)
	cdc.RegisterConcrete(&MsgRemoveBlacklistRolePermission{}, "kiraHub/MsgRemoveBlacklistRolePermission", nil)
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgWhitelistPermissions{},
		&MsgBlacklistPermissions{},

		&MsgSetNetworkProperties{},
		&MsgSetExecutionFee{},

		&MsgClaimCouncilor{},

		&MsgAssignRole{},
		&MsgCreateRole{},
		&MsgRemoveRole{},

		&MsgWhitelistRolePermission{},
		&MsgBlacklistRolePermission{},
		&MsgRemoveWhitelistRolePermission{},
		&MsgRemoveBlacklistRolePermission{},

		&MsgProposalAssignPermission{},
		&MsgProposalUpsertDataRegistry{},
		&MsgVoteProposal{},
	)

	registry.RegisterInterface(
		"kira.gov.Content",
		(*Content)(nil),
		&AssignPermissionProposal{},
		&UpsertDataRegistryProposal{},
	)

}

var (
	amino = codec.NewLegacyAmino()

	// ModuleCdc references the global x/staking module codec. Note, the codec should
	// ONLY be used in certain instances of tests and for JSON encoding as Amino is
	// still used for that purpose.
	//
	// The actual codec used for serialization should be provided to x/staking and
	// defined at the application level.
	ModuleCdc = codec.NewAminoCodec(amino)
)

func init() {
	RegisterCodec(amino)
	amino.Seal()
}
