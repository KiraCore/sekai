package types

import sdk "github.com/KiraCore/cosmos-sdk/types"

type MsgClaimValidator struct {
	Moniker   string // 64 chars max
	Website   string // 64 chars max
	Social    string // 64 chars max
	Identity  string // 64 chars max
	Comission sdk.Dec
	ValKey    sdk.ValAddress
	PubKey    sdk.AccAddress
}

func (m MsgClaimValidator) Route() string {
	panic("implement me")
}

func (m MsgClaimValidator) Type() string {
	panic("implement me")
}

func (m MsgClaimValidator) ValidateBasic() error {
	panic("implement me")
}

func (m MsgClaimValidator) GetSignBytes() []byte {
	panic("implement me")
}

func (m MsgClaimValidator) GetSigners() []sdk.AccAddress {
	panic("implement me")
}
