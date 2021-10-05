package types

type Content interface {
	// ProposalType is a string that should be unique by proposal content type in order to be able
	// to parse and apply once it passes.
	ProposalType() string

	// VotePermission returns the PermValue a user needs to have in order to be able to vote the proposal.
	VotePermission() PermValue

	// ProposalPermission returns PermValue a user needs to have in order to be able to submit the proposal.
	ProposalPermission() PermValue

	// ValidateBasic returns basic validation result for the proposal
	ValidateBasic() error

	// ProposalRoute() string
	// String() string
}
