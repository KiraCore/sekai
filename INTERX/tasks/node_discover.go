package tasks

import (
	"time"

	"github.com/KiraCore/sekai/INTERX/common"
	"github.com/KiraCore/sekai/INTERX/config"
	"github.com/KiraCore/sekai/INTERX/global"
	"github.com/KiraCore/sekai/INTERX/types"
)

var (
	NodeListResponse types.NodeListResponse
)

func getPingInfo() {

}

func NodeDiscover(isLog bool) {
	flag := make(map[string]bool)

	for {
		global.Mutex.Lock()
		NodeListResponse.Scanning = true
		global.Mutex.Unlock()
		uniqueAddresses := config.LoadUniqueAddresses()
		nodeList := []types.NodeList{}

		common.GetLogger().Info(config.Config.AddrBooks)
		common.GetLogger().Info("[node-discover] addresses = ", uniqueAddresses)

		for _, addr := range uniqueAddresses {
			ipAddr := addr.Addr.IP
			interxUrl := "http://" + ipAddr + ":" + config.Config.NodeDiscovery.DefaultInterxPort
			if config.Config.NodeDiscovery.UseHttps {
				interxUrl = "https://" + ipAddr + ":" + config.Config.NodeDiscovery.DefaultInterxPort
			}
			common.GetLogger().Info(interxUrl)

			interxNodeList := common.GetInterxNodeList(interxUrl)
			if interxNodeList != nil {
				for _, node := range interxNodeList.NodeList {
					if !flag[node.ID] {
						flag[node.ID] = true
						global.Mutex.Lock()
						nodeList = append(nodeList, node)
						global.Mutex.Unlock()
					}
				}
			}

			interxStatus := common.GetInterxStatus(interxUrl)
			if interxStatus != nil {
				newNode := types.NodeList{}
				newNode.Moniker = interxStatus.Moniker
				newNode.KiraAddr = interxStatus.KiraAddr
				newNode.Version = interxStatus.Version

				// Validator Node
				if len(interxStatus.ValidatorNodeID) > 0 {
					newNode.ID = interxStatus.ValidatorNodeID
					newNode.IP = ""
					newNode.Validator = true
					newNode.Seed = false
					global.Mutex.Lock()
					nodeList = append(nodeList, newNode)
					global.Mutex.Unlock()
				}

				newNode.IP = ipAddr
				newNode.Validator = false

				// Seed Node
				if len(interxStatus.SeedNodeID) > 0 {
					newNode.ID = interxStatus.SeedNodeID
					newNode.Seed = true
					global.Mutex.Lock()
					nodeList = append(nodeList, newNode)
					global.Mutex.Unlock()
				}

				newNode.Seed = false

				// Sentry Node
				if len(interxStatus.SentryNodeID) > 0 {
					newNode.ID = interxStatus.SentryNodeID
					global.Mutex.Lock()
					nodeList = append(nodeList, newNode)
					global.Mutex.Unlock()
				}

				// Priv Sentry Node
				if len(interxStatus.PrivSentryNodeID) > 0 {
					newNode.ID = interxStatus.PrivSentryNodeID
					global.Mutex.Lock()
					nodeList = append(nodeList, newNode)
					global.Mutex.Unlock()
				}
			}
		}

		global.Mutex.Lock()
		NodeListResponse.Scanning = false
		NodeListResponse.LastUpdate = time.Now().UTC().Unix()
		NodeListResponse.NodeList = nodeList
		global.Mutex.Unlock()

		time.Sleep(2 * time.Second)
	}
}
