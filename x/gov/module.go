package gov

import (
	"encoding/json"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"

	cli2 "github.com/KiraCore/sekai/x/gov/client/cli"
	keeper2 "github.com/KiraCore/sekai/x/gov/keeper"
	customgovtypes "github.com/KiraCore/sekai/x/gov/types"

	"github.com/KiraCore/sekai/middleware"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	types2 "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/gogo/protobuf/grpc"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	abci "github.com/tendermint/tendermint/abci/types"
)

var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

type AppModuleBasic struct{}

func (b AppModuleBasic) Name() string {
	return customgovtypes.ModuleName
}

func (b AppModuleBasic) RegisterInterfaces(registry types2.InterfaceRegistry) {
	customgovtypes.RegisterInterfaces(registry)
}

func (b AppModuleBasic) DefaultGenesis(cdc codec.JSONMarshaler) json.RawMessage {
	return cdc.MustMarshalJSON(customgovtypes.DefaultGenesis())
}

func (b AppModuleBasic) ValidateGenesis(marshaler codec.JSONMarshaler, config client.TxEncodingConfig, message json.RawMessage) error {
	return nil
}

func (b AppModuleBasic) RegisterRESTRoutes(context client.Context, router *mux.Router) {
}

func (b AppModuleBasic) RegisterGRPCRoutes(context client.Context, serveMux *runtime.ServeMux) {
}

func (b AppModuleBasic) RegisterLegacyAminoCodec(amino *codec.LegacyAmino) {
}

func (b AppModuleBasic) GetTxCmd() *cobra.Command {
	return cli2.NewTxCmd()
}

// GetQueryCmd implement query commands for this module
func (b AppModuleBasic) GetQueryCmd() *cobra.Command {
	queryCmd := &cobra.Command{
		Use:   customgovtypes.RouterKey,
		Short: "query commands for the customgov module",
	}
	queryCmd.AddCommand(
		cli2.GetCmdQueryPermissions(),
		cli2.GetCmdQueryNetworkProperties(),
		cli2.GetCmdQueryExecutionFee(),
		cli2.GetCmdQueryRolePermissions(),
		cli2.GetCmdQueryRolesByAddress(),
		cli2.GetCmdQueryProposals(),
		cli2.GetCmdQueryCouncilRegistry(),
		cli2.GetCmdQueryProposal(),
	)

	queryCmd.PersistentFlags().String("node", "tcp://localhost:26657", "<host>:<port> to Tendermint RPC interface for this chain")
	return queryCmd
}

// AppModule extends the cosmos SDK gov.
type AppModule struct {
	AppModuleBasic
	customGovKeeper keeper2.Keeper
}

func (am AppModule) RegisterInterfaces(registry types2.InterfaceRegistry) {
	customgovtypes.RegisterInterfaces(registry)
}

func (am AppModule) InitGenesis(
	ctx sdk.Context,
	cdc codec.JSONMarshaler,
	data json.RawMessage,
) []abci.ValidatorUpdate {
	var genesisState customgovtypes.GenesisState
	cdc.MustUnmarshalJSON(data, &genesisState)

	for _, actor := range genesisState.NetworkActors {
		am.customGovKeeper.SaveNetworkActor(ctx, *actor)
		for _, role := range actor.Roles {
			am.customGovKeeper.AssignRoleToActor(ctx, *actor, customgovtypes.Role(role))
		}
		for _, perm := range actor.Permissions.Whitelist {
			err := am.customGovKeeper.AddWhitelistPermission(ctx, *actor, customgovtypes.PermValue(perm))
			if err != nil {
				panic(err)
			}
		}
		// TODO when we add keeper function for managing blacklist mapping, we can just enable this
		// for _, perm := range actor.Permissions.Blacklist {
		// 	am.customGovKeeper.RemoveWhitelistPermission(ctx, *actor, customgovtypes.PermValue(perm))
		// }
	}

	for index, perm := range genesisState.Permissions {
		role := customgovtypes.Role(index)
		am.customGovKeeper.CreateRole(ctx, role)
		for _, white := range perm.Whitelist {
			err := am.customGovKeeper.WhitelistRolePermission(ctx, role, customgovtypes.PermValue(white))
			if err != nil {
				panic(err)
			}
		}
		for _, black := range perm.Blacklist {
			err := am.customGovKeeper.BlacklistRolePermission(ctx, role, customgovtypes.PermValue(black))
			if err != nil {
				panic(err)
			}
		}
	}

	am.customGovKeeper.SaveProposalID(ctx, genesisState.StartingProposalId)

	am.customGovKeeper.SetNetworkProperties(ctx, genesisState.NetworkProperties)

	for _, fee := range genesisState.ExecutionFees {
		am.customGovKeeper.SetExecutionFee(ctx, fee)
	}

	return nil
}

func (am AppModule) ExportGenesis(context sdk.Context, marshaler codec.JSONMarshaler) json.RawMessage {
	return nil
}

func (am AppModule) RegisterInvariants(registry sdk.InvariantRegistry) {}

func (am AppModule) QuerierRoute() string { return "" }

func (am AppModule) LegacyQuerierHandler(marshaler *codec.LegacyAmino) sdk.Querier {
	return nil
}

func (am AppModule) BeginBlock(context sdk.Context, block abci.RequestBeginBlock) {}

func (am AppModule) EndBlock(ctx sdk.Context, block abci.RequestEndBlock) []abci.ValidatorUpdate {
	EndBlocker(ctx, am.customGovKeeper)

	return []abci.ValidatorUpdate{}
}

func (am AppModule) Name() string {
	return customgovtypes.ModuleName
}

// Route returns the message routing key for the staking module.
func (am AppModule) Route() sdk.Route {
	return middleware.NewRoute(customgovtypes.ModuleName, NewHandler(am.customGovKeeper))
}

// RegisterQueryService registers a GRPC query service to respond to the
// module-specific GRPC queries.
func (am AppModule) RegisterQueryService(server grpc.Server) {
	querier := NewQuerier(am.customGovKeeper)
	customgovtypes.RegisterQueryServer(server, querier)
}

// NewAppModule returns a new Custom Staking module.
func NewAppModule(
	keeper keeper2.Keeper,
) AppModule {
	return AppModule{
		customGovKeeper: keeper,
	}
}
