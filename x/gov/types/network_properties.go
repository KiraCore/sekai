package types

type PropertyInfo struct {
	Name        string `json:"name"`
	Format      string `json:"format"`
	Description string `json:"description"`
}

var PropertyMetadata = []PropertyInfo{
	{
		Name:        "MinTxFee",
		Format:      "uint64",
		Description: "Minimum transaction fee on the network",
	},
	{
		Name:        "MaxTxFee",
		Format:      "uint64",
		Description: "Maximum transaction fee on the network",
	},
	{
		Name:        "VoteQuorum",
		Format:      "uint64",
		Description: "Vote quorum to meet for a proposal to pass.",
	},
	{
		Name:        "MinimumProposalEndTime",
		Format:      "uint64",
		Description: "Minimum proposal voting duration for a proposal should be live for voting.",
	},
	{
		Name:        "ProposalEnactmentTime",
		Format:      "uint64",
		Description: "Proposal enactment time.",
	},
	{
		Name:        "MinProposalEndBlocks",
		Format:      "uint64",
		Description: "Minimum number of blocks a proposal should be live for voting.",
	},
	{
		Name:        "MinProposalEnactmentBlocks",
		Format:      "uint64",
		Description: "Minimum number of blocks a proposal should be in enactment.",
	},
	{
		Name:        "EnableForeignFeePayments",
		Format:      "bool",
		Description: "Flag that describes foreign token (non-KEX) is enabled as fees.",
	},
	{
		Name:        "MischanceRankDecreaseAmount",
		Format:      "uint64",
		Description: "Rank decrease amount when a validator miss a block.",
	},
	{
		Name:        "MaxMischance",
		Format:      "uint64",
		Description: "Maximun number of sequencial miss on blocks before penalties.",
	},
	{
		Name:        "MischanceConfidence",
		Format:      "uint64",
		Description: "Number of missed blocks accepted before increasing misschance.",
	},
	{
		Name:        "InactiveRankDecreasePercent",
		Format:      "decimal",
		Description: "Percentage of rank decrease when a validator node become inactive.",
	},
	{
		Name:        "MinValidators",
		Format:      "uint64",
		Description: "Number of active validators to be an active network.",
	},
	{
		Name:        "PoorNetworkMaxBankSend",
		Format:      "uint64",
		Description: "Maximum number of tokens transferrable on poor network.",
	},
	{
		Name:        "UnjailMaxTime",
		Format:      "uint64",
		Description: "Maximum time a validator can unjail after jail.",
	},
	{
		Name:        "EnableTokenWhitelist",
		Format:      "bool",
		Description: "Flag to let only whitelisted tokens are transferrable",
	},
	{
		Name:        "EnableTokenBlacklist",
		Format:      "bool",
		Description: "Flag to prevent transfer of blacklisted tokens",
	},
	{
		Name:        "MinIdentityApprovalTip",
		Format:      "uint64",
		Description: "Minimum amount of tokens to be given for an identity record approval",
	},
	{
		Name:        "UniqueIdentityKeys",
		Format:      "string",
		Description: "Comma separated list of identity keys that should be unique across all identity keys",
	},
	{
		Name:        "UbiHardcap",
		Format:      "uint64",
		Description: "The maximum amount of tokens that can be allocated for sum of ubi records",
	},
	{
		Name:        "ValidatorsFeeShare",
		Format:      "decimal",
		Description: "The portion of fees to be given to the validator from the block fees",
	},
	{
		Name:        "InflationRate",
		Format:      "decimal",
		Description: "The rate of inflation during the inflation period",
	},
	{
		Name:        "InflationPeriod",
		Format:      "uint64",
		Description: "The duration unit for InflationRate",
	},
	{
		Name:        "UnstakingPeriod",
		Format:      "uint64",
		Description: "The unstaking duration on multistaking module",
	},
	{
		Name:        "MaxDelegators",
		Format:      "uint64",
		Description: "The maximum number of pool delegators for a single pool",
	},
	{
		Name:        "MinDelegationPushout",
		Format:      "uint64",
		Description: "The multiplier to push out a min delegation user when the maximum number of delegators filled in",
	},
	{
		Name:        "SlashingPeriod",
		Format:      "uint64",
		Description: "The period to take colluders on slash proposal",
	},
	{
		Name:        "MaxJailedPercentage",
		Format:      "decimal",
		Description: "The percentage of jails acceptable before slash proposal happens",
	},
	{
		Name:        "MaxSlashingPercentage",
		Format:      "decimal",
		Description: "The maximum slash percentage to for jail",
	},
	{
		Name:        "MinCustodyReward",
		Format:      "uint64",
		Description: "The minimum custody reward",
	},
	{
		Name:        "MaxCustodyBufferSize",
		Format:      "uint64",
		Description: "The minimum custody buffer size",
	},
	{
		Name:        "MaxCustodyTxSize",
		Format:      "uint64",
		Description: "The minimum custody transaction size",
	},
	{
		Name:        "AbstentionRankDecreaseAmount",
		Format:      "uint64",
		Description: "Rank decrease amount when a councilor does not participate in voting",
	},
	{
		Name:        "MaxAbstention",
		Format:      "uint64",
		Description: "The maximum absention count on voting for an active councilor",
	},
	{
		Name:        "MinCollectiveBond",
		Format:      "uint64",
		Description: "The minimum size of collective to be bootstrapped within bonding period",
	},
	{
		Name:        "MinCollectiveBondingTime",
		Format:      "uint64",
		Description: "The time to bootstrap minimum collectives bonds",
	},
	{
		Name:        "MaxCollectiveOutputs",
		Format:      "uint64",
		Description: "The maximum number of outputs a bonding pool could have",
	},
	{
		Name:        "MinCollectiveClaimPeriod",
		Format:      "uint64",
		Description: "The minimum acceptable collective claim period",
	},
	{
		Name:        "ValidatorRecoveryBond",
		Format:      "uint64",
		Description: "The amount of KEX to spend for issuing validator recovery token",
	},
	{
		Name:        "MaxAnnualInflation",
		Format:      "decimal",
		Description: "The maximum inflation ratio of kex by which supply can increase over the period of 1 year",
	},
}
