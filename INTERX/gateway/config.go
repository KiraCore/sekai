package gateway

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sort"

	sekaiapp "github.com/KiraCore/sekai/app"
)

// Endpoint is a struct for API endpiont.
type Endpoint struct {
	URL    string `json:"url"`
	Method string `json:"method"`
}

// APIConfig is a struct for configuration.
type APIConfig struct {
	API           Endpoint `json:"api"`
	Disable       bool     `json:"disable,omitempty"`
	RateLimit     float64  `json:"rate_limit,omitempty"`
	AuthRateLimit float64  `json:"auth_rate_limit,omitempty"`
}

// InterxConfig is a struct for configuration.
type InterxConfig []APIConfig

func readConfig() InterxConfig {
	file, _ := ioutil.ReadFile("./config.json")

	config := InterxConfig{}

	err := json.Unmarshal([]byte(file), &config)
	if err != nil {
		fmt.Println("Invalid configuration error: {}", err)
	}

	sort.SliceStable(config, func(i, j int) bool {
		if config[i].API.URL != config[j].API.URL {
			return config[i].API.URL < config[j].API.URL
		}
		return config[i].API.Method < config[j].API.Method
	})

	return config
}

var (
	config         = readConfig()
	encodingConfig = sekaiapp.MakeEncodingConfig()
	privKey        = GenEd25519PrivKey()
)
