package database

import (
	interx "github.com/KiraCore/sekai/INTERX/config"
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

func getBlockDbDriver() *db.Driver {
	driver, err := db.New(interx.GetDbCacheDir() + "block")
	if err != nil {
		panic(err)
	}

	return driver
}

// GetBlockTime is a function to get blockTime
func GetBlockTime(height int64) (int64, error) {
	DisableStdout()

	data := BlockData{}
	err := blockDb.Open(BlockData{}).Where("height", "=", height).First().AsEntity(&data)
	if err != nil {
		EnableStdout()
		return 0, err
	}

	EnableStdout()
	return data.Timestamp, nil
}

// AddBlockTime is a function to add blockTime
func AddBlockTime(height int64, timestamp int64) {
	DisableStdout()

	data := BlockData{
		Height:    height,
		Timestamp: timestamp,
	}

	_, err := GetBlockTime(height)

	if err != nil {
		err = blockDb.Open(BlockData{}).Insert(data)
		if err != nil {
			panic(err)
		}
	}

	EnableStdout()
}

var (
	blockDb *db.Driver = getBlockDbDriver()
)
