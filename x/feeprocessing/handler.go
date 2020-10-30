package feeprocessing

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"

	feeprocessingkeeper "github.com/KiraCore/sekai/x/feeprocessing/keeper"
	"github.com/KiraCore/sekai/x/feeprocessing/types"
)

// NewHandler handle custom messages
func NewHandler(fk feeprocessingkeeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		default:
			return nil, errors.Wrapf(errors.ErrUnknownRequest, "unrecognized %s message type: %T", types.ModuleName, msg)
		}
	}
}
