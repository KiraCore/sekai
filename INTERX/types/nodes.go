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
type NodePeer struct {
	ID   string `json:"id"`
	Ping int64  `json:"ping"`
}

type NodeList struct {
	ID        string     `json:"id"`
	IP        string     `json:"ip"`
	Moniker   string     `json:"moniker"`
	KiraAddr  string     `json:"kira_addr"`
	Version   string     `json:"version"`
	Seed      bool       `json:"seed"`
	Validator bool       `json:"validator"`
	Peers     []NodePeer `json:"peers"`
}

type NodeListResponse struct {
	LastUpdate int64      `json:"last_update"`
	Scanning   bool       `json:"scanning"`
	NodeList   []NodeList `json:"node_list"`
}
