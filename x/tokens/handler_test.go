package tokens_test

import (
	"bytes"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/KiraCore/sekai/app"
	simapp "github.com/KiraCore/sekai/app"
	"github.com/KiraCore/sekai/x/gov"
	"github.com/KiraCore/sekai/x/gov/types"
	govtypes "github.com/KiraCore/sekai/x/gov/types"
	tokens "github.com/KiraCore/sekai/x/tokens"
	tokenstypes "github.com/KiraCore/sekai/x/tokens/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

func TestMain(m *testing.M) {
	app.SetConfig()
	os.Exit(m.Run())
}

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

func setPermissionToAddr(t *testing.T, app *simapp.SekaiApp, ctx sdk.Context, addr sdk.AccAddress, perm types.PermValue) error {
	proposerActor := govtypes.NewDefaultActor(addr)
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
					"ETH",
					"Ethereum",
					"icon",
					6,
					[]string{"finney"},
				), nil
			},
		},
		{
			name: "lack permission test",
			constructor: func(addr sdk.AccAddress) (*tokenstypes.MsgUpsertTokenAlias, error) {
				return tokenstypes.NewMsgUpsertTokenAlias(
					addr,
					"ETH",
					"Ethereum",
					"icon",
					6,
					[]string{"finney"},
				), nil
			},
			handlerErr: "PERMISSION_UPSERT_TOKEN_ALIAS: not enough permissions",
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
			handlerErr: "PERMISSION_UPSERT_TOKEN_RATE: not enough permissions",
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

func TestHandler_CreateProposalUpsertTokenAliases_Errors(t *testing.T) {
	proposerAddr, err := sdk.AccAddressFromBech32("kira1alzyfq40zjsveat87jlg8jxetwqmr0a29sgd0f")
	require.NoError(t, err)

	tests := []struct {
		name         string
		content      govtypes.Content
		preparePerms func(t *testing.T, app *simapp.SekaiApp, ctx sdk.Context)
		expectedErr  error
	}{
		{
			"Proposer does not have Perm",
			tokenstypes.NewUpsertTokenAliasProposal(
				"BTC",
				"Bitcoin",
				"http://theicon.com",
				18,
				[]string{},
			),
			func(t *testing.T, app *simapp.SekaiApp, ctx sdk.Context) {},
			errors.Wrap(types.ErrNotEnoughPermissions, types.PermCreateUpsertTokenAliasProposal.String()),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			app := simapp.Setup(false)
			ctx := app.NewContext(false, tmproto.Header{})

			tt.preparePerms(t, app, ctx)

			handler := gov.NewHandler(app.CustomGovKeeper)
			msg, err := govtypes.NewMsgSubmitProposal(proposerAddr, "title", "some desc", tt.content)
			require.NoError(t, err)
			_, err = handler(ctx, msg)
			require.EqualError(t, err, tt.expectedErr.Error())
		})
	}
}

func TestHandler_CreateProposalUpsertTokenAliases(t *testing.T) {
	proposerAddr, err := sdk.AccAddressFromBech32("kira1alzyfq40zjsveat87jlg8jxetwqmr0a29sgd0f")
	require.NoError(t, err)

	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{
		Time: time.Now(),
	})

	// Set proposer Permissions
	proposerActor := types.NewDefaultActor(proposerAddr)
	err2 := app.CustomGovKeeper.AddWhitelistPermission(ctx, proposerActor, types.PermCreateUpsertTokenAliasProposal)
	require.NoError(t, err2)

	properties := app.CustomGovKeeper.GetNetworkProperties(ctx)
	properties.ProposalEndTime = 10
	app.CustomGovKeeper.SetNetworkProperties(ctx, properties)

	handler := gov.NewHandler(app.CustomGovKeeper)
	proposal := tokenstypes.NewUpsertTokenAliasProposal(
		"BTC",
		"Bitcoin",
		"http://sdlkfjalsdk.es",
		18,
		[]string{
			"atom",
		},
	)
	msg, err := govtypes.NewMsgSubmitProposal(proposerAddr, "title", "some desc", proposal)
	require.NoError(t, err)
	res, err := handler(
		ctx,
		msg,
	)
	require.NoError(t, err)
	expData, _ := proto.Marshal(&govtypes.MsgSubmitProposalResponse{ProposalID: 1})
	require.Equal(t, expData, res.Data)

	savedProposal, found := app.CustomGovKeeper.GetProposal(ctx, 1)
	require.True(t, found)

	expectedSavedProposal, err := types.NewProposal(
		1,
		"title",
		"some desc",
		tokenstypes.NewUpsertTokenAliasProposal(
			"BTC",
			"Bitcoin",
			"http://sdlkfjalsdk.es",
			18,
			[]string{
				"atom",
			},
		),
		ctx.BlockTime(),
		ctx.BlockTime().Add(time.Second*time.Duration(properties.ProposalEndTime)),
		ctx.BlockTime().Add(time.Second*time.Duration(properties.ProposalEndTime)+
			time.Second*time.Duration(properties.ProposalEnactmentTime),
		),
		ctx.BlockHeight()+2,
		ctx.BlockHeight()+3,
	)
	require.NoError(t, err)
	require.Equal(t, expectedSavedProposal, savedProposal)

	// Next proposal ID is increased.
	id := app.CustomGovKeeper.GetNextProposalID(ctx)
	require.Equal(t, uint64(2), id)

	// Is not on finished active proposals.
	iterator := app.CustomGovKeeper.GetActiveProposalsWithFinishedVotingEndTimeIterator(ctx, ctx.BlockTime())
	require.False(t, iterator.Valid())

	ctx = ctx.WithBlockTime(ctx.BlockTime().Add(time.Minute * 10))
	iterator = app.CustomGovKeeper.GetActiveProposalsWithFinishedVotingEndTimeIterator(ctx, ctx.BlockTime())
	require.True(t, iterator.Valid())
}

func TestHandler_CreateProposalUpsertTokenRates_Errors(t *testing.T) {
	proposerAddr, err := sdk.AccAddressFromBech32("kira1alzyfq40zjsveat87jlg8jxetwqmr0a29sgd0f")
	require.NoError(t, err)

	tests := []struct {
		name         string
		content      govtypes.Content
		preparePerms func(t *testing.T, app *simapp.SekaiApp, ctx sdk.Context)
		expectedErr  error
	}{
		{
			"Proposer does not have Perm",
			tokenstypes.NewUpsertTokenRatesProposal(
				"btc",
				sdk.NewDec(1234),
				false,
			),
			func(t *testing.T, app *simapp.SekaiApp, ctx sdk.Context) {},
			errors.Wrap(types.ErrNotEnoughPermissions, types.PermCreateUpsertTokenRateProposal.String()),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			app := simapp.Setup(false)
			ctx := app.NewContext(false, tmproto.Header{})

			tt.preparePerms(t, app, ctx)

			handler := gov.NewHandler(app.CustomGovKeeper)
			msg, err := govtypes.NewMsgSubmitProposal(proposerAddr, "title", "some desc", tt.content)
			require.NoError(t, err)
			_, err = handler(ctx, msg)
			require.EqualError(t, err, tt.expectedErr.Error())
		})
	}
}

func TestHandler_CreateProposalUpsertTokenRates(t *testing.T) {
	proposerAddr, err := sdk.AccAddressFromBech32("kira1alzyfq40zjsveat87jlg8jxetwqmr0a29sgd0f")
	require.NoError(t, err)

	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{
		Time: time.Now(),
	})

	// Set proposer Permissions
	proposerActor := types.NewDefaultActor(proposerAddr)
	err2 := app.CustomGovKeeper.AddWhitelistPermission(ctx, proposerActor, types.PermCreateUpsertTokenRateProposal)
	require.NoError(t, err2)

	properties := app.CustomGovKeeper.GetNetworkProperties(ctx)
	properties.ProposalEndTime = 10
	app.CustomGovKeeper.SetNetworkProperties(ctx, properties)

	handler := gov.NewHandler(app.CustomGovKeeper)
	proposal := tokenstypes.NewUpsertTokenRatesProposal(
		"btc",
		sdk.NewDec(1234),
		false,
	)
	msg, err := govtypes.NewMsgSubmitProposal(proposerAddr, "title", "some desc", proposal)
	require.NoError(t, err)
	res, err := handler(
		ctx,
		msg,
	)
	require.NoError(t, err)
	expData, _ := proto.Marshal(&govtypes.MsgSubmitProposalResponse{ProposalID: 1})
	require.Equal(t, expData, res.Data)

	savedProposal, found := app.CustomGovKeeper.GetProposal(ctx, 1)
	require.True(t, found)

	expectedSavedProposal, err := types.NewProposal(
		1,
		"title",
		"some desc",
		tokenstypes.NewUpsertTokenRatesProposal(
			"btc",
			sdk.NewDec(1234),
			false,
		),
		ctx.BlockTime(),
		ctx.BlockTime().Add(time.Second*time.Duration(properties.ProposalEndTime)),
		ctx.BlockTime().Add(time.Second*time.Duration(properties.ProposalEndTime)+
			time.Second*time.Duration(properties.ProposalEnactmentTime),
		),
		ctx.BlockHeight()+2,
		ctx.BlockHeight()+3,
	)
	require.NoError(t, err)
	require.Equal(t, expectedSavedProposal, savedProposal)

	// Next proposal ID is increased.
	id := app.CustomGovKeeper.GetNextProposalID(ctx)
	require.Equal(t, uint64(2), id)

	// Is not on finished active proposals.
	iterator := app.CustomGovKeeper.GetActiveProposalsWithFinishedVotingEndTimeIterator(ctx, ctx.BlockTime())
	require.False(t, iterator.Valid())

	ctx = ctx.WithBlockTime(ctx.BlockTime().Add(time.Minute * 10))
	iterator = app.CustomGovKeeper.GetActiveProposalsWithFinishedVotingEndTimeIterator(ctx, ctx.BlockTime())
	require.True(t, iterator.Valid())
}

func TestHandler_CreateProposalTokensWhiteBlackChange(t *testing.T) {
	proposerAddr, err := sdk.AccAddressFromBech32("kira1alzyfq40zjsveat87jlg8jxetwqmr0a29sgd0f")
	require.NoError(t, err)

	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{
		Time: time.Now(),
	})

	// Set proposer Permissions
	proposerActor := types.NewDefaultActor(proposerAddr)
	err = app.CustomGovKeeper.AddWhitelistPermission(ctx, proposerActor, types.PermCreateTokensWhiteBlackChangeProposal)
	require.NoError(t, err)

	properties := app.CustomGovKeeper.GetNetworkProperties(ctx)
	properties.ProposalEndTime = 10
	app.CustomGovKeeper.SetNetworkProperties(ctx, properties)

	handler := gov.NewHandler(app.CustomGovKeeper)
	proposal := tokenstypes.NewTokensWhiteBlackChangeProposal(
		false,
		true,
		[]string{"atom"},
	)
	msg, err := govtypes.NewMsgSubmitProposal(proposerAddr, "title", "some desc", proposal)
	require.NoError(t, err)
	res, err := handler(
		ctx,
		msg,
	)
	require.NoError(t, err)
	expData, _ := proto.Marshal(&govtypes.MsgSubmitProposalResponse{ProposalID: 1})
	require.Equal(t, expData, res.Data)

	savedProposal, found := app.CustomGovKeeper.GetProposal(ctx, 1)
	require.True(t, found)

	expectedSavedProposal, err := types.NewProposal(
		1,
		"title",
		"some desc",
		proposal,
		ctx.BlockTime(),
		ctx.BlockTime().Add(time.Second*time.Duration(properties.ProposalEndTime)),
		ctx.BlockTime().Add(time.Second*time.Duration(properties.ProposalEndTime)+
			time.Second*time.Duration(properties.ProposalEnactmentTime),
		),
		ctx.BlockHeight()+2,
		ctx.BlockHeight()+3,
	)
	require.NoError(t, err)
	require.Equal(t, expectedSavedProposal, savedProposal)

	// Next proposal ID is increased.
	id := app.CustomGovKeeper.GetNextProposalID(ctx)
	require.Equal(t, uint64(2), id)

	// Is not on finished active proposals.
	iterator := app.CustomGovKeeper.GetActiveProposalsWithFinishedVotingEndTimeIterator(ctx, ctx.BlockTime())
	require.False(t, iterator.Valid())

	ctx = ctx.WithBlockTime(ctx.BlockTime().Add(time.Minute * 10))
	iterator = app.CustomGovKeeper.GetActiveProposalsWithFinishedVotingEndTimeIterator(ctx, ctx.BlockTime())
	require.True(t, iterator.Valid())
}
