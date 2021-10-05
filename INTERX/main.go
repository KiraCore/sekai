package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/KiraCore/sekai/INTERX/common"
	"github.com/KiraCore/sekai/INTERX/config"
	"github.com/KiraCore/sekai/INTERX/gateway"
	_ "github.com/KiraCore/sekai/INTERX/statik"
	"github.com/tyler-smith/go-bip39"
	"google.golang.org/grpc/grpclog"
)

func printUsage() {
	fmt.Println("Interx Daemon (server)")
	fmt.Println()

	fmt.Println("Usage:")
	fmt.Printf("    interxd [command]\n")
	fmt.Println()

	fmt.Println("Available Commands:")
	fmt.Printf("    init	:	Generate interx configuration file.\n")
	fmt.Printf("    start	:	Start interx with configuration file.\n")
	fmt.Println()

	fmt.Println("Flags:")
	fmt.Printf("    -h, --help	:	help for interxd.\n")
	fmt.Println()

	fmt.Println("Use \"interxd [command] --help\" for more information about a command.")
}

func main() {
	initCommand := flag.NewFlagSet("init", flag.ExitOnError)
	startCommand := flag.NewFlagSet("start", flag.ExitOnError)

	entropy, _ := bip39.NewEntropy(256)
	signingMnemonic, _ := bip39.NewMnemonic(entropy)
	faucetMnemonic, _ := bip39.NewMnemonic(entropy)

	initConfigFilePtr := initCommand.String("config", "./config.json", "The interx configuration path.")
	initVersion := initCommand.String("version", "0.1.0", "The interxd version")
	initServeHTTPS := initCommand.Bool("serve_https", false, "http or https.")
	initGrpcPtr := initCommand.String("grpc", "dns:///0.0.0.0:9090", "The grpc endpoint of the sekaid.")
	initRPCPtr := initCommand.String("rpc", "http://0.0.0.0:26657", "The rpc endpoint of the sekaid.")

	initNodeType := initCommand.String("node_type", "seed", "The node type.")
	initSentryNodeId := initCommand.String("sentry_node_id", "", "The sentry node id.")
	initSnapshotNodeId := initCommand.String("snapshot_node_id", "", "The snapshot node id.")
	initValidatorNodeId := initCommand.String("validator_node_id", "", "The validator node id.")
	initSeedNodeId := initCommand.String("seed_node_id", "", "The seed node id.")

	initPortPtr := initCommand.String("port", "11000", "The interx port.")
	initSigningMnemonicPtr := initCommand.String("signing_mnemonic", signingMnemonic, "The interx signing mnemonic file path or seeds.")

	initSyncStatus := initCommand.Int64("status_sync", 5, "The time in seconds and INTERX syncs node status.")
	initHaltedAvgBlockTimes := initCommand.Int64("halted_avg_block_times", 10, "This will be used for checking consensus halted.")

	initCacheDirPtr := initCommand.String("cache_dir", "cache", "The interx cache directory path.")
	initMaxCacheSize := initCommand.String("max_cache_size", "2GB", "The maximum cache size.")
	initCachingDuration := initCommand.Int64("caching_duration", 5, "The caching clear duration in seconds.")
	initMaxDownloadSize := initCommand.String("download_file_size_limitation", "10MB", "The maximum download file size.")

	initFaucetMnemonicPtr := initCommand.String("faucet_mnemonic", faucetMnemonic, "The interx faucet mnemonic file path or seeds.")
	initFaucetTimeLimit := initCommand.Int64("faucet_time_limit", 20, "The claim time limitation in seconds.")

	initFaucetAmounts := initCommand.String("faucet_amounts", "100000stake,100000ukex,100000validatortoken", "The faucet amount for each asset.")
	initFaucetMinimumAmounts := initCommand.String("faucet_minimum_amounts", "1000stake,1000ukex,1000validatortoken", "The minimum faucet amount for each asset.")
	initFeeAmounts := initCommand.String("fee_amounts", "stake 1000ukex,ukex 1000ukex,validatortoken 1000ukex", "The fee amount for each denom. `stake 1000ukex` means it will use `1000ukex` for `stake` assets transfer.")

	initAddrBook := initCommand.String("addrbook", "addrbook.json", "The address books")
	initNodeKey := initCommand.String("node_key", "node_key.json", "The node key file path")
	initTxModes := initCommand.String("tx_modes", "sync,async,block", "The allowed transaction modes")

	initNodeDiscoveryUseHttps := initCommand.Bool("node_discovery_use_https", false, "The option to use https in node discovery")
	initNodeDiscoveryInterxPort := initCommand.String("node_discovery_interx_port", "11000", "The default interx port to be used in node discovery")
	initNodeDiscoveryTendermintPort := initCommand.String("node_discovery_tendermint_port", "26657", "The default tendermint port to be used in node discovery")
	initNodeDiscoveryTimeout := initCommand.String("node_discovery_timeout", "3s", "The connection timeout to be used in node discovery")

	startConfigPtr := startCommand.String("config", "./config.json", "The interx configurtion path. (Required)")

	flag.Usage = printUsage
	flag.Parse()

	// Verify that a subcommand has been provided
	// os.Arg[0] is the main command
	// os.Arg[1] will be the subcommand
	if len(os.Args) >= 2 {

		// Switch on the subcommand
		// Parse the flags for appropriate FlagSet
		// FlagSet.Parse() requires a set of arguments to parse as input
		// os.Args[2:] will be all arguments starting after the subcommand at os.Args[1]
		switch os.Args[1] {
		case "init":
			initCommand.Parse(os.Args[2:])

			if initCommand.Parsed() {
				// Check which subcommand was Parsed using the FlagSet.Parsed() function. Handle each case accordingly.
				// FlagSet.Parse() will evaluate to false if no flags were parsed (i.e. the user did not provide any flags)
				config.InitConfig(
					*initVersion,
					*initConfigFilePtr,
					*initServeHTTPS,
					*initGrpcPtr,
					*initRPCPtr,
					*initNodeType,
					*initSentryNodeId,
					*initSnapshotNodeId,
					*initValidatorNodeId,
					*initSeedNodeId,
					*initPortPtr,
					*initSigningMnemonicPtr,
					*initSyncStatus,
					*initHaltedAvgBlockTimes,
					*initCacheDirPtr,
					*initMaxCacheSize,
					*initCachingDuration,
					*initMaxDownloadSize,
					*initFaucetMnemonicPtr,
					*initFaucetTimeLimit,
					*initFaucetAmounts,
					*initFaucetMinimumAmounts,
					*initFeeAmounts,
					*initAddrBook,
					*initTxModes,
					*initNodeDiscoveryUseHttps,
					*initNodeDiscoveryInterxPort,
					*initNodeDiscoveryTendermintPort,
					*initNodeDiscoveryTimeout,
					*initNodeKey,
				)

				fmt.Printf("Created interx configuration file: %s\n", *initConfigFilePtr)
				return
			}
		case "start":
			startCommand.Parse(os.Args[2:])

			if startCommand.Parsed() {
				// Check which subcommand was Parsed using the FlagSet.Parsed() function. Handle each case accordingly.
				// FlagSet.Parse() will evaluate to false if no flags were parsed (i.e. the user did not provide any flags)
				configFilePath := *startConfigPtr
				fmt.Println("configFilePath", configFilePath)

				// Adds gRPC internal logs. This is quite verbose, so adjust as desired!
				log := common.GetLogger()
				grpclog.SetLoggerV2(log)

				err := gateway.Run(configFilePath, log)

				log.Fatalln(err)
				return
			}
		default:
			fmt.Println("init or start command is available.")
			os.Exit(1)
		}
	}

	printUsage()

}
