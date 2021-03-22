package common

import (
	"time"

	"github.com/KiraCore/sekai/INTERX/config"
)

type BlockHeightTime struct {
	Height    int64 `json:"height"`
	BlockTime int64 `json:"timestamp"`
}

var (
	N                 int               = 0
	LatestNBlockTimes []BlockHeightTime = make([]BlockHeightTime, 0)
)

func GetAverageBlockTime() float64 {
	var total int64 = 0
	for _, block := range LatestNBlockTimes {
		total += block.BlockTime
	}

	GetLogger().Infof("[GetAverageBlockTime] %v", LatestNBlockTimes)

	return float64(total) / float64(len(LatestNBlockTimes))
}

func AddNewBlock(height int64) {
	timestamp, err := GetBlockTime(config.Config.RPC, height)
	if err != nil {
		GetLogger().Errorf("[AddNewBlock] Can't get block: %d", height)
		return
	}

	if len(LatestNBlockTimes) > 0 && LatestNBlockTimes[len(LatestNBlockTimes)-1].Height >= height {
		// not a new block
		GetLogger().Errorf("[AddNewBlock] not a new block: %d", height)
		return
	}

	prevBlockTimestamp, err := GetBlockTime(config.Config.RPC, height-1)
	if err != nil {
		GetLogger().Errorf("[AddNewBlock] Can't get block: %d", height-1)
		return
	}

	if len(LatestNBlockTimes) > 0 && timestamp-prevBlockTimestamp >= int64(GetAverageBlockTime()*3) {
		// a block just after a halt
		GetLogger().Errorf("[AddNewBlock] block just after a halt: %d", height)
		return
	}

	GetLogger().Infof("[GetAverageBlockTime] %v", LatestNBlockTimes)
	// insert new block
	LatestNBlockTimes = append(LatestNBlockTimes, BlockHeightTime{
		Height:    height,
		BlockTime: timestamp - prevBlockTimestamp,
	})

	GetLogger().Infof("[GetAverageBlockTime] %v", LatestNBlockTimes)
	if len(LatestNBlockTimes) > N {
		GetLogger().Infof("[GetAverageBlockTime] %v", LatestNBlockTimes)
		LatestNBlockTimes = LatestNBlockTimes[len(LatestNBlockTimes)-N:]
	}
}

func UpdateN(_N int) {
	if N > _N {
		LatestNBlockTimes = LatestNBlockTimes[N-_N:]
		N = _N
		return
	}

	for N < _N {
		var current = NodeStatus.Block - 1
		if len(LatestNBlockTimes) > 0 {
			current = LatestNBlockTimes[0].Height - 1
		}

		currentBlockTimestamp, err := GetBlockTime(config.Config.RPC, current)
		if err != nil {
			GetLogger().Errorf("[UpdateN] Can't get block: %d", current)
			return
		}

		prevBlockTimestamp, err := GetBlockTime(config.Config.RPC, current-1)
		if err != nil {
			GetLogger().Errorf("[UpdateN] Can't get block: %d", current-1)
			return
		}

		// insert new block
		LatestNBlockTimes = append(LatestNBlockTimes, BlockHeightTime{
			Height:    current,
			BlockTime: currentBlockTimestamp - prevBlockTimestamp,
		})

		N++
	}
}

func IsConsensusStopped(validatorCount int) bool {
	blockHeight := NodeStatus.Block

	if blockHeight <= 1 {
		GetLogger().Errorf("[UpdateN] block <= 1: %d", blockHeight)
		return false
	}

	var n int = int(blockHeight - 1)
	if n > validatorCount {
		n = validatorCount
	}

	UpdateN(n)

	blockHeight = NodeStatus.Block
	blockTime, _ := time.Parse(time.RFC3339, NodeStatus.Blocktime)

	return float64(time.Now().Unix()-blockTime.Unix()) >= GetAverageBlockTime()*3
}
