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
