//noalias
package types

// Slashing module event types
const (
	EventTypeInactivate = "inactivate"
	EventTypeLiveness   = "liveness"

	AttributeKeyAddress      = "address"
	AttributeKeyHeight       = "height"
	AttributeKeyPower        = "power"
	AttributeKeyReason       = "reason"
	AttributeKeyInactivated  = "inactivated"
	AttributeKeyMissedBlocks = "missed_blocks"

	AttributeValueDoubleSign       = "double_sign"
	AttributeValueMissingSignature = "missing_signature"
	AttributeValueCategory         = ModuleName
)
