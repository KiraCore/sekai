package tasks

// RunTasks is a function to run threads.
func RunTasks(rpcAddr string) {
	go SyncStatus(rpcAddr, false)
	go CacheHeaderCheck(rpcAddr, false)
	go CacheDataCheck(rpcAddr, false)
	go CacheMaxSizeCheck(false)
}

// RunHostingTasks is a function to run threads for file hosting.
func RunHostingTasks() {
	go DataReferenceCheck(true)
}
