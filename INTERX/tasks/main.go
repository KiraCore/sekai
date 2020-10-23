package tasks

// RunTasks is a function to run threads.
func RunTasks(rpcAddr string) {
	go SyncStatus(rpcAddr)
	go CacheHeaderCheck(rpcAddr)
	go CacheDataCheck(rpcAddr)
}
