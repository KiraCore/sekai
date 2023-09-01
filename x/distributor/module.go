package distributor

import (
	"context"
	"encoding/json"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"

	distributorcli "github.com/KiraCore/sekai/x/distributor/client/cli"
	distributorkeeper "github.com/KiraCore/sekai/x/distributor/keeper"
	"github.com/KiraCore/sekai/x/distributor/types"
	distributortypes "github.com/KiraCore/sekai/x/distributor/types"

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
	return distributortypes.ModuleName
}

func (b AppModuleBasic) RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	distributortypes.RegisterInterfaces(registry)
}

func (b AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	return cdc.MustMarshalJSON(distributortypes.DefaultGenesis())
}

func (b AppModuleBasic) ValidateGenesis(marshaler codec.JSONCodec, config client.TxEncodingConfig, message json.RawMessage) error {
	return nil
}

func (b AppModuleBasic) RegisterRESTRoutes(clientCtx client.Context, router *mux.Router) {
}

func (b AppModuleBasic) RegisterGRPCRoutes(clientCtx client.Context, serveMux *runtime.ServeMux) {
	distributortypes.RegisterQueryHandlerClient(context.Background(), serveMux, types.NewQueryClient(clientCtx))
}

func (b AppModuleBasic) RegisterLegacyAminoCodec(amino *codec.LegacyAmino) {
	distributortypes.RegisterCodec(amino)
}

func (b AppModuleBasic) GetTxCmd() *cobra.Command {
	return distributorcli.NewTxCmd()
}

// GetQueryCmd implement query commands for this module
func (b AppModuleBasic) GetQueryCmd() *cobra.Command {
	return distributorcli.NewQueryCmd()
}

// AppModule for distributor management
type AppModule struct {
	AppModuleBasic
	distributorKeeper distributorkeeper.Keeper
	customGovKeeper   distributortypes.CustomGovKeeper
}

// RegisterQueryService registers a GRPC query service to respond to the
// module-specific GRPC queries.
func (am AppModule) RegisterServices(cfg module.Configurator) {
	distributortypes.RegisterMsgServer(cfg.MsgServer(), distributorkeeper.NewMsgServerImpl(am.distributorKeeper, am.customGovKeeper))
	querier := distributorkeeper.NewQuerier(am.distributorKeeper)
	distributortypes.RegisterQueryServer(cfg.QueryServer(), querier)
}

func (am AppModule) RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	distributortypes.RegisterInterfaces(registry)
}

func (am AppModule) InitGenesis(
	ctx sdk.Context,
	cdc codec.JSONCodec,
	data json.RawMessage,
) []abci.ValidatorUpdate {
	var genesisState distributortypes.GenesisState
	cdc.MustUnmarshalJSON(data, &genesisState)

	am.distributorKeeper.SetFeesTreasury(ctx, genesisState.FeesTreasury)
	am.distributorKeeper.SetSnapPeriod(ctx, genesisState.SnapPeriod)
	am.distributorKeeper.SetYearStartSnapshot(ctx, genesisState.YearStartSnapshot)
	am.distributorKeeper.SetPeriodicSnapshot(ctx, genesisState.PeriodicSnapshot)

	for _, vote := range genesisState.ValidatorVotes {
		consAddr, err := sdk.ConsAddressFromBech32(vote.ConsAddr)
		if err != nil {
			panic(err)
		}
		am.distributorKeeper.SetValidatorVote(ctx, consAddr, vote.Height)
	}

	if genesisState.PreviousProposer != "" {
		consAddr, err := sdk.ConsAddressFromBech32(genesisState.PreviousProposer)
		if err != nil {
			panic(err)
		}
		am.distributorKeeper.SetPreviousProposerConsAddr(ctx, consAddr)
	}

	return nil
}

func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	var genesisState distributortypes.GenesisState
	genesisState.FeesTreasury = am.distributorKeeper.GetFeesTreasury(ctx)
	genesisState.SnapPeriod = am.distributorKeeper.GetSnapPeriod(ctx)
	genesisState.ValidatorVotes = am.distributorKeeper.GetAllValidatorVotes(ctx)
	genesisState.PreviousProposer = am.distributorKeeper.GetPreviousProposerConsAddr(ctx).String()
	genesisState.YearStartSnapshot = am.distributorKeeper.GetYearStartSnapshot(ctx)
	genesisState.PeriodicSnapshot = am.distributorKeeper.GetPeriodicSnapshot(ctx)
	return cdc.MustMarshalJSON(&genesisState)
}

// ConsensusVersion implements AppModule/ConsensusVersion.
func (AppModule) ConsensusVersion() uint64 { return 1 }

func (am AppModule) RegisterInvariants(registry sdk.InvariantRegistry) {}

func (am AppModule) QuerierRoute() string {
	return distributortypes.QuerierRoute
}

func (am AppModule) BeginBlock(clientCtx sdk.Context, req abci.RequestBeginBlock) {
	am.distributorKeeper.BeginBlocker(clientCtx, req)
}

func (am AppModule) EndBlock(ctx sdk.Context, block abci.RequestEndBlock) []abci.ValidatorUpdate {
	am.distributorKeeper.EndBlocker(ctx)
	return nil
}

func (am AppModule) Name() string {
	return distributortypes.ModuleName
}

// NewAppModule returns a new Custom Staking module.
func NewAppModule(
	keeper distributorkeeper.Keeper,
	customGovKeeper distributortypes.CustomGovKeeper,
) AppModule {
	return AppModule{
		distributorKeeper: keeper,
		customGovKeeper:   customGovKeeper,
	}
}
