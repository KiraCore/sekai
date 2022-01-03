//noalias
package types

// Slashing module event types
const (
	EventTypeInactivate = "inactivate"
	EventTypeJail       = "jail"
	EventTypeLiveness   = "liveness"

	AttributeKeyAddress          = "address"
	AttributeKeyHeight           = "height"
	AttributeKeyPower            = "power"
	AttributeKeyReason           = "reason"
	AttributeKeyInactivated      = "inactivated"
	AttributeKeyJailed           = "jailed"
	AttributeKeyMischance        = "mischance"
	AttributeKeyLastPresentBlock = "last_present_block"
	AttributeKeyMissedBlocks     = "missed_blocks"
	AttributeKeyProducedBlocks   = "produced_blocks"
	AttributeKeyDescription      = "description"

	AttributeValueDoubleSign       = "double_sign"
	AttributeValueMissingSignature = "missing_signature"
	AttributeValueCategory         = ModuleName
)
