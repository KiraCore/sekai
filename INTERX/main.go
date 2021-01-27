package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/KiraCore/sekai/INTERX/common"
	"github.com/KiraCore/sekai/INTERX/config"
	"github.com/KiraCore/sekai/INTERX/gateway"
	_ "github.com/KiraCore/sekai/INTERX/statik"
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

	// helpPtr := flag.Bool("help", false, "Show interxd cli instructions.")

	initGrpcPtr := initCommand.String("grpc", "dns:///0.0.0.0:9090", "The grpc endpoint of the sekaid.")
	initRpcPtr := initCommand.String("rpc", "http://0.0.0.0:26657", "The rpc endpoint of the sekaid.")
	initPortPtr := initCommand.String("port", "11000", "The interx port.")
	initSigningMnemonicPtr := initCommand.String("signing_mnemonic", "swap exercise equip shoot mad inside floor wheel loan visual stereo build frozen always bulb naive subway foster marine erosion shuffle flee action there", "The interx signing mnemonic file path or seeds.")
	initFaucetMnemonicPtr := initCommand.String("faucet_mnemonic", "equip exercise shoot mad inside floor wheel loan visual stereo build frozen potato always bulb naive subway foster marine erosion shuffle flee action there", "The interx faucet mnemonic file path or seeds.")
	initConfigFilePtr := initCommand.String("config", "./config.json", "The interx configuration path.")

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
				config.InitConfig(*initConfigFilePtr, *initGrpcPtr, *initRpcPtr, *initPortPtr, *initSigningMnemonicPtr, *initFaucetMnemonicPtr)
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
