package signerkey

import sdk "github.com/KiraCore/cosmos-sdk/types"

// HandleMessage handles upsertSignerKey message
func HandleMessage(context sdk.Context, keeper Keeper, message Message) (*sdk.Result, error) {
	err := keeper.UpsertSignerKey(context, message.PubKey, message.KeyType, message.Permissions, message.Curator)
	return &sdk.Result{}, err
}
