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
}

type InterxNode struct {
	ID      string `json:"id"`
	IP      string `json:"ip"`
	Ping    int64  `json:"ping"`
	Moniker string `json:"moniker"`
	Faucet  string `json:"faucet"`
	Type    string `json:"type"`
	Version string `json:"version"`
}

type SnapNode struct {
	IP       string `json:"ip"`
	Port     uint16 `json:"port"`
	Size     int    `json:"size"`
	Checksum string `json:"checksum"`
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
