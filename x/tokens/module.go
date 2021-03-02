package tokens

import (
	"context"
	"encoding/json"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"

	"github.com/KiraCore/sekai/middleware"
	tokenscli "github.com/KiraCore/sekai/x/tokens/client/cli"
	tokenskeeper "github.com/KiraCore/sekai/x/tokens/keeper"
	"github.com/KiraCore/sekai/x/tokens/types"
	tokenstypes "github.com/KiraCore/sekai/x/tokens/types"

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
	return tokenstypes.ModuleName
}

func (b AppModuleBasic) RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	tokenstypes.RegisterInterfaces(registry)
}

func (b AppModuleBasic) DefaultGenesis(cdc codec.JSONMarshaler) json.RawMessage {
	return cdc.MustMarshalJSON(tokenstypes.DefaultGenesis())
}

func (b AppModuleBasic) ValidateGenesis(marshaler codec.JSONMarshaler, config client.TxEncodingConfig, message json.RawMessage) error {
	return nil
}

func (b AppModuleBasic) RegisterRESTRoutes(clientCtx client.Context, router *mux.Router) {
}

func (b AppModuleBasic) RegisterGRPCRoutes(clientCtx client.Context, serveMux *runtime.ServeMux) {
	tokenstypes.RegisterQueryHandlerClient(context.Background(), serveMux, types.NewQueryClient(clientCtx))
}

func (b AppModuleBasic) RegisterLegacyAminoCodec(amino *codec.LegacyAmino) {
	tokenstypes.RegisterCodec(amino)
}

func (b AppModuleBasic) GetTxCmd() *cobra.Command {
	return tokenscli.NewTxCmd()
}

// GetQueryCmd implement query commands for this module
func (b AppModuleBasic) GetQueryCmd() *cobra.Command {
	queryCmd := &cobra.Command{
		Use:   tokenstypes.RouterKey,
		Short: "query commands for the customgov module",
	}
	queryCmd.AddCommand(
		tokenscli.GetCmdQueryTokenAlias(),
		tokenscli.GetCmdQueryAllTokenAliases(),
		tokenscli.GetCmdQueryTokenAliasesByDenom(),
		tokenscli.GetCmdQueryTokenRate(),
		tokenscli.GetCmdQueryAllTokenRates(),
		tokenscli.GetCmdQueryTokenRatesByDenom(),
	)

	queryCmd.PersistentFlags().String("node", "tcp://localhost:26657", "<host>:<port> to Tendermint RPC interface for this chain")
	return queryCmd
}

// AppModule for tokens management
type AppModule struct {
	AppModuleBasic
	tokensKeeper    tokenskeeper.Keeper
	customGovKeeper tokenstypes.CustomGovKeeper
}

// RegisterQueryService registers a GRPC query service to respond to the
// module-specific GRPC queries.
func (am AppModule) RegisterServices(cfg module.Configurator) {
	tokenstypes.RegisterMsgServer(cfg.MsgServer(), tokenskeeper.NewMsgServerImpl(am.tokensKeeper, am.customGovKeeper))
	querier := tokenskeeper.NewQuerier(am.tokensKeeper)
	tokenstypes.RegisterQueryServer(cfg.QueryServer(), querier)
}

func (am AppModule) RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	tokenstypes.RegisterInterfaces(registry)
}

func (am AppModule) InitGenesis(
	ctx sdk.Context,
	cdc codec.JSONMarshaler,
	data json.RawMessage,
) []abci.ValidatorUpdate {
	var genesisState tokenstypes.GenesisState
	cdc.MustUnmarshalJSON(data, &genesisState)

	for _, alias := range genesisState.Aliases {
		am.tokensKeeper.UpsertTokenAlias(ctx, *alias)
	}

	for _, rate := range genesisState.Rates {
		am.tokensKeeper.UpsertTokenRate(ctx, *rate)
	}

	return nil
}

func (am AppModule) ExportGenesis(clientCtx sdk.Context, marshaler codec.JSONMarshaler) json.RawMessage {
	return nil
}

func (am AppModule) RegisterInvariants(registry sdk.InvariantRegistry) {}

func (am AppModule) QuerierRoute() string {
	return tokenstypes.QuerierRoute
}

// LegacyQuerierHandler returns the staking module sdk.Querier.
func (am AppModule) LegacyQuerierHandler(legacyQuerierCdc *codec.LegacyAmino) sdk.Querier {
	return nil
}

func (am AppModule) BeginBlock(clientCtx sdk.Context, block abci.RequestBeginBlock) {}

func (am AppModule) EndBlock(ctx sdk.Context, block abci.RequestEndBlock) []abci.ValidatorUpdate {
	return nil
}

func (am AppModule) Name() string {
	return tokenstypes.ModuleName
}

// Route returns the message routing key for the staking module.
func (am AppModule) Route() sdk.Route {
	return middleware.NewRoute(tokenstypes.ModuleName, NewHandler(am.tokensKeeper, am.customGovKeeper))
}

// NewAppModule returns a new Custom Staking module.
func NewAppModule(
	keeper tokenskeeper.Keeper,
	customGovKeeper tokenstypes.CustomGovKeeper,
) AppModule {
	return AppModule{
		tokensKeeper:    keeper,
		customGovKeeper: customGovKeeper,
	}
}
