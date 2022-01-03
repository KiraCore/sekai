package types

import (
	"fmt"
	"time"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Default parameter namespace
const (
	DefaultDowntimeInactiveDuration = 60 * 10 * time.Second
)

// Parameter store keys
var (
	KeyDowntimeInactiveDuration = []byte("DowntimeInactiveDuration")
)

// ParamKeyTable for slashing module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params object
func NewParams(
	downtimeInactiveDuration time.Duration,
) Params {

	return Params{
		DowntimeInactiveDuration: downtimeInactiveDuration,
	}
}

// ParamSetPairs - Implements params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyDowntimeInactiveDuration, &p.DowntimeInactiveDuration, validateDowntimeInactiveDuration),
	}
}

// DefaultParams defines the parameters for this module
func DefaultParams() Params {
	return NewParams(
		DefaultDowntimeInactiveDuration,
	)
}

func validateDowntimeInactiveDuration(i interface{}) error {
	v, ok := i.(time.Duration)
	if !ok {
		return fmt.Errorf("invalid downtime_inactive_duration parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("downtime inactive duration must be positive: %s", v)
	}

	return nil
}
