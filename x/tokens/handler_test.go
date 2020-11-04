package tokens_test

import (
	"bytes"
	"strconv"
	"strings"
	"testing"

	"github.com/KiraCore/sekai/simapp"
	"github.com/KiraCore/sekai/x/gov/types"
	customgovtypes "github.com/KiraCore/sekai/x/gov/types"
	tokens "github.com/KiraCore/sekai/x/tokens"
	tokenstypes "github.com/KiraCore/sekai/x/tokens/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

func NewAccountByIndex(accNum int) sdk.AccAddress {
	var buffer bytes.Buffer
	i := accNum + 100
	numString := strconv.Itoa(i)
	buffer.WriteString("A58856F0FD53BF058B4909A21AEC019107BA6") //base address string

	buffer.WriteString(numString) //adding on final two digits to make addresses unique
	res, _ := sdk.AccAddressFromHex(buffer.String())
	bech := res.String()
	addr, _ := simapp.TestAddr(buffer.String(), bech)
	buffer.Reset()
	return addr
}

func setPermissionToAddr(t *testing.T, app *simapp.SimApp, ctx sdk.Context, addr sdk.AccAddress, perm types.PermValue) error {
	proposerActor := customgovtypes.NewDefaultActor(addr)
	err := proposerActor.Permissions.AddToWhitelist(perm)
	require.NoError(t, err)

	app.CustomGovKeeper.SaveNetworkActor(ctx, proposerActor)

	return nil
}

func TestNewHandler_MsgUpsertTokenAlias(t *testing.T) {

	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})
	handler := tokens.NewHandler(app.TokensKeeper, app.CustomGovKeeper)

	tests := []struct {
		name        string
		constructor func(sdk.AccAddress) (*tokenstypes.MsgUpsertTokenAlias, error)
		handlerErr  string
	}{
		{
			name: "good permission test",
			constructor: func(addr sdk.AccAddress) (*tokenstypes.MsgUpsertTokenAlias, error) {
				err := setPermissionToAddr(t, app, ctx, addr, types.PermUpsertTokenAlias)
				require.NoError(t, err)
				return tokenstypes.NewMsgUpsertTokenAlias(
					addr,
					0, 0,
					[]tokenstypes.VoteType{tokenstypes.VoteType_no, tokenstypes.VoteType_yes},
					"ETH",
					"Ethereum",
					"icon",
					6,
					[]string{"finney"},
					tokenstypes.ProposalStatus_active,
				), nil
			},
		},
		{
			name: "lack permission test",
			constructor: func(addr sdk.AccAddress) (*tokenstypes.MsgUpsertTokenAlias, error) {
				return tokenstypes.NewMsgUpsertTokenAlias(
					addr,
					0, 0,
					[]tokenstypes.VoteType{tokenstypes.VoteType_no, tokenstypes.VoteType_yes},
					"ETH",
					"Ethereum",
					"icon",
					6,
					[]string{"finney"},
					tokenstypes.ProposalStatus_active,
				), nil
			},
			handlerErr: "PermUpsertTokenAlias: not enough permissions",
		},
	}
	for i, tt := range tests {
		addr := NewAccountByIndex(i)
		theMsg, err := tt.constructor(addr)
		require.NoError(t, err)

		_, err = handler(ctx, theMsg)
		if len(tt.handlerErr) != 0 {
			require.Error(t, err)
			require.Contains(t, err.Error(), tt.handlerErr)
		} else {
			require.NoError(t, err)

			// test various query commands
			alias := app.TokensKeeper.GetTokenAlias(ctx, theMsg.Symbol)
			require.True(t, alias != nil)
			aliasesAll := app.TokensKeeper.ListTokenAlias(ctx)
			require.True(t, len(aliasesAll) > 0)
			aliasesByDenom := app.TokensKeeper.GetTokenAliasesByDenom(ctx, theMsg.Denoms)
			require.True(t, aliasesByDenom[theMsg.Denoms[0]] != nil)

			// try different alias for same denom
			theMsg.Symbol += "V2"
			_, err = handler(ctx, theMsg)
			require.Error(t, err)
			require.True(t, strings.Contains(err.Error(), "denom is already registered"))
		}
	}
}

func TestNewHandler_MsgUpsertTokenRate(t *testing.T) {

	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})
	handler := tokens.NewHandler(app.TokensKeeper, app.CustomGovKeeper)

	tests := []struct {
		name        string
		constructor func(sdk.AccAddress) (*tokenstypes.MsgUpsertTokenRate, error)
		handlerErr  string
	}{
		{
			name: "good permission test",
			constructor: func(addr sdk.AccAddress) (*tokenstypes.MsgUpsertTokenRate, error) {
				err := setPermissionToAddr(t, app, ctx, addr, types.PermUpsertTokenRate)
				require.NoError(t, err)
				return tokenstypes.NewMsgUpsertTokenRate(
					addr,
					"finney", sdk.NewDecWithPrec(1, 3), // 0.001
					true,
				), nil
			},
		},
		{
			name: "lack permission test",
			constructor: func(addr sdk.AccAddress) (*tokenstypes.MsgUpsertTokenRate, error) {
				return tokenstypes.NewMsgUpsertTokenRate(
					addr,
					"finney", sdk.NewDecWithPrec(1, 3), // 0.001
					true,
				), nil
			},
			handlerErr: "PermUpsertTokenRate: not enough permissions",
		},
		{
			name: "negative rate value test",
			constructor: func(addr sdk.AccAddress) (*tokenstypes.MsgUpsertTokenRate, error) {
				return tokenstypes.NewMsgUpsertTokenRate(
					addr,
					"finney", sdk.NewDec(-1), // -1
					true,
				), nil
			},
			handlerErr: "rate should be positive",
		},
		{
			name: "bond denom rate change test",
			constructor: func(addr sdk.AccAddress) (*tokenstypes.MsgUpsertTokenRate, error) {
				err := setPermissionToAddr(t, app, ctx, addr, types.PermUpsertTokenRate)
				require.NoError(t, err)
				return tokenstypes.NewMsgUpsertTokenRate(
					addr,
					"ukex", sdk.NewDec(10),
					true,
				), nil
			},
			handlerErr: "bond denom rate is read-only",
		},
	}
	for i, tt := range tests {
		addr := NewAccountByIndex(i)
		theMsg, err := tt.constructor(addr)
		require.NoError(t, err)

		_, err = handler(ctx, theMsg)
		if len(tt.handlerErr) != 0 {
			require.Error(t, err)
			require.Contains(t, err.Error(), tt.handlerErr)
		} else {
			require.NoError(t, err)

			// test various query commands
			rate := app.TokensKeeper.GetTokenRate(ctx, theMsg.Denom)
			require.True(t, rate != nil)
			ratesAll := app.TokensKeeper.ListTokenRate(ctx)
			require.True(t, len(ratesAll) > 0)
			ratesByDenom := app.TokensKeeper.GetTokenRatesByDenom(ctx, []string{theMsg.Denom})
			require.True(t, ratesByDenom[theMsg.Denom] != nil)
		}
	}
}
