package database

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

// GetBlockTime is a function to get blockTime
func GetBlockTime(height int64) (int64, error) {
	data := BlockData{}
	err := database.Open(BlockData{}).Where("height", "=", height).First().AsEntity(&data)
	if err != nil {
		return 0, err
	}

	return data.Timestamp, nil
}

// AddBlockTime is a function to add blockTime
func AddBlockTime(height int64, timestamp int64) {
	data := BlockData{
		Height:    height,
		Timestamp: timestamp,
	}

	_, err := GetBlockTime(height)

	if err != nil {
		err = database.Insert(data)
		if err != nil {
			panic(err)
		}
	}
}
