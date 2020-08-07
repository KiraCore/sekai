package types

import (
	"fmt"

	sdk "github.com/KiraCore/cosmos-sdk/types"
)

var _ sdk.Msg = &MsgClaimValidator{}

func NewMsgClaimValidator(
	moniker string,
	website string,
	social string,
	identity string,
	comission sdk.Dec,
	valKey sdk.ValAddress,
	pubKey sdk.AccAddress,
) (*MsgClaimValidator, error) {
	if valKey == nil {
		return nil, fmt.Errorf("validator not set")
	}

	if pubKey == nil {
		return nil, fmt.Errorf("public key not set")
	}

	return &MsgClaimValidator{
		Moniker:    moniker,
		Website:    website,
		Social:     social,
		Identity:   identity,
		Commission: comission,
		ValKey:     valKey,
		PubKey:     pubKey,
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
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m MsgClaimValidator) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{
		sdk.AccAddress(m.ValKey),
	}
}
