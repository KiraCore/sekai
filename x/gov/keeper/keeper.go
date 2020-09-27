package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	KeyProposalKeyIDPrefix = []byte{0x00}

	KeyPrefixPermissionsRegistry       = []byte{0x20}
	KeyPrefixCouncilorIdentityRegistry = []byte{0x21}
	KeyPrefixActors                    = []byte{0x22}
)

type Keeper struct {
	cdc      codec.BinaryMarshaler
	storeKey sdk.StoreKey
}

func NewKeeper(storeKey sdk.StoreKey, cdc codec.BinaryMarshaler) Keeper {
	return Keeper{cdc: cdc, storeKey: storeKey}
}
