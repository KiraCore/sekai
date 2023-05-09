package cli

import (
	"fmt"
	"strconv"
	"strings"

	govtypes "github.com/KiraCore/sekai/x/gov/types"
	"github.com/KiraCore/sekai/x/layer2/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
)

const (
	FlagTitle              = "title"
	FlagDescription        = "description"
	FlagDappName           = "dapp-name"
	FlagDappDescription    = "dapp-description"
	FlagDenom              = "denom"
	FlagWebsite            = "website"
	FlagLogo               = "logo"
	FlagSocial             = "social"
	FlagDocs               = "docs"
	FlagControllerRoles    = "controller-roles"
	FlagControllerAccounts = "controller-accounts"
	FlagBinaryInfo         = "binary-info"
	FlagLpPoolConfig       = "lp-pool-config"
	FlagIssuranceConfig    = "issurance-config"
	FlagUpdateTimeMax      = "update-time-max"
	FlagExecutorsMin       = "executors-min"
	FlagExecutorsMax       = "executors-max"
	FlagVerifiersMin       = "verifiers-min"
	FlagDappStatus         = "dapp-status"
	FlagBond               = "bond"
	FlagVoteQuorum         = "vote-quorum"
	FlagVotePeriod         = "vote-period"
	FlagVoteEnactment      = "vote-enactment"
	FlagAddr               = "addr"
	FlagAmount             = "amount"
)

// NewTxCmd returns a root CLI command handler for all x/bank transaction commands.
func NewTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Tokens sub commands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(
		GetTxCreateDappProposalCmd(),
		GetTxBondDappProposalCmd(),
		GetTxReclaimDappBondProposalCmd(),
		GetTxJoinDappVerifierWithBondCmd(),
		GetTxExitDappCmd(),
		GetTxExecuteDappTxCmd(),
		GetTxDenounceLeaderTxCmd(),
		GetTxTransitionDappCmd(),
		GetTxApproveDappTransitionCmd(),
		GetTxRejectDappTransitionCmd(),
		GetTxProposalJoinDappCmd(),
		GetTxProposalUpsertDappCmd(),
	)

	return txCmd
}

// GetTxCreateDappProposalCmd implement cli command for MsgCreateDappProposal
func GetTxCreateDappProposalCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-dapp-proposal",
		Short: "Submit a proposal to launch a dapp",
		Args:  cobra.MinimumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			name, err := cmd.Flags().GetString(FlagDappName)
			if err != nil {
				return fmt.Errorf("invalid name: %w", err)
			}
			denom, err := cmd.Flags().GetString(FlagDenom)
			if err != nil {
				return fmt.Errorf("invalid denom: %w", err)
			}
			description, err := cmd.Flags().GetString(FlagDappDescription)
			if err != nil {
				return fmt.Errorf("invalid description: %w", err)
			}
			website, err := cmd.Flags().GetString(FlagWebsite)
			if err != nil {
				return fmt.Errorf("invalid website: %w", err)
			}
			logo, err := cmd.Flags().GetString(FlagLogo)
			if err != nil {
				return fmt.Errorf("invalid logo: %w", err)
			}
			social, err := cmd.Flags().GetString(FlagSocial)
			if err != nil {
				return fmt.Errorf("invalid social: %w", err)
			}
			docs, err := cmd.Flags().GetString(FlagDocs)
			if err != nil {
				return fmt.Errorf("invalid docs: %w", err)
			}
			bondStr, err := cmd.Flags().GetString(FlagBond)
			if err != nil {
				return fmt.Errorf("invalid bonds: %w", err)
			}
			bond, err := sdk.ParseCoinNormalized(bondStr)
			if err != nil {
				return fmt.Errorf("invalid bonds: %w", err)
			}

			voteQuorum, err := cmd.Flags().GetUint64(FlagVoteQuorum)
			if err != nil {
				return fmt.Errorf("invalid vote quorum: %w", err)
			}
			votePeriod, err := cmd.Flags().GetUint64(FlagVotePeriod)
			if err != nil {
				return fmt.Errorf("invalid vote period: %w", err)
			}
			voteEnactment, err := cmd.Flags().GetUint64(FlagVoteEnactment)
			if err != nil {
				return fmt.Errorf("invalid vote period: %w", err)
			}

			ctrlRolesStr, err := cmd.Flags().GetString(FlagControllerRoles)
			if err != nil {
				return fmt.Errorf("invalid controller roles: %w", err)
			}
			ctrlAccountsStr, err := cmd.Flags().GetString(FlagControllerAccounts)
			if err != nil {
				return fmt.Errorf("invalid controller accounts: %w", err)
			}

			ctrlRoles, ctrlAccounts, err := parseRolesAccounts(ctrlRolesStr, ctrlAccountsStr)
			if err != nil {
				return fmt.Errorf("invalid controller role/accounts: %w", err)
			}

			issuranceStr, err := cmd.Flags().GetString(FlagIssuranceConfig)
			if err != nil {
				return fmt.Errorf("invalid issurance config: %w", err)
			}
			issurance := types.IssuranceConfig{}
			err = clientCtx.Codec.UnmarshalJSON([]byte(issuranceStr), &issurance)
			if err != nil {
				return fmt.Errorf("invalid issurance config: %w", err)
			}

			lpPoolConfigStr, err := cmd.Flags().GetString(FlagLpPoolConfig)
			if err != nil {
				return fmt.Errorf("invalid lp pool config: %w", err)
			}
			lpPoolConfig := types.LpPoolConfig{}
			err = clientCtx.Codec.UnmarshalJSON([]byte(lpPoolConfigStr), &lpPoolConfig)
			if err != nil {
				return fmt.Errorf("invalid lp pool config: %w", err)
			}

			binaryInfoStr, err := cmd.Flags().GetString(FlagBinaryInfo)
			if err != nil {
				return fmt.Errorf("invalid binaryInfo: %w", err)
			}
			binaryInfo := types.BinaryInfo{}
			err = clientCtx.Codec.UnmarshalJSON([]byte(binaryInfoStr), &binaryInfo)
			if err != nil {
				return fmt.Errorf("invalid binaryInfo: %w", err)
			}

			statusStr, err := cmd.Flags().GetString(FlagDappStatus)
			if err != nil {
				return fmt.Errorf("invalid status: %w", err)
			}

			updateMaxTime, err := cmd.Flags().GetUint64(FlagUpdateTimeMax)
			if err != nil {
				return fmt.Errorf("invalid updateMaxTime: %w", err)
			}
			executorsMin, err := cmd.Flags().GetUint64(FlagExecutorsMin)
			if err != nil {
				return fmt.Errorf("invalid executorsMin: %w", err)
			}
			executorsMax, err := cmd.Flags().GetUint64(FlagExecutorsMax)
			if err != nil {
				return fmt.Errorf("invalid executorsMax: %w", err)
			}
			verifiersMin, err := cmd.Flags().GetUint64(FlagVerifiersMin)
			if err != nil {
				return fmt.Errorf("invalid verifiersMin: %w", err)
			}

			msg := &types.MsgCreateDappProposal{
				Sender: clientCtx.GetFromAddress().String(),
				Bond:   bond,
				Dapp: types.Dapp{
					Name:        name,
					Denom:       denom,
					Description: description,
					Website:     website,
					Logo:        logo,
					Social:      social,
					Docs:        docs,
					Controllers: types.Controllers{
						Whitelist: types.AccountRange{
							Roles:     ctrlRoles,
							Addresses: ctrlAccounts,
						},
					},
					Bin:           []types.BinaryInfo{binaryInfo},
					Pool:          lpPoolConfig,
					Issurance:     issurance,
					UpdateTimeMax: updateMaxTime,
					ExecutorsMin:  executorsMin,
					ExecutorsMax:  executorsMax,
					VerifiersMin:  verifiersMin,
					Status:        types.DappStatus(types.SessionStatus_value[statusStr]),
					VoteQuorum:    voteQuorum,
					VotePeriod:    votePeriod,
					VoteEnactment: voteEnactment,
				},
			}

			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	_ = cmd.MarkFlagRequired(flags.FlagFrom)

	cmd.Flags().String(FlagDappName, "", "Dapp name")
	cmd.Flags().String(FlagDenom, "", "Dapp denom")
	cmd.Flags().String(FlagDappDescription, "", "Dapp description")
	cmd.Flags().String(FlagWebsite, "", "Dapp website")
	cmd.Flags().String(FlagLogo, "", "Dapp logo")
	cmd.Flags().String(FlagSocial, "", "Dapp social")
	cmd.Flags().String(FlagDocs, "", "Dapp docs")
	cmd.Flags().String(FlagBond, "", "Initial bond deposit for dapp")
	cmd.Flags().Uint64(FlagVoteQuorum, 0, "vote quorum of the dapp")
	cmd.Flags().Uint64(FlagVotePeriod, 0, "vote period of the dapp")
	cmd.Flags().Uint64(FlagVoteEnactment, 0, "vote enactment of the dapp")
	cmd.Flags().String(FlagControllerRoles, "", "controller roles on the dapp.")
	cmd.Flags().String(FlagControllerAccounts, "", "controller accounts on the dapp.")
	cmd.Flags().String(FlagIssuranceConfig, "{}", "dapp issurance config.")
	cmd.Flags().String(FlagLpPoolConfig, "{}", "dapp lp config.")
	cmd.Flags().String(FlagBinaryInfo, "{}", "dapp binary info.")
	cmd.Flags().String(FlagDappStatus, "{}", "dapp status.")
	cmd.Flags().Uint64(FlagUpdateTimeMax, 0, "dapp update time max")
	cmd.Flags().Uint64(FlagExecutorsMin, 0, "dapp executors min")
	cmd.Flags().Uint64(FlagExecutorsMax, 0, "dapp executors max")
	cmd.Flags().Uint64(FlagVerifiersMin, 0, "dapp verifiers min")

	return cmd
}

// GetTxBondDappProposalCmd implement cli command for MsgBondDappProposal
func GetTxBondDappProposalCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bond-dapp-proposal",
		Short: "Bond on dapp proposal",
		Args:  cobra.MinimumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			name, err := cmd.Flags().GetString(FlagDappName)
			if err != nil {
				return fmt.Errorf("invalid name: %w", err)
			}
			bondStr, err := cmd.Flags().GetString(FlagBond)
			if err != nil {
				return fmt.Errorf("invalid bonds: %w", err)
			}
			bond, err := sdk.ParseCoinNormalized(bondStr)
			if err != nil {
				return fmt.Errorf("invalid bonds: %w", err)
			}

			msg := &types.MsgBondDappProposal{
				Sender:   clientCtx.GetFromAddress().String(),
				DappName: name,
				Bond:     bond,
			}

			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	_ = cmd.MarkFlagRequired(flags.FlagFrom)

	cmd.Flags().String(FlagDappName, "", "Dapp name")
	cmd.Flags().String(FlagBond, "", "Initial bond deposit for dapp")

	return cmd
}

// GetTxReclaimDappBondProposalCmd implement cli command for MsgReclaimDappBondProposal
func GetTxReclaimDappBondProposalCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reclaim-dapp-proposal",
		Short: "Reclaim from dapp proposal",
		Args:  cobra.MinimumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			name, err := cmd.Flags().GetString(FlagDappName)
			if err != nil {
				return fmt.Errorf("invalid name: %w", err)
			}
			bondStr, err := cmd.Flags().GetString(FlagBond)
			if err != nil {
				return fmt.Errorf("invalid bonds: %w", err)
			}
			bond, err := sdk.ParseCoinNormalized(bondStr)
			if err != nil {
				return fmt.Errorf("invalid bonds: %w", err)
			}

			msg := &types.MsgReclaimDappBondProposal{
				Sender:   clientCtx.GetFromAddress().String(),
				DappName: name,
				Bond:     bond,
			}

			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	_ = cmd.MarkFlagRequired(flags.FlagFrom)

	cmd.Flags().String(FlagDappName, "", "Dapp name")
	cmd.Flags().String(FlagBond, "", "Initial bond deposit for dapp")

	return cmd
}

// GetTxJoinDappVerifierWithBondCmd implement cli command for MsgJoinDappVerifierWithBond
func GetTxJoinDappVerifierWithBondCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "join-dapp-verifier-with-bond [dapp-name] [interx]",
		Short: "Join dapp verifier with bond",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := &types.MsgJoinDappVerifierWithBond{
				Sender:   clientCtx.GetFromAddress().String(),
				DappName: args[0],
				Interx:   args[1],
			}

			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	_ = cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

// GetTxExitDappCmd implement cli command for MsgExitDapp
func GetTxExitDappCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "exit-dapp [dapp-name]",
		Short: "Send request to exit dapp",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := &types.MsgExitDapp{
				Sender:   clientCtx.GetFromAddress().String(),
				DappName: args[0],
			}

			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	_ = cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

// GetTxExecuteDappTxCmd implement cli command for MsgExecuteDappTx
func GetTxExecuteDappTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "execute-dapp-tx [dapp-name] [gateway]",
		Short: "Send signal to start dapp execution",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := &types.MsgExecuteDappTx{
				Sender:   clientCtx.GetFromAddress().String(),
				DappName: args[0],
				Gateway:  args[1],
			}

			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	_ = cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

// GetTxDenounceLeaderTxCmd implement cli command for MsgDenounceLeaderTx
func GetTxDenounceLeaderTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "denounce-leader [dapp-name] [leader] [denounce-txt] [version]",
		Short: "Send leader denounce transaction",
		Args:  cobra.MinimumNArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := &types.MsgDenounceLeaderTx{
				Sender:       clientCtx.GetFromAddress().String(),
				DappName:     args[0],
				Leader:       args[1],
				DenounceText: args[2],
				Version:      args[3],
			}

			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	_ = cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

// GetTxTransitionDappCmd implement cli command for MsgTransitionDappTx
func GetTxTransitionDappCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "transition-dapp [dapp-name] [status-hash] [version]",
		Short: "Send dapp transition message",
		Args:  cobra.MinimumNArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := &types.MsgTransitionDappTx{
				Sender:     clientCtx.GetFromAddress().String(),
				DappName:   args[0],
				StatusHash: args[1],
				Version:    args[2],
			}

			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	_ = cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

// GetTxApproveDappTransitionCmd implement cli command for MsgApproveDappTransitionTx
func GetTxApproveDappTransitionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "approve-dapp-transition [dapp-name] [version]",
		Short: "Send dapp transition approval message",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := &types.MsgApproveDappTransitionTx{
				Sender:   clientCtx.GetFromAddress().String(),
				DappName: args[0],
				Version:  args[1],
			}

			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	_ = cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

// GetTxRejectDappTransitionCmd implement cli command for MsgRejectDappTransitionTx
func GetTxRejectDappTransitionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reject-dapp-transition [dapp-name] [version]",
		Short: "Send dapp transition reject message",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := &types.MsgRejectDappTransitionTx{
				Sender:   clientCtx.GetFromAddress().String(),
				DappName: args[0],
				Version:  args[1],
			}

			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	_ = cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

// GetTxProposalJoinDappCmd implement cli command for ProposalJoinDapp
func GetTxProposalJoinDappCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "proposal-join-dapp [dapp-name] [is-executor] [is-verifier] [interx]",
		Short: "Create a proposal to join dapp",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			title, err := cmd.Flags().GetString(FlagTitle)
			if err != nil {
				return fmt.Errorf("invalid title: %w", err)
			}
			description, err := cmd.Flags().GetString(FlagDescription)
			if err != nil {
				return fmt.Errorf("invalid description: %w", err)
			}

			executor, err := strconv.ParseBool(args[1])
			if err != nil {
				return fmt.Errorf("invalid executor flag: %w", err)
			}

			verifier, err := strconv.ParseBool(args[2])
			if err != nil {
				return fmt.Errorf("invalid verifier flag: %w", err)
			}

			msg, err := govtypes.NewMsgSubmitProposal(
				clientCtx.FromAddress,
				title,
				description,
				&types.ProposalJoinDapp{
					Sender:   clientCtx.GetFromAddress().String(),
					DappName: args[0],
					Executor: executor,
					Verifier: verifier,
					Interx:   args[3],
				},
			)
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(FlagTitle, "", "The title of the proposal.")
	cmd.MarkFlagRequired(FlagTitle)
	cmd.Flags().String(FlagDescription, "", "The description of the proposal, it can be a url, some text, etc.")
	cmd.MarkFlagRequired(FlagDescription)

	flags.AddTxFlagsToCmd(cmd)
	_ = cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

// GetTxProposalUpsertDappCmd implement cli command for ProposalUpsertDapp
func GetTxProposalUpsertDappCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "proposal-upsert-dapp",
		Short: "Create a proposal to upsert dapp",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			title, err := cmd.Flags().GetString(FlagTitle)
			if err != nil {
				return fmt.Errorf("invalid title: %w", err)
			}
			description, err := cmd.Flags().GetString(FlagDescription)
			if err != nil {
				return fmt.Errorf("invalid description: %w", err)
			}

			name, err := cmd.Flags().GetString(FlagDappName)
			if err != nil {
				return fmt.Errorf("invalid name: %w", err)
			}
			denom, err := cmd.Flags().GetString(FlagDenom)
			if err != nil {
				return fmt.Errorf("invalid denom: %w", err)
			}
			dappDescription, err := cmd.Flags().GetString(FlagDappDescription)
			if err != nil {
				return fmt.Errorf("invalid description: %w", err)
			}
			website, err := cmd.Flags().GetString(FlagWebsite)
			if err != nil {
				return fmt.Errorf("invalid website: %w", err)
			}
			logo, err := cmd.Flags().GetString(FlagLogo)
			if err != nil {
				return fmt.Errorf("invalid logo: %w", err)
			}
			social, err := cmd.Flags().GetString(FlagSocial)
			if err != nil {
				return fmt.Errorf("invalid social: %w", err)
			}
			docs, err := cmd.Flags().GetString(FlagDocs)
			if err != nil {
				return fmt.Errorf("invalid docs: %w", err)
			}

			ctrlRolesStr, err := cmd.Flags().GetString(FlagControllerRoles)
			if err != nil {
				return fmt.Errorf("invalid controller roles: %w", err)
			}
			ctrlAccountsStr, err := cmd.Flags().GetString(FlagControllerAccounts)
			if err != nil {
				return fmt.Errorf("invalid controller accounts: %w", err)
			}

			ctrlRoles, ctrlAccounts, err := parseRolesAccounts(ctrlRolesStr, ctrlAccountsStr)
			if err != nil {
				return fmt.Errorf("invalid controller role/accounts: %w", err)
			}

			issuranceStr, err := cmd.Flags().GetString(FlagIssuranceConfig)
			if err != nil {
				return fmt.Errorf("invalid issurance config: %w", err)
			}
			issurance := types.IssuranceConfig{}
			err = clientCtx.Codec.UnmarshalJSON([]byte(issuranceStr), &issurance)
			if err != nil {
				return fmt.Errorf("invalid issurance config: %w", err)
			}

			lpPoolConfigStr, err := cmd.Flags().GetString(FlagLpPoolConfig)
			if err != nil {
				return fmt.Errorf("invalid lp pool config: %w", err)
			}
			lpPoolConfig := types.LpPoolConfig{}
			err = clientCtx.Codec.UnmarshalJSON([]byte(lpPoolConfigStr), &lpPoolConfig)
			if err != nil {
				return fmt.Errorf("invalid lp pool config: %w", err)
			}

			binaryInfoStr, err := cmd.Flags().GetString(FlagBinaryInfo)
			if err != nil {
				return fmt.Errorf("invalid binaryInfo: %w", err)
			}
			binaryInfo := types.BinaryInfo{}
			err = clientCtx.Codec.UnmarshalJSON([]byte(binaryInfoStr), &binaryInfo)
			if err != nil {
				return fmt.Errorf("invalid binaryInfo: %w", err)
			}

			statusStr, err := cmd.Flags().GetString(FlagDappStatus)
			if err != nil {
				return fmt.Errorf("invalid status: %w", err)
			}

			updateMaxTime, err := cmd.Flags().GetUint64(FlagUpdateTimeMax)
			if err != nil {
				return fmt.Errorf("invalid updateMaxTime: %w", err)
			}
			executorsMin, err := cmd.Flags().GetUint64(FlagExecutorsMin)
			if err != nil {
				return fmt.Errorf("invalid executorsMin: %w", err)
			}
			executorsMax, err := cmd.Flags().GetUint64(FlagExecutorsMax)
			if err != nil {
				return fmt.Errorf("invalid executorsMax: %w", err)
			}
			verifiersMin, err := cmd.Flags().GetUint64(FlagVerifiersMin)
			if err != nil {
				return fmt.Errorf("invalid verifiersMin: %w", err)
			}
			voteQuorum, err := cmd.Flags().GetUint64(FlagVoteQuorum)
			if err != nil {
				return fmt.Errorf("invalid vote quorum: %w", err)
			}
			votePeriod, err := cmd.Flags().GetUint64(FlagVotePeriod)
			if err != nil {
				return fmt.Errorf("invalid vote period: %w", err)
			}
			voteEnactment, err := cmd.Flags().GetUint64(FlagVoteEnactment)
			if err != nil {
				return fmt.Errorf("invalid vote period: %w", err)
			}

			msg, err := govtypes.NewMsgSubmitProposal(
				clientCtx.FromAddress,
				title,
				description,
				&types.ProposalUpsertDapp{
					Sender: clientCtx.GetFromAddress().String(),
					Dapp: types.Dapp{
						Name:        name,
						Denom:       denom,
						Description: dappDescription,
						Website:     website,
						Logo:        logo,
						Social:      social,
						Docs:        docs,
						Controllers: types.Controllers{
							Whitelist: types.AccountRange{
								Roles:     ctrlRoles,
								Addresses: ctrlAccounts,
							},
						},
						Bin:           []types.BinaryInfo{binaryInfo},
						Pool:          lpPoolConfig,
						Issurance:     issurance,
						UpdateTimeMax: updateMaxTime,
						ExecutorsMin:  executorsMin,
						ExecutorsMax:  executorsMax,
						VerifiersMin:  verifiersMin,
						Status:        types.DappStatus(types.SessionStatus_value[statusStr]),
						VoteQuorum:    voteQuorum,
						VotePeriod:    votePeriod,
						VoteEnactment: voteEnactment,
					},
				},
			)
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(FlagTitle, "", "The title of the proposal.")
	cmd.MarkFlagRequired(FlagTitle)
	cmd.Flags().String(FlagDescription, "", "The description of the proposal, it can be a url, some text, etc.")
	cmd.MarkFlagRequired(FlagDescription)

	cmd.Flags().String(FlagDappName, "", "Dapp name")
	cmd.Flags().String(FlagDenom, "", "Dapp denom")
	cmd.Flags().String(FlagDappDescription, "", "Dapp description")
	cmd.Flags().String(FlagWebsite, "", "Dapp website")
	cmd.Flags().String(FlagLogo, "", "Dapp logo")
	cmd.Flags().String(FlagSocial, "", "Dapp social")
	cmd.Flags().String(FlagDocs, "", "Dapp docs")
	cmd.Flags().String(FlagBond, "", "Initial bond deposit for dapp")
	cmd.Flags().Uint64(FlagVoteQuorum, 0, "vote quorum of the dapp")
	cmd.Flags().Uint64(FlagVotePeriod, 0, "vote period of the dapp")
	cmd.Flags().Uint64(FlagVoteEnactment, 0, "vote enactment of the dapp")
	cmd.Flags().String(FlagControllerRoles, "", "controller roles on the dapp.")
	cmd.Flags().String(FlagControllerAccounts, "", "controller accounts on the dapp.")
	cmd.Flags().String(FlagIssuranceConfig, "{}", "dapp issurance config.")
	cmd.Flags().String(FlagLpPoolConfig, "{}", "dapp lp config.")
	cmd.Flags().String(FlagBinaryInfo, "{}", "dapp binary info.")
	cmd.Flags().String(FlagDappStatus, "{}", "dapp status.")
	cmd.Flags().Uint64(FlagUpdateTimeMax, 0, "dapp update time max")
	cmd.Flags().Uint64(FlagExecutorsMin, 0, "dapp executors min")
	cmd.Flags().Uint64(FlagExecutorsMax, 0, "dapp executors max")
	cmd.Flags().Uint64(FlagVerifiersMin, 0, "dapp verifiers min")

	flags.AddTxFlagsToCmd(cmd)
	_ = cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

func parseRolesAccounts(rolesStr, accountsStr string) ([]uint64, []string, error) {
	roleStrArr := strings.Split(rolesStr, ",")
	if rolesStr == "" {
		roleStrArr = []string{}
	}
	roles := []uint64{}
	for _, roleStr := range roleStrArr {
		role, err := strconv.Atoi(roleStr)
		if err != nil {
			return []uint64{}, []string{}, fmt.Errorf("invalid role: %w", err)
		}
		roles = append(roles, uint64(role))
	}
	accounts := strings.Split(accountsStr, ",")
	if accountsStr == "" {
		accounts = []string{}
	}
	return roles, accounts, nil
}

// TODO: for step2
//   rpc RedeemDappPoolTx(MsgRedeemDappPoolTx) returns (MsgRedeemDappPoolTxResponse);
//   rpc SwapDappPoolTx(MsgSwapDappPoolTx) returns (MsgSwapDappPoolTxResponse);
//   rpc ConvertDappPoolTx(MsgConvertDappPoolTx) returns (MsgConvertDappPoolTxResponse);
//   rpc PauseDappTx(MsgPauseDappTx) returns (MsgPauseDappTxResponse);
//   rpc UnPauseDappTx(MsgUnPauseDappTx) returns (MsgUnPauseDappTxResponse);
//   rpc ReactivateDappTx(MsgReactivateDappTx) returns (MsgReactivateDappTxResponse);
//   rpc TransferDappTx(MsgTransferDappTx) returns (MsgTransferDappTxResponse);
//   rpc MintCreateFtTx(MsgMintCreateFtTx) returns (MsgMintCreateFtTxResponse);
//   rpc MintCreateNftTx(MsgMintCreateNftTx) returns (MsgMintCreateNftTxResponse);
//   rpc MintIssueTx(MsgMintIssueTx) returns (MsgMintIssueTxResponse);
//   rpc MintBurnTx(MsgMintBurnTx) returns (MsgMintBurnTxResponse);
