package types

// constants
var (
	ModuleName = "distributor"
	// RouterKey to be used for routing msgs
	RouterKey    = ModuleName
	QuerierRoute = ModuleName

	FeesCollectedKey       = []byte("fees_collected")
	FeesTreasuryKey        = []byte("fees_treasury")
	SnapPeriodKey          = []byte("snap_period")
	ProposerKey            = []byte("proposer_key")
	PrefixKeyValidatorVote = []byte("validator_vote_prefix")
)
