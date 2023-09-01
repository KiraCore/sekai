package custody

import (
	"context"
	"encoding/json"

	custodycli "github.com/KiraCore/sekai/x/custody/client/cli"
	custodykeeper "github.com/KiraCore/sekai/x/custody/keeper"
	"github.com/KiraCore/sekai/x/custody/types"
	custodytypes "github.com/KiraCore/sekai/x/custody/types"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"
)

var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

type AppModuleBasic struct{}

func (AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *runtime.ServeMux) {}

func (b AppModuleBasic) Name() string {
	return custodytypes.ModuleName
}

func (b AppModuleBasic) RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	custodytypes.RegisterInterfaces(registry)
}

func (b AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	return json.RawMessage("{}")
}

func (b AppModuleBasic) ValidateGenesis(marshaler codec.JSONCodec, config client.TxEncodingConfig, message json.RawMessage) error {
	return nil
}

func (b AppModuleBasic) RegisterRESTRoutes(clientCtx client.Context, router *mux.Router) {}

func (b AppModuleBasic) RegisterGRPCRoutes(clientCtx client.Context, serveMux *runtime.ServeMux) {
	_ = custodytypes.RegisterQueryHandlerClient(context.Background(), serveMux, types.NewQueryClient(clientCtx))
}

func (b AppModuleBasic) RegisterLegacyAminoCodec(amino *codec.LegacyAmino) {
	custodytypes.RegisterCodec(amino)
}

func (b AppModuleBasic) GetTxCmd() *cobra.Command {
	return custodycli.NewTxCmd()
}

func (b AppModuleBasic) GetQueryCmd() *cobra.Command {
	return custodycli.NewQueryCmd()
}

type AppModule struct {
	AppModuleBasic
	custodyKeeper   custodykeeper.Keeper
	customGovKeeper custodytypes.CustomGovKeeper
	bankKeeper      custodytypes.BankKeeper
}

func (am AppModule) RegisterServices(cfg module.Configurator) {
	custodytypes.RegisterMsgServer(cfg.MsgServer(), custodykeeper.NewMsgServerImpl(am.custodyKeeper, am.customGovKeeper, am.bankKeeper))
	querier := custodykeeper.NewQuerier(am.custodyKeeper)
	custodytypes.RegisterQueryServer(cfg.QueryServer(), querier)
}

func (am AppModule) RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	custodytypes.RegisterInterfaces(registry)
}

func (am AppModule) InitGenesis(
	ctx sdk.Context,
	cdc codec.JSONCodec,
	data json.RawMessage,
) []abci.ValidatorUpdate {
	//var genesisState custodytypes.GenesisState
	//cdc.MustUnmarshalJSON(data, &genesisState)
	//
	//am.custodyKeeper.SetMaxCustodyBufferSize(ctx, genesisState.MaxCustodyBufferSize)
	//am.custodyKeeper.SetMaxCustodyTxSize(ctx, genesisState.MaxCustodyTxSize)

	return nil
}

func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	//var genesisState custodytypes.GenesisState
	//genesisState.MaxCustodyBufferSize = am.custodyKeeper.GetMaxCustodyBufferSize(ctx)
	//genesisState.MaxCustodyTxSize = am.custodyKeeper.GetMaxCustodyTxSize(ctx)
	return nil
}

func (AppModule) ConsensusVersion() uint64 { return 1 }

func (am AppModule) RegisterInvariants(registry sdk.InvariantRegistry) {}

func (am AppModule) QuerierRoute() string {
	return custodytypes.QuerierRoute
}

func (am AppModule) BeginBlock(clientCtx sdk.Context, block abci.RequestBeginBlock) {}

func (am AppModule) EndBlock(ctx sdk.Context, block abci.RequestEndBlock) []abci.ValidatorUpdate {
	return nil
}

func (am AppModule) Name() string {
	return custodytypes.ModuleName
}

func (am AppModule) CheckTx() string {
	return custodytypes.ModuleName
}

func NewAppModule(keeper custodykeeper.Keeper, customGovKeeper custodytypes.CustomGovKeeper, bankKeeper custodytypes.BankKeeper) AppModule {
	return AppModule{
		custodyKeeper:   keeper,
		customGovKeeper: customGovKeeper,
		bankKeeper:      bankKeeper,
	}
}
