package database

import (
	"time"

	"github.com/KiraCore/sekai/INTERX/config"
	"github.com/sonyarouje/simdb/db"
)

// ReferenceData is a struct for reference details.
type ReferenceData struct {
	Key           string    `json:"key"`
	URL           string    `json:"url"`
	ContentLength int64     `json:"content_length"`
	LastModified  time.Time `json:"last_modified"`
	FilePath      string    `json:"file_path"`
}

// ID is a field for facuet claim struct.
func (c ReferenceData) ID() (jsonField string, value interface{}) {
	value = c.Key
	jsonField = "key"
	return
}

func LoadReferenceDbDriver() {
	DisableStdout()
	driver, _ := db.New(config.GetDbCacheDir() + "/ref")
	EnableStdout()

	refDb = driver
}

// GetAllReferences is a function to get all references
func GetAllReferences() ([]ReferenceData, error) {
	if refDb == nil {
		panic("cache dir not set")
	}

	var references []ReferenceData

	DisableStdout()
	err := refDb.Open(ReferenceData{}).Get().AsEntity(&references)
	EnableStdout()

	return references, err
}

// GetReference is a function to get reference by key
func GetReference(key string) (ReferenceData, error) {
	if refDb == nil {
		panic("cache dir not set")
	}

	data := ReferenceData{}

	DisableStdout()
	err := refDb.Open(ReferenceData{}).Where("key", "=", key).First().AsEntity(&data)
	EnableStdout()

	return data, err
}

// AddReference is a function to add reference
func AddReference(key string, url string, contentLength int64, lastModified time.Time, filepath string) {
	if refDb == nil {
		panic("cache dir not set")
	}

	data := ReferenceData{
		Key:           key,
		URL:           url,
		ContentLength: contentLength,
		LastModified:  lastModified,
		FilePath:      filepath,
	}

	_, err := GetReference(key)

	if err == nil {
		DisableStdout()
		err := refDb.Open(ReferenceData{}).Update(data)
		EnableStdout()

		if err != nil {
			panic(err)
		}
	} else {
		DisableStdout()
		err := refDb.Open(ReferenceData{}).Insert(data)
		EnableStdout()

		if err != nil {
			panic(err)
		}
	}

}

var (
	refDb *db.Driver
)
