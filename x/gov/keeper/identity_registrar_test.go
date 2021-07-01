package keeper_test

import (
	"testing"
	"time"

	"github.com/KiraCore/sekai/simapp"
	"github.com/KiraCore/sekai/x/gov/keeper"
	"github.com/KiraCore/sekai/x/gov/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

func TestKeeper_CheckIfWithinAddressArray(t *testing.T) {
	addr1 := sdk.AccAddress("foo1________________")
	addr2 := sdk.AccAddress("foo2________________")
	addr3 := sdk.AccAddress("foo3________________")
	addr4 := sdk.AccAddress("foo4________________")

	array := []sdk.AccAddress{addr1, addr2, addr3}

	require.True(t, keeper.CheckIfWithinAddressArray(addr1, array))
	require.True(t, keeper.CheckIfWithinAddressArray(addr2, array))
	require.True(t, keeper.CheckIfWithinAddressArray(addr3, array))
	require.False(t, keeper.CheckIfWithinAddressArray(addr4, array))
}

func TestKeeper_LastIdentityRecordId(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	lastRecordId := app.CustomGovKeeper.GetLastIdentityRecordId(ctx)
	require.Equal(t, lastRecordId, uint64(0))

	app.CustomGovKeeper.SetLastIdentityRecordId(ctx, 5)

	lastRecordId = app.CustomGovKeeper.GetLastIdentityRecordId(ctx)
	require.Equal(t, lastRecordId, uint64(5))
}

func TestKeeper_LastIdRecordVerifyRequestId(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	lastRecordId := app.CustomGovKeeper.GetLastIdRecordVerifyRequestId(ctx)
	require.Equal(t, lastRecordId, uint64(0))

	app.CustomGovKeeper.SetLastIdRecordVerifyRequestId(ctx, 5)

	lastRecordId = app.CustomGovKeeper.GetLastIdRecordVerifyRequestId(ctx)
	require.Equal(t, lastRecordId, uint64(5))
}

func TestKeeper_IdentityRecordBasicFlow(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	// try to get non existent record
	record := app.CustomGovKeeper.GetIdentityRecord(ctx, 1)
	require.Nil(t, record)

	// create a new record and check if set correctly
	addr1 := sdk.AccAddress("foo1________________")
	addr2 := sdk.AccAddress("foo2________________")
	addr3 := sdk.AccAddress("foo3________________")
	infos := make(map[string]string)
	infos["key"] = "value"
	newRecord := types.IdentityRecord{
		Id:        1,
		Address:   addr1,
		Infos:     infos,
		Date:      time.Now().UTC(),
		Verifiers: []sdk.AccAddress{addr2, addr3},
	}
	app.CustomGovKeeper.SetIdentityRecord(ctx, newRecord)
	record = app.CustomGovKeeper.GetIdentityRecord(ctx, 1)
	require.NotNil(t, record)
	require.Equal(t, *record, newRecord)

	// check no panics
	app.CustomGovKeeper.DeleteIdentityRecord(ctx, 0)

	// remove existing id and check
	app.CustomGovKeeper.DeleteIdentityRecord(ctx, 1)
	record = app.CustomGovKeeper.GetIdentityRecord(ctx, 1)
	require.Nil(t, record)
}

func TestKeeper_IdentityRecordAddEditRemove(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	// create a new record and check if set correctly
	addr1 := sdk.AccAddress("foo1________________")
	addr2 := sdk.AccAddress("foo2________________")
	addr3 := sdk.AccAddress("foo3________________")
	infos := make(map[string]string)
	infos["key"] = "value"
	now := time.Now().UTC()
	recordId := app.CustomGovKeeper.CreateIdentityRecord(ctx, addr1, infos, now)
	require.Equal(t, recordId, uint64(1))

	record := app.CustomGovKeeper.GetIdentityRecord(ctx, 1)
	require.NotNil(t, record)
	expectedRecord := types.IdentityRecord{
		Id:      1,
		Address: addr1,
		Infos:   infos,
		Date:    now,
	}
	require.Equal(t, *record, expectedRecord)

	recordId = app.CustomGovKeeper.CreateIdentityRecord(ctx, addr1, infos, now)
	require.Equal(t, recordId, uint64(2))
	recordId = app.CustomGovKeeper.CreateIdentityRecord(ctx, addr2, infos, now)
	require.Equal(t, recordId, uint64(3))

	records := app.CustomGovKeeper.GetAllIdentityRecords(ctx)
	require.Len(t, records, 3)
	records = app.CustomGovKeeper.GetIdRecordsByAddress(ctx, addr1)
	require.Len(t, records, 2)
	records = app.CustomGovKeeper.GetIdRecordsByAddress(ctx, addr2)
	require.Len(t, records, 1)

	// remove existing id and check
	app.CustomGovKeeper.DeleteIdentityRecord(ctx, 2)
	records = app.CustomGovKeeper.GetAllIdentityRecords(ctx)
	require.Len(t, records, 2)
	records = app.CustomGovKeeper.GetIdRecordsByAddress(ctx, addr1)
	require.Len(t, records, 1)
	records = app.CustomGovKeeper.GetIdRecordsByAddress(ctx, addr2)
	require.Len(t, records, 1)

	infos["key1"] = "value1"
	now = now.Add(time.Second)

	// try editing with other owner
	err := app.CustomGovKeeper.EditIdentityRecord(ctx, 1, addr3, infos, now)
	require.Error(t, err)

	// try editing deleted record
	err = app.CustomGovKeeper.EditIdentityRecord(ctx, 2, addr2, infos, now)
	require.Error(t, err)

	// set verifier of identity record
	app.CustomGovKeeper.SetIdentityRecord(ctx, types.IdentityRecord{
		Id:        1,
		Address:   addr1,
		Infos:     infos,
		Date:      now,
		Verifiers: []sdk.AccAddress{addr3},
	})

	err = app.CustomGovKeeper.EditIdentityRecord(ctx, 1, addr1, infos, now)
	require.NoError(t, err)
	record = app.CustomGovKeeper.GetIdentityRecord(ctx, 1)
	require.NotNil(t, record)
	require.Equal(t, *record, types.IdentityRecord{
		Id:      1,
		Address: addr1,
		Infos:   infos,
		Date:    now,
	})
	records = app.CustomGovKeeper.GetAllIdentityRecords(ctx)
	require.Len(t, records, 2)
	records = app.CustomGovKeeper.GetIdRecordsByAddress(ctx, addr1)
	require.Len(t, records, 1)
	records = app.CustomGovKeeper.GetIdRecordsByAddress(ctx, addr2)
	require.Len(t, records, 1)
}

func TestKeeper_IdentityRecordApproveFlow(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	// create a new record and check if set correctly
	addr1 := sdk.AccAddress("foo1________________")
	addr2 := sdk.AccAddress("foo2________________")
	addr3 := sdk.AccAddress("foo3________________")
	app.BankKeeper.SetBalance(ctx, addr1, sdk.NewInt64Coin(sdk.DefaultBondDenom, 10000))
	app.BankKeeper.SetBalance(ctx, addr2, sdk.NewInt64Coin(sdk.DefaultBondDenom, 10000))
	app.BankKeeper.SetBalance(ctx, addr3, sdk.NewInt64Coin(sdk.DefaultBondDenom, 10000))

	infos := make(map[string]string)
	infos["key"] = "value"
	now := time.Now().UTC()
	recordId := app.CustomGovKeeper.CreateIdentityRecord(ctx, addr1, infos, now)
	require.Equal(t, recordId, uint64(1))

	record := app.CustomGovKeeper.GetIdentityRecord(ctx, 1)
	require.NotNil(t, record)
	expectedRecord := types.IdentityRecord{
		Id:      1,
		Address: addr1,
		Infos:   infos,
		Date:    now,
	}
	require.Equal(t, *record, expectedRecord)
	recordId = app.CustomGovKeeper.CreateIdentityRecord(ctx, addr1, infos, now)
	require.Equal(t, recordId, uint64(2))
	recordId = app.CustomGovKeeper.CreateIdentityRecord(ctx, addr2, infos, now)
	require.Equal(t, recordId, uint64(3))

	// bigger tip than balance
	ctxCache, _ := ctx.CacheContext()
	reqId, err := app.CustomGovKeeper.RequestIdentityRecordsVerify(ctxCache, addr1, addr3, []uint64{1}, sdk.NewInt64Coin(sdk.DefaultBondDenom, 1000000000))
	require.Equal(t, reqId, uint64(0))
	require.Error(t, err)

	// request id record 1 to addr3 by addr1
	reqId, err = app.CustomGovKeeper.RequestIdentityRecordsVerify(ctx, addr1, addr3, []uint64{1}, sdk.NewInt64Coin(sdk.DefaultBondDenom, 10))
	require.Equal(t, reqId, uint64(1))
	require.NoError(t, err)
	err = app.CustomGovKeeper.ApproveIdentityRecords(ctx, addr3, 1)
	require.NoError(t, err)
	record = app.CustomGovKeeper.GetIdentityRecord(ctx, 1)
	require.NotNil(t, record)
	require.Equal(t, *record, types.IdentityRecord{
		Id:        1,
		Address:   addr1,
		Infos:     infos,
		Date:      now,
		Verifiers: []sdk.AccAddress{addr3},
	})

	// request id record 1 to addr3 by addr1 again
	reqId, err = app.CustomGovKeeper.RequestIdentityRecordsVerify(ctx, addr1, addr3, []uint64{1}, sdk.NewInt64Coin(sdk.DefaultBondDenom, 10))
	require.Equal(t, reqId, uint64(2))
	require.NoError(t, err)
	err = app.CustomGovKeeper.ApproveIdentityRecords(ctx, addr3, 2)
	require.NoError(t, err)
	record = app.CustomGovKeeper.GetIdentityRecord(ctx, 1)
	require.NotNil(t, record)
	require.Equal(t, *record, types.IdentityRecord{
		Id:        1,
		Address:   addr1,
		Infos:     infos,
		Date:      now,
		Verifiers: []sdk.AccAddress{addr3},
	})

	// request id record 1 and 3 to addr3 by addr1 again
	reqId, err = app.CustomGovKeeper.RequestIdentityRecordsVerify(ctx, addr1, addr3, []uint64{1, 3}, sdk.NewInt64Coin(sdk.DefaultBondDenom, 10))
	require.Equal(t, reqId, uint64(3))
	require.NoError(t, err)
	err = app.CustomGovKeeper.ApproveIdentityRecords(ctx, addr3, 3)
	require.NoError(t, err)
	record = app.CustomGovKeeper.GetIdentityRecord(ctx, 1)
	require.NotNil(t, record)
	require.Equal(t, *record, types.IdentityRecord{
		Id:        1,
		Address:   addr1,
		Infos:     infos,
		Date:      now,
		Verifiers: []sdk.AccAddress{addr3},
	})
	record = app.CustomGovKeeper.GetIdentityRecord(ctx, 3)
	require.NotNil(t, record)
	require.Equal(t, *record, types.IdentityRecord{
		Id:        3,
		Address:   addr2,
		Infos:     infos,
		Date:      now,
		Verifiers: []sdk.AccAddress{addr3},
	})

	// request id record 2 to addr3 by addr2
	reqId, err = app.CustomGovKeeper.RequestIdentityRecordsVerify(ctx, addr2, addr3, []uint64{2}, sdk.NewInt64Coin(sdk.DefaultBondDenom, 10))
	require.Equal(t, reqId, uint64(4))
	require.NoError(t, err)
	err = app.CustomGovKeeper.ApproveIdentityRecords(ctx, addr3, 4)
	require.NoError(t, err)
	record = app.CustomGovKeeper.GetIdentityRecord(ctx, 2)
	require.NotNil(t, record)
	require.Equal(t, *record, types.IdentityRecord{
		Id:        2,
		Address:   addr1,
		Infos:     infos,
		Date:      now,
		Verifiers: []sdk.AccAddress{addr3},
	})

	// request non-exist identity record
	ctxCache, _ = ctx.CacheContext()
	reqId, err = app.CustomGovKeeper.RequestIdentityRecordsVerify(ctxCache, addr2, addr3, []uint64{5}, sdk.NewInt64Coin(sdk.DefaultBondDenom, 10))
	require.Equal(t, reqId, uint64(5))
	require.Error(t, err)

	// approve with non-approver
	reqId, err = app.CustomGovKeeper.RequestIdentityRecordsVerify(ctx, addr1, addr3, []uint64{1}, sdk.NewInt64Coin(sdk.DefaultBondDenom, 10))
	require.Equal(t, reqId, uint64(5))
	require.NoError(t, err)
	ctxCache, _ = ctx.CacheContext()
	err = app.CustomGovKeeper.ApproveIdentityRecords(ctxCache, addr2, 5)
	require.Error(t, err)

	// approve not existing request id
	ctxCache, _ = ctx.CacheContext()
	err = app.CustomGovKeeper.ApproveIdentityRecords(ctxCache, addr2, 0xFFFFF)
	require.Error(t, err)

	// try edit and check if verification records all gone
	infos["key1"] = "value1"
	err = app.CustomGovKeeper.EditIdentityRecord(ctx, 1, addr1, infos, now)
	require.NoError(t, err)
	record = app.CustomGovKeeper.GetIdentityRecord(ctx, 1)
	require.NotNil(t, record)
	require.Equal(t, *record, types.IdentityRecord{
		Id:      1,
		Address: addr1,
		Infos:   infos,
		Date:    now,
	})

	// TODO: add msg validate basic test
	// reqId, err := app.CustomGovKeeper.RequestIdentityRecordsVerify(ctx, addr1, addr3, []uint64{1}, sdk.Coin{})
	// require.Equal(t, reqId, uint64(0))
	// require.NoError(t, err)
	// ctxCache, _ := ctx.CacheContext()
	// reqId, err = app.CustomGovKeeper.RequestIdentityRecordsVerify(ctxCache, addr1, addr3, []uint64{1}, sdk.NewInt64Coin(sdk.DefaultBondDenom, 0))
	// require.Equal(t, reqId, uint64(0))
	// require.Error(t, err)
	// TODO: add msg validate basic test
	// reqId, err = app.CustomGovKeeper.RequestIdentityRecordsVerify(ctx, addr1, addr3, []uint64{}, sdk.NewInt64Coin(sdk.DefaultBondDenom, 10))
	// require.Equal(t, reqId, uint64(2))
	// require.NoError(t, err)
	// TODO: add validate basic test for empty verifier
	// reqId, err = app.CustomGovKeeper.RequestIdentityRecordsVerify(ctx, addr2, sdk.AccAddress{}, []uint64{1}, sdk.NewInt64Coin(sdk.DefaultBondDenom, 10))
	// require.Equal(t, reqId, uint64(5))
	// require.NoError(t, err)

	// TODO: add coin balance movement check on request, approve, cancel process
	// TODO: ask non-tip request exist and if so, implement

	// TODO: check queries for verify requests
	// GetIdRecordsVerifyRequest
	// GetIdRecordsVerifyRequestsByRequester
	// GetIdRecordsVerifyRequestsByApprover
	// GetAllIdRecordsVerifyRequests

	// TODO: try deleting verifyRequest after creation
	// RequestIdentityRecordsVerify
	// DeleteIdRecordsVerifyRequest
	// CancelIdentityRecordsVerifyRequest

	// TODO: try deleting IdRecord after request creation
}

// func (k Keeper) RequestIdentityRecordsVerify(ctx sdk.Context, address, verifier sdk.AccAddress, recordIds []uint64, tip sdk.Coin) (uint64, error) {
// func (k Keeper) GetIdRecordsVerifyRequest(ctx sdk.Context, requestId uint64) *types.IdentityRecordsVerify {
// func (k Keeper) DeleteIdRecordsVerifyRequest(ctx sdk.Context, requestId uint64) {
// func (k Keeper) ApproveIdentityRecords(ctx sdk.Context, verifier sdk.AccAddress, requestId uint64) error {
// func (k Keeper) CancelIdentityRecordsVerifyRequest(ctx sdk.Context, executor sdk.AccAddress, requestId uint64) error {
// func (k Keeper) GetIdRecordsVerifyRequestsByRequester(ctx sdk.Context, requester sdk.AccAddress) []types.IdentityRecordsVerify {
// func (k Keeper) GetIdRecordsVerifyRequestsByApprover(ctx sdk.Context, requester sdk.AccAddress) []types.IdentityRecordsVerify {
// func (k Keeper) GetAllIdRecordsVerifyRequests(ctx sdk.Context) []types.IdentityRecordsVerify {
