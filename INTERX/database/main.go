package database

import (
	"time"

	interx "github.com/KiraCore/sekai/INTERX/config"
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

func getDbDriver() *db.Driver {
	driver, err := db.New(".db")
	if err != nil {
		panic(err)
	}

	return driver
}

func isExist(address string) bool {
	data := FaucetClaim{}
	err := database.Open(FaucetClaim{}).Where("address", "=", address).First().AsEntity(&data)
	if err != nil {
		return false
	}

	return true
}

func getClaim(address string) time.Time {
	data := FaucetClaim{}
	err := database.Open(FaucetClaim{}).Where("address", "=", address).First().AsEntity(&data)
	if err != nil {
		panic(err)
	}

	return data.Claim
}

// GetClaimTimeLeft is a function to get left time for next claim
func GetClaimTimeLeft(address string) int64 {
	if !isExist(address) {
		return 0
	}

	diff := time.Now().Unix() - getClaim(address).Unix()

	if diff > interx.Config.Faucet.TimeLimit {
		return 0
	}

	return interx.Config.Faucet.TimeLimit - diff
}

// AddNewClaim is a function to add current claim time
func AddNewClaim(address string, claim time.Time) {
	data := FaucetClaim{
		Address: address,
		Claim:   claim,
	}

	if isExist(address) {
		err := database.Update(data)
		if err != nil {
			panic(err)
		}
	} else {
		err := database.Insert(data)
		if err != nil {
			panic(err)
		}
	}
}

var (
	database *db.Driver = getDbDriver()
)
