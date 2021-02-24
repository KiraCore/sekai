package database

import (
	"time"

	"github.com/KiraCore/sekai/INTERX/config"
	"github.com/sonyarouje/simdb/db"
)

// FaucetClaim is a struct for facuet claim.
type FaucetClaim struct {
	Address string    `json:"address"`
	Claim   time.Time `json:"claim"`
}

// ID is a field for facuet claim struct.
func (c FaucetClaim) ID() (jsonField string, value interface{}) {
	value = c.Address
	jsonField = "address"
	return
}

func LoadFaucetDbDriver() {
	DisableStdout()
	driver, _ := db.New(config.GetDbCacheDir() + "/faucet")
	EnableStdout()

	faucetDb = driver
}

func isClaimExist(address string) bool {
	if faucetDb == nil {
		panic("cache dir not set")
	}

	data := FaucetClaim{}

	DisableStdout()
	err := faucetDb.Open(FaucetClaim{}).Where("address", "=", address).First().AsEntity(&data)
	EnableStdout()

	if err != nil {
		return false
	}

	return true
}

func getClaim(address string) time.Time {
	if faucetDb == nil {
		panic("cache dir not set")
	}

	data := FaucetClaim{}

	DisableStdout()
	err := faucetDb.Open(FaucetClaim{}).Where("address", "=", address).First().AsEntity(&data)
	EnableStdout()

	if err != nil {
		panic(err)
	}

	return data.Claim
}

// GetClaimTimeLeft is a function to get left time for next claim
func GetClaimTimeLeft(address string) int64 {
	if faucetDb == nil {
		panic("cache dir not set")
	}

	if !isClaimExist(address) {
		return 0
	}

	diff := time.Now().Unix() - getClaim(address).Unix()

	if diff > config.Config.Faucet.TimeLimit {
		return 0
	}

	return config.Config.Faucet.TimeLimit - diff
}

// AddNewClaim is a function to add current claim time
func AddNewClaim(address string, claim time.Time) {
	if faucetDb == nil {
		panic("cache dir not set")
	}

	data := FaucetClaim{
		Address: address,
		Claim:   claim,
	}

	exists := isClaimExist(address)

	if exists {
		DisableStdout()
		err := faucetDb.Open(FaucetClaim{}).Update(data)
		EnableStdout()

		if err != nil {
			panic(err)
		}
	} else {
		DisableStdout()
		err := faucetDb.Open(FaucetClaim{}).Insert(data)
		EnableStdout()

		if err != nil {
			panic(err)
		}
	}
}

var (
	faucetDb *db.Driver
)
