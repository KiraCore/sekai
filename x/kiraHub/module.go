package kiraHub

import (
	"encoding/json"

	"github.com/KiraCore/cosmos-sdk/codec/types"

	"github.com/KiraCore/cosmos-sdk/client"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"

	constants "github.com/KiraCore/sekai/x/kiraHub/constants"
	"github.com/gogo/protobuf/grpc"
	abciTypes "github.com/tendermint/tendermint/abci/types"

	"github.com/KiraCore/cosmos-sdk/codec"
	sdkTypes "github.com/KiraCore/cosmos-sdk/types"
	"github.com/KiraCore/cosmos-sdk/types/module"
)

var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

type AppModuleBasic struct {
}

func (b AppModuleBasic) RegisterInterfaces(registry types.InterfaceRegistry) {
	panic("implement me")
}

func (AppModuleBasic) Name() string {
	return constants.ModuleName
}

func (AppModuleBasic) RegisterCodec(codec *codec.Codec) {
	RegisterCodec(codec)
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
	panic("implement me")
}

func (b AppModuleBasic) GetQueryCmd() *cobra.Command {
	panic("implement me")
}

type AppModule struct {
	AppModuleBasic
	keeper Keeper
}

func NewAppModule(keeper Keeper) AppModule {
	return AppModule{keeper: keeper}
}
func (AppModule) Name() string {
	return ModuleName
}
func (appModule AppModule) RegisterInterfaces(registry types.InterfaceRegistry) {
	panic("implement me")
}

func (appModule AppModule) ValidateGenesis(marshaler codec.JSONMarshaler, config client.TxEncodingConfig, message json.RawMessage) error {
	panic("implement me")
}

func (appModule AppModule) GetTxCmd() *cobra.Command {
	panic("implement me")
}

func (appModule AppModule) GetQueryCmd() *cobra.Command {
	panic("implement me")
}

func (appModule AppModule) Route() sdkTypes.Route {
	return sdkTypes.NewRoute(ModuleName, NewHandler(appModule.keeper))
}

func (appModule AppModule) LegacyQuerierHandler(marshaler codec.JSONMarshaler) sdkTypes.Querier {
	panic("implement me")
}

func (appModule AppModule) RegisterQueryService(server grpc.Server) {
	panic("implement me")
}

func (appModule AppModule) RegisterInvariants(_ sdkTypes.InvariantRegistry) {}

func (appModule AppModule) NewHandler() sdkTypes.Handler {
	return NewHandler(appModule.keeper)
}
func (AppModule) QuerierRoute() string {
	return QuerierRoute
}
func (appModule AppModule) NewQuerierHandler() sdkTypes.Querier {
	return NewQuerier(appModule.keeper)
}
func (appModule AppModule) InitGenesis(context sdkTypes.Context, jsonMarshaler codec.JSONMarshaler, rawMessage json.RawMessage) []abciTypes.ValidatorUpdate {
	var genesisState GenesisState
	jsonMarshaler.MustUnmarshalJSON(rawMessage, &genesisState)
	InitializeGenesisState(context, appModule.keeper, genesisState)
	return []abciTypes.ValidatorUpdate{}
}
func (appModule AppModule) ExportGenesis(context sdkTypes.Context, jsonMarshaler codec.JSONMarshaler) json.RawMessage {
	gs := ExportGenesis(context, appModule.keeper)
	return jsonMarshaler.MustMarshalJSON(gs)
}
func (AppModule) BeginBlock(_ sdkTypes.Context, _ abciTypes.RequestBeginBlock) {}

func (AppModule) EndBlock(_ sdkTypes.Context, _ abciTypes.RequestEndBlock) []abciTypes.ValidatorUpdate {
	return []abciTypes.ValidatorUpdate{}
}
