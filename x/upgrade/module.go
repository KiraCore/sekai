package upgrade

import (
	"encoding/json"

	"github.com/KiraCore/sekai/middleware"
	"github.com/KiraCore/sekai/x/upgrade/client/cli"
	"github.com/KiraCore/sekai/x/upgrade/keeper"
	"github.com/KiraCore/sekai/x/upgrade/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
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

func (b AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage { // InitGenesis is ignored, no sense in serializing future upgrades
	return cdc.MustMarshalJSON(types.DefaultGenesis())
}

func (b AppModuleBasic) ValidateGenesis(marshaler codec.JSONCodec, config client.TxEncodingConfig, message json.RawMessage) error {
	return nil
}

func (b AppModuleBasic) RegisterRESTRoutes(clientCtx client.Context, router *mux.Router) {
}

func (b AppModuleBasic) RegisterGRPCRoutes(clientCtx client.Context, serveMux *runtime.ServeMux) {
	//types.RegisterQueryHandlerClient(context.Background(), serveMux, types.NewQueryClient(clientCtx))
}

func (b AppModuleBasic) RegisterLegacyAminoCodec(amino *codec.LegacyAmino) {
	types.RegisterCodec(amino)
}

func (b AppModuleBasic) GetTxCmd() *cobra.Command {
	return cli.GetTxCmd()
}

// GetQueryCmd implement query commands for this module
func (b AppModuleBasic) GetQueryCmd() *cobra.Command {
	return cli.GetQueryCmd()
}

// AppModule for tokens management
type AppModule struct {
	AppModuleBasic
	upgradeKeeper   keeper.Keeper
	customGovKeeper types.CustomGovKeeper
}

// RegisterQueryService registers a GRPC query service to respond to the
// module-specific GRPC queries.
func (am AppModule) RegisterServices(cfg module.Configurator) {
	types.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServerImpl(am.upgradeKeeper, am.customGovKeeper))
	querier := keeper.NewQuerier(am.upgradeKeeper)
	types.RegisterQueryServer(cfg.QueryServer(), querier)
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

	if genesisState.CurrentPlan != nil {
		am.upgradeKeeper.SaveCurrentPlan(ctx, *genesisState.CurrentPlan)
	}

	if genesisState.NextPlan != nil {
		am.upgradeKeeper.SaveNextPlan(ctx, *genesisState.NextPlan)
	}

	return nil
}

func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	currentPlan, err := am.upgradeKeeper.GetCurrentPlan(ctx)
	if err != nil {
		panic(err)
	}
	nextPlan, err := am.upgradeKeeper.GetNextPlan(ctx)
	if err != nil {
		panic(err)
	}
	genesisState := types.GenesisState{
		CurrentPlan: currentPlan,
		NextPlan:    nextPlan,
	}

	return cdc.MustMarshalJSON(&genesisState)
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

func (am AppModule) BeginBlock(clientCtx sdk.Context, block abci.RequestBeginBlock) {
	BeginBlocker(am.upgradeKeeper, clientCtx, block)
}

func (am AppModule) EndBlock(ctx sdk.Context, block abci.RequestEndBlock) []abci.ValidatorUpdate {
	return nil
}

func (am AppModule) Name() string {
	return types.ModuleName
}

// Route returns the message routing key for the staking module.
func (am AppModule) Route() sdk.Route {
	return middleware.NewRoute(types.ModuleName, NewHandler(am.upgradeKeeper, am.customGovKeeper))
}

// NewAppModule returns a new Custom Staking module.
func NewAppModule(
	keeper keeper.Keeper,
	customGovKeeper types.CustomGovKeeper,
) AppModule {
	return AppModule{
		upgradeKeeper:   keeper,
		customGovKeeper: customGovKeeper,
	}
}
