package types

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	tmprotocrypto "github.com/tendermint/tendermint/proto/tendermint/crypto"
)

// NewValidator generates new Validator.
func NewValidator(valKey sdk.ValAddress, pubKey cryptotypes.PubKey) (Validator, error) {

	pkAny, err := codectypes.NewAnyWithValue(pubKey)
	if err != nil {
		return Validator{}, err
	}

	v := Validator{
		ValKey: valKey,
		PubKey: pkAny,
		Status: Active,
	}

	err = v.Validate()
	if err != nil {
		return v, err
	}

	return v, nil
}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (v Validator) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	var pubkey cryptotypes.PubKey
	return unpacker.UnpackAny(v.PubKey, &pubkey)
}

// Validate validates if a validator is correct.
func (v Validator) Validate() error {
	return nil
}

// GetConsAddr extracts Consensus key address
func (v Validator) GetConsAddr() sdk.ConsAddress {
	return sdk.ConsAddress(v.GetConsPubKey().Address())
}

// IsInactivated returns if validator is inactivated
func (v Validator) IsInactivated() bool {
	return v.Status == Inactive
}

// IsPaused returns if validator is paused
func (v Validator) IsPaused() bool {
	return v.Status == Paused
}

// IsActive returns if validator is active
func (v Validator) IsActive() bool {
	return v.Status == Active
}

// IsJailed returns if validator is jailed
func (v Validator) IsJailed() bool {
	return v.Status == Jailed
}

// GetConsPubKey returns the validator PubKey as a cryptotypes.PubKey.
func (v Validator) GetConsPubKey() cryptotypes.PubKey {
	pk, ok := v.PubKey.GetCachedValue().(cryptotypes.PubKey)
	if !ok {
		panic("invalid key")
	}

	return pk
}

// ConsPubKey returns the validator PubKey as a cryptotypes.PubKey.
func (v Validator) ConsPubKey() (cryptotypes.PubKey, error) {
	pk, ok := v.PubKey.GetCachedValue().(cryptotypes.PubKey)
	if !ok {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidType, "expecting cryptotypes.PubKey, got %T", pk)
	}

	return pk, nil
}

// TmConsPubKey casts Validator.ConsensusPubkey to crypto.PubKey
func (v Validator) TmConsPubKey() (tmprotocrypto.PublicKey, error) {
	pk, err := v.ConsPubKey()
	if err != nil {
		return tmprotocrypto.PublicKey{}, err
	}

	tmPk, err := cryptocodec.ToTmProtoPublicKey(pk)
	if err != nil {
		return tmprotocrypto.PublicKey{}, err
	}

	return tmPk, nil
}

// ConsensusPower gets the consensus-engine power. Aa reduction of 10^6 from
// validator tokens is applied
func (v Validator) ConsensusPower() int64 {
	if v.IsActive() {
		return 1
	}

	return 0
}
