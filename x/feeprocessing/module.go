package feeprocessing

import (
	"encoding/json"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"

	feeprocessingkeeper "github.com/KiraCore/sekai/x/feeprocessing/keeper"
	feeprocessingtypes "github.com/KiraCore/sekai/x/feeprocessing/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
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

func (b AppModuleBasic) RegisterGRPCGatewayRoutes(context client.Context, serveMux *runtime.ServeMux) {
}

func (b AppModuleBasic) Name() string {
	return feeprocessingtypes.ModuleName
}

func (b AppModuleBasic) RegisterInterfaces(registry codectypes.InterfaceRegistry) {
}

func (b AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	return json.RawMessage("{}")
}

func (b AppModuleBasic) ValidateGenesis(marshaler codec.JSONCodec, config client.TxEncodingConfig, message json.RawMessage) error {
	return nil
}

func (b AppModuleBasic) RegisterRESTRoutes(context client.Context, router *mux.Router) {
}

func (b AppModuleBasic) RegisterGRPCRoutes(context client.Context, serveMux *runtime.ServeMux) {
}

func (b AppModuleBasic) RegisterLegacyAminoCodec(amino *codec.LegacyAmino) {
}

func (b AppModuleBasic) GetTxCmd() *cobra.Command {
	return nil
}

// GetQueryCmd implement query commands for this module
func (b AppModuleBasic) GetQueryCmd() *cobra.Command {
	return nil
}

// AppModule extends the cosmos SDK gov.
type AppModule struct {
	AppModuleBasic
	keeper feeprocessingkeeper.Keeper
}

func (am AppModule) RegisterServices(configurator module.Configurator) {
}

func (am AppModule) RegisterInterfaces(registry codectypes.InterfaceRegistry) {
}

func (am AppModule) InitGenesis(
	ctx sdk.Context,
	cdc codec.JSONCodec,
	data json.RawMessage,
) []abci.ValidatorUpdate {
	return nil
}

func (am AppModule) ExportGenesis(context sdk.Context, marshaler codec.JSONCodec) json.RawMessage {
	return nil
}

// ConsensusVersion implements AppModule/ConsensusVersion.
func (AppModule) ConsensusVersion() uint64 { return 1 }

func (am AppModule) RegisterInvariants(registry sdk.InvariantRegistry) {}

func (am AppModule) QuerierRoute() string { return "" }

func (am AppModule) LegacyQuerierHandler(marshaler *codec.LegacyAmino) sdk.Querier {
	return nil
}

func (am AppModule) BeginBlock(context sdk.Context, block abci.RequestBeginBlock) {}

func (am AppModule) EndBlock(ctx sdk.Context, block abci.RequestEndBlock) []abci.ValidatorUpdate {
	return EndBlocker(ctx, am.keeper)
}

func (am AppModule) Name() string {
	return feeprocessingtypes.ModuleName
}

// Route returns the message routing key for the staking module.
func (am AppModule) Route() sdk.Route {
	return sdk.NewRoute(feeprocessingtypes.ModuleName, NewHandler(am.keeper))
}

// RegisterQueryService registers a GRPC query service to respond to the
// module-specific GRPC queries.
func (am AppModule) RegisterQueryService(server grpc.Server) {
}

// NewAppModule returns a new Custom Staking module.
func NewAppModule(
	keeper feeprocessingkeeper.Keeper,
) AppModule {
	return AppModule{
		keeper: keeper,
	}
}
