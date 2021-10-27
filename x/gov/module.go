package gov

import (
	"context"
	"encoding/json"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"

	"github.com/KiraCore/sekai/middleware"
	"github.com/KiraCore/sekai/x/gov/client/cli"
	"github.com/KiraCore/sekai/x/gov/keeper"
	"github.com/KiraCore/sekai/x/gov/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	abci "github.com/tendermint/tendermint/abci/types"
)

var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

type AppModuleBasic struct{}

// RegisterGRPCGatewayRoutes registers the gRPC Gateway routes for the staking module.
func (AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *runtime.ServeMux) {
}

func (b AppModuleBasic) Name() string {
	return types.ModuleName
}

func (b AppModuleBasic) RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	types.RegisterInterfaces(registry)
}

func (b AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	return cdc.MustMarshalJSON(types.DefaultGenesis())
}

func (b AppModuleBasic) ValidateGenesis(marshaler codec.JSONCodec, config client.TxEncodingConfig, message json.RawMessage) error {
	return nil
}

func (b AppModuleBasic) RegisterRESTRoutes(context client.Context, router *mux.Router) {}

func (b AppModuleBasic) RegisterGRPCRoutes(clientCtx client.Context, serveMux *runtime.ServeMux) {
	types.RegisterQueryHandlerClient(context.Background(), serveMux, types.NewQueryClient(clientCtx))
}

func (b AppModuleBasic) RegisterLegacyAminoCodec(amino *codec.LegacyAmino) {}

func (b AppModuleBasic) GetTxCmd() *cobra.Command {
	return cli.NewTxCmd()
}

// GetQueryCmd implement query commands for this module
func (b AppModuleBasic) GetQueryCmd() *cobra.Command {
	queryCmd := &cobra.Command{
		Use:   types.RouterKey,
		Short: "query commands for the customgov module",
	}
	queryCmd.AddCommand(
		cli.GetCmdQueryPermissions(),
		cli.GetCmdQueryNetworkProperties(),
		cli.GetCmdQueryExecutionFee(),
		cli.GetCmdQueryPoorNetworkMessages(),
		cli.GetCmdQueryRolePermissions(),
		cli.GetCmdQueryRolesByAddress(),
		cli.GetCmdQueryProposals(),
		cli.GetCmdQueryCouncilRegistry(),
		cli.GetCmdQueryProposal(),
		cli.GetCmdQueryVote(),
		cli.GetCmdQueryVotes(),
		cli.GetCmdQueryWhitelistedProposalVoters(),
		cli.GetCmdQueryIdentityRecord(),
		cli.GetCmdQueryIdentityRecordByAddress(),
		cli.GetCmdQueryAllIdentityRecords(),
		cli.GetCmdQueryIdentityRecordVerifyRequest(),
		cli.GetCmdQueryIdentityRecordVerifyRequestsByRequester(),
		cli.GetCmdQueryIdentityRecordVerifyRequestsByApprover(),
		cli.GetCmdQueryAllIdentityRecordVerifyRequests(),
		cli.GetCmdQueryAllDataReferenceKeys(),
	)

	return queryCmd
}

// AppModule extends the cosmos SDK gov.
type AppModule struct {
	AppModuleBasic
	customGovKeeper keeper.Keeper
}

// RegisterServices registers a GRPC query service to respond to the
// module-specific GRPC queries.
func (am AppModule) RegisterServices(cfg module.Configurator) {
	types.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServerImpl(am.customGovKeeper))
	types.RegisterQueryServer(cfg.QueryServer(), am.customGovKeeper)
}

func (am AppModule) RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	types.RegisterInterfaces(registry)
}

func (am AppModule) InitGenesis(
	ctx sdk.Context,
	cdc codec.JSONCodec,
	data json.RawMessage,
) []abci.ValidatorUpdate {
	var genesisState types.GenesisState
	cdc.MustUnmarshalJSON(data, &genesisState)

	InitGenesis(ctx, am.customGovKeeper, genesisState)
	return nil
}

func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	gs := ExportGenesis(ctx, am.customGovKeeper)
	return cdc.MustMarshalJSON(gs)
}

// ConsensusVersion implements AppModule/ConsensusVersion.
func (AppModule) ConsensusVersion() uint64 { return 1 }

func (am AppModule) RegisterInvariants(registry sdk.InvariantRegistry) {}

func (am AppModule) QuerierRoute() string {
	return types.QuerierRoute
}

// LegacyQuerierHandler returns the staking module sdk.Querier.
func (am AppModule) LegacyQuerierHandler(legacyQuerierCdc *codec.LegacyAmino) sdk.Querier {
	return nil
}

func (am AppModule) BeginBlock(context sdk.Context, block abci.RequestBeginBlock) {}

func (am AppModule) EndBlock(ctx sdk.Context, block abci.RequestEndBlock) []abci.ValidatorUpdate {
	EndBlocker(ctx, am.customGovKeeper)

	return []abci.ValidatorUpdate{}
}

func (am AppModule) Name() string {
	return types.ModuleName
}

// Route returns the message routing key for the staking module.
func (am AppModule) Route() sdk.Route {
	return middleware.NewRoute(types.ModuleName, NewHandler(am.customGovKeeper))
}

// NewAppModule returns a new Custom Staking module.
func NewAppModule(
	keeper keeper.Keeper,
) AppModule {
	return AppModule{
		customGovKeeper: keeper,
	}
}
