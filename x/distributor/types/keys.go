package types

// constants
var (
	ModuleName = "distributor"
	// RouterKey to be used for routing msgs
	RouterKey    = ModuleName
	QuerierRoute = ModuleName

	FeesTreasuryKey        = []byte("fees_treasury")
	SnapPeriodKey          = []byte("snap_period")
	ProposerKey            = []byte("proposer_key")
	PrefixKeyValidatorVote = []byte("validator_vote_prefix")
	KeyYearStartSnapshot   = []byte("year_start_snapshot")
	KeyPeriodicSnapshot    = []byte("periodic_snapshot")
)
