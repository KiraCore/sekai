package types

import (
	"errors"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.Msg = &MsgUpsertSignerKey{}

// NewMsgUpsertSignerKey create a new message to upsert signer key
func NewMsgUpsertSignerKey(
	pubKey string,
	keyType SignerKeyType,
	expiryTime int64,
	enabled bool,
	Data string,
	permissions []int64,
	curator sdk.AccAddress,
) (*MsgUpsertSignerKey, error) {
	return &MsgUpsertSignerKey{
		PubKey:      pubKey,
		KeyType:     keyType,
		ExpiryTime:  expiryTime,
		Enabled:     enabled,
		Permissions: permissions,
		Curator:     curator,
	}, nil
}

// Route returns module route to find appropriate handler
func (message MsgUpsertSignerKey) Route() string { return ModuleName }

// Type returns message type to differentiate with other messages on amino codec
func (message MsgUpsertSignerKey) Type() string { return UpsertSignerKeyTransaction }

// ValidateBasic returns basic validation error
func (message MsgUpsertSignerKey) ValidateBasic() error {
	// TODO: validate pubkey encoding by key type
	// TODO: validate permissions set

	if time.Now().Unix() > message.ExpiryTime {
		return errors.New("expiry time is invalid: now > expiryTime")
	}

	// validate if curator not empty
	if message.Curator.Empty() {
		return errors.New("curator shouldn't be empty")
	}
	return nil
}

// GetSignBytes return sorted sign bytes
func (message MsgUpsertSignerKey) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(message))
}

// GetSigners returns signer to sign the message
func (message MsgUpsertSignerKey) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{message.Curator}
}
