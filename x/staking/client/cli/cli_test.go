package cli_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	dbm "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/store/types"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/KiraCore/sekai/app"
	"github.com/KiraCore/sekai/simapp"
	"github.com/KiraCore/sekai/testutil/network"
	"github.com/KiraCore/sekai/x/staking/client/cli"
	customtypes "github.com/KiraCore/sekai/x/staking/types"
)

type IntegrationTestSuite struct {
	suite.Suite

	cfg     network.Config
	network *network.Network
}

func (s *IntegrationTestSuite) SetupSuite() {
	app.SetConfig()
	s.T().Log("setting up integration test suite")

	cfg := network.DefaultConfig()
	encodingConfig := simapp.MakeEncodingConfig()
	cfg.Codec = encodingConfig.Marshaler
	cfg.TxConfig = encodingConfig.TxConfig

	cfg.NumValidators = 1

	cfg.AppConstructor = func(val network.Validator) servertypes.Application {
		return app.NewInitApp(
			val.Ctx.Logger, dbm.NewMemDB(), nil, true, make(map[int64]bool), val.Ctx.Config.RootDir, 0,
			app.MakeEncodingConfig(),
			baseapp.SetPruning(types.NewPruningOptionsFromString(val.AppConfig.Pruning)),
			baseapp.SetMinGasPrices(val.AppConfig.MinGasPrices),
		)
	}

	s.cfg = cfg
	s.network = network.New(s.T(), cfg)

	_, err := s.network.WaitForHeight(1)
	s.Require().NoError(err)
}

func (s *IntegrationTestSuite) TearDownSuite() {
	s.T().Log("tearing down integration test suite")
	s.network.Cleanup()
}

func (s *IntegrationTestSuite) TestClaimValidatorSet_AndQueriers() {
	val := s.network.Validators[0]

	cmd := cli.GetCmdQueryValidatorByAddress()
	cmd.SetArgs(
		[]string{
			fmt.Sprintf("--%s=%s", cli.FlagValAddr, val.ValAddress.String()),
		},
	)

	_, out := testutil.ApplyMockIO(cmd)
	clientCtx := val.ClientCtx.WithOutput(out).WithOutputFormat("json")

	ctx := context.Background()
	ctx = context.WithValue(ctx, client.ClientContextKey, &clientCtx)

	err := cmd.ExecuteContext(ctx)
	s.Require().NoError(err)

	var respValidator customtypes.Validator
	clientCtx.JSONMarshaler.MustUnmarshalJSON(out.Bytes(), &respValidator)

	s.Require().Equal(val.Moniker, respValidator.Moniker)
	s.Require().Equal("the Website", respValidator.Website)
	s.Require().Equal("The social", respValidator.Social)
	s.Require().Equal("The Identity", respValidator.Identity)
	s.Require().Equal(sdk.NewDec(1), respValidator.Commission)
	s.Require().Equal(val.ValAddress, respValidator.ValKey)

	pubkey, err := sdk.GetPubKeyFromBech32(sdk.Bech32PubKeyTypeConsPub, respValidator.PubKey)
	s.Require().NoError(err)
	s.Require().Equal(val.PubKey, pubkey)

	// Query by Acc Addrs.
	cmd = cli.GetCmdQueryValidatorByAddress()
	cmd.SetArgs(
		[]string{
			fmt.Sprintf("--%s=%s", cli.FlagAddr, val.Address.String()),
		},
	)

	out.Reset()

	clientCtx = clientCtx.WithOutputFormat("json")
	err = cmd.ExecuteContext(ctx)
	s.Require().NoError(err)

	clientCtx.JSONMarshaler.MustUnmarshalJSON(out.Bytes(), &respValidator)

	s.Require().Equal(val.Moniker, respValidator.Moniker)
	s.Require().Equal("the Website", respValidator.Website)
	s.Require().Equal("The social", respValidator.Social)
	s.Require().Equal("The Identity", respValidator.Identity)
	s.Require().Equal(sdk.NewDec(1), respValidator.Commission)
	s.Require().Equal(val.ValAddress, respValidator.ValKey)

	pubkey, err = sdk.GetPubKeyFromBech32(sdk.Bech32PubKeyTypeConsPub, respValidator.PubKey)
	s.Require().NoError(err)
	s.Require().Equal(val.PubKey, pubkey)

	// Query by moniker.
	cmd = cli.GetCmdQueryValidatorByAddress()
	cmd.SetArgs(
		[]string{
			fmt.Sprintf("--%s=%s", cli.FlagMoniker, val.Moniker),
		},
	)

	out.Reset()

	clientCtx = clientCtx.WithOutputFormat("json")
	err = cmd.ExecuteContext(ctx)
	s.Require().NoError(err)

	clientCtx.JSONMarshaler.MustUnmarshalJSON(out.Bytes(), &respValidator)

	s.Require().Equal(val.Moniker, respValidator.Moniker)
	s.Require().Equal("the Website", respValidator.Website)
	s.Require().Equal("The social", respValidator.Social)
	s.Require().Equal("The Identity", respValidator.Identity)
	s.Require().Equal(sdk.NewDec(1), respValidator.Commission)
	s.Require().Equal(val.ValAddress, respValidator.ValKey)

	pubkey, err = sdk.GetPubKeyFromBech32(sdk.Bech32PubKeyTypeConsPub, respValidator.PubKey)
	s.Require().NoError(err)
	s.Require().Equal(val.PubKey, pubkey)
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}
