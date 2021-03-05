package gov

// VoteOption enumerates the valid vote options for a given governance proposal.
type VoteOption int32

const (
	// VOTE_OPTION_UNSPECIFIED defines a no-op vote option.
	VoteOption_VOTE_OPTION_UNSPECIFIED VoteOption = 0
	// VOTE_OPTION_YES defines a yes vote option.
	VoteOption_VOTE_OPTION_YES VoteOption = 1
	// VOTE_OPTION_ABSTAIN defines an abstain vote option.
	VoteOption_VOTE_OPTION_ABSTAIN VoteOption = 2
	// VOTE_OPTION_NO defines a no vote option.
	VoteOption_VOTE_OPTION_NO VoteOption = 3
	// VOTE_OPTION_NO_WITH_VETO defines a no with veto vote option.
	VoteOption_VOTE_OPTION_NO_WITH_VETO VoteOption = 4
)

// Enum value maps for VoteOption.
var (
	VoteOption_name = map[int32]string{
		0: "VOTE_OPTION_UNSPECIFIED",
		1: "VOTE_OPTION_YES",
		2: "VOTE_OPTION_ABSTAIN",
		3: "VOTE_OPTION_NO",
		4: "VOTE_OPTION_NO_WITH_VETO",
	}
	VoteOption_value = map[string]int32{
		"VOTE_OPTION_UNSPECIFIED":  0,
		"VOTE_OPTION_YES":          1,
		"VOTE_OPTION_ABSTAIN":      2,
		"VOTE_OPTION_NO":           3,
		"VOTE_OPTION_NO_WITH_VETO": 4,
	}
)

type Vote struct {
	ProposalID uint64 `json:"proposal_id"`
	Voter      string `json:"voter"`
	Option     string `json:"option"`
}
