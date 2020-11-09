package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ExecutionStatus describes msg execution status
type ExecutionStatus struct {
	MsgType  string
	FeePayer sdk.AccAddress
	Success  bool
}
