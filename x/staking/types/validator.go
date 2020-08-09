package types

import (
	sdk "github.com/KiraCore/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"
)

// NewValidator generates new Validator.
func NewValidator(moniker string, website string, social string, identity string, comission sdk.Dec, valKey sdk.ValAddress, pubKey crypto.PubKey) (Validator, error) {
	var pkStr string
	if pubKey != nil {
		pkStr = sdk.MustBech32ifyPubKey(sdk.Bech32PubKeyTypeConsPub, pubKey)
	}

	v := Validator{
		Moniker:    moniker,
		Website:    website,
		Social:     social,
		Identity:   identity,
		Commission: comission,
		ValKey:     valKey,
		PubKey:     pkStr,
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
