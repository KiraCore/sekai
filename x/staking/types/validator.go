package types

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
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

// GetConsPubKey returns the validator PubKey as a cryptotypes.PubKey.
func (v Validator) GetConsPubKey() crypto.PubKey {
	pk, ok := v.PubKey.GetCachedValue().(crypto.PubKey)
	if !ok {
		panic("invalid key")
	}

	return pk
}

// GetConsAddr extracts Consensus key address
func (v Validator) GetConsAddr() sdk.ConsAddress {
	return sdk.ConsAddress(v.GetConsPubKey().Address())
}

// IsInactivated returns if validator is inactivated
func (v Validator) IsInactivated() bool {
	return v.Inactivated
}

// IsPaused returns if validator is paused
func (v Validator) IsPaused() bool {
	return v.Paused
}

// TmConsPubKey casts Validator.ConsensusPubkey to crypto.PubKey
func (v Validator) TmConsPubKey() (crypto.PubKey, error) {
	pk, ok := v.PubKey.GetCachedValue().(crypto.PubKey)
	if !ok {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidType, "Expecting crypto.PubKey, got %T", pk)
	}

	// The way things are refactored now, v.ConsensusPubkey is sometimes a TM
	// ed25519 pubkey, sometimes our own ed25519 pubkey. This is very ugly and
	// inconsistent.
	// Luckily, here we coerce it into a TM ed25519 pubkey always, as this
	// pubkey will be passed into TM (eg calling encoding.PubKeyToProto).
	if intoTmPk, ok := pk.(types.IntoTmPubKey); ok {
		return intoTmPk.AsTmPubKey(), nil
	}

	return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidPubKey, "Logic error: ConsensusPubkey must be an SDK key and SDK PubKey types must be convertible to tendermint PubKey; got: %T", pk)
}
