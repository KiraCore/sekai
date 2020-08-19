package kiraHub

import (
	"encoding/json"

	cdcTypes "github.com/cosmos/cosmos-sdk/codec/types"

	"github.com/cosmos/cosmos-sdk/client"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"

	"github.com/KiraCore/sekai/x/kiraHub/client/cli"
	"github.com/KiraCore/sekai/x/kiraHub/keeper"
	"github.com/KiraCore/sekai/x/kiraHub/types"
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

type AppModuleBasic struct {
}

func (b AppModuleBasic) RegisterInterfaces(registry cdcTypes.InterfaceRegistry) {
	types.RegisterInterfaces(registry)
}

func (AppModuleBasic) Name() string {
	return types.ModuleName
}

func (AppModuleBasic) RegisterCodec(codec *codec.LegacyAmino) {
	types.RegisterCodec(codec)
}

func (AppModuleBasic) DefaultGenesis(jsonMarshaler codec.JSONMarshaler) json.RawMessage {
	return jsonMarshaler.MustMarshalJSON(DefaultGenesisState())
}

func (AppModuleBasic) RegisterRESTRoutes(cliContext client.Context, router *mux.Router) {
	RegisterRESTRoutes(cliContext, router)
}

func (b AppModuleBasic) ValidateGenesis(marshaler codec.JSONMarshaler, config client.TxEncodingConfig, message json.RawMessage) error {
	var genesisState GenesisState
	Error := marshaler.UnmarshalJSON(message, &genesisState)
	if Error != nil {
		return Error
	}
	return ValidateGenesis(genesisState)
}

func (AppModuleBasic) GetTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:   types.RouterKey,
		Short: "kiraHub transaction subcommands",
	}

	txCmd.AddCommand(
		cli.CreateOrderBook(),
		cli.CreateOrder(),
		cli.UpsertSignerKey())

	txCmd.PersistentFlags().String("node", "tcp://localhost:26657", "<host>:<port> to Tendermint RPC interface for this chain")
	txCmd.PersistentFlags().String("keyring-backend", "os", "Select keyring's backend (os|file|test)")
	txCmd.PersistentFlags().String("from", "", "Name or address of private key with which to sign")
	txCmd.PersistentFlags().String("broadcast-mode", "sync", "Transaction broadcasting mode (sync|async|block)")

	return txCmd
}

func (b AppModuleBasic) GetQueryCmd() *cobra.Command {
	queryCmd := &cobra.Command{
		Use:   types.RouterKey,
		Short: "query commands for the kiraHub module",
	}
	queryCmd.AddCommand(
		cli.GetOrderBooksCmd(),
		cli.GetOrderBooksByTPCmd(),
		cli.GetOrdersCmd(),
		cli.ListSignerKeysCmd())

	queryCmd.PersistentFlags().String("node", "tcp://localhost:26657", "<host>:<port> to Tendermint RPC interface for this chain")
	return queryCmd
}

type AppModule struct {
	AppModuleBasic
	keeper keeper.Keeper
}

func (am AppModule) ValidateGenesis(marshaler codec.JSONMarshaler, config client.TxEncodingConfig, message json.RawMessage) error {
	return nil
}

func (am AppModule) RegisterInterfaces(registry cdcTypes.InterfaceRegistry) {
	panic("implement me")
}

func (am AppModule) NewHandler() sdkTypes.Handler {
	return NewHandler(am.keeper)
}

func (am AppModule) InitGenesis(context sdkTypes.Context, jsonMarshaler codec.JSONMarshaler, rawMessage json.RawMessage) []abciTypes.ValidatorUpdate {
	var genesisState GenesisState
	jsonMarshaler.MustUnmarshalJSON(rawMessage, &genesisState)
	InitializeGenesisState(context, am.keeper, genesisState)
	return []abciTypes.ValidatorUpdate{}
}
func (am AppModule) ExportGenesis(context sdkTypes.Context, jsonMarshaler codec.JSONMarshaler) json.RawMessage {
	gs := ExportGenesis(context, am.keeper)
	return jsonMarshaler.MustMarshalJSON(gs)
}

func (am AppModule) RegisterInvariants(_ sdkTypes.InvariantRegistry) {}

func (am AppModule) QuerierRoute() string {
	return QuerierRoute
}

func (am AppModule) LegacyQuerierHandler(marshaler codec.JSONMarshaler) sdkTypes.Querier {
	return nil
}

func (am AppModule) BeginBlock(_ sdkTypes.Context, _ abciTypes.RequestBeginBlock) {}

func (am AppModule) EndBlock(_ sdkTypes.Context, _ abciTypes.RequestEndBlock) []abciTypes.ValidatorUpdate {
	return []abciTypes.ValidatorUpdate{}
}

func (am AppModule) Name() string {
	return ModuleName
}

func (am AppModule) Route() sdkTypes.Route {
	return sdkTypes.NewRoute(ModuleName, NewHandler(am.keeper))
}

func (am AppModule) RegisterQueryService(server grpc.Server) {
	querier := NewQuerier(am.keeper)
	types.RegisterQueryServer(server, querier)
}

func NewAppModule(keeper keeper.Keeper) AppModule {
	return AppModule{keeper: keeper}
}
