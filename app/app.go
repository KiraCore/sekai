package app

import (
	"io"
	"net/http"
	"os"
	"path/filepath"

	customante "github.com/KiraCore/sekai/app/ante"
	"github.com/KiraCore/sekai/middleware"
	"github.com/KiraCore/sekai/x/evidence"
	evidencekeeper "github.com/KiraCore/sekai/x/evidence/keeper"
	evidencetypes "github.com/KiraCore/sekai/x/evidence/types"
	"github.com/KiraCore/sekai/x/feeprocessing"
	feeprocessingkeeper "github.com/KiraCore/sekai/x/feeprocessing/keeper"
	feeprocessingtypes "github.com/KiraCore/sekai/x/feeprocessing/types"
	"github.com/KiraCore/sekai/x/genutil"
	genutiltypes "github.com/KiraCore/sekai/x/genutil/types"
	customgov "github.com/KiraCore/sekai/x/gov"
	customgovkeeper "github.com/KiraCore/sekai/x/gov/keeper"
	customgovtypes "github.com/KiraCore/sekai/x/gov/types"
	customslashing "github.com/KiraCore/sekai/x/slashing"
	customslashingkeeper "github.com/KiraCore/sekai/x/slashing/keeper"
	customslashingtypes "github.com/KiraCore/sekai/x/slashing/types"
	customstaking "github.com/KiraCore/sekai/x/staking"
	customstakingkeeper "github.com/KiraCore/sekai/x/staking/keeper"
	customstakingtypes "github.com/KiraCore/sekai/x/staking/types"
	"github.com/KiraCore/sekai/x/tokens"
	tokenskeeper "github.com/KiraCore/sekai/x/tokens/keeper"
	tokenstypes "github.com/KiraCore/sekai/x/tokens/types"
	"github.com/KiraCore/sekai/x/upgrade"
	upgradekeeper "github.com/KiraCore/sekai/x/upgrade/keeper"
	upgradetypes "github.com/KiraCore/sekai/x/upgrade/types"
	bam "github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/server/api"
	"github.com/cosmos/cosmos-sdk/server/config"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	simappparams "github.com/cosmos/cosmos-sdk/simapp/params"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	authrest "github.com/cosmos/cosmos-sdk/x/auth/client/rest"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	"github.com/cosmos/cosmos-sdk/x/auth/simulation"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/gorilla/mux"
	"github.com/rakyll/statik/fs"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	tmos "github.com/tendermint/tendermint/libs/os"
	dbm "github.com/tendermint/tm-db"
)

const appName = "Sekai"

var (
	// DefaultNodeHome default home directories for the application daemon
	DefaultNodeHome string

	// ModuleBasics defines the module BasicManager is in charge of setting up basic,
	// non-dependant module elements, such as codec registration
	// and genesis verification.
	ModuleBasics = module.NewBasicManager(
		auth.AppModuleBasic{},
		genutil.AppModuleBasic{},
		bank.AppModuleBasic{},
		params.AppModuleBasic{},
		upgrade.AppModuleBasic{},
		customslashing.AppModuleBasic{},
		customstaking.AppModuleBasic{},
		customgov.AppModuleBasic{},
		evidence.AppModuleBasic{},
		tokens.AppModuleBasic{},
		feeprocessing.AppModuleBasic{},
	)

	// module account permissions
	maccPerms = map[string][]string{
		authtypes.FeeCollectorName: nil,
	}

	// module accounts that are allowed to receive tokens
	allowedReceivingModAcc = map[string]bool{}
)

// NewApp extended ABCI application
type SekaiApp struct {
	*bam.BaseApp
	cdc               *codec.LegacyAmino
	appCodec          codec.Marshaler
	interfaceRegistry types.InterfaceRegistry

	invCheckPeriod uint

	// keys to access the substores
	keys    map[string]*sdk.KVStoreKey
	tKeys   map[string]*sdk.TransientStoreKey
	memKeys map[string]*sdk.MemoryStoreKey

	// keepers
	accountKeeper authkeeper.AccountKeeper
	bankKeeper    bankkeeper.Keeper
	upgradeKeeper upgradekeeper.Keeper
	paramsKeeper  paramskeeper.Keeper

	customSlashingKeeper customslashingkeeper.Keeper
	customStakingKeeper  customstakingkeeper.Keeper
	customGovKeeper      customgovkeeper.Keeper
	tokensKeeper         tokenskeeper.Keeper
	feeprocessingKeeper  feeprocessingkeeper.Keeper
	evidenceKeeper       evidencekeeper.Keeper

	// Module Manager
	mm *module.Manager

	// simulation manager
	sm *module.SimulationManager
}

func init() {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	DefaultNodeHome = filepath.Join(userHomeDir, ".simapp")
}

// NewSimApp returns a reference to an initialized SimApp.
func NewInitApp(
	logger log.Logger, db dbm.DB, traceStore io.Writer, loadLatest bool, skipUpgradeHeights map[int64]bool,
	homePath string, invCheckPeriod uint, encodingConfig simappparams.EncodingConfig, appOpts servertypes.AppOptions, baseAppOptions ...func(*bam.BaseApp),
) *SekaiApp {
	// TODO: Remove cdc in favor of appCodec once all modules are migrated.
	appCodec := encodingConfig.Marshaler
	cdc := encodingConfig.Amino
	interfaceRegistry := encodingConfig.InterfaceRegistry

	// BaseApp handles interactions with Tendermint through the ABCI protocol
	bApp := bam.NewBaseApp(appName, logger, db, encodingConfig.TxConfig.TxDecoder(), baseAppOptions...)
	bApp.SetCommitMultiStoreTracer(traceStore)
	bApp.SetAppVersion(version.Version)
	bApp.SetInterfaceRegistry(interfaceRegistry)

	keys := sdk.NewKVStoreKeys(
		authtypes.StoreKey,
		banktypes.StoreKey,
		paramstypes.StoreKey,
		upgradetypes.StoreKey,
		customslashingtypes.ModuleName,
		customstakingtypes.ModuleName,
		customgovtypes.ModuleName,
		tokenstypes.ModuleName,
		feeprocessingtypes.ModuleName,
		evidencetypes.StoreKey,
	)
	tKeys := sdk.NewTransientStoreKeys(paramstypes.TStoreKey)

	// Here you initialize your application with the store keys it requires
	app := &SekaiApp{
		BaseApp:           bApp,
		cdc:               cdc,
		appCodec:          appCodec,
		interfaceRegistry: interfaceRegistry,
		invCheckPeriod:    invCheckPeriod,
		keys:              keys,
		tKeys:             tKeys,
	}

	app.paramsKeeper = initParamsKeeper(appCodec, cdc, keys[paramstypes.StoreKey], tKeys[paramstypes.TStoreKey])

	// set the BaseApp's parameter store
	app.SetParamStore(app.paramsKeeper.Subspace(bam.Paramspace).WithKeyTable(paramskeeper.ConsensusParamsKeyTable()))

	// The AccountKeeper handles address -> account lookups
	app.accountKeeper = authkeeper.NewAccountKeeper(
		appCodec, keys[authtypes.StoreKey], app.GetSubspace(authtypes.ModuleName), authtypes.ProtoBaseAccount, maccPerms,
	)
	app.bankKeeper = bankkeeper.NewBaseKeeper(
		appCodec, keys[banktypes.StoreKey], app.accountKeeper, app.GetSubspace(banktypes.ModuleName), app.BlockedAddrs(),
	)
	app.upgradeKeeper = upgradekeeper.NewKeeper(keys[upgradetypes.StoreKey], appCodec)

	app.customGovKeeper = customgovkeeper.NewKeeper(keys[customgovtypes.ModuleName], appCodec, app.bankKeeper)
	customStakingKeeper := customstakingkeeper.NewKeeper(keys[customstakingtypes.ModuleName], cdc, app.customGovKeeper)
	app.customSlashingKeeper = customslashingkeeper.NewKeeper(
		appCodec, keys[customslashingtypes.StoreKey], &customStakingKeeper, app.customGovKeeper, app.GetSubspace(customslashingtypes.ModuleName),
	)
	app.tokensKeeper = tokenskeeper.NewKeeper(keys[tokenstypes.ModuleName], appCodec)
	// NOTE: customStakingKeeper above is passed by reference, so that it will contain these hooks
	app.customStakingKeeper = *customStakingKeeper.SetHooks(
		customstakingtypes.NewMultiStakingHooks(app.customSlashingKeeper.Hooks()),
	)

	app.feeprocessingKeeper = feeprocessingkeeper.NewKeeper(keys[feeprocessingtypes.ModuleName], appCodec, app.bankKeeper, app.tokensKeeper, app.customGovKeeper)

	// create evidence keeper with router
	evidenceKeeper := evidencekeeper.NewKeeper(
		appCodec, keys[evidencetypes.StoreKey], &app.customStakingKeeper, app.customSlashingKeeper,
	)
	// If evidence needs to be handled for the app, set routes in router here and seal
	app.evidenceKeeper = *evidenceKeeper

	/****  Module Options ****/

	// NOTE: Any module instantiated in the module manager that is later modified
	// must be passed by reference here.
	app.mm = module.NewManager(
		auth.NewAppModule(appCodec, app.accountKeeper, simulation.RandomGenesisAccounts),
		genutil.NewAppModule(
			app.accountKeeper, app.customStakingKeeper, app.BaseApp.DeliverTx,
			encodingConfig.TxConfig,
		),
		bank.NewAppModule(appCodec, app.bankKeeper, app.accountKeeper),
		upgrade.NewAppModule(app.upgradeKeeper, app.customGovKeeper),
		params.NewAppModule(app.paramsKeeper),
		customslashing.NewAppModule(appCodec, app.customSlashingKeeper, app.accountKeeper, app.bankKeeper, app.customStakingKeeper),
		customstaking.NewAppModule(app.customStakingKeeper, app.customGovKeeper),
		customgov.NewAppModule(app.customGovKeeper, customgov.NewProposalRouter(
			[]customgov.ProposalHandler{
				customgov.NewApplyAssignPermissionProposalHandler(app.customGovKeeper),
				customgov.NewApplySetNetworkPropertyProposalHandler(app.customGovKeeper),
				customgov.NewApplyUpsertDataRegistryProposalHandler(app.customGovKeeper),
				customgov.NewApplySetPoorNetworkMessagesProposalHandler(app.customGovKeeper),
				tokens.NewApplyUpsertTokenAliasProposalHandler(app.tokensKeeper),
				tokens.NewApplyUpsertTokenRatesProposalHandler(app.tokensKeeper),
				tokens.NewApplyWhiteBlackChangeProposalHandler(app.tokensKeeper),
				customstaking.NewApplyUnjailValidatorProposalHandler(app.customStakingKeeper),
				customslashing.NewApplyResetWholeValidatorRankProposalHandler(app.customSlashingKeeper),
				customgov.NewApplyCreateRoleProposalHandler(app.customGovKeeper),
				upgrade.NewApplySoftwareUpgradeProposalHandler(app.upgradeKeeper),
			},
		)),
		tokens.NewAppModule(app.tokensKeeper, app.customGovKeeper),
		feeprocessing.NewAppModule(app.feeprocessingKeeper),
		evidence.NewAppModule(app.evidenceKeeper),
	)

	// During begin block slashing happens after distr.BeginBlocker so that
	// there is nothing left over in the validator fee pool, so as to keep the
	// CanWithdrawInvariant invariant.
	app.mm.SetOrderBeginBlockers(
		upgradetypes.ModuleName, customslashingtypes.ModuleName,
		evidencetypes.ModuleName, customstakingtypes.ModuleName,
	)
	app.mm.SetOrderEndBlockers(
		customgovtypes.ModuleName,
		customstakingtypes.ModuleName,
		feeprocessingtypes.ModuleName,
	)

	// NOTE: The genutils moodule must occur after staking so that pools are
	// properly initialized with tokens from genesis accounts.
	// NOTE: Capability module must occur first so that it can initialize any capabilities
	// so that other modules that want to create or claim capabilities afterwards in InitChain
	// can do so safely.
	app.mm.SetOrderInitGenesis(
		authtypes.ModuleName,
		banktypes.ModuleName,
		customstakingtypes.ModuleName,
		customslashingtypes.ModuleName,
		customgovtypes.ModuleName,
		tokenstypes.ModuleName,
		feeprocessingtypes.ModuleName,
		genutiltypes.ModuleName,
		evidencetypes.ModuleName,
	)

	app.mm.RegisterRoutes(app.Router(), app.QueryRouter(), encodingConfig.Amino)
	app.mm.RegisterServices(module.NewConfigurator(app.MsgServiceRouter(), app.GRPCQueryRouter()))

	// add test gRPC service for testing gRPC queries in isolation
	testdata.RegisterQueryServer(app.GRPCQueryRouter(), testdata.QueryImpl{})

	// create the simulation manager and define the order of the modules for deterministic simulations
	//
	// NOTE: this is not required apps that don't use the simulator for fuzz testing
	// transactions
	app.sm = module.NewSimulationManager(
		auth.NewAppModule(appCodec, app.accountKeeper, simulation.RandomGenesisAccounts),
		bank.NewAppModule(appCodec, app.bankKeeper, app.accountKeeper),
		customslashing.NewAppModule(appCodec, app.customSlashingKeeper, app.accountKeeper, app.bankKeeper, app.customStakingKeeper),
		params.NewAppModule(app.paramsKeeper),
		evidence.NewAppModule(app.evidenceKeeper),
	)

	app.sm.RegisterStoreDecoders()

	// initialize stores
	app.MountKVStores(keys)
	app.MountTransientStores(tKeys)

	// initialize BaseApp
	app.SetInitChainer(app.InitChainer)
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetAnteHandler(
		customante.NewAnteHandler(
			app.customStakingKeeper,
			app.customGovKeeper,
			app.tokensKeeper,
			app.feeprocessingKeeper,
			app.accountKeeper,
			app.bankKeeper,
			ante.DefaultSigVerificationGasConsumer,
			encodingConfig.TxConfig.SignModeHandler(),
		),
	)
	app.SetEndBlocker(app.EndBlocker)

	if loadLatest {
		if err := app.LoadLatestVersion(); err != nil {
			tmos.Exit(err.Error())
		}
	}

	middleware.SetKeepers(app.customGovKeeper, app.feeprocessingKeeper)

	return app
}

// Name returns the name of the App
func (app *SekaiApp) Name() string { return app.BaseApp.Name() }

// BeginBlocker application updates every begin block
func (app *SekaiApp) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return app.mm.BeginBlock(ctx, req)
}

// EndBlocker application updates every end block
func (app *SekaiApp) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	return app.mm.EndBlock(ctx, req)
}

// InitChainer application update at chain initialization
func (app *SekaiApp) InitChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	var genesisState GenesisState
	app.cdc.MustUnmarshalJSON(req.AppStateBytes, &genesisState)
	return app.mm.InitGenesis(ctx, app.appCodec, genesisState)
}

// LoadHeight loads a particular height
func (app *SekaiApp) LoadHeight(height int64) error {
	return app.LoadVersion(height)
}

// ModuleAccountAddrs returns all the app's module account addresses.
func (app *SekaiApp) ModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[authtypes.NewModuleAddress(acc).String()] = true
	}

	return modAccAddrs
}

// BlockedAddrs returns all the app's module account addresses that are not
// allowed to receive external tokens.
func (app *SekaiApp) BlockedAddrs() map[string]bool {
	blockedAddrs := make(map[string]bool)
	for acc := range maccPerms {
		blockedAddrs[authtypes.NewModuleAddress(acc).String()] = !allowedReceivingModAcc[acc]
	}

	return blockedAddrs
}

// Codec returns SimApp's codec.
//
// NOTE: This is solely to be used for testing purposes as it may be desirable
// for modules to register their own custom testing types.
func (app *SekaiApp) LegacyAmino() *codec.LegacyAmino {
	return app.cdc
}

// AppCodec returns SimApp's app codec.
//
// NOTE: This is solely to be used for testing purposes as it may be desirable
// for modules to register their own custom testing types.
func (app *SekaiApp) AppCodec() codec.Marshaler {
	return app.appCodec
}

// InterfaceRegistry returns SimApp's InterfaceRegistry
func (app *SekaiApp) InterfaceRegistry() types.InterfaceRegistry {
	return app.interfaceRegistry
}

// GetKey returns the KVStoreKey for the provided store key.
//
// NOTE: This is solely to be used for testing purposes.
func (app *SekaiApp) GetKey(storeKey string) *sdk.KVStoreKey {
	return app.keys[storeKey]
}

// GetTKey returns the TransientStoreKey for the provided store key.
//
// NOTE: This is solely to be used for testing purposes.
func (app *SekaiApp) GetTKey(storeKey string) *sdk.TransientStoreKey {
	return app.tKeys[storeKey]
}

// GetMemKey returns the MemStoreKey for the provided mem key.
//
// NOTE: This is solely used for testing purposes.
func (app *SekaiApp) GetMemKey(storeKey string) *sdk.MemoryStoreKey {
	return app.memKeys[storeKey]
}

// GetSubspace returns a param subspace for a given module name.
//
// NOTE: This is solely to be used for testing purposes.
func (app *SekaiApp) GetSubspace(moduleName string) paramstypes.Subspace {
	subspace, _ := app.paramsKeeper.GetSubspace(moduleName)
	return subspace
}

// SimulationManager implements the SimulationApp interface
func (app *SekaiApp) SimulationManager() *module.SimulationManager {
	return app.sm
}

// RegisterAPIRoutes registers all application module routes with the provided
// API server.
func (app *SekaiApp) RegisterAPIRoutes(apiSvr *api.Server, apiConfig config.APIConfig) {
	clientCtx := apiSvr.ClientCtx
	rpc.RegisterRoutes(clientCtx, apiSvr.Router)
	authrest.RegisterTxRoutes(clientCtx, apiSvr.Router)
	// Register new tx routes from grpc-gateway.
	authtx.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	// Register legacy and grpc-gateway routes for all modules.
	ModuleBasics.RegisterRESTRoutes(clientCtx, apiSvr.Router)
	ModuleBasics.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	// register swagger API from root so that other applications can override easily
	if apiConfig.Swagger {
		RegisterSwaggerAPI(clientCtx, apiSvr.Router)
	}
}

// RegisterTxService implements the Application.RegisterTxService method.
func (app *SekaiApp) RegisterTxService(clientCtx client.Context) {
	authtx.RegisterTxService(app.BaseApp.GRPCQueryRouter(), clientCtx, app.BaseApp.Simulate, app.interfaceRegistry)
}

// RegisterTendermintService implements the Application.RegisterTendermintService method.
func (app *SekaiApp) RegisterTendermintService(clientCtx client.Context) {
	tmservice.RegisterTendermintService(app.BaseApp.GRPCQueryRouter(), clientCtx, app.interfaceRegistry)
}

// RegisterSwaggerAPI registers swagger route with API Server
func RegisterSwaggerAPI(ctx client.Context, rtr *mux.Router) {
	statikFS, err := fs.New()
	if err != nil {
		panic(err)
	}

	staticServer := http.FileServer(statikFS)
	rtr.PathPrefix("/swagger/").Handler(http.StripPrefix("/swagger/", staticServer))
}

// GetMaccPerms returns a copy of the module account permissions
func GetMaccPerms() map[string][]string {
	dupMaccPerms := make(map[string][]string)
	for k, v := range maccPerms {
		dupMaccPerms[k] = v
	}
	return dupMaccPerms
}

// initParamsKeeper init params keeper and its subspaces
func initParamsKeeper(appCodec codec.BinaryMarshaler, legacyAmino *codec.LegacyAmino, key, tkey sdk.StoreKey) paramskeeper.Keeper {
	paramsKeeper := paramskeeper.NewKeeper(appCodec, legacyAmino, key, tkey)

	paramsKeeper.Subspace(authtypes.ModuleName)
	paramsKeeper.Subspace(banktypes.ModuleName)
	paramsKeeper.Subspace(customslashingtypes.ModuleName)

	return paramsKeeper
}
