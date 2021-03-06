package cli

import (
	"fmt"

	"github.com/KiraCore/sekai/x/genutil"
	customgovtypes "github.com/KiraCore/sekai/x/gov/types"
	customstakingtypes "github.com/KiraCore/sekai/x/staking/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking/client/cli"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const (
	FlagComission   = "commission"
	FlagValKey      = "validator-key"
	FlagDescription = "description"
)

func GetTxClaimValidatorCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "claim-validator-seat",
		Short: "Claim validator seat to become a Validator",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			serverCtx := server.GetServerContextFromCmd(cmd)

			moniker, _ := cmd.Flags().GetString(FlagMoniker)
			comission, _ := cmd.Flags().GetString(FlagComission)

			var (
				valPubKey cryptotypes.PubKey
			)
			if valPubKeyString, _ := cmd.Flags().GetString(cli.FlagPubKey); valPubKeyString != "" {
				valPubKey, err = types.GetPubKeyFromBech32(types.Bech32PubKeyTypeConsPub, valPubKeyString)
				if err != nil {
					return errors.Wrap(err, "failed to get consensus node public key")
				}
			} else {
				_, valPubKey, err = genutil.InitializeNodeValidatorFiles(serverCtx.Config)
				if err != nil {
					return errors.Wrap(err, "failed to initialize node validator files")
				}
			}

			comm, err := types.NewDecFromStr(comission)
			val := types.ValAddress(clientCtx.GetFromAddress())

			msg, err := customstakingtypes.NewMsgClaimValidator(moniker, comm, val, valPubKey)
			if err != nil {
				return fmt.Errorf("error creating tx: %w", err)
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	AddValidatorFlags(cmd)
	flags.AddTxFlagsToCmd(cmd)

	_ = cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

// GetTxProposalUnjailValidatorCmd implement cli command for MsgUpsertTokenAlias
func GetTxProposalUnjailValidatorCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "proposal-unjail-validator hash reference",
		Short: "Creates an proposal to unjail validator (the from address is the validator)",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			hash := args[0]
			reference := args[1]

			description, err := cmd.Flags().GetString(FlagDescription)
			if err != nil {
				return fmt.Errorf("invalid description: %w", err)
			}

			msg, err := customgovtypes.NewMsgSubmitProposal(
				clientCtx.FromAddress,
				description,
				customstakingtypes.NewUnjailValidatorProposal(clientCtx.FromAddress, hash, reference),
			)
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(FlagDescription, "", "The description of the proposal, it can be a url, some text, etc.")
	cmd.MarkFlagRequired(FlagDescription)

	flags.AddTxFlagsToCmd(cmd)
	_ = cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}
