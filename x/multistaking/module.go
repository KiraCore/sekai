package multistaking

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"

	govkeeper "github.com/KiraCore/sekai/x/gov/keeper"

	"github.com/KiraCore/sekai/x/multistaking/keeper"
	"github.com/KiraCore/sekai/x/multistaking/types"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"

	"github.com/KiraCore/sekai/x/multistaking/client/cli"
	stakingtypes "github.com/KiraCore/sekai/x/multistaking/types"
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
	return stakingtypes.ModuleName
}

func (b AppModuleBasic) RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	stakingtypes.RegisterInterfaces(registry)
}

func (b AppModuleBasic) DefaultGenesis(marshaler codec.JSONCodec) json.RawMessage {
	return nil
}

func (b AppModuleBasic) ValidateGenesis(marshaler codec.JSONCodec, config client.TxEncodingConfig, message json.RawMessage) error {
	return nil
}

func (b AppModuleBasic) RegisterGRPCRoutes(clientCtx client.Context, serveMux *runtime.ServeMux) {
	stakingtypes.RegisterQueryHandlerClient(context.Background(), serveMux, types.NewQueryClient(clientCtx))
}

func (b AppModuleBasic) RegisterRESTRoutes(context client.Context, router *mux.Router) {
}

func (b AppModuleBasic) GetTxCmd() *cobra.Command {
	return cli.NewTxCmd()
}

func (b AppModuleBasic) RegisterLegacyAminoCodec(amino *codec.LegacyAmino) {
	stakingtypes.RegisterCodec(amino)
}

func (b AppModuleBasic) GetQueryCmd() *cobra.Command {
	return cli.GetQueryCmd()
}

// AppModule extends the cosmos SDK staking.
type AppModule struct {
	AppModuleBasic
	keeper              keeper.Keeper
	bankKeeper          types.BankKeeper
	customStakingKeeper types.StakingKeeper
	customGovKeeper     govkeeper.Keeper
}

// RegisterQueryService registers a GRPC query service to respond to the
// module-specific GRPC queries.
func (am AppModule) RegisterServices(cfg module.Configurator) {
	stakingtypes.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServerImpl(am.keeper, am.bankKeeper, am.customGovKeeper, am.customStakingKeeper))
	stakingtypes.RegisterQueryServer(cfg.QueryServer(), keeper.NewQuerier(am.keeper))
}

func (am AppModule) RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	stakingtypes.RegisterInterfaces(registry)
}

func (am AppModule) InitGenesis(
	ctx sdk.Context,
	cdc codec.JSONCodec,
	// keeper keeper.Keeper,
	data json.RawMessage,
) []abci.ValidatorUpdate {
	var gs types.GenesisState
	err := cdc.UnmarshalJSON(data, &gs)
	if err != nil {
		panic(fmt.Sprintf("failed to unmarshal %s genesis state: %s", types.ModuleName, err))
	}

	InitGenesis(ctx, am.keeper, &gs)
	return []abci.ValidatorUpdate{}
}

func (am AppModule) ExportGenesis(context sdk.Context, marshaler codec.JSONCodec) json.RawMessage {
	gs := ExportGenesis(context, am.keeper)
	return marshaler.MustMarshalJSON(gs)
}

// ConsensusVersion implements AppModule/ConsensusVersion.
func (AppModule) ConsensusVersion() uint64 { return 1 }

func (am AppModule) RegisterInvariants(registry sdk.InvariantRegistry) {}

func (am AppModule) QuerierRoute() string {
	return types.QuerierRoute
}

func (am AppModule) BeginBlock(context sdk.Context, block abci.RequestBeginBlock) {}

func (am AppModule) EndBlock(ctx sdk.Context, block abci.RequestEndBlock) []abci.ValidatorUpdate {
	return EndBlocker(ctx, am.keeper)
}

func (am AppModule) Name() string {
	return stakingtypes.ModuleName
}

// NewAppModule returns a new multistaking module.
func NewAppModule(
	keeper keeper.Keeper,
	bankKeeper types.BankKeeper,
	govKeeper govkeeper.Keeper,
	sk types.StakingKeeper,
) AppModule {
	return AppModule{
		keeper:              keeper,
		bankKeeper:          bankKeeper,
		customGovKeeper:     govKeeper,
		customStakingKeeper: sk,
	}
}
