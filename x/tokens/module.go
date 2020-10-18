package tokens

import (
	"encoding/json"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"

	cli2 "github.com/KiraCore/sekai/x/tokens/client/cli"
	keeper2 "github.com/KiraCore/sekai/x/tokens/keeper"
	tokenstypes "github.com/KiraCore/sekai/x/tokens/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	types2 "github.com/cosmos/cosmos-sdk/codec/types"
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

func (b AppModuleBasic) Name() string {
	return tokenstypes.ModuleName
}

func (b AppModuleBasic) RegisterInterfaces(registry types2.InterfaceRegistry) {
	tokenstypes.RegisterInterfaces(registry)
}

func (b AppModuleBasic) DefaultGenesis(cdc codec.JSONMarshaler) json.RawMessage {
	return cdc.MustMarshalJSON(tokenstypes.DefaultGenesis())
}

func (b AppModuleBasic) ValidateGenesis(marshaler codec.JSONMarshaler, config client.TxEncodingConfig, message json.RawMessage) error {
	return nil
}

func (b AppModuleBasic) RegisterRESTRoutes(context client.Context, router *mux.Router) {
}

func (b AppModuleBasic) RegisterGRPCRoutes(context client.Context, serveMux *runtime.ServeMux) {
}

func (b AppModuleBasic) RegisterLegacyAminoCodec(amino *codec.LegacyAmino) {
	tokenstypes.RegisterCodec(amino)
}

func (b AppModuleBasic) GetTxCmd() *cobra.Command {
	return cli2.NewTxCmd()
}

// GetQueryCmd implement query commands for this module
func (b AppModuleBasic) GetQueryCmd() *cobra.Command {
	queryCmd := &cobra.Command{
		Use:   tokenstypes.RouterKey,
		Short: "query commands for the customgov module",
	}
	queryCmd.AddCommand(
		cli2.GetCmdQueryTokenAlias(),
		cli2.GetCmdQueryAllTokenAliases(),
		cli2.GetCmdQueryTokenAliasesByDenom(),
		cli2.GetCmdQueryTokenRate(),
		cli2.GetCmdQueryAllTokenRates(),
		cli2.GetCmdQueryTokenRatesByDenom(),
	)

	queryCmd.PersistentFlags().String("node", "tcp://localhost:26657", "<host>:<port> to Tendermint RPC interface for this chain")
	return queryCmd
}

// AppModule for tokens management
type AppModule struct {
	AppModuleBasic
	tokensKeeper    keeper2.Keeper
	customGovKeeper tokenstypes.CustomGovKeeper
}

func (am AppModule) RegisterInterfaces(registry types2.InterfaceRegistry) {
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

func (am AppModule) ExportGenesis(context sdk.Context, marshaler codec.JSONMarshaler) json.RawMessage {
	return nil
}

func (am AppModule) RegisterInvariants(registry sdk.InvariantRegistry) {}

func (am AppModule) QuerierRoute() string { return "" }

func (am AppModule) LegacyQuerierHandler(marshaler *codec.LegacyAmino) sdk.Querier {
	return nil
}

func (am AppModule) BeginBlock(context sdk.Context, block abci.RequestBeginBlock) {}

func (am AppModule) EndBlock(ctx sdk.Context, block abci.RequestEndBlock) []abci.ValidatorUpdate {
	return nil
}

func (am AppModule) Name() string {
	return tokenstypes.ModuleName
}

// Route returns the message routing key for the staking module.
func (am AppModule) Route() sdk.Route {
	return sdk.NewRoute(tokenstypes.ModuleName, NewHandler(am.tokensKeeper, am.customGovKeeper))
}

// RegisterQueryService registers a GRPC query service to respond to the
// module-specific GRPC queries.
func (am AppModule) RegisterQueryService(server grpc.Server) {
	querier := NewQuerier(am.tokensKeeper)
	tokenstypes.RegisterQueryServer(server, querier)
}

// NewAppModule returns a new Custom Staking module.
func NewAppModule(
	keeper keeper2.Keeper,
	customGovKeeper tokenstypes.CustomGovKeeper,
) AppModule {
	return AppModule{
		tokensKeeper:    keeper,
		customGovKeeper: customGovKeeper,
	}
}
