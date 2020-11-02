package types

// ModuleName describes name of module
const ModuleName = "feeprocessing"

// RouterKey to be used for routing msgs
const RouterKey = ModuleName

// constants
var (
	KeyFeePaymentHistory = []byte("fee_payment_history")
	KeyExecutionStatus   = []byte("execution_status")
)
