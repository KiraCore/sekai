package handlers

import (
	"github.com/KiraCore/sekai/x/kiraHub/keeper"
	"github.com/KiraCore/sekai/x/kiraHub/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// HandlerMsgUpsertSignerKey handles upsertSignerKey message
func HandlerMsgUpsertSignerKey(context sdk.Context, keeper keeper.Keeper, message *types.MsgUpsertSignerKey) (*sdk.Result, error) {
	err := keeper.UpsertSignerKey(context, message.PubKey, message.KeyType, message.ExpiryTime, message.Permissions, message.Curator)
	return &sdk.Result{}, err
}
