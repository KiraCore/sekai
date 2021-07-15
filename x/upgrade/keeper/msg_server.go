package keeper

import (
	"github.com/KiraCore/sekai/x/upgrade/types"
)

type msgServer struct {
	keeper Keeper
	cgk    types.CustomGovKeeper
}

// NewMsgServerImpl returns an implementation of the upgrade MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper, cgk types.CustomGovKeeper) types.MsgServer {
	return &msgServer{
		keeper: keeper,
		cgk:    cgk,
	}
}
