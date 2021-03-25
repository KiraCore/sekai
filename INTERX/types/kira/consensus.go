package kira

type Validator struct {
	Address string `json:"address"`
	PubKey  struct {
		Type  string `json:"type"`
		Value string `json:"value"`
	} `json:"pub_key"`
	VotingPower string `json:"voting_power"`

	ProposerPriority string `json:"proposer_priority"`
}

type ValidatorSet struct {
	// NOTE: persisted via reflect, must be exported.
	Validators []Validator `json:"validators"`
	Proposer   Validator   `json:"proposer"`
}

type HeightVoteSet struct {
	Round              int64    `json:"round"`
	Prevotes           []string `json:"prevotes"`
	PrevotesBitArray   string   `json:"prevotes_bit_array"`
	Precommits         []string `json:"precommits"`
	PrecommitsBitArray string   `json:"precommits_bit_array"`
}

type RoundState struct {
	Validators ValidatorSet    `json:"validators"`
	Votes      []HeightVoteSet `json:"votes"`
}

type ResultDumpConsensusState struct {
	RoundState RoundState `json:"round_state"`
}
