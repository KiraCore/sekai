package handlers

import (
	sdk "github.com/KiraCore/cosmos-sdk/types"
	"github.com/KiraCore/sekai/x/kiraHub/keeper"
	"github.com/KiraCore/sekai/x/kiraHub/types"
)

// HandlerMsgUpsertSignerKey handles upsertSignerKey message
func HandlerMsgUpsertSignerKey(context sdk.Context, keeper keeper.Keeper, message *types.MsgUpsertSignerKey) (*sdk.Result, error) {
	err := keeper.UpsertSignerKey(context, message.PubKey, message.KeyType, message.ExpiryTime, message.Permissions, message.Curator)
	return &sdk.Result{}, err
}
