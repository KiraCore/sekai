package keeper_test

import (
	"strings"
	"testing"
	"time"

	simapp "github.com/KiraCore/sekai/app"
	"github.com/KiraCore/sekai/x/gov/keeper"
	"github.com/KiraCore/sekai/x/gov/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

func TestKeeper_ValidateIdentityRecordKey(t *testing.T) {
	require.False(t, keeper.ValidateIdentityRecordKey("_abc"))
	require.False(t, keeper.ValidateIdentityRecordKey("1abc"))
	require.True(t, keeper.ValidateIdentityRecordKey("ab_a"))
	require.True(t, keeper.ValidateIdentityRecordKey("ab_1a"))
	require.True(t, keeper.ValidateIdentityRecordKey("aa_Aa"))
}

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
	record := app.CustomGovKeeper.GetIdentityRecordById(ctx, 1)
	require.Nil(t, record)

	// create a new record and check if set correctly
	addr1 := sdk.AccAddress("foo1________________")
	addr2 := sdk.AccAddress("foo2________________")
	addr3 := sdk.AccAddress("foo3________________")
	newRecord := types.IdentityRecord{
		Id:        1,
		Address:   addr1.String(),
		Key:       "key",
		Value:     "value",
		Date:      time.Now().UTC(),
		Verifiers: []string{addr2.String(), addr3.String()},
	}
	app.CustomGovKeeper.SetIdentityRecord(ctx, newRecord)
	record = app.CustomGovKeeper.GetIdentityRecordById(ctx, 1)
	require.NotNil(t, record)
	require.Equal(t, *record, newRecord)

	// check no panics
	app.CustomGovKeeper.DeleteIdentityRecordById(ctx, 0)

	// remove existing id and check
	app.CustomGovKeeper.DeleteIdentityRecordById(ctx, 1)
	record = app.CustomGovKeeper.GetIdentityRecordById(ctx, 1)
	require.Nil(t, record)

	// check automatic conversion to lowercase
	uppercaseRecord := types.IdentityRecord{
		Id:        2,
		Address:   addr1.String(),
		Key:       "MyKey",
		Value:     "value",
		Date:      time.Now().UTC(),
		Verifiers: []string{addr2.String(), addr3.String()},
	}
	app.CustomGovKeeper.SetIdentityRecord(ctx, uppercaseRecord)
	record = app.CustomGovKeeper.GetIdentityRecordById(ctx, 2)
	require.NotNil(t, record)
	require.Equal(t, record.Key, "mykey")

	// try to get via uppercase key
	recordId := app.CustomGovKeeper.GetIdentityRecordIdByAddressKey(ctx, addr1, "MYKEY")
	require.Equal(t, recordId, uint64(2))

	// try to get by key
	recordId = app.CustomGovKeeper.GetIdentityRecordIdByAddressKey(ctx, addr1, "_key")
	require.Equal(t, recordId, uint64(0))

	// check invalid key set
	invalidRecord := types.IdentityRecord{
		Id:        1,
		Address:   addr1.String(),
		Key:       "_key",
		Value:     "value",
		Date:      time.Now().UTC(),
		Verifiers: []string{addr2.String(), addr3.String()},
	}
	require.Panics(t, func() {
		app.CustomGovKeeper.SetIdentityRecord(ctx, invalidRecord)
	})
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
	ctx = ctx.WithBlockTime(now)
	err := app.CustomGovKeeper.RegisterIdentityRecords(ctx, addr1, types.WrapInfos(infos))
	require.NoError(t, err)

	record := app.CustomGovKeeper.GetIdentityRecordById(ctx, 1)
	require.NotNil(t, record)
	expectedRecord := types.IdentityRecord{
		Id:      1,
		Address: addr1.String(),
		Key:     "key",
		Value:   "value",
		Date:    now,
	}
	require.Equal(t, *record, expectedRecord)

	err = app.CustomGovKeeper.RegisterIdentityRecords(ctx, addr2, types.WrapInfos(infos))
	require.NoError(t, err)

	records := app.CustomGovKeeper.GetAllIdentityRecords(ctx)
	require.Len(t, records, 2)
	records = app.CustomGovKeeper.GetIdRecordsByAddress(ctx, addr1)
	require.NotNil(t, records)
	records = app.CustomGovKeeper.GetIdRecordsByAddress(ctx, addr2)
	require.NotNil(t, records)

	// remove existing id and check
	app.CustomGovKeeper.DeleteIdentityRecords(ctx, addr2, []string{})
	records = app.CustomGovKeeper.GetIdRecordsByAddress(ctx, addr2)
	require.Len(t, records, 0)
	records = app.CustomGovKeeper.GetIdRecordsByAddress(ctx, addr1)
	require.Len(t, records, 1)
	records = app.CustomGovKeeper.GetAllIdentityRecords(ctx)
	require.Len(t, records, 1)

	infos["key1"] = "value1"
	now = now.Add(time.Second)
	ctx = ctx.WithBlockTime(now)

	// try deleting one key
	err = app.CustomGovKeeper.DeleteIdentityRecords(ctx, addr3, []string{"key1"})
	require.NoError(t, err)

	// set verifier of identity record
	app.CustomGovKeeper.RegisterIdentityRecords(ctx, addr1, types.WrapInfos(infos))

	err = app.CustomGovKeeper.RegisterIdentityRecords(ctx, addr1, types.WrapInfos(infos))
	require.NoError(t, err)
	record = app.CustomGovKeeper.GetIdentityRecordById(ctx, 1)
	require.NotNil(t, record)
	require.Equal(t, *record, types.IdentityRecord{
		Id:      1,
		Address: addr1.String(),
		Key:     "key",
		Value:   "value",
		Date:    now,
	})
	records = app.CustomGovKeeper.GetAllIdentityRecords(ctx)
	require.Len(t, records, 2)
	records = app.CustomGovKeeper.GetIdRecordsByAddress(ctx, addr1)
	require.NotNil(t, record)
	records = app.CustomGovKeeper.GetIdRecordsByAddress(ctx, addr2)
	require.Len(t, records, 0)

	// check identity records by address
	records = app.CustomGovKeeper.GetIdRecordsByAddress(ctx, addr1)
	require.Len(t, records, 2)

	// check identity records by address and keys
	records, err = app.CustomGovKeeper.GetIdRecordsByAddressAndKeys(ctx, addr1, []string{})
	require.NoError(t, err)
	require.Len(t, records, 2)

	records, err = app.CustomGovKeeper.GetIdRecordsByAddressAndKeys(ctx, addr1, []string{"key"})
	require.NoError(t, err)
	require.Len(t, records, 1)

	records, err = app.CustomGovKeeper.GetIdRecordsByAddressAndKeys(ctx, addr1, []string{"invalidkey"})
	require.Error(t, err)
}

func TestKeeper_TryLongMonikerField(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	// create a new record and check if set correctly
	addr1 := sdk.AccAddress("foo1________________")
	infos := make(map[string]string)
	infos["moniker"] = strings.Repeat("A", 33)
	now := time.Now().UTC()
	ctx = ctx.WithBlockTime(now)
	err := app.CustomGovKeeper.RegisterIdentityRecords(ctx, addr1, types.WrapInfos(infos))
	require.Error(t, err)
}

func TestKeeper_TrySameMonikerField(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	// create a new record and check if set correctly
	addr1 := sdk.AccAddress("foo1________________")
	addr2 := sdk.AccAddress("foo2________________")
	infos := make(map[string]string)
	infos["moniker"] = "AAA"
	now := time.Now().UTC()
	ctx = ctx.WithBlockTime(now)
	err := app.CustomGovKeeper.RegisterIdentityRecords(ctx, addr1, types.WrapInfos(infos))
	require.NoError(t, err)
	err = app.CustomGovKeeper.RegisterIdentityRecords(ctx, addr1, types.WrapInfos(infos))
	require.NoError(t, err)
	err = app.CustomGovKeeper.RegisterIdentityRecords(ctx, addr2, types.WrapInfos(infos))
	require.Error(t, err)
	infos["moniker"] = "AAA2"
	err = app.CustomGovKeeper.RegisterIdentityRecords(ctx, addr2, types.WrapInfos(infos))
	require.NoError(t, err)
}

func TestKeeper_TryUniqueIdentityKeysSet(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	// create a new record and check if set correctly
	now := time.Now().UTC()
	ctx = ctx.WithBlockTime(now)
	err := app.CustomGovKeeper.SetNetworkProperty(ctx, types.UniqueIdentityKeys, types.NetworkPropertyValue{StrValue: "moniker,email"})
	require.NoError(t, err)
	networkProperties := app.CustomGovKeeper.GetNetworkProperties(ctx)
	require.NotNil(t, networkProperties)

	require.Equal(t, networkProperties.UniqueIdentityKeys, "moniker,email")
}

func TestKeeper_IdentityKeysManagement(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	// create a new record and check if set correctly
	addr1 := sdk.AccAddress("foo1________________")
	infos := make(map[string]string)
	infos["MyKey"] = "MyValue"
	infos["Nike"] = "MyNike"
	now := time.Now().UTC()
	ctx = ctx.WithBlockTime(now)
	err := app.CustomGovKeeper.RegisterIdentityRecords(ctx, addr1, types.WrapInfos(infos))
	require.NoError(t, err)

	record := app.CustomGovKeeper.GetIdentityRecordById(ctx, 1)
	require.NotNil(t, record)
	expectedRecord := types.IdentityRecord{
		Id:      1,
		Address: addr1.String(),
		Key:     "mykey",
		Value:   "MyValue",
		Date:    now,
	}
	require.Equal(t, *record, expectedRecord)

	record = app.CustomGovKeeper.GetIdentityRecordById(ctx, 2)
	require.NotNil(t, record)
	expectedRecord = types.IdentityRecord{
		Id:      2,
		Address: addr1.String(),
		Key:     "nike",
		Value:   "MyNike",
		Date:    now,
	}
	require.Equal(t, *record, expectedRecord)

	// check invalid key involved registration
	infos["1Nike"] = "MyNike"
	err = app.CustomGovKeeper.RegisterIdentityRecords(ctx, addr1, types.WrapInfos(infos))
	require.Error(t, err)

	records, err := app.CustomGovKeeper.GetIdRecordsByAddressAndKeys(ctx, addr1, []string{"MyKey", "nike"})
	require.NoError(t, err)
	require.Len(t, records, 2)

	records, err = app.CustomGovKeeper.GetIdRecordsByAddressAndKeys(ctx, addr1, []string{"MyKey", "nike", "A"})
	require.Error(t, err)

	records, err = app.CustomGovKeeper.GetIdRecordsByAddressAndKeys(ctx, addr1, []string{"MyKey", "nike", "_"})
	require.Error(t, err)

	// get address from identity record key value
	addrs := app.CustomGovKeeper.GetAddressesByIdRecordKey(ctx, "mykey", "MyValue")
	require.Len(t, addrs, 1)

	addrs = app.CustomGovKeeper.GetAddressesByIdRecordKey(ctx, "mykey", "MyValue2")
	require.Len(t, addrs, 0)

	// delete by uppercase key and check if deleted correctly
	err = app.CustomGovKeeper.DeleteIdentityRecords(ctx, addr1, []string{"myKey"})
	require.NoError(t, err)

	records, err = app.CustomGovKeeper.GetIdRecordsByAddressAndKeys(ctx, addr1, []string{"MyKey"})
	require.Error(t, err)

	// test for moniker field deletion is not enabled
	err = app.CustomGovKeeper.RegisterIdentityRecords(ctx, addr1, []types.IdentityInfoEntry{
		{
			Key:  "moniker",
			Info: "node0",
		},
	})
	require.NoError(t, err)
	err = app.CustomGovKeeper.DeleteIdentityRecords(ctx, addr1, []string{"moniker"})
	require.Error(t, err)
}

func TestKeeper_IdentityRecordApproveFlow(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	// create a new record and check if set correctly
	addr1 := sdk.AccAddress("foo1________________")
	addr2 := sdk.AccAddress("foo2________________")
	addr3 := sdk.AccAddress("foo3________________")
	addr4 := sdk.AccAddress("foo4________________")
	initCoins := sdk.Coins{sdk.NewInt64Coin(sdk.DefaultBondDenom, 10000)}
	app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, initCoins)
	app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, initCoins)
	app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, initCoins)
	app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, addr1, initCoins)
	app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, addr2, initCoins)
	app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, addr3, initCoins)

	infos := make(map[string]string)
	infos["key"] = "value"
	now := time.Now().UTC()
	ctx = ctx.WithBlockTime(now)
	err := app.CustomGovKeeper.RegisterIdentityRecords(ctx, addr1, types.WrapInfos(infos))
	require.NoError(t, err)

	record := app.CustomGovKeeper.GetIdentityRecordById(ctx, 1)
	require.NotNil(t, record)
	expectedRecord := types.IdentityRecord{
		Id:      1,
		Address: addr1.String(),
		Key:     "key",
		Value:   "value",
		Date:    now,
	}
	require.Equal(t, *record, expectedRecord)

	err = app.CustomGovKeeper.RegisterIdentityRecords(ctx, addr2, types.WrapInfos(infos))
	require.NoError(t, err)

	// bigger tip than balance
	ctxCache, _ := ctx.CacheContext()
	reqId, err := app.CustomGovKeeper.RequestIdentityRecordsVerify(ctxCache, addr1, addr3, []uint64{1}, sdk.NewInt64Coin(sdk.DefaultBondDenom, 1000000000))
	require.Equal(t, reqId, uint64(0))
	require.Error(t, err)

	// test smaller tip than minimum tip
	ctxCache, _ = ctx.CacheContext()
	reqId, err = app.CustomGovKeeper.RequestIdentityRecordsVerify(ctxCache, addr1, addr3, []uint64{1}, sdk.NewInt64Coin(sdk.DefaultBondDenom, 199))
	require.Error(t, err)
	require.Equal(t, reqId, uint64(1))

	// request id record 1 to addr3 by addr1
	reqId, err = app.CustomGovKeeper.RequestIdentityRecordsVerify(ctx, addr1, addr3, []uint64{1}, sdk.NewInt64Coin(sdk.DefaultBondDenom, 200))
	require.Equal(t, reqId, uint64(1))
	require.NoError(t, err)
	request := app.CustomGovKeeper.GetIdRecordsVerifyRequest(ctx, 1)
	require.NotNil(t, request)
	require.Equal(t, *request, types.IdentityRecordsVerify{
		Id:                 1,
		Address:            addr1.String(),
		Verifier:           addr3.String(),
		RecordIds:          []uint64{1},
		Tip:                sdk.NewInt64Coin(sdk.DefaultBondDenom, 200),
		LastRecordEditDate: now,
	})
	coins := app.BankKeeper.GetAllBalances(ctx, addr1)
	require.Equal(t, coins, sdk.Coins{sdk.NewInt64Coin(sdk.DefaultBondDenom, 9800)})
	coins = app.BankKeeper.GetAllBalances(ctx, app.AccountKeeper.GetModuleAddress(types.ModuleName))
	require.Equal(t, coins, sdk.Coins{sdk.NewInt64Coin(sdk.DefaultBondDenom, 200)})
	err = app.CustomGovKeeper.HandleIdentityRecordsVerifyRequest(ctx, addr3, 1, true)
	require.NoError(t, err)
	coins = app.BankKeeper.GetAllBalances(ctx, addr3)
	require.Equal(t, coins, sdk.Coins{sdk.NewInt64Coin(sdk.DefaultBondDenom, 10200)})
	coins = app.BankKeeper.GetAllBalances(ctx, app.AccountKeeper.GetModuleAddress(types.ModuleName))
	require.Equal(t, coins, sdk.Coins{})
	record = app.CustomGovKeeper.GetIdentityRecordById(ctx, 1)
	require.NotNil(t, record)
	require.Equal(t, *record, types.IdentityRecord{
		Id:        1,
		Address:   addr1.String(),
		Key:       "key",
		Value:     "value",
		Date:      now,
		Verifiers: []string{addr3.String()},
	})
	request = app.CustomGovKeeper.GetIdRecordsVerifyRequest(ctx, 1)
	require.Nil(t, request)

	// request id record 1 to addr3 by addr1 again
	reqId, err = app.CustomGovKeeper.RequestIdentityRecordsVerify(ctx, addr1, addr3, []uint64{1}, sdk.NewInt64Coin(sdk.DefaultBondDenom, 200))
	require.Equal(t, reqId, uint64(2))
	require.NoError(t, err)
	coins = app.BankKeeper.GetAllBalances(ctx, addr1)
	require.Equal(t, coins, sdk.Coins{sdk.NewInt64Coin(sdk.DefaultBondDenom, 9600)})
	err = app.CustomGovKeeper.HandleIdentityRecordsVerifyRequest(ctx, addr3, 2, true)
	require.NoError(t, err)
	coins = app.BankKeeper.GetAllBalances(ctx, addr3)
	require.Equal(t, coins, sdk.Coins{sdk.NewInt64Coin(sdk.DefaultBondDenom, 10400)})
	record = app.CustomGovKeeper.GetIdentityRecordById(ctx, 1)
	require.NotNil(t, record)
	require.Equal(t, *record, types.IdentityRecord{
		Id:        1,
		Address:   addr1.String(),
		Key:       "key",
		Value:     "value",
		Date:      now,
		Verifiers: []string{addr3.String()},
	})

	// request id record 1 and 2 to addr3 by addr1 again
	reqId, err = app.CustomGovKeeper.RequestIdentityRecordsVerify(ctx, addr1, addr3, []uint64{1, 2}, sdk.NewInt64Coin(sdk.DefaultBondDenom, 200))
	require.Equal(t, reqId, uint64(3))
	require.NoError(t, err)
	coins = app.BankKeeper.GetAllBalances(ctx, addr1)
	require.Equal(t, coins, sdk.Coins{sdk.NewInt64Coin(sdk.DefaultBondDenom, 9400)})
	err = app.CustomGovKeeper.HandleIdentityRecordsVerifyRequest(ctx, addr3, 3, true)
	require.NoError(t, err)
	coins = app.BankKeeper.GetAllBalances(ctx, addr3)
	require.Equal(t, coins, sdk.Coins{sdk.NewInt64Coin(sdk.DefaultBondDenom, 10600)})
	record = app.CustomGovKeeper.GetIdentityRecordById(ctx, 1)
	require.NotNil(t, record)
	require.Equal(t, *record, types.IdentityRecord{
		Id:        1,
		Address:   addr1.String(),
		Key:       "key",
		Value:     "value",
		Date:      now,
		Verifiers: []string{addr3.String()},
	})
	record = app.CustomGovKeeper.GetIdentityRecordById(ctx, 2)
	require.NotNil(t, record)
	require.Equal(t, *record, types.IdentityRecord{
		Id:        2,
		Address:   addr2.String(),
		Key:       "key",
		Value:     "value",
		Date:      now,
		Verifiers: []string{addr3.String()},
	})

	// request id record 2 to addr3 by addr2
	reqId, err = app.CustomGovKeeper.RequestIdentityRecordsVerify(ctx, addr2, addr3, []uint64{2}, sdk.NewInt64Coin(sdk.DefaultBondDenom, 200))
	require.Equal(t, reqId, uint64(4))
	require.NoError(t, err)
	err = app.CustomGovKeeper.HandleIdentityRecordsVerifyRequest(ctx, addr3, 4, true)
	require.NoError(t, err)
	record = app.CustomGovKeeper.GetIdentityRecordById(ctx, 2)
	require.NotNil(t, record)
	require.Equal(t, *record, types.IdentityRecord{
		Id:        2,
		Address:   addr2.String(),
		Key:       "key",
		Value:     "value",
		Date:      now,
		Verifiers: []string{addr3.String()},
	})

	// request non-exist identity record
	ctxCache, _ = ctx.CacheContext()
	reqId, err = app.CustomGovKeeper.RequestIdentityRecordsVerify(ctxCache, addr2, addr3, []uint64{5}, sdk.NewInt64Coin(sdk.DefaultBondDenom, 200))
	require.Equal(t, reqId, uint64(5))
	require.Error(t, err)

	// approve with non-approver
	reqId, err = app.CustomGovKeeper.RequestIdentityRecordsVerify(ctx, addr1, addr3, []uint64{1}, sdk.NewInt64Coin(sdk.DefaultBondDenom, 200))
	require.Equal(t, reqId, uint64(5))
	require.NoError(t, err)
	ctxCache, _ = ctx.CacheContext()
	err = app.CustomGovKeeper.HandleIdentityRecordsVerifyRequest(ctxCache, addr2, 5, true)
	require.Error(t, err)

	// approve not existing request id
	ctxCache, _ = ctx.CacheContext()
	err = app.CustomGovKeeper.HandleIdentityRecordsVerifyRequest(ctxCache, addr2, 0xFFFFF, true)
	require.Error(t, err)

	// try edit and check if verification records all gone
	infos["key1"] = "value1"
	err = app.CustomGovKeeper.DeleteIdentityRecords(ctx, addr1, []string{"key", "key1"})
	require.NoError(t, err)
	record = app.CustomGovKeeper.GetIdentityRecordById(ctx, 1)
	require.Nil(t, record)

	// check get queries
	requests := app.CustomGovKeeper.GetIdRecordsVerifyRequestsByRequester(ctx, addr1)
	require.Len(t, requests, 1)
	requests = app.CustomGovKeeper.GetIdRecordsVerifyRequestsByApprover(ctx, addr1)
	require.Len(t, requests, 0)
	requests = app.CustomGovKeeper.GetAllIdRecordsVerifyRequests(ctx)
	require.Len(t, requests, 1)

	// remove all and query again
	app.CustomGovKeeper.DeleteIdRecordsVerifyRequest(ctx, 5)
	requests = app.CustomGovKeeper.GetAllIdRecordsVerifyRequests(ctx)
	require.Len(t, requests, 0)

	// try to cancel request and check coin moves correctly
	reqId, err = app.CustomGovKeeper.RequestIdentityRecordsVerify(ctx, addr2, addr3, []uint64{2}, sdk.NewInt64Coin(sdk.DefaultBondDenom, 200))
	require.Equal(t, reqId, uint64(6))
	require.NoError(t, err)
	coins = app.BankKeeper.GetAllBalances(ctx, addr2)
	require.Equal(t, coins, sdk.Coins{sdk.NewInt64Coin(sdk.DefaultBondDenom, 9600)})
	cacheCtx, _ := ctx.CacheContext()
	err = app.CustomGovKeeper.CancelIdentityRecordsVerifyRequest(cacheCtx, addr3, 6)
	require.Error(t, err)
	err = app.CustomGovKeeper.CancelIdentityRecordsVerifyRequest(ctx, addr2, 6)
	require.NoError(t, err)
	request = app.CustomGovKeeper.GetIdRecordsVerifyRequest(ctx, 6)
	require.Nil(t, request)
	coins = app.BankKeeper.GetAllBalances(ctx, addr2)
	require.Equal(t, coins, sdk.Coins{sdk.NewInt64Coin(sdk.DefaultBondDenom, 9800)})

	// try deleting request after request creation
	reqId, err = app.CustomGovKeeper.RequestIdentityRecordsVerify(ctx, addr2, addr3, []uint64{2}, sdk.NewInt64Coin(sdk.DefaultBondDenom, 200))
	require.Equal(t, reqId, uint64(7))
	app.CustomGovKeeper.DeleteIdRecordsVerifyRequest(ctx, 7)
	request = app.CustomGovKeeper.GetIdRecordsVerifyRequest(ctx, 7)
	require.Nil(t, request)
	requests = app.CustomGovKeeper.GetIdRecordsVerifyRequestsByRequester(ctx, addr2)
	require.Len(t, requests, 0)
	requests = app.CustomGovKeeper.GetIdRecordsVerifyRequestsByApprover(ctx, addr3)
	require.Len(t, requests, 0)

	// check automatic reject if record is edited after raising verification request
	reqId, err = app.CustomGovKeeper.RequestIdentityRecordsVerify(ctx, addr2, addr4, []uint64{2}, sdk.NewInt64Coin(sdk.DefaultBondDenom, 200))
	require.Equal(t, reqId, uint64(8))
	ctx = ctx.WithBlockTime(now.Add(time.Second))
	app.CustomGovKeeper.RegisterIdentityRecords(ctx, addr2, types.WrapInfos(infos))
	ctx, _ = ctx.CacheContext()
	err = app.CustomGovKeeper.HandleIdentityRecordsVerifyRequest(ctx, addr4, 8, true)
	require.NoError(t, err)
	record = app.CustomGovKeeper.GetIdentityRecordById(ctx, 2)
	require.NotNil(t, record)
	require.False(t, keeper.CheckIfWithinStringArray(addr4.String(), record.Verifiers))
	coins = app.BankKeeper.GetAllBalances(ctx, addr4)
	require.Equal(t, coins, sdk.Coins{sdk.NewInt64Coin(sdk.DefaultBondDenom, 200)})

	// try deleting id record after request creation
	reqId, err = app.CustomGovKeeper.RequestIdentityRecordsVerify(ctx, addr2, addr3, []uint64{2}, sdk.NewInt64Coin(sdk.DefaultBondDenom, 200))
	require.Equal(t, reqId, uint64(9))
	app.CustomGovKeeper.DeleteIdentityRecords(ctx, addr2, []string{})
	cacheCtx, _ = ctx.CacheContext()
	err = app.CustomGovKeeper.HandleIdentityRecordsVerifyRequest(cacheCtx, addr3, 9, true)
	require.Error(t, err)
}
