package ubi

import (
	"context"
	"encoding/json"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"

	ubicli "github.com/KiraCore/sekai/x/ubi/client/cli"
	ubikeeper "github.com/KiraCore/sekai/x/ubi/keeper"
	"github.com/KiraCore/sekai/x/ubi/types"
	ubitypes "github.com/KiraCore/sekai/x/ubi/types"

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
	return ubitypes.ModuleName
}

func (b AppModuleBasic) RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	ubitypes.RegisterInterfaces(registry)
}

func (b AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	return cdc.MustMarshalJSON(ubitypes.DefaultGenesis())
}

func (b AppModuleBasic) ValidateGenesis(marshaler codec.JSONCodec, config client.TxEncodingConfig, message json.RawMessage) error {
	return nil
}

func (b AppModuleBasic) RegisterRESTRoutes(clientCtx client.Context, router *mux.Router) {
}

func (b AppModuleBasic) RegisterGRPCRoutes(clientCtx client.Context, serveMux *runtime.ServeMux) {
	ubitypes.RegisterQueryHandlerClient(context.Background(), serveMux, types.NewQueryClient(clientCtx))
}

func (b AppModuleBasic) RegisterLegacyAminoCodec(amino *codec.LegacyAmino) {
	ubitypes.RegisterCodec(amino)
}

func (b AppModuleBasic) GetTxCmd() *cobra.Command {
	return ubicli.NewTxCmd()
}

// GetQueryCmd implement query commands for this module
func (b AppModuleBasic) GetQueryCmd() *cobra.Command {
	return ubicli.NewQueryCmd()
}

// AppModule for ubi management
type AppModule struct {
	AppModuleBasic
	ubiKeeper       ubikeeper.Keeper
	customGovKeeper ubitypes.CustomGovKeeper
}

// RegisterQueryService registers a GRPC query service to respond to the
// module-specific GRPC queries.
func (am AppModule) RegisterServices(cfg module.Configurator) {
	ubitypes.RegisterMsgServer(cfg.MsgServer(), ubikeeper.NewMsgServerImpl(am.ubiKeeper, am.customGovKeeper))
	querier := ubikeeper.NewQuerier(am.ubiKeeper)
	ubitypes.RegisterQueryServer(cfg.QueryServer(), querier)
}

func (am AppModule) RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	ubitypes.RegisterInterfaces(registry)
}

func (am AppModule) InitGenesis(
	ctx sdk.Context,
	cdc codec.JSONCodec,
	data json.RawMessage,
) []abci.ValidatorUpdate {
	var genesisState ubitypes.GenesisState
	cdc.MustUnmarshalJSON(data, &genesisState)

	for _, record := range genesisState.UbiRecords {
		am.ubiKeeper.SetUBIRecord(ctx, record)
	}

	return nil
}

func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	var genesisState ubitypes.GenesisState
	genesisState.UbiRecords = am.ubiKeeper.GetUBIRecords(ctx)
	return cdc.MustMarshalJSON(&genesisState)
}

// ConsensusVersion implements AppModule/ConsensusVersion.
func (AppModule) ConsensusVersion() uint64 { return 1 }

func (am AppModule) RegisterInvariants(registry sdk.InvariantRegistry) {}

func (am AppModule) QuerierRoute() string {
	return ubitypes.QuerierRoute
}

func (am AppModule) BeginBlock(clientCtx sdk.Context, block abci.RequestBeginBlock) {}

func (am AppModule) EndBlock(ctx sdk.Context, block abci.RequestEndBlock) []abci.ValidatorUpdate {
	EndBlocker(ctx, am.ubiKeeper)

	return nil
}

func (am AppModule) Name() string {
	return ubitypes.ModuleName
}

// NewAppModule returns a new Custom Staking module.
func NewAppModule(
	keeper ubikeeper.Keeper,
	customGovKeeper ubitypes.CustomGovKeeper,
) AppModule {
	return AppModule{
		ubiKeeper:       keeper,
		customGovKeeper: customGovKeeper,
	}
}
