package tasks

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/KiraCore/sekai/INTERX/common"
	"github.com/KiraCore/sekai/INTERX/config"
	"github.com/KiraCore/sekai/INTERX/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"

	sekaitypes "github.com/KiraCore/sekai/types"
)

var (
	AllValidators types.AllValidators
)

const (
	// Undefined status
	Undefined string = "UNDEFINED"
	// Active status
	Active string = "ACTIVE"
	// Inactive status
	Inactive string = "INACTIVE"
	// Paused status
	Paused string = "PAUSED"
	// Jailed status
	Jailed string = "JAILED"
)

func ToString(data interface{}) string {
	out, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	return string(out)
}

func QueryValidators(gwCosmosmux *runtime.ServeMux, gatewayAddr string) error {
	type ValidatorsResponse = struct {
		Validators []types.QueryValidator `json:"validators,omitempty"`
		Actors     []string               `json:"actors,omitempty"`
		Pagination interface{}            `json:"pagination,omitempty"`
	}

	result := ValidatorsResponse{}

	limit := sekaitypes.PageIterationLimit - 1
	offset := 0
	for {
		validatorsQueryRequest, _ := http.NewRequest("GET", "http://"+gatewayAddr+config.QueryValidators+"?pagination.offset="+strconv.Itoa(offset)+"&pagination.limit="+strconv.Itoa(limit), nil)

		validatorsQueryResponse, failure, _ := common.ServeGRPC(validatorsQueryRequest, gwCosmosmux)

		if validatorsQueryResponse == nil {
			return errors.New(ToString(failure))
		}

		byteData, err := json.Marshal(validatorsQueryResponse)
		if err != nil {
			return err
		}

		subResult := ValidatorsResponse{}
		err = json.Unmarshal(byteData, &subResult)
		if err != nil {
			return err
		}

		if len(subResult.Validators) == 0 {
			break
		}

		result.Actors = subResult.Actors
		result.Validators = append(result.Validators, subResult.Validators...)

		offset += limit
	}

	type ValidatorInfoResponse = struct {
		ValValidatorInfos []types.ValidatorSigningInfo `json:"info,omitempty"`
	}
	validatorInfosResponse := ValidatorInfoResponse{}

	offset = 0
	for {
		validatorInfosQueryRequest, _ := http.NewRequest("GET", "http://"+gatewayAddr+config.QueryValidatorInfos+"?pagination.offset="+strconv.Itoa(offset)+"&pagination.limit="+strconv.Itoa(limit), nil)

		validatorInfosQueryResponse, failure, _ := common.ServeGRPC(validatorInfosQueryRequest, gwCosmosmux)

		if validatorInfosQueryResponse == nil {
			return errors.New(ToString(failure))
		}

		byteData, err := json.Marshal(validatorInfosQueryResponse)
		if err != nil {
			return err
		}

		subResult := ValidatorInfoResponse{}
		err = json.Unmarshal(byteData, &subResult)
		if err != nil {
			return err
		}

		if len(subResult.ValValidatorInfos) == 0 {
			break
		}

		validatorInfosResponse.ValValidatorInfos = append(validatorInfosResponse.ValValidatorInfos, subResult.ValValidatorInfos...)

		offset += limit
	}

	for index, validator := range result.Validators {

		pubkeyHexString := validator.Pubkey[14 : len(validator.Pubkey)-1]
		bytes, _ := hex.DecodeString(pubkeyHexString)
		pubkey := ed25519.PubKey{
			Key: bytes,
		}
		address := sdk.ConsAddress(pubkey.Address()).String()

		var valSigningInfo types.ValidatorSigningInfo
		for _, signingInfo := range validatorInfosResponse.ValValidatorInfos {
			if signingInfo.Address == address {
				valSigningInfo = signingInfo
				break
			}
		}

		for _, record := range result.Validators[index].Identity {
			if record.Key == "logo" || record.Key == "avatar" {
				result.Validators[index].Logo = record.Value
			} else if record.Key == "description" {
				result.Validators[index].Description = record.Value
			} else if record.Key == "website" {
				result.Validators[index].Website = record.Value
			} else if record.Key == "social" {
				result.Validators[index].Social = record.Value
			} else if record.Key == "contact" {
				result.Validators[index].Contact = record.Value
			} else if record.Key == "validator_node_id" {
				result.Validators[index].Validator_node_id = record.Value
			} else if record.Key == "sentry_node_id" {
				result.Validators[index].Sentry_node_id = record.Value
			}
		}

		result.Validators[index].Identity = nil
		result.Validators[index].StartHeight = valSigningInfo.StartHeight
		result.Validators[index].InactiveUntil = valSigningInfo.InactiveUntil
		result.Validators[index].Mischance = valSigningInfo.Mischance
		result.Validators[index].MischanceConfidence = valSigningInfo.MischanceConfidence
		result.Validators[index].LastPresentBlock = valSigningInfo.LastPresentBlock
		result.Validators[index].MissedBlocksCounter = valSigningInfo.MissedBlocksCounter
		result.Validators[index].ProducedBlocksCounter = valSigningInfo.ProducedBlocksCounter
	}

	sort.Sort(types.QueryValidators(result.Validators))
	for index := range result.Validators {
		result.Validators[index].Top = index + 1
	}

	allValidators := types.AllValidators{}

	allValidators.Validators = result.Validators
	allValidators.Waiting = make([]string, 0)
	for _, actor := range result.Actors {
		isWaiting := true
		for _, validator := range result.Validators {
			if validator.Address == actor {
				isWaiting = false
				break
			}
		}

		if isWaiting {
			allValidators.Waiting = append(allValidators.Waiting, actor)
		}
	}

	allValidators.Status.TotalValidators = len(result.Validators)
	allValidators.Status.WaitingValidators = len(allValidators.Waiting)

	allValidators.Status.ActiveValidators = 0
	allValidators.Status.PausedValidators = 0
	allValidators.Status.InactiveValidators = 0
	allValidators.Status.JailedValidators = 0
	for _, validator := range result.Validators {
		if validator.Status == Active {
			allValidators.Status.ActiveValidators++
		}
		if validator.Status == Inactive {
			allValidators.Status.InactiveValidators++
		}
		if validator.Status == Paused {
			allValidators.Status.PausedValidators++
		}
		if validator.Status == Jailed {
			allValidators.Status.JailedValidators++
		}
	}

	AllValidators = allValidators

	// common.GetLogger().Info(AllValidators)

	return nil
}

func SyncValidators(gwCosmosmux *runtime.ServeMux, gatewayAddr string, isLog bool) {
	lastBlock := int64(0)
	for {
		if common.NodeStatus.Block != lastBlock {
			err := QueryValidators(gwCosmosmux, gatewayAddr)

			if err != nil && isLog {
				common.GetLogger().Error("[sync-validators] Failed to query validators: ", err)
			}

			lastBlock = common.NodeStatus.Block
		}

		time.Sleep(1 * time.Second)
	}
}
