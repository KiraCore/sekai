package tasks

import (
	"time"

	"github.com/KiraCore/sekai/INTERX/common"
	"github.com/KiraCore/sekai/INTERX/config"
	"github.com/KiraCore/sekai/INTERX/types"
)

var (
	NodeList   []types.NodeList
	Scanning   bool
	LastUpdate int64
)

func NodeDiscover(isLog bool) {
	for {
		Scanning = true
		uniqueAddresses := config.LoadUniqueAddresses()

		common.GetLogger().Info(config.Config.AddrBooks)
		common.GetLogger().Info("[node-discover] addresses = ", uniqueAddresses)

		for _, addr := range uniqueAddresses {
			ipAddr := addr.Addr.IP
			interxUrl := "http://" + ipAddr + ":" + config.Config.NodeDiscovery.DefaultInterxPort
			if config.Config.NodeDiscovery.UseHttps {
				interxUrl = "https://" + ipAddr + ":" + config.Config.NodeDiscovery.DefaultInterxPort
			}

			interxStatus := common.GetInterxStatus(interxUrl)

			if interxStatus == nil {

			}

			common.GetLogger().Info(interxUrl)
		}

		Scanning = false
		LastUpdate = time.Now().UTC().Unix()

		time.Sleep(2 * time.Second)
	}
}
