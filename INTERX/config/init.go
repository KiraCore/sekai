package config

import "fmt"

// InitConfig is a function to load interx configurations from a given file
func InitConfig(configFilePath string, grpc string, rpc string, port string, signingMnemonic string, faucetMnemonic string) {
	fmt.Println(configFilePath, grpc, rpc, port)
	Config = InterxConfig{}
}
