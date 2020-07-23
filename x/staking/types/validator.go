package types

import (
	sdk "github.com/KiraCore/cosmos-sdk/types"
)

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
func NewValidator(moniker string, website string, social string, identity string, comission sdk.Dec, valKey sdk.ValAddress, pubKey sdk.AccAddress) (Validator, error) {
	v := Validator{
		Moniker:   moniker,
		Website:   website,
		Social:    social,
		Identity:  identity,
		Comission: comission,
		ValKey:    valKey,
		PubKey:    pubKey,
	}

	err := v.Validate()
	if err != nil {
		return v, err
	}

	return v, nil
}

// Validate validates if a validator is correct.
func (v Validator) Validate() error {
	if len(v.Moniker) > 64 {
		return ErrInvalidMonikerLength
	}

	if len(v.Website) > 64 {
		return ErrInvalidWebsiteLength
	}

	if len(v.Social) > 64 {
		return ErrInvalidSocialLength
	}

	if len(v.Identity) > 64 {
		return ErrInvalidIdentityLength
	}

	return nil
}

type ValidatorIdentityRegistry struct {
}
