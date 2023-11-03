package basket

import (
	"context"
	"encoding/json"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"

	basketcli "github.com/KiraCore/sekai/x/basket/client/cli"
	basketkeeper "github.com/KiraCore/sekai/x/basket/keeper"
	"github.com/KiraCore/sekai/x/basket/types"
	baskettypes "github.com/KiraCore/sekai/x/basket/types"

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
	return baskettypes.ModuleName
}

func (b AppModuleBasic) RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	baskettypes.RegisterInterfaces(registry)
}

func (b AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	return cdc.MustMarshalJSON(baskettypes.DefaultGenesis())
}

func (b AppModuleBasic) ValidateGenesis(marshaler codec.JSONCodec, config client.TxEncodingConfig, message json.RawMessage) error {
	return nil
}

func (b AppModuleBasic) RegisterRESTRoutes(clientCtx client.Context, router *mux.Router) {
}

func (b AppModuleBasic) RegisterGRPCRoutes(clientCtx client.Context, serveMux *runtime.ServeMux) {
	baskettypes.RegisterQueryHandlerClient(context.Background(), serveMux, types.NewQueryClient(clientCtx))
}

func (b AppModuleBasic) RegisterLegacyAminoCodec(amino *codec.LegacyAmino) {
	baskettypes.RegisterCodec(amino)
}

func (b AppModuleBasic) GetTxCmd() *cobra.Command {
	return basketcli.NewTxCmd()
}

// GetQueryCmd implement query commands for this module
func (b AppModuleBasic) GetQueryCmd() *cobra.Command {
	return basketcli.NewQueryCmd()
}

// AppModule for basket management
type AppModule struct {
	AppModuleBasic
	basketKeeper    basketkeeper.Keeper
	customGovKeeper baskettypes.CustomGovKeeper
}

// RegisterQueryService registers a GRPC query service to respond to the
// module-specific GRPC queries.
func (am AppModule) RegisterServices(cfg module.Configurator) {
	baskettypes.RegisterMsgServer(cfg.MsgServer(), basketkeeper.NewMsgServerImpl(am.basketKeeper, am.customGovKeeper))
	querier := basketkeeper.NewQuerier(am.basketKeeper)
	baskettypes.RegisterQueryServer(cfg.QueryServer(), querier)
}

func (am AppModule) RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	baskettypes.RegisterInterfaces(registry)
}

func (am AppModule) InitGenesis(
	ctx sdk.Context,
	cdc codec.JSONCodec,
	data json.RawMessage,
) []abci.ValidatorUpdate {
	var genesisState baskettypes.GenesisState
	cdc.MustUnmarshalJSON(data, &genesisState)

	am.basketKeeper.SetLastBasketId(ctx, genesisState.LastBasketId)
	for _, basket := range genesisState.Baskets {
		am.basketKeeper.SetBasket(ctx, basket)
	}

	for _, amount := range genesisState.HistoricalMints {
		am.basketKeeper.SetMintAmount(ctx, time.Unix(int64(amount.Time), 0), amount.BasketId, amount.Amount)
	}

	for _, amount := range genesisState.HistoricalBurns {
		am.basketKeeper.SetBurnAmount(ctx, time.Unix(int64(amount.Time), 0), amount.BasketId, amount.Amount)
	}

	for _, amount := range genesisState.HistoricalSwaps {
		am.basketKeeper.SetSwapAmount(ctx, time.Unix(int64(amount.Time), 0), amount.BasketId, amount.Amount)
	}

	return nil
}

func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	genesisState := baskettypes.GenesisState{
		Baskets:         am.basketKeeper.GetAllBaskets(ctx),
		LastBasketId:    am.basketKeeper.GetLastBasketId(ctx),
		HistoricalMints: am.basketKeeper.GetAllMintAmounts(ctx),
		HistoricalBurns: am.basketKeeper.GetAllBurnAmounts(ctx),
		HistoricalSwaps: am.basketKeeper.GetAllSwapAmounts(ctx),
	}
	return cdc.MustMarshalJSON(&genesisState)
}

// ConsensusVersion implements AppModule/ConsensusVersion.
func (AppModule) ConsensusVersion() uint64 { return 1 }

func (am AppModule) RegisterInvariants(registry sdk.InvariantRegistry) {}

func (am AppModule) QuerierRoute() string {
	return baskettypes.QuerierRoute
}

func (am AppModule) BeginBlock(clientCtx sdk.Context, req abci.RequestBeginBlock) {
	BeginBlocker(clientCtx, req, am.basketKeeper)
}

func (am AppModule) EndBlock(ctx sdk.Context, block abci.RequestEndBlock) []abci.ValidatorUpdate {
	EndBlocker(ctx, am.basketKeeper)
	return nil
}

func (am AppModule) Name() string {
	return baskettypes.ModuleName
}

// NewAppModule returns a new Custom Staking module.
func NewAppModule(
	keeper basketkeeper.Keeper,
	customGovKeeper baskettypes.CustomGovKeeper,
) AppModule {
	return AppModule{
		basketKeeper:    keeper,
		customGovKeeper: customGovKeeper,
	}
}
