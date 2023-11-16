package collectives

import (
	"context"
	"encoding/json"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"

	collectivescli "github.com/KiraCore/sekai/x/collectives/client/cli"
	collectiveskeeper "github.com/KiraCore/sekai/x/collectives/keeper"
	"github.com/KiraCore/sekai/x/collectives/types"
	collectivestypes "github.com/KiraCore/sekai/x/collectives/types"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
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
	return collectivestypes.ModuleName
}

func (b AppModuleBasic) RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	collectivestypes.RegisterInterfaces(registry)
}

func (b AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	return cdc.MustMarshalJSON(collectivestypes.DefaultGenesis())
}

func (b AppModuleBasic) ValidateGenesis(marshaler codec.JSONCodec, config client.TxEncodingConfig, message json.RawMessage) error {
	return nil
}

func (b AppModuleBasic) RegisterRESTRoutes(clientCtx client.Context, router *mux.Router) {
}

func (b AppModuleBasic) RegisterGRPCRoutes(clientCtx client.Context, serveMux *runtime.ServeMux) {
	collectivestypes.RegisterQueryHandlerClient(context.Background(), serveMux, types.NewQueryClient(clientCtx))
}

func (b AppModuleBasic) RegisterLegacyAminoCodec(amino *codec.LegacyAmino) {
	collectivestypes.RegisterCodec(amino)
}

func (b AppModuleBasic) GetTxCmd() *cobra.Command {
	return collectivescli.NewTxCmd()
}

// GetQueryCmd implement query commands for this module
func (b AppModuleBasic) GetQueryCmd() *cobra.Command {
	return collectivescli.NewQueryCmd()
}

// AppModule for collectives management
type AppModule struct {
	AppModuleBasic
	collectivesKeeper collectiveskeeper.Keeper
}

// RegisterQueryService registers a GRPC query service to respond to the
// module-specific GRPC queries.
func (am AppModule) RegisterServices(cfg module.Configurator) {
	collectivestypes.RegisterMsgServer(cfg.MsgServer(), collectiveskeeper.NewMsgServerImpl(am.collectivesKeeper))
	querier := collectiveskeeper.NewQuerier(am.collectivesKeeper)
	collectivestypes.RegisterQueryServer(cfg.QueryServer(), querier)
}

func (am AppModule) RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	collectivestypes.RegisterInterfaces(registry)
}

func (am AppModule) InitGenesis(
	ctx sdk.Context,
	cdc codec.JSONCodec,
	data json.RawMessage,
) []abci.ValidatorUpdate {
	var genesisState collectivestypes.GenesisState
	cdc.MustUnmarshalJSON(data, &genesisState)

	return nil
}

func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	genesisState := collectivestypes.GenesisState{}
	return cdc.MustMarshalJSON(&genesisState)
}

// ConsensusVersion implements AppModule/ConsensusVersion.
func (AppModule) ConsensusVersion() uint64 { return 1 }

func (am AppModule) RegisterInvariants(registry sdk.InvariantRegistry) {}

func (am AppModule) QuerierRoute() string {
	return collectivestypes.QuerierRoute
}

func (am AppModule) BeginBlock(clientCtx sdk.Context, req abci.RequestBeginBlock) {
	am.collectivesKeeper.BeginBlocker(clientCtx)
}

func (am AppModule) EndBlock(ctx sdk.Context, block abci.RequestEndBlock) []abci.ValidatorUpdate {
	am.collectivesKeeper.EndBlocker(ctx)
	return nil
}

func (am AppModule) Name() string {
	return collectivestypes.ModuleName
}

// NewAppModule returns a new Custom Staking module.
func NewAppModule(
	keeper collectiveskeeper.Keeper,
) AppModule {
	return AppModule{
		collectivesKeeper: keeper,
	}
}
