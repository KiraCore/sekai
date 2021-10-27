package common

import (
	"time"

	"github.com/KiraCore/sekai/INTERX/config"
)

type BlockHeightTime struct {
	Height    int64   `json:"height"`
	BlockTime float64 `json:"timestamp"`
}

var (
	N                 int               = 0
	LatestNBlockTimes []BlockHeightTime = make([]BlockHeightTime, 0)
)

func GetAverageBlockTime() float64 {
	var total float64 = 0
	for _, block := range LatestNBlockTimes {
		total += block.BlockTime
	}

	if total == 0 {
		return 0
	}

	// GetLogger().Infof("[GetAverageBlockTime] %v", LatestNBlockTimes)

	return total / float64(len(LatestNBlockTimes))
}

func AddNewBlock(height int64, timestamp int64) {
	if len(LatestNBlockTimes) > 0 && LatestNBlockTimes[len(LatestNBlockTimes)-1].Height >= height {
		// not a new block
		GetLogger().Errorf("[AddNewBlock] not a new block: %d", height)
		return
	}
	prevBlockTimestamp, err := GetBlockNanoTime(config.Config.RPC, height-1)
	if err != nil {
		GetLogger().Errorf("[AddNewBlock] Can't get block: %d", height-1)
		return
	}

	var timespan float64 = (float64(timestamp) - float64(prevBlockTimestamp)) / 1e9

	if len(LatestNBlockTimes) > 0 && timespan >= GetAverageBlockTime()*float64(config.Config.Block.HaltedAvgBlockTimes) {
		// a block just after a halt
		GetLogger().Errorf("[AddNewBlock] block just after a halt: %d, timestamp: %f, average: %f", height, timespan, GetAverageBlockTime())
		return
	}

	// insert new block
	LatestNBlockTimes = append(LatestNBlockTimes, BlockHeightTime{
		Height:    height,
		BlockTime: timespan,
	})

	if len(LatestNBlockTimes) > N {
		LatestNBlockTimes = LatestNBlockTimes[len(LatestNBlockTimes)-N:]
	}
}

func UpdateN(_N int) {
	if N > _N {
		LatestNBlockTimes = LatestNBlockTimes[N-_N:]
		N = _N
		return
	}

	var current = NodeStatus.Block - 1
	if len(LatestNBlockTimes) > 0 {
		current = LatestNBlockTimes[0].Height - 1
	}

	for N < _N {

		currentBlockTimestamp, err := GetBlockNanoTime(config.Config.RPC, current)
		if err != nil {
			GetLogger().Errorf("[UpdateN] Can't get block: %d", current)
			return
		}

		prevBlockTimestamp, err := GetBlockNanoTime(config.Config.RPC, current-1)
		if err != nil {
			GetLogger().Errorf("[UpdateN] Can't get block: %d", current-1)
			return
		}

		// insert new block
		LatestNBlockTimes = append(
			[]BlockHeightTime{
				{
					Height:    current,
					BlockTime: (float64(currentBlockTimestamp) - float64(prevBlockTimestamp)) / 1e9,
				},
			},
			LatestNBlockTimes...,
		)

		N++
		current = current - 1
	}
}

func IsConsensusStopped(validatorCount int) bool {
	blockHeight := NodeStatus.Block
	blockTime, _ := time.Parse(time.RFC3339, NodeStatus.Blocktime)

	if blockHeight <= 1 {
		GetLogger().Errorf("[UpdateN] block <= 1: %d", blockHeight)
		return false
	}

	var n int = int(blockHeight - 1)
	if n > validatorCount {
		n = validatorCount
	}

	UpdateN(n)

	if float64(time.Now().Unix()-blockTime.Unix()) < GetAverageBlockTime()*float64(config.Config.Block.HaltedAvgBlockTimes) {
		return false
	}

	_, err := GetBlockNanoTime(config.Config.RPC, blockHeight+1)

	return err != nil
}
