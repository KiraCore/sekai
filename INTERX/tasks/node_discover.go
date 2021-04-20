package tasks

import (
	"time"

	"github.com/KiraCore/sekai/INTERX/common"
	"github.com/KiraCore/sekai/INTERX/config"
	"github.com/KiraCore/sekai/INTERX/types"
)

var (
	NodeList []types.KnownAddress
)

func NodeDiscover(isLog bool) {
	for {
		uniqueAddresses := config.LoadUniqueAddresses()
		NodeList = uniqueAddresses
		common.GetLogger().Info(config.Config.AddrBooks)

		common.GetLogger().Info("[node-discover] addresses = ", uniqueAddresses)

		time.Sleep(2 * time.Second)
	}
}
