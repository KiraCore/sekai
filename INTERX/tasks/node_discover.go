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

func NodeDiscover(rpcAddr string, isLog bool) {
	initPrivateIps()

	for {
		PubP2PNodeListResponse.Scanning = true
		PrivP2PNodeListResponse.Scanning = true
		InterxP2PNodeListResponse.Scanning = true
		SnapNodeListResponse.Scanning = true

		PubP2PNodeListResponse.NodeList = []types.P2PNode{}
		PrivP2PNodeListResponse.NodeList = []types.P2PNode{}
		InterxP2PNodeListResponse.NodeList = []types.InterxNode{}
		SnapNodeListResponse.NodeList = []types.SnapNode{}

		isIpInList := make(map[string]bool)   // check if ip is already queried
		isPrivNodeID := make(map[string]bool) // check if ip is already queried

		uniqueIPAddresses := config.LoadUniqueIPAddresses()

		for i := 0; i < len(uniqueIPAddresses); i++ {
			isIpInList[uniqueIPAddresses[i]] = true
		}

		isLocalPeer := make(map[string]bool)
		localPeers, _ := QueryPeers(rpcAddr)
		for _, peer := range localPeers {
			isLocalPeer[string(peer.NodeInfo.ID())] = true
		}

		peersFromIP := make(map[string]([]tmTypes.Peer))

		index := 0
		for index < len(uniqueIPAddresses) {
			ipAddr := uniqueIPAddresses[index]
			index++

			if isLog {
				common.GetLogger().Info("[node-discovery] ", ipAddr)
			}

			kiraStatus, err := QueryKiraStatus(getRPCAddress(ipAddr))
			if err != nil {
				continue
			}

			nodeInfo := types.P2PNode{}
			nodeInfo.ID = string(kiraStatus.NodeInfo.ID())
			nodeInfo.IP = ipAddr
			nodeInfo.Port = getPort(kiraStatus.NodeInfo.ListenAddr)
			nodeInfo.Peers = []string{}

			if _, ok := isLocalPeer[nodeInfo.ID]; ok {
				nodeInfo.Connected = true
			}

			if _, ok := peersFromIP[ipAddr]; !ok {
				peersFromIP[ipAddr], _ = QueryPeers(getRPCAddress(ipAddr))
			}

			peers := peersFromIP[ipAddr]
			for _, peer := range peers {
				nodeInfo.Peers = append(nodeInfo.Peers, string(peer.NodeInfo.ID()))

				if isPrivateIP(peer.RemoteIP) {
					privNodeInfo := types.P2PNode{}
					privNodeInfo.ID = string(peer.NodeInfo.ID())
					privNodeInfo.IP = peer.RemoteIP
					privNodeInfo.Port = getPort(peer.NodeInfo.ListenAddr)
					privNodeInfo.Peers = []string{}
					privNodeInfo.Peers = append(privNodeInfo.Peers, nodeInfo.ID)

					if _, ok := isLocalPeer[privNodeInfo.ID]; ok {
						privNodeInfo.Connected = true
					}

					if _, ok := isPrivNodeID[privNodeInfo.ID]; !ok {
						PrivP2PNodeListResponse.NodeList = append(PrivP2PNodeListResponse.NodeList, privNodeInfo)
						isPrivNodeID[privNodeInfo.ID] = true
					}
				} else {
					if _, ok := isIpInList[peer.RemoteIP]; !ok {
						uniqueIPAddresses = append(uniqueIPAddresses, peer.RemoteIP)
						isIpInList[peer.RemoteIP] = true
					}
				}
			}

			PubP2PNodeListResponse.NodeList = append(PubP2PNodeListResponse.NodeList, nodeInfo)
		}

		lastUpdate := time.Now().Unix()

		PubP2PNodeListResponse.LastUpdate = lastUpdate
		PrivP2PNodeListResponse.LastUpdate = lastUpdate
		InterxP2PNodeListResponse.LastUpdate = lastUpdate
		SnapNodeListResponse.LastUpdate = lastUpdate

		PubP2PNodeListResponse.Scanning = false
		PrivP2PNodeListResponse.Scanning = false
		InterxP2PNodeListResponse.Scanning = false
		SnapNodeListResponse.Scanning = false

		if isLog {
			common.GetLogger().Info("[node-discovery] finished!")
		}

		time.Sleep(24 * time.Hour)
	}
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

func getPort(listenAddr string) uint16 {
	u, _ := url.Parse(listenAddr)
	portNumber, _ := strconv.ParseUint(u.Port(), 10, 16)

	return uint16(portNumber)
}
