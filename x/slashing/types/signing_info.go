package types

import (
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewValidatorSigningInfo creates a new ValidatorSigningInfo instance
//nolint:interfacer
func NewValidatorSigningInfo(
	consAddr sdk.ConsAddress, startHeight int64,
	inactivatedUntil time.Time,
	mischance, missedBlocksCounter, producedBlocksCounter int64,
) ValidatorSigningInfo {

	return ValidatorSigningInfo{
		Address:               consAddr.String(),
		StartHeight:           startHeight,
		InactiveUntil:         inactivatedUntil,
		Mischance:             mischance,
		MissedBlocksCounter:   missedBlocksCounter,
		ProducedBlocksCounter: producedBlocksCounter,
	}
}

// String implements the stringer interface for ValidatorSigningInfo
func (i ValidatorSigningInfo) String() string {
	return fmt.Sprintf(`Validator Signing Info:
  Address:                %s
  Start Height:           %d
  Inactivated Until:      %v
  Mischance:              %d
  Missed Blocks Counter:  %d
  Produced Blocks Counter: %d`,
		i.Address, i.StartHeight, i.InactiveUntil,
		i.Mischance, i.MissedBlocksCounter, i.ProducedBlocksCounter)
}

// unmarshal a validator signing info from a store value
func UnmarshalValSigningInfo(cdc codec.Codec, value []byte) (signingInfo ValidatorSigningInfo, err error) {
	err = cdc.Unmarshal(value, &signingInfo)
	return signingInfo, err
}
