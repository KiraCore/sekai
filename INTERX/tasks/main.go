package tasks

import (
	"os"

	interx "github.com/KiraCore/sekai/INTERX/config"
)

// RunTasks is a function to run threads.
func RunTasks(rpcAddr string) {
	if _, err := os.Stat(interx.Config.CacheDir); os.IsNotExist(err) {
		os.Mkdir(interx.Config.CacheDir, os.ModePerm)
	}
	if _, err := os.Stat(interx.Config.CacheDir + "/response"); os.IsNotExist(err) {
		os.Mkdir(interx.Config.CacheDir+"/response", os.ModePerm)
	}
	if _, err := os.Stat(interx.Config.CacheDir + "/reference"); os.IsNotExist(err) {
		os.Mkdir(interx.Config.CacheDir+"/reference", os.ModePerm)
	}

	go SyncStatus(rpcAddr, false)
	go CacheHeaderCheck(rpcAddr, false)
	go CacheDataCheck(rpcAddr, false)
	go CacheMaxSizeCheck(false)
	go DataReferenceCheck(true)
}
