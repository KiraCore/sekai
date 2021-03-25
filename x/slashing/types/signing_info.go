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
	condAddr sdk.ConsAddress, startHeight, indexOffset int64,
	inactivatedUntil time.Time, tombstoned bool,
	mischance, missedBlocksCounter, producedBlocksCounter int64,
) ValidatorSigningInfo {

	return ValidatorSigningInfo{
		Address:               condAddr.String(),
		StartHeight:           startHeight,
		IndexOffset:           indexOffset,
		InactiveUntil:         inactivatedUntil,
		Tombstoned:            tombstoned,
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
  Index Offset:           %d
  Inactivated Until:      %v
  Tombstoned:             %t
  Mischance:              %d
  Missed Blocks Counter:  %d
  Produced Blocks Counter: %d`,
		i.Address, i.StartHeight, i.IndexOffset, i.InactiveUntil,
		i.Tombstoned, i.Mischance, i.MissedBlocksCounter, i.ProducedBlocksCounter)
}

// unmarshal a validator signing info from a store value
func UnmarshalValSigningInfo(cdc codec.Marshaler, value []byte) (signingInfo ValidatorSigningInfo, err error) {
	err = cdc.UnmarshalBinaryBare(value, &signingInfo)
	return signingInfo, err
}
