package types

type Content interface {
	ProposalType() string
	VotePermission() PermValue
	ProposePermission() PermValue
}
