package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type zeroGasMeter struct{}

// NewZeroGasMeter returns a reference to a new zeroGasMeter.
func NewZeroGasMeter() sdk.GasMeter {
	return &zeroGasMeter{}
}

func (g *zeroGasMeter) GasConsumed() sdk.Gas {
	return 0
}

func (g *zeroGasMeter) GasConsumedToLimit() sdk.Gas {
	return 1
}

func (g *zeroGasMeter) Limit() sdk.Gas {
	return 0
}

func (g *zeroGasMeter) ConsumeGas(amount sdk.Gas, descriptor string) {
}

func (g *zeroGasMeter) RefundGas(amount sdk.Gas, descriptor string) {
}

func (g *zeroGasMeter) IsPastLimit() bool {
	return false
}

func (g *zeroGasMeter) IsOutOfGas() bool {
	return false
}

func (g *zeroGasMeter) String() string {
	return fmt.Sprintf("InfiniteGasMeter:\n  consumed: %d", 0)
}
