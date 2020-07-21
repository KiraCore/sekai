package types

type MsgClaimValidator struct {
	Moniker   string // 64 chars max
	Website   string // 64 chars max
	Social    string // 64 chars max
	Identity  string // 64 chars max
	Comission string
	ValKey    string
	PubKey    string
}
