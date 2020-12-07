package types

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"
)

// NewValidator generates new Validator.
func NewValidator(moniker string, website string, social string,
	identity string, comission sdk.Dec, valKey sdk.ValAddress, pubKey crypto.PubKey) (Validator, error) {

	pkAny, err := codectypes.PackAny(pubKey)
	if err != nil {
		return Validator{}, err
	}

	v := Validator{
		Moniker:    moniker,
		Website:    website,
		Social:     social,
		Identity:   identity,
		Commission: comission,
		ValKey:     valKey,
		PubKey:     pkAny,
	}

	err = v.Validate()
	if err != nil {
		return v, err
	}

	return v, nil
}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (v Validator) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	var pubkey crypto.PubKey
	return unpacker.UnpackAny(v.PubKey, &pubkey)
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

func (v Validator) GetConsPubKey() crypto.PubKey {
	pk, ok := v.PubKey.GetCachedValue().(crypto.PubKey)
	if !ok {
		panic("invalid key")
	}

	return pk
}
