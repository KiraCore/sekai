package gateway

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	sekaiapp "github.com/KiraCore/sekai/app"
)

// PRCConfig is a struct to be used for PRC configuration
type PRCConfig struct {
	Disable       bool    `json:"disable"`
	RateLimit     float64 `json:"rate_limit,omitempty"`
	AuthRateLimit float64 `json:"auth_rate_limit,omitempty"`
}

func readConfig() map[string]map[string]PRCConfig {
	file, _ := ioutil.ReadFile("./config.json")

	config := map[string]map[string]PRCConfig{}

	err := json.Unmarshal([]byte(file), &config)
	if err != nil {
		fmt.Println("Invalid configuration error: {}", err)
	}

	fmt.Println(config)

	return config
}

var (
	config         = readConfig()
	encodingConfig = sekaiapp.MakeEncodingConfig()
	privKey        = GenEd25519PrivKey()
)
