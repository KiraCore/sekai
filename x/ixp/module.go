package ixp

import (
	"encoding/json"

	cdcTypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"

	"github.com/cosmos/cosmos-sdk/client"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"

	"github.com/KiraCore/sekai/x/ixp/client/cli"
	"github.com/KiraCore/sekai/x/ixp/keeper"
	"github.com/KiraCore/sekai/x/ixp/types"
	"github.com/gogo/protobuf/grpc"
	abciTypes "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
)

var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

// AppModuleBasic describe basic utilities for module
type AppModuleBasic struct {
}

// RegisterInterfaces register interfaces used in the module
func (b AppModuleBasic) RegisterInterfaces(registry cdcTypes.InterfaceRegistry) {
	types.RegisterInterfaces(registry)
}

// Name returns module name
func (AppModuleBasic) Name() string {
	return types.ModuleName
}

// RegisterCodec register codec for app module
func (AppModuleBasic) RegisterCodec(codec *codec.LegacyAmino) {
	types.RegisterCodec(codec)
}

// DefaultGenesis returns default genesis json data for this module
func (AppModuleBasic) DefaultGenesis(jsonMarshaler codec.JSONMarshaler) json.RawMessage {
	return nil
	// return jsonMarshaler.MustMarshalJSON(DefaultGenesisState())
}

// RegisterRESTRoutes register REST routes
func (AppModuleBasic) RegisterRESTRoutes(cliContext client.Context, router *mux.Router) {
	RegisterRESTRoutes(cliContext, router)
}

// RegisterGRPCRoutes register GRPC routes
func (b AppModuleBasic) RegisterGRPCRoutes(context client.Context, serveMux *runtime.ServeMux) {
}

// RegisterGRPCRoutes register LegacyAminoCodec
func (AppModuleBasic) RegisterLegacyAminoCodec(amino *codec.LegacyAmino) {
}

// ValidateGenesis validate genesis for this module
func (b AppModuleBasic) ValidateGenesis(marshaler codec.JSONMarshaler, config client.TxEncodingConfig, message json.RawMessage) error {
	// var genesisState types.GenesisState
	// Error := marshaler.UnmarshalJSON(message, &genesisState)
	// if Error != nil {
	// 	return Error
	// }
	// return types.ValidateGenesis(genesisState)
	return nil
}

// GetTxCmd implement transaction commands for this module
func (AppModuleBasic) GetTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:   types.RouterKey,
		Short: "ixp transaction subcommands",
	}

	txCmd.AddCommand(
		cli.CreateOrderBook(),
		cli.CreateOrder(),
		cli.CancelOrder(),
		cli.UpsertSignerKey())

	txCmd.PersistentFlags().String("node", "tcp://localhost:26657", "<host>:<port> to Tendermint RPC interface for this chain")
	txCmd.PersistentFlags().String("keyring-backend", "os", "Select keyring's backend (os|file|test)")
	txCmd.PersistentFlags().String("from", "", "Name or address of private key with which to sign")
	txCmd.PersistentFlags().String("broadcast-mode", "sync", "Transaction broadcasting mode (sync|async|block)")

	return txCmd
}

// GetQueryCmd implement query commands for this module
func (b AppModuleBasic) GetQueryCmd() *cobra.Command {
	queryCmd := &cobra.Command{
		Use:   types.RouterKey,
		Short: "query commands for the ixp module",
	}
	queryCmd.AddCommand(
		cli.GetOrderBooksCmd(),
		cli.GetOrderBooksByTradingPairCmd(),
		cli.GetOrdersCmd(),
		cli.GetSignerKeysCmd())

	queryCmd.PersistentFlags().String("node", "tcp://localhost:26657", "<host>:<port> to Tendermint RPC interface for this chain")
	return queryCmd
}

// AppModule struct describes module specific functions
type AppModule struct {
	AppModuleBasic
	keeper keeper.Keeper
}

// ValidateGenesis validate genesis for this module
func (am AppModule) ValidateGenesis(marshaler codec.JSONMarshaler, config client.TxEncodingConfig, message json.RawMessage) error {
	return nil
}

// RegisterInterfaces register interfaces used in the module
func (am AppModule) RegisterInterfaces(registry cdcTypes.InterfaceRegistry) {
	types.RegisterInterfaces(registry)
}

// NewHandler returns messages handler for this module
func (am AppModule) NewHandler() sdkTypes.Handler {
	return NewHandler(am.keeper)
}

// InitGenesis returns genesis state for this module
func (am AppModule) InitGenesis(context sdkTypes.Context, jsonMarshaler codec.JSONMarshaler, rawMessage json.RawMessage) []abciTypes.ValidatorUpdate {
	var genesisState types.GenesisState
	jsonMarshaler.MustUnmarshalJSON(rawMessage, &genesisState)
	return []abciTypes.ValidatorUpdate{}
}

// ExportGenesis returns json format of genesis data
func (am AppModule) ExportGenesis(context sdkTypes.Context, jsonMarshaler codec.JSONMarshaler) json.RawMessage {
	return nil
}

// RegisterInvariants register invariants
func (am AppModule) RegisterInvariants(_ sdkTypes.InvariantRegistry) {}

// QuerierRoute return querier route string value
func (am AppModule) QuerierRoute() string {
	return QuerierRoute
}

// LegacyQuerierHandler returns querier handler
func (am AppModule) LegacyQuerierHandler(marshaler *codec.LegacyAmino) sdkTypes.Querier {
	return nil
}

// BeginBlock do action for begin block
func (am AppModule) BeginBlock(_ sdkTypes.Context, _ abciTypes.RequestBeginBlock) {}

// EndBlock do action for end block
func (am AppModule) EndBlock(_ sdkTypes.Context, _ abciTypes.RequestEndBlock) []abciTypes.ValidatorUpdate {
	return []abciTypes.ValidatorUpdate{}
}

// Name returns module's name
func (am AppModule) Name() string {
	return ModuleName
}

// Route returns module route
func (am AppModule) Route() sdkTypes.Route {
	return sdkTypes.NewRoute(ModuleName, NewHandler(am.keeper))
}

// RegisterQueryService register query server
func (am AppModule) RegisterQueryService(server grpc.Server) {
	querier := NewQuerier(am.keeper)
	types.RegisterQueryServer(server, querier)
}

// NewAppModule returns new app module instance
func NewAppModule(keeper keeper.Keeper) AppModule {
	return AppModule{keeper: keeper}
}
