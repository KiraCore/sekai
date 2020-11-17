package types

type Content interface {
	// ProposalType is a string that should be unique by proposal content type in order to be able
	// to parse and apply once it passes.
	ProposalType() string

	// VotePermission returns the PermValue a user needs to have in order to be able to vote the proposal.
	VotePermission() PermValue
}
