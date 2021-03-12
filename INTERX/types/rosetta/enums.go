package rosetta

type CoinAction string

const (
	CoinCreated CoinAction = "coin_created"
	CoinSpent   CoinAction = "coin_spent"
)

type Direction string

const (
	Forward  Direction = "forward"
	Backward Direction = "backward"
)

type SignatureType string

const (
	Ecdsa           SignatureType = "ecdsa"
	EcdsaRecovery   SignatureType = "ecdsa_recovery"
	Ed25519         SignatureType = "ed25519"
	Schnorr1        SignatureType = "schnorr_1"
	SchnorrPoseidon SignatureType = "schnorr_poseidon"
)

type CurveType string

const (
	Secp256k1    CurveType = "secp256k1"
	Secp256r1    CurveType = "secp256r1"
	Edwards25519 CurveType = "edwards25519"
	Tweedle      CurveType = "tweedle"
)

type ExemptionType string

const (
	GreaterOrEqual ExemptionType = "greater_or_equal"
	LessOrEqual    ExemptionType = "less_or_equal"
	Dynamic        ExemptionType = "dynamic"
)

type BlockEventType string

const (
	BlockAdded   BlockEventType = "block_added"
	BlockRemoved BlockEventType = "block_removed"
)

type Operator string

const (
	OperatorOR  Operator = "or"
	OperatorAnd Operator = "and"
)
