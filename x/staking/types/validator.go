package types

import sdk "github.com/KiraCore/cosmos-sdk/types"

type Validator struct {
	Moniker  string
	Website  string
	Social   string
	Identity string

	Comission sdk.Dec

	ValKey sdk.ValAddress
	PubKey sdk.AccAddress
}

// NewValidator generates new Validator.
func NewValidator(moniker string, website string, social string, identity string, comission sdk.Dec, valKey sdk.ValAddress, pubKey sdk.AccAddress) (*Validator, error) {
	return &Validator{Moniker: moniker, Website: website, Social: social, Identity: identity, Comission: comission, ValKey: valKey, PubKey: pubKey}, nil
}

type ValidatorIdentityRegistry struct {
}
