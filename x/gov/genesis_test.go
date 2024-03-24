package gov_test

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"

	simapp "github.com/KiraCore/sekai/app"
	kiratypes "github.com/KiraCore/sekai/types"
	"github.com/KiraCore/sekai/x/gov"
	"github.com/KiraCore/sekai/x/gov/keeper"
	"github.com/KiraCore/sekai/x/gov/types"
	dbm "github.com/cometbft/cometbft-db"
	"github.com/cometbft/cometbft/libs/log"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/store"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

// test for gov export / init genesis process
func TestSimappExportGenesis(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.NewContext(false, tmproto.Header{})

	genesis := gov.ExportGenesis(ctx, app.CustomGovKeeper)
	bz, err := app.AppCodec().MarshalJSON(genesis)
	require.NoError(t, err)
	buffer := new(bytes.Buffer)
	err = json.Compact(buffer, []byte(`{	
  "default_denom": "ukex",	
  "bech32_prefix": "kira",
  "starting_proposal_id": "1",
  "next_role_id": "3",
  "roles": [
    {
      "id": 1,
      "sid": "sudo",
      "description": "Sudo role"
    },
    {
      "id": 2,
      "sid": "validator",
      "description": "Validator role"
    }
  ],
  "role_permissions": {
    "1": {
      "blacklist": [],
      "whitelist": [
        1,
        2,
        3,
        6,
        8,
        9,
        12,
        13,
        10,
        11,
        14,
        15,
        18,
        19,
        20,
        21,
        22,
        23,
        31,
        32,
        24,
        25,
        16,
        17,
        4,
        5,
        26,
        27,
        28,
        29,
        30,
        33,
        34,
        35,
        36,
        37,
        38,
        39,
        40,
        41,
        42,
        43,
        44,
        45,
        46,
        47,
        48,
        49,
        50,
        51,
        52,
        53,
        54,
        55,
        56,
        57,
        58,
        59,
        60,
        61,
        62,
        63,
        64,
        65,
        66
      ]
    },
    "2": {
      "blacklist": [],
      "whitelist": [
        2
      ]
    }
  },
  "network_actors": [],
  "network_properties": {
    "min_tx_fee": "100",
    "max_tx_fee": "1000000",
    "vote_quorum": "33",
    "minimum_proposal_end_time": "300",
    "proposal_enactment_time": "300",
    "min_proposal_end_blocks": "2",
    "min_proposal_enactment_blocks": "1",
    "enable_foreign_fee_payments": true,
    "mischance_rank_decrease_amount": "10",
    "max_mischance": "110",
    "mischance_confidence": "10",
    "inactive_rank_decrease_percent": "0.500000000000000000",
    "min_validators": "1",
    "poor_network_max_bank_send": "1000000",
    "unjail_max_time": "600",
    "enable_token_whitelist": false,
    "enable_token_blacklist": true,
    "min_identity_approval_tip": "200",
    "unique_identity_keys": "moniker,username",
    "ubi_hardcap": "6000000",
    "validators_fee_share": "0.500000000000000000",
    "inflation_rate": "0.180000000000000000",
    "inflation_period": "31557600",
    "unstaking_period": "2629800",
    "max_delegators": "100",
    "min_delegation_pushout": "10",
    "slashing_period": "3600",
    "max_jailed_percentage": "0.250000000000000000",
    "max_slashing_percentage": "0.010000000000000000",
    "min_custody_reward": "200",
    "max_custody_buffer_size": "10",
    "max_custody_tx_size": "8192",
    "abstention_rank_decrease_amount": "1",
    "max_abstention": "2",
    "min_collective_bond": "100000",
    "min_collective_bonding_time": "86400",
    "max_collective_outputs": "10",
    "min_collective_claim_period": "14400",
    "validator_recovery_bond": "300000",
    "max_annual_inflation": "0.350000000000000000",
    "max_proposal_title_size": "128",
    "max_proposal_description_size": "1024",
    "max_proposal_poll_option_size": "64",
    "max_proposal_poll_option_count": "128",
    "max_proposal_reference_size": "512",
    "max_proposal_checksum_size": "128",
  	"min_dapp_bond": "1000000",	
    "max_dapp_bond": "10000000",	
    "dapp_liquidation_threshold": "0",
    "dapp_liquidation_period": "0",
    "dapp_bond_duration": "604800",
    "dapp_verifier_bond": "0.001000000000000000",
    "dapp_auto_denounce_time": "60",
    "dapp_mischance_rank_decrease_amount": "1",
    "dapp_max_mischance": "10",
    "dapp_inactive_rank_decrease_percent": "10",
    "dapp_pool_slippage_default": "0.100000000000000000",
    "minting_ft_fee": "100000000000000",
    "minting_nft_fee": "100000000000000",
    "veto_threshold": "33.400000000000000000",
    "autocompound_interval_num_blocks": "17280",
	"bridge_address": "test",
    "bridge_cosmos_ethereum_exchange_rate": 10,
    "bridge_ethereum_cosmos_exchange_rate": 0.100000000000000000,
  },
  "execution_fees": [
    {
      "transaction_type": "activate",
      "execution_fee": "100",
      "failure_fee": "1000",
      "timeout": "10",
      "default_parameters": "0"
    },
    {
      "transaction_type": "claim-councilor",
      "execution_fee": "100",
      "failure_fee": "1",
      "timeout": "10",
      "default_parameters": "0"
    },
    {
      "transaction_type": "claim-proposal-type-x",
      "execution_fee": "100",
      "failure_fee": "1",
      "timeout": "10",
      "default_parameters": "0"
    },
    {
      "transaction_type": "claim-validator",
      "execution_fee": "100",
      "failure_fee": "1",
      "timeout": "10",
      "default_parameters": "0"
    },
    {
      "transaction_type": "pause",
      "execution_fee": "100",
      "failure_fee": "100",
      "timeout": "10",
      "default_parameters": "0"
    },
    {
      "transaction_type": "submit-proposal-type-x",
      "execution_fee": "10",
      "failure_fee": "1",
      "timeout": "10",
      "default_parameters": "0"
    },
    {
      "transaction_type": "unpause",
      "execution_fee": "100",
      "failure_fee": "100",
      "timeout": "10",
      "default_parameters": "0"
    },
    {
      "transaction_type": "upsert-token-alias",
      "execution_fee": "100",
      "failure_fee": "1",
      "timeout": "10",
      "default_parameters": "0"
    },
    {
      "transaction_type": "veto-proposal-type-x",
      "execution_fee": "100",
      "failure_fee": "1",
      "timeout": "10",
      "default_parameters": "0"
    },
    {
      "transaction_type": "vote-proposal-type-x",
      "execution_fee": "100",
      "failure_fee": "1",
      "timeout": "10",
      "default_parameters": "0"
    }
  ],
  "poor_network_messages": {
    "messages": [
      "submit-proposal",
      "set-network-properties",
      "vote-proposal",
      "claim-councilor",
      "whitelist-permissions",
      "blacklist-permissions",
      "create-role",
      "assign-role",
      "unassign-role",
      "whitelist-role-permission",
      "blacklist-role-permission",
      "remove-whitelist-role-permission",
      "remove-blacklist-role-permission",
      "claim-validator",
      "activate",
      "pause",
      "unpause",
      "register-identity-records",
      "edit-identity-record",
      "request-identity-records-verify",
      "handle-identity-records-verify-request",
      "cancel-identity-records-verify-request"
    ]
  },
  "proposals": [],
  "votes": [],
  "data_registry": {},
  "identity_records": [],
  "last_identity_record_id": "0",
  "id_records_verify_requests": [],
  "last_id_record_verify_request_id": "0",
  "proposal_durations": {}
}`))
	require.NoError(t, err)
	require.Equal(t, string(bz), buffer.String())
}

func TestExportInitGenesis(t *testing.T) {
	// Initialize store keys
	keyGovernance := sdk.NewKVStoreKey(types.ModuleName)

	// Initialize memory database and mount stores on it
	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(keyGovernance, storetypes.StoreTypeIAVL, nil)
	ms.LoadLatestVersion()

	ctx := sdk.NewContext(ms, tmproto.Header{
		Height: 1,
		Time:   time.Date(2020, time.March, 1, 1, 0, 0, 0, time.UTC),
	}, false, log.TestingLogger())

	k := keeper.NewKeeper(keyGovernance, simapp.MakeEncodingConfig().Marshaler, nil)

	genState := types.GenesisState{
		DefaultDenom: "ukex",
		Bech32Prefix: "kira",
		NextRoleId:   3,
		Roles: []types.Role{
			{
				Id:          uint32(types.RoleSudo),
				Sid:         "sudo",
				Description: "Sudo role",
			},
			{
				Id:          uint32(types.RoleValidator),
				Sid:         "validator",
				Description: "Validator role",
			},
		},
		RolePermissions: map[uint64]*types.Permissions{
			uint64(types.RoleSudo): types.NewPermissions([]types.PermValue{
				types.PermSetPermissions,
				types.PermClaimValidator,
				types.PermClaimCouncilor,
				types.PermUpsertTokenAlias,
			}, nil),
		},
		StartingProposalId: 1,
		NetworkProperties: &types.NetworkProperties{
			MinTxFee:                         100,
			MaxTxFee:                         1000000,
			VoteQuorum:                       33,
			MinimumProposalEndTime:           300, // 300 seconds / 5 mins
			ProposalEnactmentTime:            300, // 300 seconds / 5 mins
			MinProposalEndBlocks:             2,
			MinProposalEnactmentBlocks:       1,
			MischanceRankDecreaseAmount:      1,
			MaxMischance:                     1,
			InactiveRankDecreasePercent:      sdk.NewDecWithPrec(2, 2),
			MinValidators:                    1,
			PoorNetworkMaxBankSend:           1,
			EnableForeignFeePayments:         true,
			MinIdentityApprovalTip:           200,
			UniqueIdentityKeys:               "moniker,username",
			UbiHardcap:                       6000_000,
			ValidatorsFeeShare:               sdk.NewDecWithPrec(50, 2), // 50%
			InflationRate:                    sdk.NewDecWithPrec(18, 2), // 18%
			InflationPeriod:                  31557600,                  // 1 year
			UnstakingPeriod:                  2629800,                   // 1 month
			MaxDelegators:                    100,
			MinDelegationPushout:             10,
			SlashingPeriod:                   3600,
			MaxJailedPercentage:              sdk.NewDecWithPrec(25, 2),
			MaxSlashingPercentage:            sdk.NewDecWithPrec(1, 2),
			MinCustodyReward:                 200,
			MaxCustodyBufferSize:             10,
			MaxCustodyTxSize:                 8192,
			AbstentionRankDecreaseAmount:     1,
			MaxAbstention:                    2,
			MinCollectiveBond:                100_000, // in KEX
			MinCollectiveBondingTime:         86400,   // in seconds
			MaxCollectiveOutputs:             10,
			MinCollectiveClaimPeriod:         14400,                     // 4hrs
			ValidatorRecoveryBond:            300000,                    // 300k KEX
			MaxAnnualInflation:               sdk.NewDecWithPrec(35, 2), // 35%
			MaxProposalTitleSize:             128,
			MaxProposalDescriptionSize:       1024,
			MaxProposalPollOptionSize:        64,
			MaxProposalPollOptionCount:       128,
			MaxProposalReferenceSize:         512,
			MaxProposalChecksumSize:          128,
			MinDappBond:                      1000000,
			MaxDappBond:                      10000000,
			DappBondDuration:                 604800,
			DappVerifierBond:                 sdk.NewDecWithPrec(1, 3), //0.1%
			DappAutoDenounceTime:             60,                       // 60s
			DappMischanceRankDecreaseAmount:  1,
			DappMaxMischance:                 10,
			DappInactiveRankDecreasePercent:  10,
			DappPoolSlippageDefault:          sdk.NewDecWithPrec(1, 1), // 10%
			MintingFtFee:                     100_000_000_000_000,
			MintingNftFee:                    100_000_000_000_000,
			VetoThreshold:                    sdk.NewDecWithPrec(3340, 2), // 33.40%
			AutocompoundIntervalNumBlocks:    17280,
			BridgeAddress:                    "test",
			BridgeCosmosEthereumExchangeRate: sdk.NewDec(10),
			BridgeEthereumCosmosExchangeRate: sdk.NewDecWithPrec(1, 1),
		},
		ExecutionFees: []types.ExecutionFee{
			{
				TransactionType:   "claim-validator-seat",
				ExecutionFee:      10,
				FailureFee:        1,
				Timeout:           10,
				DefaultParameters: 0,
			},
		},
		PoorNetworkMessages: &types.AllowedMessages{
			Messages: []string{
				kiratypes.MsgTypeSetNetworkProperties,
			},
		},
	}
	err := gov.InitGenesis(ctx, k, genState)
	require.NoError(t, err)

	genesis := gov.ExportGenesis(ctx, k)
	bz, err := simapp.MakeEncodingConfig().Marshaler.MarshalJSON(genesis)
	require.NoError(t, err)
	buffer := new(bytes.Buffer)
	err = json.Compact(buffer, []byte(`	{
  "default_denom": "ukex",	
  "bech32_prefix": "kira",
  "starting_proposal_id": "1",	
  "next_role_id": "3",	
  "roles": [	
    {	
      "id": 1,	
      "sid": "sudo",	
      "description": "Sudo role"	
    },	
    {	
      "id": 2,	
      "sid": "validator",	
      "description": "Validator role"	
    }	
  ],	
  "role_permissions": {	
    "1": {	
      "blacklist": [],	
      "whitelist": [	
        1,	
        2,	
        3,	
        6	
      ]	
    },	
    "2": {	
      "blacklist": [],	
      "whitelist": []	
    }	
  },	
  "network_actors": [],	
  "network_properties": {	
    "min_tx_fee": "100",	
    "max_tx_fee": "1000000",	
    "vote_quorum": "33",	
    "minimum_proposal_end_time": "300",	
    "proposal_enactment_time": "300",	
    "min_proposal_end_blocks": "2",	
    "min_proposal_enactment_blocks": "1",	
    "enable_foreign_fee_payments": true,	
    "mischance_rank_decrease_amount": "1",	
    "max_mischance": "1",	
    "mischance_confidence": "0",	
    "inactive_rank_decrease_percent": "0.020000000000000000",	
    "min_validators": "1",	
    "poor_network_max_bank_send": "1",	
    "unjail_max_time": "0",	
    "enable_token_whitelist": false,	
    "enable_token_blacklist": false,	
    "min_identity_approval_tip": "200",	
    "unique_identity_keys": "moniker,username",	
    "ubi_hardcap": "6000000",	
    "validators_fee_share": "0.500000000000000000",	
    "inflation_rate": "0.180000000000000000",	
    "inflation_period": "31557600",	
    "unstaking_period": "2629800",	
    "max_delegators": "100",	
    "min_delegation_pushout": "10",	
    "slashing_period": "3600",	
    "max_jailed_percentage": "0.250000000000000000",	
    "max_slashing_percentage": "0.010000000000000000",	
    "min_custody_reward": "200",	
    "max_custody_buffer_size": "10",	
    "max_custody_tx_size": "8192",	
    "abstention_rank_decrease_amount": "1",	
    "max_abstention": "2",	
    "min_collective_bond": "100000",	
    "min_collective_bonding_time": "86400",	
    "max_collective_outputs": "10",	
    "min_collective_claim_period": "14400",	
    "validator_recovery_bond": "300000",	
    "max_annual_inflation": "0.350000000000000000",
    "max_proposal_title_size": "128",
    "max_proposal_description_size": "1024",
    "max_proposal_poll_option_size": "64",
    "max_proposal_poll_option_count": "128",
    "max_proposal_reference_size": "512",
    "max_proposal_checksum_size": "128",
    "min_dapp_bond": "1000000",	
    "max_dapp_bond": "10000000",	
    "dapp_liquidation_threshold": "0",
    "dapp_liquidation_period": "0",
    "dapp_bond_duration": "604800",
    "dapp_verifier_bond": "0.001000000000000000",
    "dapp_auto_denounce_time": "60",
    "dapp_mischance_rank_decrease_amount": "1",
    "dapp_max_mischance": "10",
    "dapp_inactive_rank_decrease_percent": "10",
    "dapp_pool_slippage_default": "0.100000000000000000",
    "minting_ft_fee": "100000000000000",
    "minting_nft_fee": "100000000000000",
    "veto_threshold": "33.400000000000000000",
    "autocompound_interval_num_blocks": "17280",
	"bridge_address": "test",
    "bridge_cosmos_ethereum_exchange_rate": 10,
    "bridge_ethereum_cosmos_exchange_rate": 0.100000000000000000,
  },	
  "execution_fees": [	
    {	
      "transaction_type": "claim-validator-seat",	
      "execution_fee": "10",	
      "failure_fee": "1",	
      "timeout": "10",	
      "default_parameters": "0"	
    }	
  ],	
  "poor_network_messages": {	
    "messages": [	
      "set-network-properties"	
    ]	
  },	
  "proposals": [],	
  "votes": [],	
  "data_registry": {},	
  "identity_records": [],	
  "last_identity_record_id": "0",	
  "id_records_verify_requests": [],	
  "last_id_record_verify_request_id": "0",	
  "proposal_durations": {}	
}`))
	require.NoError(t, err)
	require.Equal(t, string(bz), buffer.String())
}
