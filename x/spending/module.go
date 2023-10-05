package spending

import (
	"context"
	"encoding/json"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"

	spendingcli "github.com/KiraCore/sekai/x/spending/client/cli"
	spendingkeeper "github.com/KiraCore/sekai/x/spending/keeper"
	"github.com/KiraCore/sekai/x/spending/types"
	spendingtypes "github.com/KiraCore/sekai/x/spending/types"

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
	return spendingtypes.ModuleName
}

func (b AppModuleBasic) RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	spendingtypes.RegisterInterfaces(registry)
}

func (b AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	return cdc.MustMarshalJSON(spendingtypes.DefaultGenesis())
}

func (b AppModuleBasic) ValidateGenesis(marshaler codec.JSONCodec, config client.TxEncodingConfig, message json.RawMessage) error {
	return nil
}

func (b AppModuleBasic) RegisterRESTRoutes(clientCtx client.Context, router *mux.Router) {
}

func (b AppModuleBasic) RegisterGRPCRoutes(clientCtx client.Context, serveMux *runtime.ServeMux) {
	spendingtypes.RegisterQueryHandlerClient(context.Background(), serveMux, types.NewQueryClient(clientCtx))
}

func (b AppModuleBasic) RegisterLegacyAminoCodec(amino *codec.LegacyAmino) {
	spendingtypes.RegisterCodec(amino)
}

func (b AppModuleBasic) GetTxCmd() *cobra.Command {
	return spendingcli.NewTxCmd()
}

// GetQueryCmd implement query commands for this module
func (b AppModuleBasic) GetQueryCmd() *cobra.Command {
	return spendingcli.NewQueryCmd()
}

// AppModule for tokens management
type AppModule struct {
	AppModuleBasic
	spendingKeeper  spendingkeeper.Keeper
	customGovKeeper spendingtypes.CustomGovKeeper
	bankKeeper      spendingtypes.BankKeeper
}

// RegisterQueryService registers a GRPC query service to respond to the
// module-specific GRPC queries.
func (am AppModule) RegisterServices(cfg module.Configurator) {
	spendingtypes.RegisterMsgServer(cfg.MsgServer(), spendingkeeper.NewMsgServerImpl(am.spendingKeeper, am.customGovKeeper, am.bankKeeper))
	querier := spendingkeeper.NewQuerier(am.spendingKeeper, am.customGovKeeper)
	spendingtypes.RegisterQueryServer(cfg.QueryServer(), querier)
}

func (am AppModule) RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	spendingtypes.RegisterInterfaces(registry)
}

func (am AppModule) InitGenesis(
	ctx sdk.Context,
	cdc codec.JSONCodec,
	data json.RawMessage,
) []abci.ValidatorUpdate {
	var genesisState spendingtypes.GenesisState
	cdc.MustUnmarshalJSON(data, &genesisState)

	for _, pool := range genesisState.Pools {
		am.spendingKeeper.SetSpendingPool(ctx, pool)
	}

	for _, claimInfo := range genesisState.Claims {
		am.spendingKeeper.SetClaimInfo(ctx, claimInfo)
	}
	return nil
}

func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	var genesisState spendingtypes.GenesisState
	genesisState.Pools = am.spendingKeeper.GetAllSpendingPools(ctx)
	genesisState.Claims = am.spendingKeeper.GetAllClaimInfos(ctx)
	return cdc.MustMarshalJSON(&genesisState)
}

// ConsensusVersion implements AppModule/ConsensusVersion.
func (AppModule) ConsensusVersion() uint64 { return 1 }

func (am AppModule) RegisterInvariants(registry sdk.InvariantRegistry) {}

func (am AppModule) QuerierRoute() string {
	return spendingtypes.QuerierRoute
}

func (am AppModule) BeginBlock(ctx sdk.Context, block abci.RequestBeginBlock) {
	am.spendingKeeper.BeginBlocker(ctx)
}

func (am AppModule) EndBlock(ctx sdk.Context, block abci.RequestEndBlock) []abci.ValidatorUpdate {
	am.spendingKeeper.EndBlocker(ctx)
	return nil
}

func (am AppModule) Name() string {
	return spendingtypes.ModuleName
}

// NewAppModule returns a new Custom Staking module.
func NewAppModule(
	keeper spendingkeeper.Keeper,
	customGovKeeper spendingtypes.CustomGovKeeper,
	bankKeeper spendingtypes.BankKeeper,
) AppModule {
	return AppModule{
		spendingKeeper:  keeper,
		customGovKeeper: customGovKeeper,
		bankKeeper:      bankKeeper,
	}
}
