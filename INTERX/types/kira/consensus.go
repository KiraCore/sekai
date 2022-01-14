package kira

import (
	"time"

	tmConsTypes "github.com/tendermint/tendermint/consensus/types"
)

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

type LastCommit struct {
	Votes         []string `json:"votes"`
	VotesBitArray string   `json:"votes_bit_array"`
}

type RoundState struct {
	Height    string                    `json:"height"`
	Round     int64                     `json:"round"`
	Step      tmConsTypes.RoundStepType `json:"step"`
	StartTime time.Time                 `json:"start_time"`

	CommitTime time.Time       `json:"commit_time"`
	Validators ValidatorSet    `json:"validators"`
	Votes      []HeightVoteSet `json:"votes"`
	LastCommit LastCommit      `json:"last_commit"`

	TriggeredTimeoutPrecommit bool `json:"triggered_timeout_precommit"`
}

type ResultDumpConsensusState struct {
	RoundState RoundState `json:"round_state"`
}

// Interx Response
type ConsensusResponse struct {
	Height    string    `json:"height"` // Height we are working on
	Round     int64     `json:"round"`
	Step      string    `json:"step"`
	StartTime time.Time `json:"start_time"`

	CommitTime      time.Time  `json:"commit_time"`
	LastCommit      LastCommit `json:"last_commit"`
	ConsensusHealth string     `json:"consensus_health"`

	TriggeredTimeoutPrecommit bool `json:"triggered_timeout_precommit"`

	Proposer         string   `json:"proposer"`
	Precommits       []string `json:"precommits"`
	Prevotes         []string `json:"prevotes"`
	Noncommits       []string `json:"noncommits"`
	AverageBlockTime float64  `json:"average_block_time"`
	ConsensusStopped bool     `json:"consensus_stopped"`
}
