package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	registerPermissionsCodec(cdc)
	registerRolesCodec(cdc)
	registerCouncilorCodec(cdc)
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

		&MsgClaimCouncilor{},

		&MsgAssignRole{},
		&MsgCreateRole{},
		&MsgRemoveRole{},

		&MsgWhitelistRolePermission{},
		&MsgBlacklistRolePermission{},
		&MsgRemoveWhitelistRolePermission{},
		&MsgRemoveBlacklistRolePermission{},
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
