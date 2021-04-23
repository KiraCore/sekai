package database

import (
	"github.com/KiraCore/sekai/INTERX/config"
	"github.com/sonyarouje/simdb/db"
)

// BlockData is a struct for block details.
type BlockData struct {
	Height    int64 `json:"height"`
	Timestamp int64 `json:"timestamp"`
}

// ID is a field for facuet claim struct.
func (c BlockData) ID() (jsonField string, value interface{}) {
	value = c.Height
	jsonField = "height"
	return
}

func LoadBlockDbDriver() {
	DisableStdout()
	driver, _ := db.New(config.GetDbCacheDir() + "/block")
	EnableStdout()

	blockDb = driver
}

// GetBlockTime is a function to get blockTime
func GetBlockTime(height int64) (int64, error) {
	if blockDb == nil {
		panic("cache dir not set")
	}

	data := BlockData{}

	DisableStdout()
	err := blockDb.Open(BlockData{}).Where("height", "=", height).First().AsEntity(&data)
	EnableStdout()

	if err != nil {
		return 0, err
	}

	return data.Timestamp, nil
}

// AddBlockTime is a function to add blockTime
func AddBlockTime(height int64, timestamp int64) {
	if blockDb == nil {
		panic("cache dir not set")
	}

	data := BlockData{
		Height:    height,
		Timestamp: timestamp,
	}

	_, err := GetBlockTime(height)

	if err != nil {
		DisableStdout()
		err = blockDb.Open(BlockData{}).Insert(data)
		EnableStdout()

		if err != nil {
			panic(err)
		}

	}
}

var (
	blockDb *db.Driver
)
