package cli_test

import (
	"fmt"

	"github.com/KiraCore/sekai/x/gov/client/cli"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/testutil"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (s IntegrationTestSuite) TestTxRegisterIdentityRecords() {
	val := s.network.Validators[0]
	cmd := cli.GetTxRegisterIdentityRecords()

	infosFile := testutil.WriteToNewTempFile(s.T(), `
		{
			"key1": "value1",
			"key2": "value2"
		}
	`)

	clientCtx := val.ClientCtx.WithOutputFormat("json")
	out, err := clitestutil.ExecTestCLICmd(
		clientCtx,
		cmd,
		[]string{
			fmt.Sprintf("--%s=%s", cli.FlagInfosFile, infosFile.Name()),
			fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
			fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
			fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
			fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(100))).String()),
		},
	)
	s.Require().NoError(err)

	err = s.network.WaitForNextBlock()
	s.Require().NoError(err)
	fmt.Println("out", out)
}

func (s IntegrationTestSuite) TestTxEditIdentityRecord() {
	val := s.network.Validators[0]
	cmd := cli.GetTxDeleteIdentityRecords()

	clientCtx := val.ClientCtx.WithOutputFormat("json")
	out, err := clitestutil.ExecTestCLICmd(
		clientCtx,
		cmd,
		[]string{
			fmt.Sprintf("--%s=%s", cli.FlagKeys, "key1,key2"),
			fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
			fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
			fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
			fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(100))).String()),
		},
	)
	s.Require().NoError(err)

	err = s.network.WaitForNextBlock()
	s.Require().NoError(err)
	fmt.Println("out", out)
}

func (s IntegrationTestSuite) TestTxRequestIdentityRecordsVerify() {
	val := s.network.Validators[0]
	cmd := cli.GetTxRequestIdentityRecordsVerify()

	clientCtx := val.ClientCtx.WithOutputFormat("json")
	out, err := clitestutil.ExecTestCLICmd(
		clientCtx,
		cmd,
		[]string{
			fmt.Sprintf("--%s=%s", cli.FlagVerifier, val.Address.String()),
			fmt.Sprintf("--%s=%s", cli.FlagRecordIds, "1"),
			fmt.Sprintf("--%s=%s", cli.FlagTip, "10stake"),
			fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
			fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
			fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
			fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(100))).String()),
		},
	)
	s.Require().NoError(err)

	err = s.network.WaitForNextBlock()
	s.Require().NoError(err)
	fmt.Println("out", out)
}

func (s IntegrationTestSuite) TestTxHandleIdentityRecordsVerifyRequest() {
	val := s.network.Validators[0]
	cmd := cli.GetTxHandleIdentityRecordsVerifyRequest()

	clientCtx := val.ClientCtx.WithOutputFormat("json")
	out, err := clitestutil.ExecTestCLICmd(
		clientCtx,
		cmd,
		[]string{
			"1",
			fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
			fmt.Sprintf("--%s=true", cli.FlagApprove),
			fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
			fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
			fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(100))).String()),
		},
	)
	s.Require().NoError(err)

	err = s.network.WaitForNextBlock()
	s.Require().NoError(err)
	fmt.Println("out", out)
}

func (s IntegrationTestSuite) TestTxCancelIdentityRecordsVerifyRequest() {
	val := s.network.Validators[0]
	cmd := cli.GetTxCancelIdentityRecordsVerifyRequest()

	clientCtx := val.ClientCtx.WithOutputFormat("json")
	out, err := clitestutil.ExecTestCLICmd(
		clientCtx,
		cmd,
		[]string{
			"1",
			fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
			fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
			fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
			fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(100))).String()),
		},
	)
	s.Require().NoError(err)

	err = s.network.WaitForNextBlock()
	s.Require().NoError(err)
	fmt.Println("out", out)
}

func (s IntegrationTestSuite) TestCmdQueryIdentityRecord() {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx
	cmd := cli.GetCmdQueryIdentityRecord()
	_, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, []string{
		fmt.Sprintf("%d", 1),
	})
	s.Require().NoError(err)
}

func (s IntegrationTestSuite) TestCmdQueryIdentityRecordByAddress() {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx
	cmd := cli.GetCmdQueryIdentityRecordByAddress()
	_, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, []string{
		fmt.Sprintf("%s", val.Address.String()),
	})
	s.Require().NoError(err)
}

func (s IntegrationTestSuite) TestCmdQueryAllIdentityRecords() {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx
	cmd := cli.GetCmdQueryAllIdentityRecords()
	_, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, []string{})
	s.Require().NoError(err)
}

func (s IntegrationTestSuite) TestCmdQueryIdentityRecordVerifyRequest() {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx
	cmd := cli.GetCmdQueryIdentityRecordVerifyRequest()
	_, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, []string{
		fmt.Sprintf("%d", 1),
	})
	s.Require().NoError(err)
}

func (s IntegrationTestSuite) TestCmdQueryIdentityRecordVerifyRequestsByRequester() {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx
	cmd := cli.GetCmdQueryIdentityRecordVerifyRequestsByRequester()
	_, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, []string{
		fmt.Sprintf("%s", val.Address.String()),
	})
	s.Require().NoError(err)
}

func (s IntegrationTestSuite) TestCmdQueryIdentityRecordVerifyRequestsByApprover() {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx
	cmd := cli.GetCmdQueryIdentityRecordVerifyRequestsByApprover()
	_, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, []string{
		fmt.Sprintf("%s", val.Address.String()),
	})
	s.Require().NoError(err)
}

func (s IntegrationTestSuite) TestCmdQueryAllIdentityRecordVerifyRequests() {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx
	cmd := cli.GetCmdQueryAllIdentityRecordVerifyRequests()
	_, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, []string{})
	s.Require().NoError(err)
}
