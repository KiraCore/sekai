package tasks

import (
	"github.com/KiraCore/sekai/INTERX/types"
)

var (
	PubP2PNodeListResponse    types.P2PNodeListResponse
	PrivP2PNodeListResponse   types.P2PNodeListResponse
	InterxP2PNodeListResponse types.InterxNodeListResponse
	SnapNodeListResponse      types.SnapNodeListResponse
)

func getP2PNodeInfo(ipAddr string) {

}

func NodeDiscover(isLog bool) {
	/*
		flag := make(map[string]bool)

		for {
			global.Mutex.Lock()
			PubP2PNodeListResponse.Scanning = true
			PrivP2PNodeListResponse.Scanning = true
			global.Mutex.Unlock()
			uniqueIPAddresses := config.LoadUniqueIPAddresses()
			pubNodeList := []types.P2PNode{}
			privNodeList := []types.P2PNode{}

			common.GetLogger().Info(config.Config.AddrBooks)
			common.GetLogger().Info("[node-discover] addresses = ", uniqueIPAddresses)

			// matchPubIPAddr := regexp.MustCompile("^([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])(?<!172\.(16|17|18|19|20|21|22|23|24|25|26|27|28|29|30|31))(?<!127)(?<!^10)(?<!^0)\.([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])(?<!192\.168)(?<!172\.(16|17|18|19|20|21|22|23|24|25|26|27|28|29|30|31))\.([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])(?<!\.255$)(?<!\b255.255.255.0\b)(?<!\b255.255.255.242\b)$")
			matchPubIPAddr := regexp.MustCompile("")

			for _, ipAddr := range uniqueIPAddresses {
				interxUrl := "http://" + ipAddr + ":" + config.Config.NodeDiscovery.DefaultInterxPort
				if config.Config.NodeDiscovery.UseHttps {
					interxUrl = "https://" + ipAddr + ":" + config.Config.NodeDiscovery.DefaultInterxPort
				}

				common.GetLogger().Info(interxUrl)

				info := getP2PNodeInfo(ipAddr)

				if matchPubIPAddr.MatchString(ipAddr) {
					// interxNodeList := common.GetInterxNodeList(interxUrl)
					// if interxNodeList != nil {
					// 	for _, node := range interxNodeList.NodeList {
					// 		if !flag[node.ID] {
					// 			flag[node.ID] = true
					// 			global.Mutex.Lock()
					// 			nodeList = append(nodeList, node)
					// 			global.Mutex.Unlock()
					// 		}
					// 	}
					// }

					interxStatus := common.GetInterxStatus(interxUrl)
					if interxStatus != nil {
						newNode := types.P2PNodeList{}
						newNode.IP = ipAddr
						// newNode.Connected = true
						// newNode.Port =
						// newNode.Ping =
						// newNode.Peers =
						// newNode.ID = interxStatus.ValidatorNodeID
						pubNodeList = append(pubNodeList, newNode)
					} else {

					}
				} else {

				}
			}

			global.Mutex.Lock()
			NodeListResponse.Scanning = false
			NodeListResponse.LastUpdate = time.Now().UTC().Unix()
			NodeListResponse.NodeList = nodeList
			global.Mutex.Unlock()

			time.Sleep(2 * time.Second)
		}*/
}
