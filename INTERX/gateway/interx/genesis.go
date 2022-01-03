package interx

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/KiraCore/sekai/INTERX/common"
	"github.com/KiraCore/sekai/INTERX/config"
	"github.com/KiraCore/sekai/INTERX/global"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	tmjson "github.com/tendermint/tendermint/libs/json"
	tmtypes "github.com/tendermint/tendermint/types"
)

// RegisterGenesisQueryRoutes registers genesis query routers.
func RegisterGenesisQueryRoutes(r *mux.Router, gwCosmosmux *runtime.ServeMux, rpcAddr string) {
	r.HandleFunc(config.QueryGenesis, QueryGenesis(rpcAddr)).Methods("GET")
	r.HandleFunc(config.QueryGenesisSum, QueryGenesisSum(rpcAddr)).Methods("GET")

	common.AddRPCMethod("GET", config.QueryGenesis, "This is an API to query genesis.", true)
	common.AddRPCMethod("GET", config.QueryGenesisSum, "This is an API to get genesis checksum.", true)
}

func genesisPath() string {
	return config.GetReferenceCacheDir() + "/genesis.json"
}

func saveGenesis(rpcAddr string) error {
	_, err := getGenesisCheckSum()
	if err == nil {
		return nil
	}

	data, _, _ := common.MakeTendermintRPCRequest(rpcAddr, "/genesis", "")

	type GenesisResponse struct {
		Genesis tmtypes.GenesisDoc `json:"genesis"`
	}

	genesis := GenesisResponse{}
	byteData, _ := json.Marshal(data)
	err = tmjson.Unmarshal(byteData, &genesis)
	if err != nil {
		return err
	}

	err = genesis.Genesis.ValidateAndComplete()
	if err != nil {
		return err
	}

	global.Mutex.Lock()
	err = genesis.Genesis.SaveAs(genesisPath())
	global.Mutex.Unlock()

	return err
}

func getGenesisCheckSum() (string, error) {
	global.Mutex.Lock()
	data, err := ioutil.ReadFile(genesisPath())
	global.Mutex.Unlock()

	if err != nil {
		return "", err
	}

	return common.GetSha256SumFromBytes(data), nil
}

func GetGenesisResults(rpcAddr string) (*tmtypes.GenesisDoc, string, error) {
	err := saveGenesis(rpcAddr)
	if err != nil {
		return nil, "", err
	}

	global.Mutex.Lock()
	data, err := ioutil.ReadFile(genesisPath())
	global.Mutex.Unlock()

	if err != nil {
		return nil, "", err
	}

	genesis := tmtypes.GenesisDoc{}
	err = tmjson.Unmarshal(data, &genesis)

	return &genesis, common.GetSha256SumFromBytes(data), err
}

// QueryGenesis is a function to query genesis.
func QueryGenesis(rpcAddr string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if saveGenesis(rpcAddr) != nil {
			common.ServeError(0, "", "interx error", http.StatusInternalServerError)
		} else {
			http.ServeFile(w, r, genesisPath())
		}
	}
}

func queryGenesisSumHandler(rpcAddr string) (interface{}, interface{}, int) {
	saveGenesis(rpcAddr)
	checksum, err := getGenesisCheckSum()
	if err != nil {
		return common.ServeError(0, "", "interx error", http.StatusInternalServerError)
	}

	type GenesisChecksumResponse struct {
		Checksum string `json:"checksum,omitempty"`
	}
	result := GenesisChecksumResponse{
		Checksum: "0x" + checksum,
	}

	return result, nil, http.StatusOK
}

// QueryGenesisSum is a function to get genesis checksum.
func QueryGenesisSum(rpcAddr string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := common.GetInterxRequest(r)
		response := common.GetResponseFormat(request, rpcAddr)
		statusCode := http.StatusOK

		response.Response, response.Error, statusCode = queryGenesisSumHandler(rpcAddr)

		common.WrapResponse(w, request, *response, statusCode, false)
	}
}
