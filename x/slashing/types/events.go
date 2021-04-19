//noalias
package types

// Slashing module event types
const (
	EventTypeInactivate = "inactivate"
	EventTypeLiveness   = "liveness"

	AttributeKeyAddress          = "address"
	AttributeKeyHeight           = "height"
	AttributeKeyPower            = "power"
	AttributeKeyReason           = "reason"
	AttributeKeyInactivated      = "inactivated"
	AttributeKeyMischance        = "mischance"
	AttributeKeyLastPresentBlock = "last_present_block"
	AttributeKeyMissedBlocks     = "missed_blocks"
	AttributeKeyProducedBlocks   = "produced_blocks"

	AttributeValueDoubleSign       = "double_sign"
	AttributeValueMissingSignature = "missing_signature"
	AttributeValueCategory         = ModuleName
)
