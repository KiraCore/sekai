package keeper

import (
	"context"
	"github.com/KiraCore/sekai/x/custody/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
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

func (s msgServer) AddToWhiteList(goCtx context.Context, msg *types.MsgAddToCustodyWhiteList) (*types.MsgAddToCustodyWhiteListResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	record := types.CustodyWhiteListRecord{
		Address:          msg.Address,
		CustodyWhiteList: s.keeper.GetCustodyWhiteListByAddress(ctx, msg.Address),
	}

	if record.CustodyWhiteList == nil {
		record.CustodyWhiteList = new(types.CustodyWhiteList)
		record.CustodyWhiteList.Addresses = map[string]bool{}
	}

	record.CustodyWhiteList.Addresses[msg.AddAddress.String()] = true
	s.keeper.AddToCustodyWhiteList(ctx, record)

	return &types.MsgAddToCustodyWhiteListResponse{}, nil
}

func (s msgServer) RemoveFromWhiteList(goCtx context.Context, msg *types.MsgRemoveFromCustodyWhiteList) (*types.MsgRemoveFromCustodyWhiteListResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	record := types.CustodyWhiteListRecord{
		Address:          msg.Address,
		CustodyWhiteList: s.keeper.GetCustodyWhiteListByAddress(ctx, msg.Address),
	}

	if record.CustodyWhiteList == nil {
		return nil, errors.Wrap(types.ErrNoWhiteLists, "Can not remove from the empty whitelist")
	}

	if !record.CustodyWhiteList.Addresses[msg.RemoveAddress.String()] {
		return nil, errors.Wrap(types.ErrNoWhiteListsElement, "Can not remove missing element from the whitelist")
	}

	record.CustodyWhiteList.Addresses[msg.RemoveAddress.String()] = false
	s.keeper.AddToCustodyWhiteList(ctx, record)

	return &types.MsgRemoveFromCustodyWhiteListResponse{}, nil
}

func (s msgServer) DropWhiteList(goCtx context.Context, msg *types.MsgDropCustodyWhiteList) (*types.MsgDropCustodyWhiteListResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	s.keeper.DropCustodyWhiteListByAddress(ctx, msg.Address)

	return &types.MsgDropCustodyWhiteListResponse{}, nil
}
