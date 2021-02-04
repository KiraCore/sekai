package tasks

// RunTasks is a function to run threads.
func RunTasks(rpcAddr string) {
	go SyncStatus(rpcAddr, false)
	go CacheHeaderCheck(rpcAddr, false)
	go CacheDataCheck(rpcAddr, false)
	go CacheMaxSizeCheck(false)
	go DataReferenceCheck(false)
}
