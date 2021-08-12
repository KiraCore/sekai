package staking

import (
	"context"
	"encoding/json"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"

	govkeeper "github.com/KiraCore/sekai/x/gov/keeper"

	"github.com/KiraCore/sekai/middleware"
	"github.com/KiraCore/sekai/x/staking/keeper"
	"github.com/KiraCore/sekai/x/staking/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/KiraCore/sekai/x/staking/client/cli"
	stakingtypes "github.com/KiraCore/sekai/x/staking/types"
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
	proposalCmd := &cobra.Command{
		Use:                        "proposal",
		Short:                      "Proposal subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	proposalCmd.AddCommand(cli.GetTxProposalUnjailValidatorCmd())

	txCommand := &cobra.Command{
		Use:                        "customstaking",
		Short:                      "staking module subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCommand.AddCommand(cli.GetTxClaimValidatorCmd(), proposalCmd)

	return txCommand
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
	customStakingKeeper keeper.Keeper
	customGovKeeper     govkeeper.Keeper
}

// RegisterQueryService registers a GRPC query service to respond to the
// module-specific GRPC queries.
func (am AppModule) RegisterServices(cfg module.Configurator) {
	stakingtypes.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServerImpl(am.customStakingKeeper, am.customGovKeeper))
	stakingtypes.RegisterQueryServer(cfg.QueryServer(), keeper.NewQuerier(am.customStakingKeeper))
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
	var genesisState types.GenesisState
	cdc.MustUnmarshalJSON(data, &genesisState)

	valUpdate := make([]abci.ValidatorUpdate, len(genesisState.Validators))

	for i, val := range genesisState.Validators {
		am.customStakingKeeper.AddValidator(ctx, val)
		am.customStakingKeeper.AfterValidatorJoined(ctx, val.GetConsAddr(), val.ValKey)

		consPk, err := val.TmConsPubKey()
		if err != nil {
			panic(err)
		}

		valUpdate[i] = abci.ValidatorUpdate{
			Power:  1,
			PubKey: consPk,
		}
	}

	return valUpdate
}

func (am AppModule) ExportGenesis(context sdk.Context, marshaler codec.JSONCodec) json.RawMessage {
	gs := ExportGenesis(context, am.customStakingKeeper)
	return marshaler.MustMarshalJSON(gs)
}

// ConsensusVersion implements AppModule/ConsensusVersion.
func (AppModule) ConsensusVersion() uint64 { return 1 }

func (am AppModule) RegisterInvariants(registry sdk.InvariantRegistry) {}

func (am AppModule) QuerierRoute() string {
	return types.QuerierRoute
}

// LegacyQuerierHandler returns the staking module sdk.Querier.
func (am AppModule) LegacyQuerierHandler(legacyQuerierCdc *codec.LegacyAmino) sdk.Querier {
	return nil
}

func (am AppModule) BeginBlock(context sdk.Context, block abci.RequestBeginBlock) {}

func (am AppModule) EndBlock(ctx sdk.Context, block abci.RequestEndBlock) []abci.ValidatorUpdate {
	return EndBlocker(ctx, am.customStakingKeeper)
}

func (am AppModule) Name() string {
	return stakingtypes.ModuleName
}

// Route returns the message routing key for the staking module.
func (am AppModule) Route() sdk.Route {
	return middleware.NewRoute(stakingtypes.ModuleName, NewHandler(am.customStakingKeeper, am.customGovKeeper))
}

// NewAppModule returns a new Custom Staking module.
func NewAppModule(
	keeper keeper.Keeper,
	govKeeper govkeeper.Keeper,
) AppModule {
	return AppModule{
		customStakingKeeper: keeper,
		customGovKeeper:     govKeeper,
	}
}
