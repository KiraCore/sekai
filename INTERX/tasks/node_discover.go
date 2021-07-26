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
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	tmjson "github.com/tendermint/tendermint/libs/json"
	"github.com/tendermint/tendermint/libs/protoio"
	"github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/p2p/conn"
	tmp2p "github.com/tendermint/tendermint/proto/tendermint/p2p"
	tmTypes "github.com/tendermint/tendermint/rpc/core/types"
	tmJsonRPCTypes "github.com/tendermint/tendermint/rpc/jsonrpc/types"
)

var (
	PubP2PNodeListResponse    types.P2PNodeListResponse
	PrivP2PNodeListResponse   types.P2PNodeListResponse
	InterxP2PNodeListResponse types.InterxNodeListResponse
	SnapNodeListResponse      types.SnapNodeListResponse
)

func timeout() time.Duration {
	timeoutDuration, err := time.ParseDuration(config.Config.NodeDiscovery.ConnectionTimeout)

	if err != nil {
		return 3 * time.Second
	}

	return timeoutDuration
}

func getTendermintRPCAddress(ipAddr string) string {
	return "http://" + ipAddr + ":" + config.Config.NodeDiscovery.DefaultTendermintPort
}

func getInterxAddress(ipAddr string) string {
	return "http://" + ipAddr + ":" + config.Config.NodeDiscovery.DefaultInterxPort
}

func QueryPeers(rpcAddr string) ([]tmTypes.Peer, error) {
	peers := []tmTypes.Peer{}

	u, err := url.Parse(rpcAddr)
	_, err = net.DialTimeout("tcp", u.Host, timeout())
	if err != nil {
		common.GetLogger().Info(err)
		return peers, err
	}

	endpoint := fmt.Sprintf("%s/net_info", rpcAddr)

	client := http.Client{
		Timeout: timeout(),
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

			kiraStatus, err := QueryKiraStatus(getTendermintRPCAddress(ipAddr))
			if err != nil {
				continue
			}

			nodeInfo := types.P2PNode{}
			nodeInfo.ID = string(kiraStatus.NodeInfo.ID())
			nodeInfo.IP = ipAddr
			nodeInfo.Port = getPort(kiraStatus.NodeInfo.ListenAddr)
			nodeInfo.Peers = []string{}

			// verify p2p node_id via p2p connect
			peerNodeInfo, ping := connect(p2p.NewNetAddressIPPort(parseIP(nodeInfo.IP), uint16(nodeInfo.Port)), timeout())
			if peerNodeInfo == nil || nodeInfo.ID != string(peerNodeInfo.ID()) {
				continue
			}

			nodeInfo.Ping = ping

			if _, ok := isLocalPeer[nodeInfo.ID]; ok {
				nodeInfo.Connected = true
			}

			if _, ok := peersFromIP[ipAddr]; !ok {
				peersFromIP[ipAddr], _ = QueryPeers(getTendermintRPCAddress(ipAddr))
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

			interxStartTime := makeTimestamp()
			// interxStatus := common.GetInterxStatus(getInterxAddress(ipAddr))
			interxStatus := common.GetInterxStatus(getInterxAddress("127.0.0.1"))

			if interxStatus != nil {
				interxEndTime := makeTimestamp()

				interxInfo := types.InterxNode{}
				interxInfo.ID = interxStatus.ID
				interxInfo.IP = ipAddr
				interxInfo.Ping = interxEndTime - interxStartTime
				interxInfo.Moniker = interxStatus.InterxInfo.Moniker
				interxInfo.Faucet = interxStatus.InterxInfo.FaucetAddr
				interxInfo.Type = interxStatus.InterxInfo.Node.NodeType
				interxInfo.Version = interxStatus.InterxInfo.Version

				InterxP2PNodeListResponse.NodeList = append(InterxP2PNodeListResponse.NodeList, interxInfo)
			}
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

func connect(
	netAddress *p2p.NetAddress,
	timeoutDuration time.Duration,
) (*p2p.DefaultNodeInfo, int64) {
	// load node_key
	privKey := ed25519.GenPrivKey()
	nodeKey := &p2p.NodeKey{
		PrivKey: privKey,
	}

	// manual handshaking
	// dial to address

	startTime := makeTimestamp()
	connection, err := netAddress.DialTimeout(timeoutDuration)
	endTime := makeTimestamp()
	connectionTime := endTime - startTime

	if err != nil {
		return nil, 0
	}

	// create secret connection
	startTime = makeTimestamp()
	secretConn, err := upgradeSecretConn(connection, timeoutDuration, nodeKey.PrivKey)
	endTime = makeTimestamp()
	if endTime-startTime > connectionTime {
		connectionTime = endTime - startTime
	}

	if err != nil {
		return nil, 0
	}

	// handshake

	startTime = makeTimestamp()
	peerNodeInfo, err := handshake(secretConn, timeoutDuration, p2p.DefaultNodeInfo{})
	endTime = makeTimestamp()
	if endTime-startTime > connectionTime {
		connectionTime = endTime - startTime
	}

	if err != nil {
		return nil, 0
	}

	return peerNodeInfo, connectionTime
}

func makeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func parseIP(host string) net.IP {
	ip := net.ParseIP(host)
	if ip == nil {
		ips, err := net.LookupIP(host)
		if err != nil {
			return nil
		}
		ip = ips[0]
	}

	return ip
}

func upgradeSecretConn(
	c net.Conn,
	timeout time.Duration,
	privKey crypto.PrivKey,
) (*conn.SecretConnection, error) {
	if err := c.SetDeadline(time.Now().Add(timeout)); err != nil {
		return nil, err
	}

	sc, err := conn.MakeSecretConnection(c, privKey)
	if err != nil {
		return nil, err
	}

	return sc, sc.SetDeadline(time.Time{})
}

func handshake(
	c net.Conn,
	timeout time.Duration,
	nodeInfo p2p.NodeInfo,
) (*p2p.DefaultNodeInfo, error) {
	if err := c.SetDeadline(time.Now().Add(timeout)); err != nil {
		return nil, err
	}

	var (
		errc = make(chan error, 2)

		pbpeerNodeInfo tmp2p.DefaultNodeInfo
		peerNodeInfo   p2p.DefaultNodeInfo
		ourNodeInfo    = nodeInfo.(p2p.DefaultNodeInfo)
	)

	go func(errc chan<- error, c net.Conn) {
		_, err := protoio.NewDelimitedWriter(c).WriteMsg(ourNodeInfo.ToProto())
		errc <- err
	}(errc, c)
	go func(errc chan<- error, c net.Conn) {
		protoReader := protoio.NewDelimitedReader(c, p2p.MaxNodeInfoSize())
		_, err := protoReader.ReadMsg(&pbpeerNodeInfo)
		errc <- err
	}(errc, c)

	for i := 0; i < cap(errc); i++ {
		err := <-errc
		if err != nil {
			return nil, err
		}
	}

	peerNodeInfo, err := p2p.DefaultNodeInfoFromToProto(&pbpeerNodeInfo)
	if err != nil {
		return nil, err
	}

	return &peerNodeInfo, c.SetDeadline(time.Time{})
}
