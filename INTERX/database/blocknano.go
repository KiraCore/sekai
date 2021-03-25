package database

import (
	"github.com/KiraCore/sekai/INTERX/config"
	"github.com/sonyarouje/simdb/db"
)

// BlockNanoData is a struct for block details.
type BlockNanoData struct {
	Height    int64 `json:"height"`
	Timestamp int64 `json:"timestamp"`
}

// ID is a field for facuet claim struct.
func (c BlockNanoData) ID() (jsonField string, value interface{}) {
	value = c.Height
	jsonField = "height"
	return
}

func LoadBlockNanoDbDriver() {
	DisableStdout()
	driver, _ := db.New(config.GetDbCacheDir() + "/blocknano")
	EnableStdout()

	blockNanoDb = driver
}

// GetBlockNanoTime is a function to get blockTime
func GetBlockNanoTime(height int64) (int64, error) {
	if blockNanoDb == nil {
		panic("cache dir not set")
	}

	data := BlockNanoData{}

	DisableStdout()
	err := blockNanoDb.Open(BlockNanoData{}).Where("height", "=", height).First().AsEntity(&data)
	EnableStdout()

	if err != nil {
		return 0, err
	}

	return data.Timestamp, nil
}

// AddBlockNanoTime is a function to add blockTime
func AddBlockNanoTime(height int64, timestamp int64) {
	if blockNanoDb == nil {
		panic("cache dir not set")
	}

	data := BlockNanoData{
		Height:    height,
		Timestamp: timestamp,
	}

	_, err := GetBlockNanoTime(height)

	if err != nil {
		DisableStdout()
		err = blockNanoDb.Open(BlockNanoData{}).Insert(data)
		EnableStdout()

		if err != nil {
			panic(err)
		}

	}
}

var (
	blockNanoDb *db.Driver
)
