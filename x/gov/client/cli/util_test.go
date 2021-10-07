package cli_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/KiraCore/sekai/testutil/network"
	"github.com/KiraCore/sekai/x/gov/client/cli"
	"github.com/KiraCore/sekai/x/gov/types"
	customstakingcli "github.com/KiraCore/sekai/x/staking/client/cli"
	tokenscli "github.com/KiraCore/sekai/x/tokens/client/cli"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankcli "github.com/cosmos/cosmos-sdk/x/bank/client/cli"
	"github.com/stretchr/testify/require"
)

// GetRolesByAddress calls the CLI command GetCmdQueryRolesByAddress and returns the roles.
func GetRolesByAddress(t *testing.T, network *network.Network, address sdk.AccAddress) []uint64 {
	val := network.Validators[0]
	clientCtx := val.ClientCtx
	cmd := cli.GetCmdQueryRolesByAddress()
	out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, []string{
		address.String(),
	})
	require.NoError(t, err)

	var roles types.RolesByAddressResponse
	err = val.ClientCtx.JSONCodec.UnmarshalJSON(out.Bytes(), &roles)
	require.NoError(t, err)

	return roles.Roles
}

// SetCouncilor calls CLI to set address in the Councilor Registry. The Validator 1 is the caller.
func (s IntegrationTestSuite) SetCouncilor(address sdk.Address) {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx

	cmd := cli.GetTxClaimCouncilorSeatCmd()
	_, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, []string{
		fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(100))).String()),
		fmt.Sprintf("--%s=%s", cli.FlagAddress, address.String()),
		fmt.Sprintf("--%s=%s", cli.FlagMoniker, val.Moniker),
	})
	s.Require().NoError(err)

	err = s.network.WaitForNextBlock()
	s.Require().NoError(err)
}

// SendValue sends Coins from A to B using CLI.
func (s IntegrationTestSuite) SendValue(cCtx client.Context, from sdk.AccAddress, to sdk.AccAddress, coin sdk.Coin) {
	cmd := bankcli.NewSendTxCmd()
	_, err := clitestutil.ExecTestCLICmd(cCtx, cmd, []string{
		from.String(),
		to.String(),
		coin.String(),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(100))).String()),
	})
	s.Require().NoError(err)

	err = s.network.WaitForNextBlock()
	s.Require().NoError(err)
}

func (s IntegrationTestSuite) WhitelistPermission(address sdk.AccAddress, perm string) {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx
	cmd := cli.GetTxSetWhitelistPermissions()
	out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, []string{
		fmt.Sprintf("--%s=%s", customstakingcli.FlagAddr, address.String()),
		fmt.Sprintf("--%s=%s", cli.FlagPermission, perm),
		fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(100))).String()),
	})
	s.Require().NoError(err)
	fmt.Println("IntegrationTestSuite::WhitelistPermission", out.String())
}

func (s IntegrationTestSuite) VoteWithValidator0(proposalID uint64, voteOption types.VoteOption) {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx
	cmd := cli.GetTxVoteProposal()
	out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, []string{
		fmt.Sprintf("%d", proposalID),
		fmt.Sprintf("%d", voteOption),
		fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(100))).String()),
	})
	s.Require().NoError(err)
	var result sdk.TxResponse
	s.Require().NoError(val.ClientCtx.JSONCodec.UnmarshalJSON(out.Bytes(), &result))
	s.Require().NotNil(result.Height)
}

func (s IntegrationTestSuite) QueryProposals() {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx
	cmd := cli.GetCmdQueryProposals()
	output, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, []string{})
	s.Require().NoError(err, string(output.Bytes()))
}

func (s IntegrationTestSuite) QueryProposalVotes(proposalID uint64) {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx
	cmd := cli.GetCmdQueryVotes()
	_, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, []string{
		fmt.Sprintf("%d", proposalID),
	})
	s.Require().NoError(err)
}

func (s IntegrationTestSuite) QueryPoorNetworkMessages() {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx
	cmd := cli.GetCmdQueryPoorNetworkMessages()
	_, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, []string{})
	s.Require().NoError(err)
}

func (s IntegrationTestSuite) SetPoorNetworkMessages(messages string) sdk.TxResponse {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx
	cmd := cli.GetTxProposalSetPoorNetworkMessages()
	out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, []string{
		messages,
		fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", cli.FlagTitle, "title"),
		fmt.Sprintf("--%s=%s", cli.FlagDescription, "some desc"),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(100))).String()),
	})
	s.Require().NoError(err)

	err = s.network.WaitForNextBlock()
	s.Require().NoError(err)

	var result sdk.TxResponse
	s.Require().NoError(val.ClientCtx.JSONCodec.UnmarshalJSON(out.Bytes(), &result))
	s.Require().NotNil(result.Height)
	return result
}

func (s IntegrationTestSuite) SetNetworkProperties(minTxFee, maxTxFee, minValidators uint64) {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx
	cmd := cli.NewTxSetNetworkProperties()
	out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, []string{
		fmt.Sprintf("--%s=%d", cli.FlagMinTxFee, minTxFee),
		fmt.Sprintf("--%s=%d", cli.FlagMaxTxFee, maxTxFee),
		fmt.Sprintf("--%s=%d", cli.FlagMinValidators, minValidators),
		fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(100))).String()),
	})
	s.Require().NoError(err)

	err = s.network.WaitForNextBlock()
	s.Require().NoError(err)

	var result sdk.TxResponse
	s.Require().NoError(val.ClientCtx.JSONCodec.UnmarshalJSON(out.Bytes(), &result))
	s.Require().NotNil(result.Height)
	s.Require().Contains(result.RawLog, "MsgSetNetworkProperties")
}

func (s IntegrationTestSuite) SetNetworkPropertyProposal(property string, value uint64) {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx
	cmd := cli.GetTxProposalSetNetworkProperty()
	out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, []string{
		property,
		fmt.Sprintf("%d", value),
		fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
		fmt.Sprintf("--%s=%s", cli.FlagTitle, "title"),
		fmt.Sprintf("--%s=%s", cli.FlagDescription, "some desc"),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(100))).String()),
	})
	s.Require().NoError(err)

	err = s.network.WaitForNextBlock()
	s.Require().NoError(err)

	var result sdk.TxResponse
	s.Require().NoError(val.ClientCtx.JSONCodec.UnmarshalJSON(out.Bytes(), &result))
	s.Require().NotNil(result.Height)
	s.Require().Contains(result.RawLog, "SetNetworkProperty")
}

func (s IntegrationTestSuite) UpsertRate(denom string, rate string, flagFeePayments bool) sdk.TxResponse {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx
	cmd := tokenscli.GetTxUpsertTokenRateCmd()
	out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, []string{
		fmt.Sprintf("--%s=%s", tokenscli.FlagDenom, denom),
		fmt.Sprintf("--%s=%s", tokenscli.FlagRate, rate),
		fmt.Sprintf("--%s=%s", tokenscli.FlagFeePayments, strconv.FormatBool(flagFeePayments)),
		fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(100))).String()),
	})
	s.Require().NoError(err)

	err = s.network.WaitForNextBlock()
	s.Require().NoError(err)

	var result sdk.TxResponse
	s.Require().NoError(val.ClientCtx.JSONCodec.UnmarshalJSON(out.Bytes(), &result))
	s.Require().NotNil(result.Height)
	return result
}
