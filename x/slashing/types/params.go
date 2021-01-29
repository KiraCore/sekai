package types

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Default parameter namespace
const (
	DefaultSignedBlocksWindow       = int64(100)
	DefaultDowntimeInactiveDuration = 60 * 10 * time.Second
)

var (
	DefaultMinSignedPerWindow = sdk.NewDecWithPrec(5, 1)
)

// Parameter store keys
var (
	KeySignedBlocksWindow       = []byte("SignedBlocksWindow")
	KeyMinSignedPerWindow       = []byte("MinSignedPerWindow")
	KeyDowntimeInactiveDuration = []byte("DowntimeInactiveDuration")
)

// ParamKeyTable for slashing module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params object
func NewParams(
	signedBlocksWindow int64, minSignedPerWindow sdk.Dec, DowntimeInactiveDuration time.Duration,
) Params {

	return Params{
		SignedBlocksWindow:       signedBlocksWindow,
		MinSignedPerWindow:       minSignedPerWindow,
		DowntimeInactiveDuration: DowntimeInactiveDuration,
	}
}

// ParamSetPairs - Implements params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeySignedBlocksWindow, &p.SignedBlocksWindow, validateSignedBlocksWindow),
		paramtypes.NewParamSetPair(KeyMinSignedPerWindow, &p.MinSignedPerWindow, validateMinSignedPerWindow),
		paramtypes.NewParamSetPair(KeyDowntimeInactiveDuration, &p.DowntimeInactiveDuration, validateDowntimeInactiveDuration),
	}
}

// DefaultParams defines the parameters for this module
func DefaultParams() Params {
	return NewParams(
		DefaultSignedBlocksWindow, DefaultMinSignedPerWindow, DefaultDowntimeInactiveDuration,
	)
}

func validateSignedBlocksWindow(i interface{}) error {
	v, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("signed blocks window must be positive: %d", v)
	}

	return nil
}

func validateMinSignedPerWindow(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNegative() {
		return fmt.Errorf("min signed per window cannot be negative: %s", v)
	}
	if v.GT(sdk.OneDec()) {
		return fmt.Errorf("min signed per window too large: %s", v)
	}

	return nil
}

func validateDowntimeInactiveDuration(i interface{}) error {
	v, ok := i.(time.Duration)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("downtime inactive duration must be positive: %s", v)
	}

	return nil
}
