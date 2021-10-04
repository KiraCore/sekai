package cli

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"

	"github.com/KiraCore/sekai/x/gov/types"
)

// Proposal flags
const (
	flagVoter = "voter"
)

// GetCmdQueryPermissions the query delegation command.
func GetCmdQueryPermissions() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "permissions [addr]",
		Short: "Query permissions of an address",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			accAddr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return errors.Wrap(err, "invalid account address")
			}

			params := &types.PermissionsByAddressRequest{ValAddr: accAddr}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.PermissionsByAddress(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res.Permissions)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryRolesByAddress the query delegation command.
func GetCmdQueryRolesByAddress() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "roles addr",
		Short: "Query roles assigned to an address",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			accAddr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return errors.Wrap(err, "invalid account address")
			}

			params := &types.RolesByAddressRequest{ValAddr: accAddr}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.RolesByAddress(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func GetCmdQueryRolePermissions() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "role-permissions arg-num",
		Short: "Query permissions of all the roles",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			roleNum, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid role number")
			}

			params := &types.RolePermissionsRequest{
				Role: roleNum,
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.RolePermissions(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res.Permissions)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryNetworkProperties implement query network properties
func GetCmdQueryNetworkProperties() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "network-properties",
		Short: "Query network properties",
		Args:  cobra.MinimumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			params := &types.NetworkPropertiesRequest{}
			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.NetworkProperties(context.Background(), params)
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryPoorNetworkMessages query for poor network messages
func GetCmdQueryPoorNetworkMessages() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "poor-network-messages",
		Short: "Query poor network messages",
		Args:  cobra.MinimumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			params := &types.PoorNetworkMessagesRequest{}
			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.PoorNetworkMessages(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdQueryExecutionFee query for execution fee by execution name
func GetCmdQueryExecutionFee() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "execution-fee [transaction_type]",
		Short: "Query execution fee by the type of transaction",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			params := &types.ExecutionFeeRequest{
				TransactionType: args[0],
			}
			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.ExecutionFee(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func GetCmdQueryCouncilRegistry() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "council-registry [--addr || --flagMoniker]",
		Short: "Query governance registry.",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			addr, err := cmd.Flags().GetString(FlagAddress)
			if err != nil {
				return err
			}

			moniker, err := cmd.Flags().GetString(FlagMoniker)
			if err != nil {
				return err
			}
			if addr == "" && moniker == "" {
				return fmt.Errorf("at least one flag (--flag or --moniker) is mandatory")
			}

			var res *types.CouncilorResponse
			if moniker != "" {
				params := &types.CouncilorByMonikerRequest{Moniker: moniker}

				queryClient := types.NewQueryClient(clientCtx)
				res, err = queryClient.CouncilorByMoniker(context.Background(), params)
				if err != nil {
					return err
				}
			} else {
				bech32, err := sdk.AccAddressFromBech32(addr)
				if err != nil {
					return fmt.Errorf("invalid address: %w", err)
				}

				params := &types.CouncilorByAddressRequest{ValAddr: bech32}

				queryClient := types.NewQueryClient(clientCtx)
				res, err = queryClient.CouncilorByAddress(context.Background(), params)
				if err != nil {
					return err
				}
			}

			return clientCtx.PrintProto(&res.Councilor)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	cmd.Flags().String(FlagAddress, "", "the address you want to query information")
	cmd.Flags().String(FlagMoniker, "", "the moniker you want to query information")

	return cmd
}

// GetCmdQueryProposals implements a query proposals command. Command to Get a
// Proposal Information.
func GetCmdQueryProposals() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "proposals",
		Short: "Query proposals with optional filters",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query for a all paginated proposals that match optional filters:

Example:
$ %s query gov proposals --depositor cosmos1skjwj5whet0lpe65qaq4rpq03hjxlwd9nf39lk
$ %s query gov proposals --voter cosmos1skjwj5whet0lpe65qaq4rpq03hjxlwd9nf39lk
`,
				version.AppName, version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			bechVoterAddr, _ := cmd.Flags().GetString(flagVoter)

			if len(bechVoterAddr) != 0 {
				_, err := sdk.AccAddressFromBech32(bechVoterAddr)
				if err != nil {
					return err
				}
			}

			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			res, err := queryClient.Proposals(
				context.Background(),
				&types.QueryProposalsRequest{
					Voter:      bechVoterAddr,
					Pagination: pageReq,
				},
			)
			if err != nil {
				return err
			}

			if len(res.GetProposals()) == 0 {
				return fmt.Errorf("no proposals found")
			}

			return clientCtx.PrintProto(res)
		},
	}

	cmd.Flags().String(flagVoter, "", "(optional) filter by proposals voted on by voted")
	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "customgov")

	return cmd
}

// GetCmdQueryProposal implements the query proposal command.
func GetCmdQueryProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "proposal [proposal-id]",
		Args:  cobra.ExactArgs(1),
		Short: "Query proposal details",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details for a proposal. You can find the
proposal-id by running "%s query gov proposals".

Example:
$ %s query gov proposal 1
`,
				version.AppName, version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			// validate that the proposal id is a uint
			proposalID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("proposal-id %s not a valid uint, please input a valid proposal-id", args[0])
			}

			// Query the proposal
			res, err := queryClient.Proposal(
				context.Background(),
				&types.QueryProposalRequest{ProposalId: proposalID},
			)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(&res.Proposal)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryVote implements the query proposal vote command. Command to Get a
// Proposal Information.
func GetCmdQueryVote() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vote [proposal-id] [voter-addr]",
		Args:  cobra.ExactArgs(2),
		Short: "Query details of a single vote",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details for a single vote on a proposal given its identifier.

Example:
$ %s query gov vote 1 kira1skjwj5whet0lpe65qaq4rpq03hjxlwd9nf39lk
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			// validate that the proposal id is a uint
			proposalID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("proposal-id %s not a valid int, please input a valid proposal-id", args[0])
			}

			// check to see if the proposal is in the store
			_, err = queryClient.Proposal(
				context.Background(),
				&types.QueryProposalRequest{ProposalId: proposalID},
			)
			if err != nil {
				return fmt.Errorf("failed to fetch proposal-id %d: %s", proposalID, err)
			}

			voterAddr, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			res, err := queryClient.Vote(
				context.Background(),
				&types.QueryVoteRequest{ProposalId: proposalID, Voter: voterAddr},
			)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(&res.Vote)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryVotes implements the command to query for proposal votes.
func GetCmdQueryVotes() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "votes [proposal-id]",
		Args:  cobra.ExactArgs(1),
		Short: "Query votes on a proposal",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query vote details for a single proposal by its identifier.

Example:
$ %[1]s query gov votes 1
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			// validate that the proposal id is a uint
			proposalID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("proposal-id %s not a valid int, please input a valid proposal-id", args[0])
			}

			res, err := queryClient.Votes(
				context.Background(),
				&types.QueryVotesRequest{ProposalId: proposalID},
			)

			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)

		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryWhitelistedProposalVoters implements the command to query for possible proposal voters.
func GetCmdQueryWhitelistedProposalVoters() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "voters [proposal-id]",
		Args:  cobra.ExactArgs(1),
		Short: "Query voters of a proposal",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query voters for a single proposal by its identifier.

Example:
$ %[1]s query gov voters 1
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			// validate that the proposal id is a uint
			proposalID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("proposal-id %s not a valid int, please input a valid proposal-id", args[0])
			}

			res, err := queryClient.WhitelistedProposalVoters(
				context.Background(),
				&types.QueryWhitelistedProposalVotersRequest{ProposalId: proposalID},
			)

			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)

		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryIdentityRecord implements the command to query identity record by id
func GetCmdQueryIdentityRecord() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "identity-record [id]",
		Args:  cobra.ExactArgs(1),
		Short: "Query identity record by id",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query identity record by id.

Example:
$ %[1]s query gov identity-record 1
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			// validate that the id is a uint
			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("id %s not a valid int, please input a valid id", args[0])
			}

			res, err := queryClient.IdentityRecord(
				context.Background(),
				&types.QueryIdentityRecordRequest{Id: id},
			)

			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryIdentityRecordByAddress implements the command to query identity records by records creator
func GetCmdQueryIdentityRecordByAddress() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "identity-records-by-addr [addr]",
		Args:  cobra.ExactArgs(1),
		Short: "Query identity records by owner",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query identity records by owner.

Example:
$ %[1]s query gov identity-records-by-addr [addr]
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			// validate address
			addr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			keysStr, err := cmd.Flags().GetString(FlagKeys)
			if err != nil {
				return err
			}

			keys := strings.Split(keysStr, ",")
			if keysStr == "" {
				keys = []string{}
			}

			res, err := queryClient.IdentityRecordsByAddress(
				context.Background(),
				&types.QueryIdentityRecordsByAddressRequest{
					Creator: addr,
					Keys:    keys,
				},
			)

			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	cmd.Flags().String(FlagKeys, "", "keys required when needs to be filtered")

	return cmd
}

// GetCmdQueryAllIdentityRecords implements the command to query all identity records
func GetCmdQueryAllIdentityRecords() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "all-identity-records",
		Args:  cobra.ExactArgs(0),
		Short: "Query all identity records",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query all identity records.

Example:
$ %[1]s query gov all-identity-records
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			res, err := queryClient.AllIdentityRecords(
				context.Background(),
				&types.QueryAllIdentityRecordsRequest{
					Pagination: pageReq,
				},
			)

			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "customgov")

	return cmd
}

// GetCmdQueryIdentityRecordVerifyRequest implements the command to query identity record verify request by id
func GetCmdQueryIdentityRecordVerifyRequest() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "identity-record-verify-request [id]",
		Args:  cobra.ExactArgs(1),
		Short: "Query identity record verify request by id",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query identity record verify request by id.

Example:
$ %[1]s query gov identity-record-verify-request 1
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			// validate that the id is a uint
			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("id %s not a valid int, please input a valid id", args[0])
			}

			res, err := queryClient.IdentityRecordVerifyRequest(
				context.Background(),
				&types.QueryIdentityVerifyRecordRequest{
					RequestId: id,
				},
			)

			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryIdentityRecordVerifyRequestsByRequester implements the command to query identity records verify requests by requester
func GetCmdQueryIdentityRecordVerifyRequestsByRequester() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "identity-record-verify-requests-by-requester [addr]",
		Args:  cobra.ExactArgs(1),
		Short: "Query identity records verify requests by requester",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query identity records verify requests by requester.

Example:
$ %[1]s query gov identity-record-verify-requests-by-requester [addr]
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			// validate address
			addr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			res, err := queryClient.IdentityRecordVerifyRequestsByRequester(
				context.Background(),
				&types.QueryIdentityRecordVerifyRequestsByRequester{
					Requester:  addr,
					Pagination: pageReq,
				},
			)

			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "customgov")

	return cmd
}

// GetCmdQueryIdentityRecordVerifyRequestsByApprover implements the command to query identity records verify requests by approver
func GetCmdQueryIdentityRecordVerifyRequestsByApprover() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "identity-record-verify-requests-by-approver [addr]",
		Args:  cobra.ExactArgs(1),
		Short: "Query identity record verify request by approver",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query identity record verify requests by approver.

Example:
$ %[1]s query gov identity-record-verify-requests-by-approver [addr]
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			// validate address
			addr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			res, err := queryClient.IdentityRecordVerifyRequestsByApprover(
				context.Background(),
				&types.QueryIdentityRecordVerifyRequestsByApprover{
					Approver:   addr,
					Pagination: pageReq,
				},
			)

			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "customgov")

	return cmd
}

// GetCmdQueryAllIdentityRecordVerifyRequests implements the command to query all identity records verify requests
func GetCmdQueryAllIdentityRecordVerifyRequests() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "all-identity-record-verify-requests",
		Args:  cobra.ExactArgs(0),
		Short: "Query all identity records verify requests",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query all identity records verify requests.

Example:
$ %[1]s query gov all-identity-record-verify-requests
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			res, err := queryClient.AllIdentityRecordVerifyRequests(
				context.Background(),
				&types.QueryAllIdentityRecordVerifyRequests{
					Pagination: pageReq,
				},
			)

			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "customgov")

	return cmd
}

// GetCmdQueryAllDataReferenceKeys implements the command to query data reference keys
func GetCmdQueryAllDataReferenceKeys() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "all-data-reference-keys",
		Args:  cobra.ExactArgs(0),
		Short: "Query all data reference keys",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query all data reference keys.

Example:
$ %[1]s query gov all-data-reference-keys
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			res, err := queryClient.AllDataReferenceKeys(
				context.Background(),
				&types.QueryDataReferenceKeysRequest{
					Pagination: pageReq,
				},
			)

			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "customgov")

	return cmd
}
