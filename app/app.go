package app

import (
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/KiraCore/sekai/x/feeprocessing"
	customgov "github.com/KiraCore/sekai/x/gov"
	customstaking "github.com/KiraCore/sekai/x/staking"
	"github.com/KiraCore/sekai/x/tokens"

	ibc "github.com/cosmos/cosmos-sdk/x/ibc/core"
	"github.com/gorilla/mux"
	"github.com/rakyll/statik/fs"
	"github.com/spf13/cast"

	"github.com/cosmos/cosmos-sdk/x/auth/simulation"

	customgovkeeper "github.com/KiraCore/sekai/x/gov/keeper"
	customgovtypes "github.com/KiraCore/sekai/x/gov/types"

	tokenskeeper "github.com/KiraCore/sekai/x/tokens/keeper"
	tokenstypes "github.com/KiraCore/sekai/x/tokens/types"

	feeprocessingkeeper "github.com/KiraCore/sekai/x/feeprocessing/keeper"
	feeprocessingtypes "github.com/KiraCore/sekai/x/feeprocessing/types"

	simappparams "github.com/cosmos/cosmos-sdk/simapp/params"

	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/server/api"
	authrest "github.com/cosmos/cosmos-sdk/x/auth/client/rest"

	customstakingkeeper "github.com/KiraCore/sekai/x/staking/keeper"

	"github.com/cosmos/cosmos-sdk/x/crisis"
	crisiskeeper "github.com/cosmos/cosmos-sdk/x/crisis/keeper"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	customstakingtypes "github.com/KiraCore/sekai/x/staking/types"

	customante "github.com/KiraCore/sekai/app/ante"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"

	"github.com/cosmos/cosmos-sdk/codec/types"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	tmos "github.com/tendermint/tendermint/libs/os"
	dbm "github.com/tendermint/tm-db"

	"github.com/KiraCore/sekai/middleware"
	customslashing "github.com/KiraCore/sekai/x/slashing"
	customslashingkeeper "github.com/KiraCore/sekai/x/slashing/keeper"
	customslashingtypes "github.com/KiraCore/sekai/x/slashing/types"
	bam "github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server/config"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/capability"
	capabilitykeeper "github.com/cosmos/cosmos-sdk/x/capability/keeper"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distrclient "github.com/cosmos/cosmos-sdk/x/distribution/client"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/evidence"
	evidencekeeper "github.com/cosmos/cosmos-sdk/x/evidence/keeper"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/cosmos/cosmos-sdk/x/gov"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	transfer "github.com/cosmos/cosmos-sdk/x/ibc/applications/transfer"
	ibctransfertypes "github.com/cosmos/cosmos-sdk/x/ibc/applications/transfer/types"
	ibchost "github.com/cosmos/cosmos-sdk/x/ibc/core/24-host"
	ibcmock "github.com/cosmos/cosmos-sdk/x/ibc/testing/mock"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramsclient "github.com/cosmos/cosmos-sdk/x/params/client"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/cosmos-sdk/x/upgrade"
	upgradeclient "github.com/cosmos/cosmos-sdk/x/upgrade/client"
	upgradekeeper "github.com/cosmos/cosmos-sdk/x/upgrade/keeper"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	//"github.com/KiraCore/sekai/x/kiraHub"
	//constants "github.com/KiraCore/sekai/x/kiraHub/constants"
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
		capability.AppModuleBasic{},
		gov.NewAppModuleBasic(
			paramsclient.ProposalHandler, distrclient.ProposalHandler, upgradeclient.ProposalHandler, upgradeclient.CancelProposalHandler,
		),
		params.AppModuleBasic{},
		crisis.AppModuleBasic{},
		slashing.AppModuleBasic{},
		customslashing.AppModuleBasic{},
		ibc.AppModuleBasic{},
		upgrade.AppModuleBasic{},
		evidence.AppModuleBasic{},
		transfer.AppModuleBasic{},
		customstaking.AppModuleBasic{},
		customgov.AppModuleBasic{},
		tokens.AppModuleBasic{},
		feeprocessing.AppModuleBasic{},
	)

	// module account permissions
	maccPerms = map[string][]string{
		authtypes.FeeCollectorName:  nil,
		distrtypes.ModuleName:       nil,
		minttypes.ModuleName:        {authtypes.Minter},
		govtypes.ModuleName:         {authtypes.Burner},
		ibctransfertypes.ModuleName: {authtypes.Minter, authtypes.Burner},
	}

	// module accounts that are allowed to receive tokens
	allowedReceivingModAcc = map[string]bool{
		distrtypes.ModuleName: true,
	}
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
	accountKeeper    authkeeper.AccountKeeper
	bankKeeper       bankkeeper.Keeper
	capabilityKeeper *capabilitykeeper.Keeper
	crisisKeeper     crisiskeeper.Keeper
	upgradeKeeper    upgradekeeper.Keeper
	paramsKeeper     paramskeeper.Keeper
	evidenceKeeper   evidencekeeper.Keeper

	customSlashingKeeper customslashingkeeper.Keeper
	customStakingKeeper  customstakingkeeper.Keeper
	customGovKeeper      customgovkeeper.Keeper
	tokensKeeper         tokenskeeper.Keeper
	feeprocessingKeeper  feeprocessingkeeper.Keeper

	// make scoped keepers public for test purposes
	scopedIBCKeeper      capabilitykeeper.ScopedKeeper
	scopedTransferKeeper capabilitykeeper.ScopedKeeper
	scopedIBCMockKeeper  capabilitykeeper.ScopedKeeper

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

	// TODO: Add the keys that module requires
	keys := sdk.NewKVStoreKeys(authtypes.StoreKey, banktypes.StoreKey,
		distrtypes.StoreKey, slashingtypes.StoreKey,
		govtypes.StoreKey, paramstypes.StoreKey, ibchost.StoreKey, upgradetypes.StoreKey,
		evidencetypes.StoreKey, ibctransfertypes.StoreKey, capabilitytypes.StoreKey,
		customslashingtypes.ModuleName,
		customstakingtypes.ModuleName, customgovtypes.ModuleName, tokenstypes.ModuleName, feeprocessingtypes.ModuleName,
	)
	tKeys := sdk.NewTransientStoreKeys(paramstypes.TStoreKey)
	memKeys := sdk.NewMemoryStoreKeys(capabilitytypes.MemStoreKey)

	// Here you initialize your application with the store keys it requires
	app := &SekaiApp{
		BaseApp:           bApp,
		cdc:               cdc,
		appCodec:          appCodec,
		interfaceRegistry: interfaceRegistry,
		invCheckPeriod:    invCheckPeriod,
		keys:              keys,
		tKeys:             tKeys,
		memKeys:           memKeys,
	}

	app.paramsKeeper = initParamsKeeper(appCodec, cdc, keys[paramstypes.StoreKey], tKeys[paramstypes.TStoreKey])

	// set the BaseApp's parameter store
	app.SetParamStore(app.paramsKeeper.Subspace(bam.Paramspace).WithKeyTable(paramskeeper.ConsensusParamsKeyTable()))

	app.capabilityKeeper = capabilitykeeper.NewKeeper(appCodec, keys[capabilitytypes.StoreKey], memKeys[capabilitytypes.MemStoreKey])
	scopedIBCKeeper := app.capabilityKeeper.ScopeToModule(ibchost.ModuleName)
	scopedTransferKeeper := app.capabilityKeeper.ScopeToModule(ibctransfertypes.ModuleName)
	// NOTE: the IBC mock keeper and application module is used only for testing core IBC. Do
	// note replicate if you do not need to test core IBC or light clients.
	scopedIBCMockKeeper := app.capabilityKeeper.ScopeToModule(ibcmock.ModuleName)

	// The AccountKeeper handles address -> account lookups
	app.accountKeeper = authkeeper.NewAccountKeeper(
		appCodec, keys[authtypes.StoreKey], app.GetSubspace(authtypes.ModuleName), authtypes.ProtoBaseAccount, maccPerms,
	)
	app.bankKeeper = bankkeeper.NewBaseKeeper(
		appCodec, keys[banktypes.StoreKey], app.accountKeeper, app.GetSubspace(banktypes.ModuleName), app.BlockedAddrs(),
	)
	app.crisisKeeper = crisiskeeper.NewKeeper(
		app.GetSubspace(crisistypes.ModuleName), invCheckPeriod, app.bankKeeper, authtypes.FeeCollectorName,
	)
	app.upgradeKeeper = upgradekeeper.NewKeeper(skipUpgradeHeights, keys[upgradetypes.StoreKey], appCodec, homePath)

	customStakingKeeper := customstakingkeeper.NewKeeper(keys[customstakingtypes.ModuleName], cdc)
	app.customSlashingKeeper = customslashingkeeper.NewKeeper(
		appCodec, keys[slashingtypes.StoreKey], &customStakingKeeper, app.GetSubspace(slashingtypes.ModuleName),
	)
	app.customGovKeeper = customgovkeeper.NewKeeper(keys[customgovtypes.ModuleName], appCodec)
	app.tokensKeeper = tokenskeeper.NewKeeper(keys[tokenstypes.ModuleName], appCodec)
	// NOTE: customStakingKeeper above is passed by reference, so that it will contain these hooks
	app.customStakingKeeper = *customStakingKeeper.SetHooks(
		customstakingtypes.NewMultiStakingHooks(app.customSlashingKeeper.Hooks()),
	)

	app.feeprocessingKeeper = feeprocessingkeeper.NewKeeper(keys[feeprocessingtypes.ModuleName], appCodec, app.bankKeeper, app.tokensKeeper, app.customGovKeeper)

	/****  Module Options ****/

	// NOTE: we may consider parsing `appOpts` inside module constructors. For the moment
	// we prefer to be more strict in what arguments the modules expect.
	var skipGenesisInvariants = cast.ToBool(appOpts.Get(crisis.FlagSkipGenesisInvariants))

	// NOTE: Any module instantiated in the module manager that is later modified
	// must be passed by reference here.
	app.mm = module.NewManager(
		auth.NewAppModule(appCodec, app.accountKeeper, simulation.RandomGenesisAccounts),
		bank.NewAppModule(appCodec, app.bankKeeper, app.accountKeeper),
		capability.NewAppModule(appCodec, *app.capabilityKeeper),
		crisis.NewAppModule(&app.crisisKeeper, skipGenesisInvariants),
		customslashing.NewAppModule(appCodec, app.customSlashingKeeper, app.accountKeeper, app.bankKeeper, app.customStakingKeeper),
		upgrade.NewAppModule(app.upgradeKeeper),
		evidence.NewAppModule(app.evidenceKeeper),
		params.NewAppModule(app.paramsKeeper),
		customstaking.NewAppModule(app.customStakingKeeper, app.customGovKeeper),
		customgov.NewAppModule(app.customGovKeeper, customgov.NewProposalRouter(
			[]customgov.ProposalHandler{
				customgov.NewApplyAssignPermissionProposalHandler(app.customGovKeeper),
				customgov.NewApplySetNetworkPropertyProposalHandler(app.customGovKeeper),
				customgov.NewApplyUpsertDataRegistryProposalHandler(app.customGovKeeper),
				tokens.NewApplyUpsertTokenAliasProposalHandler(app.tokensKeeper),
				tokens.NewApplyUpsertTokenRatesProposalHandler(app.tokensKeeper),
			},
		)),
		tokens.NewAppModule(app.tokensKeeper, app.customGovKeeper),
		feeprocessing.NewAppModule(app.feeprocessingKeeper),
	)

	// During begin block slashing happens after distr.BeginBlocker so that
	// there is nothing left over in the validator fee pool, so as to keep the
	// CanWithdrawInvariant invariant.
	app.mm.SetOrderBeginBlockers(
		upgradetypes.ModuleName,
		evidencetypes.ModuleName,
	)
	app.mm.SetOrderEndBlockers(
		crisistypes.ModuleName,
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
		capabilitytypes.ModuleName, authtypes.ModuleName /*distrtypes.ModuleName */ /*stakingtypes.ModuleName,*/, banktypes.ModuleName,
		/*slashingtypes.ModuleName, customslashingtypes.ModuleName, govtypes.ModuleName crisistypes.ModuleName,*/
		/*ibchost.ModuleName,*/
		/*genutiltypes.ModuleName, evidencetypes.ModuleName, ibctransfertypes.ModuleName, */
		customstakingtypes.ModuleName, customgovtypes.ModuleName, tokenstypes.ModuleName, feeprocessingtypes.ModuleName,
	)

	app.mm.RegisterInvariants(&app.crisisKeeper)
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
		capability.NewAppModule(appCodec, *app.capabilityKeeper),
		customslashing.NewAppModule(appCodec, app.customSlashingKeeper, app.accountKeeper, app.bankKeeper, app.customStakingKeeper),
		params.NewAppModule(app.paramsKeeper),
		evidence.NewAppModule(app.evidenceKeeper),
	)

	app.sm.RegisterStoreDecoders()

	// initialize stores
	app.MountKVStores(keys)
	app.MountTransientStores(tKeys)
	app.MountMemoryStores(memKeys)

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

	// Initialize and seal the capability keeper so all persistent capabilities
	// are loaded in-memory and prevent any further modules from creating scoped
	// sub-keepers.
	// This must be done during creation of baseapp rather than in InitChain so
	// that in-memory capabilities get regenerated on app restart
	ctx := app.BaseApp.NewUncachedContext(true, tmproto.Header{})
	app.capabilityKeeper.InitializeAndSeal(ctx)

	app.scopedIBCKeeper = scopedIBCKeeper
	app.scopedTransferKeeper = scopedTransferKeeper

	// NOTE: the IBC mock keeper and application module is used only for testing core IBC. Do
	// note replicate if you do not need to test core IBC or light clients.
	app.scopedIBCMockKeeper = scopedIBCMockKeeper

	middleware.SetKeepers(app.customGovKeeper, app.feeprocessingKeeper)

	return app
}

// MakeCodecs constructs the *std.Codec and *codec.LegacyAmino instances used by
// simapp. It is useful for tests and clients who do not want to construct the
// full simapp
func MakeCodecs() (codec.Marshaler, *codec.LegacyAmino) {
	config := MakeEncodingConfig()
	return config.Marshaler, config.Amino
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
	authtx.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCRouter)

	// Register legacy and grpc-gateway routes for all modules.
	ModuleBasics.RegisterRESTRoutes(clientCtx, apiSvr.Router)
	ModuleBasics.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCRouter)

	// register swagger API from root so that other applications can override easily
	if apiConfig.Swagger {
		RegisterSwaggerAPI(clientCtx, apiSvr.Router)
	}
}

// RegisterTxService implements the Application.RegisterTxService method.
func (app *SekaiApp) RegisterTxService(clientCtx client.Context) {
	authtx.RegisterTxService(app.BaseApp.GRPCQueryRouter(), clientCtx, app.BaseApp.Simulate, app.interfaceRegistry)
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
	paramsKeeper.Subspace(stakingtypes.ModuleName)
	paramsKeeper.Subspace(minttypes.ModuleName)
	paramsKeeper.Subspace(distrtypes.ModuleName)
	paramsKeeper.Subspace(slashingtypes.ModuleName)
	paramsKeeper.Subspace(customslashingtypes.ModuleName)
	paramsKeeper.Subspace(govtypes.ModuleName).WithKeyTable(govtypes.ParamKeyTable())
	paramsKeeper.Subspace(crisistypes.ModuleName)
	paramsKeeper.Subspace(ibctransfertypes.ModuleName)

	return paramsKeeper
}
