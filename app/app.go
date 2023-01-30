package app

import (
	"io"
	"net/http"
	"os"
	"path/filepath"

	customante "github.com/KiraCore/sekai/app/ante"
	"github.com/KiraCore/sekai/middleware"
	"github.com/KiraCore/sekai/x/basket"
	basketkeeper "github.com/KiraCore/sekai/x/basket/keeper"
	baskettypes "github.com/KiraCore/sekai/x/basket/types"
	"github.com/KiraCore/sekai/x/collectives"
	collectiveskeeper "github.com/KiraCore/sekai/x/collectives/keeper"
	collectivestypes "github.com/KiraCore/sekai/x/collectives/types"
	"github.com/KiraCore/sekai/x/custody"
	custodykeeper "github.com/KiraCore/sekai/x/custody/keeper"
	custodytypes "github.com/KiraCore/sekai/x/custody/types"
	"github.com/KiraCore/sekai/x/distributor"
	distributorkeeper "github.com/KiraCore/sekai/x/distributor/keeper"
	distributortypes "github.com/KiraCore/sekai/x/distributor/types"
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
	govtypes "github.com/KiraCore/sekai/x/gov/types"
	"github.com/KiraCore/sekai/x/multistaking"
	multistakingkeeper "github.com/KiraCore/sekai/x/multistaking/keeper"
	multistakingtypes "github.com/KiraCore/sekai/x/multistaking/types"
	recovery "github.com/KiraCore/sekai/x/recovery"
	recoverykeeper "github.com/KiraCore/sekai/x/recovery/keeper"
	recoverytypes "github.com/KiraCore/sekai/x/recovery/types"
	customslashing "github.com/KiraCore/sekai/x/slashing"
	customslashingkeeper "github.com/KiraCore/sekai/x/slashing/keeper"
	slashingtypes "github.com/KiraCore/sekai/x/slashing/types"
	"github.com/KiraCore/sekai/x/spending"
	spendingkeeper "github.com/KiraCore/sekai/x/spending/keeper"
	spendingtypes "github.com/KiraCore/sekai/x/spending/types"
	customstaking "github.com/KiraCore/sekai/x/staking"
	customstakingkeeper "github.com/KiraCore/sekai/x/staking/keeper"
	stakingtypes "github.com/KiraCore/sekai/x/staking/types"
	"github.com/KiraCore/sekai/x/tokens"
	tokenskeeper "github.com/KiraCore/sekai/x/tokens/keeper"
	tokenstypes "github.com/KiraCore/sekai/x/tokens/types"
	"github.com/KiraCore/sekai/x/ubi"
	ubikeeper "github.com/KiraCore/sekai/x/ubi/keeper"
	ubitypes "github.com/KiraCore/sekai/x/ubi/types"
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
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
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
		recovery.AppModuleBasic{},
		customstaking.AppModuleBasic{},
		customgov.AppModuleBasic{},
		spending.AppModuleBasic{},
		distributor.AppModuleBasic{},
		basket.AppModuleBasic{},
		ubi.AppModuleBasic{},
		evidence.AppModuleBasic{},
		tokens.AppModuleBasic{},
		feeprocessing.AppModuleBasic{},
		custody.AppModuleBasic{},
		multistaking.AppModuleBasic{},
		collectives.AppModuleBasic{},
	)

	// module account permissions
	maccPerms = map[string][]string{
		authtypes.FeeCollectorName:   nil,
		govtypes.ModuleName:          nil,
		minttypes.ModuleName:         {authtypes.Minter},
		spendingtypes.ModuleName:     nil,
		distributortypes.ModuleName:  nil,
		baskettypes.ModuleName:       {authtypes.Minter, authtypes.Burner},
		multistakingtypes.ModuleName: {authtypes.Burner},
		collectivestypes.ModuleName:  nil,
	}

	// module accounts that are allowed to receive tokens
	allowedReceivingModAcc = map[string]bool{}
)

// NewApp extended ABCI application
type SekaiApp struct {
	*bam.BaseApp
	cdc               *codec.LegacyAmino
	appCodec          codec.Codec
	interfaceRegistry types.InterfaceRegistry

	invCheckPeriod uint

	// keys to access the substores
	keys    map[string]*sdk.KVStoreKey
	tKeys   map[string]*sdk.TransientStoreKey
	memKeys map[string]*sdk.MemoryStoreKey

	// keepers
	AccountKeeper authkeeper.AccountKeeper
	BankKeeper    bankkeeper.Keeper
	UpgradeKeeper upgradekeeper.Keeper
	ParamsKeeper  paramskeeper.Keeper

	CustodyKeeper        custodykeeper.Keeper
	CustomGovKeeper      customgovkeeper.Keeper
	CustomStakingKeeper  customstakingkeeper.Keeper
	CustomSlashingKeeper customslashingkeeper.Keeper
	RecoveryKeeper       recoverykeeper.Keeper
	TokensKeeper         tokenskeeper.Keeper
	FeeProcessingKeeper  feeprocessingkeeper.Keeper
	EvidenceKeeper       evidencekeeper.Keeper
	SpendingKeeper       spendingkeeper.Keeper
	UbiKeeper            ubikeeper.Keeper
	DistrKeeper          distributorkeeper.Keeper
	BasketKeeper         basketkeeper.Keeper
	MultiStakingKeeper   multistakingkeeper.Keeper
	CollectivesKeeper    collectiveskeeper.Keeper

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

	DefaultNodeHome = filepath.Join(userHomeDir, ".sekaid")
}

// NewInitApp returns a reference to an initialized App.
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
	bApp.SetVersion(version.Version)
	bApp.SetInterfaceRegistry(interfaceRegistry)

	keys := sdk.NewKVStoreKeys(
		authtypes.StoreKey,
		banktypes.StoreKey,
		paramstypes.StoreKey,
		upgradetypes.StoreKey,
		recoverytypes.ModuleName,
		slashingtypes.ModuleName,
		stakingtypes.ModuleName,
		govtypes.ModuleName,
		spendingtypes.ModuleName,
		distributortypes.ModuleName,
		baskettypes.ModuleName,
		multistakingtypes.ModuleName,
		ubitypes.ModuleName,
		tokenstypes.ModuleName,
		feeprocessingtypes.ModuleName,
		evidencetypes.StoreKey,
		custodytypes.StoreKey,
		collectivestypes.ModuleName,
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

	app.ParamsKeeper = initParamsKeeper(appCodec, cdc, keys[paramstypes.StoreKey], tKeys[paramstypes.TStoreKey])

	// set the BaseApp's parameter store
	app.SetParamStore(app.ParamsKeeper.Subspace(bam.Paramspace).WithKeyTable(paramskeeper.ConsensusParamsKeyTable()))

	// The AccountKeeper handles address -> account lookups
	app.AccountKeeper = authkeeper.NewAccountKeeper(
		appCodec, keys[authtypes.StoreKey], app.GetSubspace(authtypes.ModuleName), authtypes.ProtoBaseAccount, maccPerms,
	)
	app.BankKeeper = bankkeeper.NewBaseKeeper(
		appCodec, keys[banktypes.StoreKey], app.AccountKeeper, app.GetSubspace(banktypes.ModuleName), app.BlockedAddrs(),
	)

	app.TokensKeeper = tokenskeeper.NewKeeper(keys[tokenstypes.ModuleName], appCodec)
	app.CustomGovKeeper = customgovkeeper.NewKeeper(keys[govtypes.ModuleName], appCodec, app.BankKeeper)
	customStakingKeeper := customstakingkeeper.NewKeeper(keys[stakingtypes.ModuleName], cdc, app.CustomGovKeeper)
	app.MultiStakingKeeper = multistakingkeeper.NewKeeper(keys[multistakingtypes.ModuleName], appCodec, app.BankKeeper, app.TokensKeeper, app.CustomGovKeeper, customStakingKeeper)
	app.CustomSlashingKeeper = customslashingkeeper.NewKeeper(
		appCodec,
		keys[slashingtypes.StoreKey],
		&customStakingKeeper,
		app.MultiStakingKeeper,
		app.CustomGovKeeper,
		app.GetSubspace(slashingtypes.ModuleName),
	)
	app.RecoveryKeeper = recoverykeeper.NewKeeper(
		appCodec,
		keys[slashingtypes.StoreKey],
		app.AccountKeeper,
		&customStakingKeeper,
		app.CustomGovKeeper,
		app.MultiStakingKeeper,
		app.CollectivesKeeper,
		app.SpendingKeeper,
		app.CustodyKeeper,
	)
	app.SpendingKeeper = spendingkeeper.NewKeeper(keys[spendingtypes.ModuleName], appCodec, app.BankKeeper, app.CustomGovKeeper)
	app.UbiKeeper = ubikeeper.NewKeeper(keys[ubitypes.ModuleName], appCodec, app.BankKeeper, app.SpendingKeeper)
	// NOTE: customStakingKeeper above is passed by reference, so that it will contain these hooks
	app.CustomStakingKeeper = *customStakingKeeper.SetHooks(
		stakingtypes.NewMultiStakingHooks(app.CustomSlashingKeeper.Hooks()),
	)
	app.DistrKeeper = distributorkeeper.NewKeeper(
		keys[distributortypes.ModuleName], appCodec,
		app.AccountKeeper, app.BankKeeper,
		app.CustomStakingKeeper, app.CustomGovKeeper,
		app.MultiStakingKeeper)
	app.MultiStakingKeeper.SetDistrKeeper(app.DistrKeeper)

	app.BasketKeeper = basketkeeper.NewKeeper(
		keys[baskettypes.ModuleName], appCodec,
		app.AccountKeeper, app.BankKeeper,
		app.CustomStakingKeeper, app.CustomGovKeeper,
		app.MultiStakingKeeper,
	)

	app.CollectivesKeeper = collectiveskeeper.NewKeeper(
		keys[collectivestypes.StoreKey], appCodec,
		app.BankKeeper,
		app.CustomGovKeeper,
		app.MultiStakingKeeper,
		app.TokensKeeper,
		app.SpendingKeeper,
	)

	app.UpgradeKeeper = upgradekeeper.NewKeeper(keys[upgradetypes.StoreKey], appCodec, app.CustomStakingKeeper)

	// app.upgradeKeeper.SetUpgradeHandler(
	// 	"upgrade1", func(ctx sdk.Context, plan upgradetypes.Plan) {
	// 	})

	app.FeeProcessingKeeper = feeprocessingkeeper.NewKeeper(keys[feeprocessingtypes.ModuleName], appCodec, app.BankKeeper, app.TokensKeeper, app.CustomGovKeeper)

	// create evidence keeper with router
	evidenceKeeper := evidencekeeper.NewKeeper(
		appCodec, keys[evidencetypes.StoreKey], &app.CustomStakingKeeper, app.CustomSlashingKeeper,
	)
	// If evidence needs to be handled for the app, set routes in router here and seal
	app.EvidenceKeeper = *evidenceKeeper

	app.CustodyKeeper = custodykeeper.NewKeeper(keys[custodytypes.StoreKey], appCodec, app.CustomGovKeeper, app.BankKeeper)

	proposalRouter := govtypes.NewProposalRouter(
		[]govtypes.ProposalHandler{
			customgov.NewApplyWhitelistAccountPermissionProposalHandler(app.CustomGovKeeper),
			customgov.NewApplyBlacklistAccountPermissionProposalHandler(app.CustomGovKeeper),
			customgov.NewApplyRemoveWhitelistedAccountPermissionProposalHandler(app.CustomGovKeeper),
			customgov.NewApplyRemoveBlacklistedAccountPermissionProposalHandler(app.CustomGovKeeper),
			customgov.NewApplyAssignRoleToAccountProposalHandler(app.CustomGovKeeper),
			customgov.NewApplyUnassignRoleFromAccountProposalHandler(app.CustomGovKeeper),
			customgov.NewApplySetNetworkPropertyProposalHandler(app.CustomGovKeeper),
			customgov.NewApplyUpsertDataRegistryProposalHandler(app.CustomGovKeeper),
			customgov.NewApplySetPoorNetworkMessagesProposalHandler(app.CustomGovKeeper),
			customgov.NewApplyResetWholeCouncilorRankProposalHandler(app.CustomGovKeeper),
			customgov.NewApplyJailCouncilorProposalHandler(app.CustomGovKeeper),
			tokens.NewApplyUpsertTokenAliasProposalHandler(app.TokensKeeper),
			tokens.NewApplyUpsertTokenRatesProposalHandler(app.TokensKeeper),
			tokens.NewApplyWhiteBlackChangeProposalHandler(app.TokensKeeper),
			customstaking.NewApplyUnjailValidatorProposalHandler(app.CustomStakingKeeper, app.CustomGovKeeper),
			customslashing.NewApplyResetWholeValidatorRankProposalHandler(app.CustomSlashingKeeper),
			customslashing.NewApplySlashValidatorProposalHandler(app.CustomSlashingKeeper),
			customgov.NewApplyCreateRoleProposalHandler(app.CustomGovKeeper),
			customgov.NewApplyRemoveRoleProposalHandler(app.CustomGovKeeper),
			customgov.NewApplyWhitelistRolePermissionProposalHandler(app.CustomGovKeeper),
			customgov.NewApplyBlacklistRolePermissionProposalHandler(app.CustomGovKeeper),
			customgov.NewApplyRemoveWhitelistedRolePermissionProposalHandler(app.CustomGovKeeper),
			customgov.NewApplyRemoveBlacklistedRolePermissionProposalHandler(app.CustomGovKeeper),
			customgov.NewApplySetProposalDurationsProposalHandler(app.CustomGovKeeper),
			upgrade.NewApplySoftwareUpgradeProposalHandler(app.UpgradeKeeper),
			upgrade.NewApplyCancelSoftwareUpgradeProposalHandler(app.UpgradeKeeper),
			spending.NewApplyUpdateSpendingPoolProposalHandler(app.SpendingKeeper),
			spending.NewApplySpendingPoolDistributionProposalHandler(app.SpendingKeeper, app.CustomGovKeeper),
			spending.NewApplySpendingPoolWithdrawProposalHandler(app.SpendingKeeper, app.BankKeeper),
			ubi.NewApplyUpsertUBIProposalHandler(app.UbiKeeper, app.CustomGovKeeper, app.SpendingKeeper),
			ubi.NewApplyRemoveUBIProposalHandler(app.UbiKeeper),
			basket.NewApplyCreateBasketProposalHandler(app.BasketKeeper),
			basket.NewApplyEditBasketProposalHandler(app.BasketKeeper),
			basket.NewApplyBasketWithdrawSurplusProposalHandler(app.BasketKeeper),
			collectives.NewApplyCollectiveSendDonationProposalHandler(app.CollectivesKeeper),
			collectives.NewApplyCollectiveUpdateProposalHandler(app.CollectivesKeeper),
			collectives.NewApplyCollectiveRemoveProposalHandler(app.CollectivesKeeper),
		})

	app.CustomGovKeeper.SetProposalRouter(proposalRouter)

	/****  Module Options ****/

	// NOTE: Any module instantiated in the module manager that is later modified
	// must be passed by reference here.
	app.mm = module.NewManager(
		auth.NewAppModule(appCodec, app.AccountKeeper, simulation.RandomGenesisAccounts),
		genutil.NewAppModule(
			app.AccountKeeper, app.CustomStakingKeeper, app.BaseApp.DeliverTx,
			encodingConfig.TxConfig,
		),
		bank.NewAppModule(appCodec, app.BankKeeper, app.AccountKeeper),
		upgrade.NewAppModule(app.UpgradeKeeper, app.CustomGovKeeper),
		params.NewAppModule(app.ParamsKeeper),
		customslashing.NewAppModule(appCodec, app.CustomSlashingKeeper, app.AccountKeeper, app.BankKeeper, app.CustomStakingKeeper),
		recovery.NewAppModule(appCodec, app.RecoveryKeeper, app.AccountKeeper, app.CustomStakingKeeper),
		customstaking.NewAppModule(app.CustomStakingKeeper, app.CustomGovKeeper),
		multistaking.NewAppModule(app.MultiStakingKeeper, app.BankKeeper, app.CustomGovKeeper, app.CustomStakingKeeper),
		customgov.NewAppModule(app.CustomGovKeeper),
		tokens.NewAppModule(app.TokensKeeper, app.CustomGovKeeper),
		spending.NewAppModule(app.SpendingKeeper, app.CustomGovKeeper, app.BankKeeper),
		distributor.NewAppModule(app.DistrKeeper, app.CustomGovKeeper),
		basket.NewAppModule(app.BasketKeeper, app.CustomGovKeeper),
		ubi.NewAppModule(app.UbiKeeper, app.CustomGovKeeper),
		feeprocessing.NewAppModule(app.FeeProcessingKeeper),
		evidence.NewAppModule(app.EvidenceKeeper),
		custody.NewAppModule(app.CustodyKeeper, app.CustomGovKeeper, app.BankKeeper),
		collectives.NewAppModule(app.CollectivesKeeper),
	)

	// During begin block slashing happens after distr.BeginBlocker so that
	// there is nothing left over in the validator fee pool, so as to keep the
	// CanWithdrawInvariant invariant.
	app.mm.SetOrderBeginBlockers(
		genutiltypes.ModuleName, paramstypes.ModuleName, govtypes.ModuleName, tokenstypes.ModuleName,
		authtypes.ModuleName, feeprocessingtypes.ModuleName, banktypes.ModuleName,
		upgradetypes.ModuleName, slashingtypes.ModuleName, recoverytypes.ModuleName,
		evidencetypes.ModuleName, stakingtypes.ModuleName,
		spendingtypes.ModuleName, ubitypes.ModuleName,
		distributortypes.ModuleName, multistakingtypes.ModuleName, custodytypes.ModuleName,
		baskettypes.ModuleName,
		distributortypes.ModuleName, multistakingtypes.ModuleName, custodytypes.ModuleName,
		baskettypes.ModuleName,
		collectivestypes.ModuleName,
	)
	app.mm.SetOrderEndBlockers(
		banktypes.ModuleName, upgradetypes.ModuleName, tokenstypes.ModuleName,
		evidencetypes.ModuleName, genutiltypes.ModuleName, paramstypes.ModuleName,
		slashingtypes.ModuleName, authtypes.ModuleName, recoverytypes.ModuleName,
		govtypes.ModuleName,
		stakingtypes.ModuleName,
		feeprocessingtypes.ModuleName,
		spendingtypes.ModuleName, ubitypes.ModuleName,
		distributortypes.ModuleName, multistakingtypes.ModuleName, custodytypes.ModuleName,
		baskettypes.ModuleName,
		collectivestypes.ModuleName,
	)

	// NOTE: The genutils moodule must occur after staking so that pools are
	// properly initialized with tokens from genesis accounts.
	// NOTE: Capability module must occur first so that it can initialize any capabilities
	// so that other modules that want to create or claim capabilities afterwards in InitChain
	// can do so safely.
	app.mm.SetOrderInitGenesis(
		authtypes.ModuleName,
		banktypes.ModuleName,
		govtypes.ModuleName, // staking module is using the moniker identity registrar and gov module should be initialized before
		stakingtypes.ModuleName,
		slashingtypes.ModuleName,
		recoverytypes.ModuleName,
		tokenstypes.ModuleName,
		feeprocessingtypes.ModuleName,
		genutiltypes.ModuleName,
		evidencetypes.ModuleName,
		upgradetypes.ModuleName,
		spendingtypes.ModuleName,
		ubitypes.ModuleName,
		paramstypes.ModuleName,
		distributortypes.ModuleName,
		custodytypes.ModuleName,
		multistakingtypes.ModuleName,
		baskettypes.ModuleName,
		collectivestypes.ModuleName,
	)

	app.mm.RegisterRoutes(app.Router(), app.QueryRouter(), encodingConfig.Amino)
	app.mm.RegisterServices(module.NewConfigurator(appCodec, app.MsgServiceRouter(), app.GRPCQueryRouter()))

	// add test gRPC service for testing gRPC queries in isolation
	testdata.RegisterQueryServer(app.GRPCQueryRouter(), testdata.QueryImpl{})

	// create the simulation manager and define the order of the modules for deterministic simulations
	//
	// NOTE: this is not required apps that don't use the simulator for fuzz testing
	// transactions
	app.sm = module.NewSimulationManager(
		auth.NewAppModule(appCodec, app.AccountKeeper, simulation.RandomGenesisAccounts),
		bank.NewAppModule(appCodec, app.BankKeeper, app.AccountKeeper),
		customslashing.NewAppModule(appCodec, app.CustomSlashingKeeper, app.AccountKeeper, app.BankKeeper, app.CustomStakingKeeper),
		recovery.NewAppModule(appCodec, app.RecoveryKeeper, app.AccountKeeper, app.CustomStakingKeeper),
		params.NewAppModule(app.ParamsKeeper),
		evidence.NewAppModule(app.EvidenceKeeper),
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
			app.CustomStakingKeeper,
			app.CustomGovKeeper,
			app.TokensKeeper,
			app.FeeProcessingKeeper,
			app.AccountKeeper,
			app.BankKeeper,
			ante.DefaultSigVerificationGasConsumer,
			encodingConfig.TxConfig.SignModeHandler(),
			app.CustodyKeeper,
		),
	)
	app.SetEndBlocker(app.EndBlocker)

	if loadLatest {
		if err := app.LoadLatestVersion(); err != nil {
			tmos.Exit(err.Error())
		}
	}

	middleware.SetKeepers(app.CustomGovKeeper, app.FeeProcessingKeeper)

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
func (app *SekaiApp) AppCodec() codec.Codec {
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
	subspace, _ := app.ParamsKeeper.GetSubspace(moduleName)
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
func initParamsKeeper(appCodec codec.BinaryCodec, legacyAmino *codec.LegacyAmino, key, tkey sdk.StoreKey) paramskeeper.Keeper {
	paramsKeeper := paramskeeper.NewKeeper(appCodec, legacyAmino, key, tkey)

	paramsKeeper.Subspace(authtypes.ModuleName)
	paramsKeeper.Subspace(banktypes.ModuleName)
	paramsKeeper.Subspace(stakingtypes.ModuleName)
	paramsKeeper.Subspace(slashingtypes.ModuleName)
	paramsKeeper.Subspace(multistakingtypes.ModuleName)
	paramsKeeper.Subspace(baskettypes.ModuleName)

	return paramsKeeper
}
