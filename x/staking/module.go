package staking

import (
	"encoding/json"

	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/KiraCore/cosmos-sdk/client"
	"github.com/KiraCore/cosmos-sdk/codec"
	types2 "github.com/KiraCore/cosmos-sdk/codec/types"
	sdk "github.com/KiraCore/cosmos-sdk/types"
	"github.com/KiraCore/cosmos-sdk/types/module"
	"github.com/KiraCore/sekai/x/staking/keeper"
	"github.com/KiraCore/sekai/x/staking/types"
	"github.com/gogo/protobuf/grpc"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/KiraCore/sekai/x/staking/client/cli"
	cumstomtypes "github.com/KiraCore/sekai/x/staking/types"
)

var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

type AppModuleBasic struct{}

func (b AppModuleBasic) Name() string {
	return cumstomtypes.ModuleName
}

func (b AppModuleBasic) RegisterInterfaces(registry types2.InterfaceRegistry) {
	cumstomtypes.RegisterInterfaces(registry)
}

func (b AppModuleBasic) DefaultGenesis(marshaler codec.JSONMarshaler) json.RawMessage {
	return nil
}

func (b AppModuleBasic) ValidateGenesis(marshaler codec.JSONMarshaler, config client.TxEncodingConfig, message json.RawMessage) error {
	return nil
}

func (b AppModuleBasic) RegisterRESTRoutes(context client.Context, router *mux.Router) {
}

func (b AppModuleBasic) GetTxCmd() *cobra.Command {
	return cli.GetTxClaimValidatorCmd()
}

func (b AppModuleBasic) GetQueryCmd() *cobra.Command {
	return cli.GetCmdQueryValidatorByAddress()
}

// RegisterCodec registers the staking module's types for the given codec.
func (AppModuleBasic) RegisterCodec(cdc *codec.Codec) {
	cumstomtypes.RegisterCodec(cdc)
}

// AppModule extends the cosmos SDK staking.
type AppModule struct {
	AppModuleBasic
	customStakingKeeper keeper.Keeper
}

func (am AppModule) RegisterCodec(c *codec.Codec) {
	panic("implement me")
}

func (am AppModule) RegisterInterfaces(registry types2.InterfaceRegistry) {
	panic("implement me")
}

func (am AppModule) InitGenesis(
	ctx sdk.Context,
	cdc codec.JSONMarshaler,
	data json.RawMessage,
) []abci.ValidatorUpdate {
	var genesisState types.GenesisState
	cdc.MustUnmarshalJSON(data, &genesisState)

	valUpdate := make([]abci.ValidatorUpdate, len(genesisState.Validators))

	for i, val := range genesisState.Validators {
		am.customStakingKeeper.AddValidator(ctx, val)
		valUpdate[i] = abci.ValidatorUpdate{
			Power:  1,
			PubKey: tmtypes.TM2PB.PubKey(val.GetConsPubKey()),
		}
	}

	return valUpdate
}

func (am AppModule) ExportGenesis(context sdk.Context, marshaler codec.JSONMarshaler) json.RawMessage {
	return nil
}

func (am AppModule) RegisterInvariants(registry sdk.InvariantRegistry) {}

func (am AppModule) QuerierRoute() string { return "" }

func (am AppModule) LegacyQuerierHandler(marshaler codec.JSONMarshaler) sdk.Querier {
	return nil
}

func (am AppModule) BeginBlock(context sdk.Context, block abci.RequestBeginBlock) {}

func (am AppModule) EndBlock(ctx sdk.Context, block abci.RequestEndBlock) []abci.ValidatorUpdate {
	valSet := am.customStakingKeeper.GetValidatorSet(ctx)

	valUpdate := make([]abci.ValidatorUpdate, len(valSet))

	for i, val := range valSet {
		am.customStakingKeeper.AddValidator(ctx, val)
		valUpdate[i] = abci.ValidatorUpdate{
			Power:  1,
			PubKey: tmtypes.TM2PB.PubKey(val.GetConsPubKey()),
		}
	}

	return valUpdate
}

func (am AppModule) Name() string {
	return cumstomtypes.ModuleName
}

// Route returns the message routing key for the staking module.
func (am AppModule) Route() sdk.Route {
	return sdk.NewRoute(cumstomtypes.ModuleName, NewHandler(am.customStakingKeeper))
}

// RegisterQueryService registers a GRPC query service to respond to the
// module-specific GRPC queries.
func (am AppModule) RegisterQueryService(server grpc.Server) {
	querier := NewQuerier(am.customStakingKeeper)
	cumstomtypes.RegisterQueryServer(server, querier)
}

// NewAppModule returns a new Custom Staking module.
func NewAppModule(
	keeper keeper.Keeper,
) AppModule {
	return AppModule{
		customStakingKeeper: keeper,
	}
}
