package cli_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	dbm "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/keys"
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
	s.T().Log("setting up integration test suite")

	cfg := network.DefaultConfig()
	encodingConfig := simapp.MakeEncodingConfig()
	cfg.Codec = encodingConfig.Marshaler
	cfg.TxConfig = encodingConfig.TxConfig

	cfg.NumValidators = 1

	cfg.AppConstructor = func(val network.Validator) servertypes.Application {
		return app.NewInitApp(
			val.Ctx.Logger, dbm.NewMemDB(), nil, true, make(map[int64]bool), val.Ctx.Config.RootDir, 0,
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

	cmd := cli.GetTxClaimValidatorCmd()
	_, out := testutil.ApplyMockIO(cmd)
	clientCtx := val.ClientCtx.WithOutput(out)

	ctx := context.Background()
	ctx = context.WithValue(ctx, client.ClientContextKey, &clientCtx)

	pubKey := "kiravalconspub1zcjduepqylc5k8r40azmw0xt7hjugr4mr5w2am7jw77ux5w6s8hpjxyrjjsq4xg7em"
	cmd.SetArgs(
		[]string{
			fmt.Sprintf("--%s=%s", cli.FlagMoniker, "Moniker"),
			fmt.Sprintf("--%s=%s", cli.FlagWebsite, "Website"),
			fmt.Sprintf("--%s=%s", cli.FlagSocial, "Social"),
			fmt.Sprintf("--%s=%s", cli.FlagIdentity, "Identity"),
			fmt.Sprintf("--%s=%s", cli.FlagComission, "10"),
			fmt.Sprintf("--%s=%s", keys.FlagPublicKey, pubKey),
			fmt.Sprintf("--%s=%s", cli.FlagValKey, val.ValAddress.String()),
			fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Moniker),
			fmt.Sprintf("--%s", flags.FlagSkipConfirmation),
			fmt.Sprintf("--%s=%s", flags.FlagFees, "2stake"),
		},
	)

	err := cmd.ExecuteContext(ctx)
	s.Require().NoError(err)

	height, err := s.network.LatestHeight()
	s.Require().NoError(err)

	_, err = s.network.WaitForHeight(height + 2)
	s.Require().NoError(err)

	query := cli.GetCmdQueryValidatorByAddress()
	query.SetArgs(
		[]string{
			fmt.Sprintf("--%s=%s", cli.FlagValAddr, val.ValAddress.String()),
		},
	)

	out.Reset()

	clientCtx = clientCtx.WithOutputFormat("json")
	err = query.ExecuteContext(ctx)
	s.Require().NoError(err)

	var respValidator customtypes.Validator
	clientCtx.JSONMarshaler.MustUnmarshalJSON(out.Bytes(), &respValidator)

	s.Require().Equal("Moniker", respValidator.Moniker)
	s.Require().Equal("Website", respValidator.Website)
	s.Require().Equal("Social", respValidator.Social)
	s.Require().Equal("Identity", respValidator.Identity)
	s.Require().Equal(sdk.NewDec(10), respValidator.Commission)
	s.Require().Equal(val.ValAddress, respValidator.ValKey)
	s.Require().Equal(pubKey, respValidator.PubKey)

	// Query by Acc Addrs.
	query = cli.GetCmdQueryValidatorByAddress()
	query.SetArgs(
		[]string{
			fmt.Sprintf("--%s=%s", cli.FlagAddr, val.Address.String()),
		},
	)

	out.Reset()

	clientCtx = clientCtx.WithOutputFormat("json")
	err = query.ExecuteContext(ctx)
	s.Require().NoError(err)

	clientCtx.JSONMarshaler.MustUnmarshalJSON(out.Bytes(), &respValidator)

	s.Require().Equal("Moniker", respValidator.Moniker)
	s.Require().Equal("Website", respValidator.Website)
	s.Require().Equal("Social", respValidator.Social)
	s.Require().Equal("Identity", respValidator.Identity)
	s.Require().Equal(sdk.NewDec(10), respValidator.Commission)
	s.Require().Equal(val.ValAddress, respValidator.ValKey)
	s.Require().Equal(pubKey, respValidator.PubKey)

	// Query by moniker.
	query = cli.GetCmdQueryValidatorByAddress()
	query.SetArgs(
		[]string{
			fmt.Sprintf("--%s=%s", cli.FlagMoniker, val.Moniker),
		},
	)

	out.Reset()

	clientCtx = clientCtx.WithOutputFormat("json")
	err = query.ExecuteContext(ctx)
	s.Require().NoError(err)

	clientCtx.JSONMarshaler.MustUnmarshalJSON(out.Bytes(), &respValidator)

	s.Require().Equal("Moniker", respValidator.Moniker)
	s.Require().Equal("Website", respValidator.Website)
	s.Require().Equal("Social", respValidator.Social)
	s.Require().Equal("Identity", respValidator.Identity)
	s.Require().Equal(sdk.NewDec(10), respValidator.Commission)
	s.Require().Equal(val.ValAddress, respValidator.ValKey)
	s.Require().Equal(pubKey, respValidator.PubKey)
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}
