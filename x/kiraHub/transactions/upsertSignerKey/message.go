package signerkey

import (
	"errors"

	sdk "github.com/KiraCore/cosmos-sdk/types"
	"github.com/KiraCore/sekai/types"
	constants "github.com/KiraCore/sekai/x/kiraHub/constants"
)

// Message defines parameters that user put to create this message
type Message struct {
	PubKey      string              `json:"pubkey" yaml:"pubkey" valid:"required~PubKey is required"`
	KeyType     types.SignerKeyType `json:"type" yaml:"type" valid:"required~Type is required"`
	ExpiryTime  int64               `json:"expires" yaml:"expires" valid:"required~Expires is required"`
	Enabled     bool                `json:"enabled" yaml:"enabled" valid:"required~Enabled is required"`
	Permissions []int               `json:"permissions" yaml:"permissions" valid:"required~Permissions is required"`
	Curator     sdk.AccAddress      `json:"curator"  yaml:"curator" valid:"required~Curator is required"`
}

var _ sdk.Msg = Message{}

// Route returns module route to find appropriate handler
func (message Message) Route() string { return constants.ModuleName }

// Type returns message type to differentiate with other messages on amino codec
func (message Message) Type() string { return constants.UpsertSignerKeyTransaction }

// ValidateBasic returns basic validation error
func (message Message) ValidateBasic() error {
	// TODO: validate pubkey encoding by key type
	// TODO: validate expiry time and permissions set

	// validate if curator not empty
	if message.Curator.Empty() {
		return errors.New("curator shouldn't be empty")
	}
	return nil
}

// GetSignBytes return sorted sign bytes
func (message Message) GetSignBytes() []byte {
	return sdk.MustSortJSON(PackageCodec.MustMarshalJSON(message))
}

// GetSigners returns signer to sign the message
func (message Message) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{message.Curator}
}
