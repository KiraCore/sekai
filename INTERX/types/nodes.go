package types

import (
	"time"
)

const (
	bucketTypeNew = 0x01
	bucketTypeOld = 0x02
)

type NetAddress struct {
	ID   string `json:"id"`
	IP   string `json:"ip"`
	Port uint16 `json:"port"`
}

type KnownAddress struct {
	Addr        NetAddress `json:"addr"`
	Src         NetAddress `json:"src"`
	Buckets     []int      `json:"buckets"`
	Attempts    int32      `json:"attempts"`
	BucketType  byte       `json:"bucket_type"`
	LastAttempt time.Time  `json:"last_attempt"`
	LastSuccess time.Time  `json:"last_success"`
	LastBanTime time.Time  `json:"last_ban_time"`
}

type AddrBookJSON struct {
	Key   string         `json:"key"`
	Addrs []KnownAddress `json:"addrs"`
}

type P2PNode struct {
	ID        string   `json:"id"`
	IP        string   `json:"ip"`
	Port      uint16   `json:"port"`
	Ping      int64    `json:"ping"`
	Connected bool     `json:"connected"`
	Peers     []string `json:"peers"`
	Alive     bool     `json:"alive"`
	Synced    bool     `json:"synced"`
}

type InterxNode struct {
	ID      string `json:"id"`
	IP      string `json:"ip"`
	Ping    int64  `json:"ping"`
	Moniker string `json:"moniker"`
	Faucet  string `json:"faucet"`
	Type    string `json:"type"`
	Version string `json:"version"`
	Alive   bool   `json:"alive"`
	Synced  bool   `json:"synced"`
}

type SnapNode struct {
	IP       string `json:"ip"`
	Port     uint16 `json:"port"`
	Size     int64  `json:"size"`
	Checksum string `json:"checksum"`
	Alive    bool   `json:"alive"`
	Synced   bool   `json:"synced"`
}

type P2PNodeListResponse struct {
	LastUpdate int64     `json:"last_update"`
	Scanning   bool      `json:"scanning"`
	NodeList   []P2PNode `json:"node_list"`
}

type InterxNodeListResponse struct {
	LastUpdate int64        `json:"last_update"`
	Scanning   bool         `json:"scanning"`
	NodeList   []InterxNode `json:"node_list"`
}

type SnapNodeListResponse struct {
	LastUpdate int64      `json:"last_update"`
	Scanning   bool       `json:"scanning"`
	NodeList   []SnapNode `json:"node_list"`
}

type P2PNodes []P2PNode

func (s P2PNodes) Len() int {
	return len(s)
}
func (s P2PNodes) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s P2PNodes) Less(i, j int) bool {
	if s[i].Connected != s[j].Connected {
		if s[j].Connected {
			return false
		}
		if s[i].Connected {
			return true
		}
	}

	if s[i].Ping != s[j].Ping {
		return s[i].Ping < s[j].Ping
	}

	return false
}

type InterxNodes []InterxNode

func (s InterxNodes) Len() int {
	return len(s)
}
func (s InterxNodes) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s InterxNodes) Less(i, j int) bool {
	if s[i].Ping != s[j].Ping {
		return s[i].Ping < s[j].Ping
	}

	return false
}

type SnapNodes []SnapNode

func (s SnapNodes) Len() int {
	return len(s)
}
func (s SnapNodes) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s SnapNodes) Less(i, j int) bool {
	if s[i].Size != s[j].Size {
		return s[i].Size > s[j].Size
	}

	return false
}
