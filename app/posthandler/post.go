package posthandler

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	feeprocessingkeeper "github.com/KiraCore/sekai/x/feeprocessing/keeper"
	customgovkeeper "github.com/KiraCore/sekai/x/gov/keeper"
)

// HandlerOptions are the options required for constructing a default SDK PostHandler.
type HandlerOptions struct {
	customGovKeeper     customgovkeeper.Keeper
	feeprocessingKeeper feeprocessingkeeper.Keeper
}

// NewPostHandler returns an empty PostHandler chain.
func NewPostHandler(options HandlerOptions) (sdk.PostHandler, error) {
	postDecorators := []sdk.PostDecorator{
		NewExecutionDecorator(options.customGovKeeper, options.feeprocessingKeeper),
	}

	return sdk.ChainPostDecorators(postDecorators...), nil
}
