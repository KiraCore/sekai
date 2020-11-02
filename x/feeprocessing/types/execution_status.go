package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ExecutionStatus describes msg execution status
type ExecutionStatus struct {
	Msg     sdk.Msg
	Success bool
}
