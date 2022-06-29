package keeper

import (
	"context"
	"github.com/KiraCore/sekai/x/custody/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type msgServer struct {
	keeper Keeper
}

// NewMsgServerImpl returns an implementation of the bank MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{
		keeper: keeper,
	}
}

var _ types.MsgServer = msgServer{}

func (s msgServer) CreateCustody(goCtx context.Context, msg *types.MsgCreteCustodyRecord) (*types.MsgCreteCustodyRecordResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	record := types.CustodyRecord{
		Address:         msg.Address,
		CustodySettings: msg.CustodySettings,
	}

	s.keeper.SetCustodyRecord(ctx, record)

	return &types.MsgCreteCustodyRecordResponse{}, nil
}
