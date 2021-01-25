package interx

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/KiraCore/sekai/INTERX/common"
	interx "github.com/KiraCore/sekai/INTERX/config"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	tmjson "github.com/tendermint/tendermint/libs/json"
	tmtypes "github.com/tendermint/tendermint/types"
)

// RegisterGenesisQueryRoutes registers genesis query routers.
func RegisterGenesisQueryRoutes(r *mux.Router, gwCosmosmux *runtime.ServeMux, rpcAddr string) {
	r.HandleFunc(common.QueryGenesis, QueryGenesis(rpcAddr)).Methods("GET")
	r.HandleFunc(common.QueryGenesisSum, QueryGenesisSum(rpcAddr)).Methods("GET")

	common.AddRPCMethod("GET", common.QueryGenesis, "This is an API to query genesis.", true)
	common.AddRPCMethod("GET", common.QueryGenesisSum, "This is an API to get genesis checksum.", true)
}

var genesisPath = interx.GetReferenceCacheDir() + "/genesis.json"

func saveGenesis(rpcAddr string) error {
	_, err := getGenesis()
	if err == nil {
		return nil
	}

	data, _, _ := common.MakeGetRequest(rpcAddr, "/genesis", "")

	type GenesisResponse struct {
		Genesis tmtypes.GenesisDoc `json:"genesis"`
	}

	genesis := GenesisResponse{}
	byteData, err := json.Marshal(data)
	err = tmjson.Unmarshal(byteData, &genesis)
	if err != nil {
		return err
	}

	fmt.Println(genesis.Genesis)
	err = genesis.Genesis.ValidateAndComplete()
	if err != nil {
		return err
	}

	common.Mutex.Lock()
	err = genesis.Genesis.SaveAs(genesisPath)
	common.Mutex.Unlock()

	return err
}

func getGenesis() (string, error) {
	common.Mutex.Lock()
	data, err := ioutil.ReadFile(genesisPath)
	common.Mutex.Unlock()

	if err != nil {
		return "", err
	}

	return common.GetBlake2bHashFromBytes(data), nil
}

// QueryGenesis is a function to query genesis.
func QueryGenesis(rpcAddr string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if saveGenesis(rpcAddr) != nil {
			common.ServeError(0, "", "interx error", http.StatusInternalServerError)
		} else {
			http.ServeFile(w, r, genesisPath)
		}
	}
}

func queryGenesisSumHandler(rpcAddr string) (interface{}, interface{}, int) {
	saveGenesis(rpcAddr)
	checksum, err := getGenesis()
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
