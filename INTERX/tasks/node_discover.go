package tasks

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/KiraCore/sekai/INTERX/common"
	"github.com/KiraCore/sekai/INTERX/config"
	"github.com/KiraCore/sekai/INTERX/types"
	tmjson "github.com/tendermint/tendermint/libs/json"
	tmTypes "github.com/tendermint/tendermint/rpc/core/types"
	tmJsonRPCTypes "github.com/tendermint/tendermint/rpc/jsonrpc/types"
)

var (
	PubP2PNodeListResponse    types.P2PNodeListResponse
	PrivP2PNodeListResponse   types.P2PNodeListResponse
	InterxP2PNodeListResponse types.InterxNodeListResponse
	SnapNodeListResponse      types.SnapNodeListResponse
)

const TENDERMINT_PORT = "26657"
const TIMEOUT = 3 * time.Second

func getRPCAddress(ipAddr string) string {
	return "http://" + ipAddr + ":" + TENDERMINT_PORT
}

func QueryPeers(rpcAddr string) ([]tmTypes.Peer, error) {
	peers := []tmTypes.Peer{}

	u, err := url.Parse(rpcAddr)
	_, err = net.DialTimeout("tcp", u.Host, TIMEOUT)
	if err != nil {
		common.GetLogger().Info(err)
		return peers, err
	}

	endpoint := fmt.Sprintf("%s/net_info", rpcAddr)

	client := http.Client{
		Timeout: TIMEOUT,
	}

	resp, err := client.Get(endpoint)
	if err != nil {
		return peers, err
	}
	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)

	response := new(tmJsonRPCTypes.RPCResponse)

	if err := json.Unmarshal(respBody, response); err != nil {
		return peers, err
	}

	if response.Error != nil {
		return peers, errors.New(fmt.Sprint(response.Error))
	}

	result := new(tmTypes.ResultNetInfo)
	if err := tmjson.Unmarshal(response.Result, result); err != nil {
		return peers, err
	}

	peers = result.Peers

	return peers, nil
}

func QueryKiraStatus(rpcAddr string) (tmTypes.ResultStatus, error) {
	result := tmTypes.ResultStatus{}

	endpoint := fmt.Sprintf("%s/status", rpcAddr)

	resp, err := http.Get(endpoint)
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)

	response := new(tmJsonRPCTypes.RPCResponse)

	if err := json.Unmarshal(respBody, response); err != nil {
		return result, err
	}

	if response.Error != nil {
		return result, errors.New(fmt.Sprint(response.Error))
	}

	if err := tmjson.Unmarshal(response.Result, &result); err != nil {
		return result, err
	}

	return result, nil
}

func getP2PNodeInfo(ipAddr string) {

}

func NodeDiscover(rpcAddr string, isLog bool) {
	initPrivateIps()

	PubP2PNodeListResponse.Scanning = true
	PrivP2PNodeListResponse.Scanning = true
	InterxP2PNodeListResponse.Scanning = true
	SnapNodeListResponse.Scanning = true

	for {
		pubP2PNodes := []types.P2PNode{}
		privP2PNodes := []types.P2PNode{}
		interxNodes := []types.InterxNode{}
		snapNodes := []types.SnapNode{}

		isQueriedIP := make(map[string]bool) // check if ip is already queried
		// isQueriedID := make(map[string]bool) // check if node id is already queried

		uniqueIPAddresses := config.LoadUniqueIPAddresses()

		isLocalPeer := make(map[string]bool)
		localPeers, _ := QueryPeers(rpcAddr)
		// localKiraStatus, _ := QueryKiraStatus(rpcAddr)
		for _, peer := range localPeers {
			isLocalPeer[string(peer.NodeInfo.ID())] = true
		}

		peersFromIP := make(map[string]([]tmTypes.Peer))
		// kiraStatusFromIP := make(map[string](tmTypes.ResultStatus))

		index := 0
		for index < len(uniqueIPAddresses) {
			ipAddr := uniqueIPAddresses[index]
			index++

			if _, ok := isQueriedIP[ipAddr]; ok {
				continue
			}
			isQueriedIP[ipAddr] = true

			common.GetLogger().Info(ipAddr)

			if _, ok := peersFromIP[ipAddr]; !ok {
				peersFromIP[ipAddr], _ = QueryPeers(getRPCAddress(ipAddr))
			} else {
				continue
			}

			// if _, ok := kiraStatusFromIP[ipAddr]; !ok {
			// 	kiraStatusFromIP[ipAddr], _ = QueryKiraStatus(getRPCAddress(ipAddr))
			// }

			peers := peersFromIP[ipAddr]
			// kiraStatus := kiraStatusFromIP[ipAddr]

			for _, peer := range peers {
				isPrivate := isPrivateIP(peer.RemoteIP)

				u, err := url.Parse(peer.NodeInfo.ListenAddr)
				portNumber, err := strconv.ParseUint(u.Port(), 10, 16)

				if err != nil {
					common.GetLogger().Info(peer.NodeInfo.ListenAddr, err)
				}

				nodeInfo := types.P2PNode{}
				nodeInfo.ID = string(peer.NodeInfo.ID())
				nodeInfo.IP = peer.RemoteIP
				nodeInfo.Port = uint16(portNumber)

				if _, ok := isLocalPeer[nodeInfo.ID]; ok {
					nodeInfo.Connected = true
				}

				// ping & check out node_id

				if isPrivate {
					privP2PNodes = append(privP2PNodes, nodeInfo)
				} else {
					pubP2PNodes = append(pubP2PNodes, nodeInfo)

					if _, ok := isQueriedIP[nodeInfo.IP]; !ok {
						uniqueIPAddresses = append(uniqueIPAddresses, nodeInfo.IP)
					}
				}
			}
		}

		lastUpdate := time.Now().Unix()

		PubP2PNodeListResponse.NodeList = pubP2PNodes
		PrivP2PNodeListResponse.NodeList = privP2PNodes
		InterxP2PNodeListResponse.NodeList = interxNodes
		SnapNodeListResponse.NodeList = snapNodes

		PubP2PNodeListResponse.LastUpdate = lastUpdate
		PrivP2PNodeListResponse.LastUpdate = lastUpdate
		InterxP2PNodeListResponse.LastUpdate = lastUpdate
		SnapNodeListResponse.LastUpdate = lastUpdate

		PubP2PNodeListResponse.Scanning = false
		PrivP2PNodeListResponse.Scanning = false
		InterxP2PNodeListResponse.Scanning = false
		SnapNodeListResponse.Scanning = false
		time.Sleep(30 * time.Minute)
	}
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

var privateIPBlocks []*net.IPNet

func initPrivateIps() {
	for _, cidr := range []string{
		"127.0.0.0/8",    // IPv4 loopback
		"10.0.0.0/8",     // RFC1918
		"172.16.0.0/12",  // RFC1918
		"192.168.0.0/16", // RFC1918
		"169.254.0.0/16", // RFC3927 link-local
		"::1/128",        // IPv6 loopback
		"fe80::/10",      // IPv6 link-local
		"fc00::/7",       // IPv6 unique local addr
	} {
		_, block, err := net.ParseCIDR(cidr)
		if err != nil {
			panic(fmt.Errorf("parse error on %q: %v", cidr, err))
		}
		privateIPBlocks = append(privateIPBlocks, block)
	}
}

func isPrivateIP(ipAddr string) bool {
	ip := net.ParseIP(ipAddr)

	if ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() {
		return true
	}

	for _, block := range privateIPBlocks {
		if block.Contains(ip) {
			return true
		}
	}
	return false
}
