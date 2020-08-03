package types

import (
	"fmt"

	sdk "github.com/KiraCore/cosmos-sdk/types"
)

var _ sdk.Msg = MsgClaimValidator{}

type MsgClaimValidator struct {
	Moniker   string // 64 chars max
	Website   string // 64 chars max
	Social    string // 64 chars max
	Identity  string // 64 chars max
	Comission sdk.Dec
	ValKey    sdk.ValAddress
	PubKey    sdk.AccAddress
}

func NewMsgClaimValidator(
	moniker string,
	website string,
	social string,
	identity string,
	comission sdk.Dec,
	valKey sdk.ValAddress,
	pubKey sdk.AccAddress,
) (MsgClaimValidator, error) {
	if valKey == nil {
		return MsgClaimValidator{}, fmt.Errorf("validator not set")
	}

	if pubKey == nil {
		return MsgClaimValidator{}, fmt.Errorf("public key not set")
	}

	return MsgClaimValidator{
		Moniker:   moniker,
		Website:   website,
		Social:    social,
		Identity:  identity,
		Comission: comission,
		ValKey:    valKey,
		PubKey:    pubKey,
	}, nil
}

func (m MsgClaimValidator) Route() string {
	return ModuleName
}

func (m MsgClaimValidator) Type() string {
	return ClaimValidator
}

func (m MsgClaimValidator) ValidateBasic() error {
	return nil
}

func (m MsgClaimValidator) GetSignBytes() []byte {
	panic("implement me")
}

func (m MsgClaimValidator) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{
		sdk.AccAddress(m.ValKey),
	}
}

func (m MsgClaimValidator) Reset() {
	panic("implement me")
}

func (m MsgClaimValidator) String() string {
	panic("implement me")
}

func (m MsgClaimValidator) ProtoMessage() {
	panic("implement me")
}
