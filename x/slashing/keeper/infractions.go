package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/crypto"

	"github.com/KiraCore/sekai/x/slashing/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// HandleValidatorSignature handles a validator signature, must be called once per validator per block.
func (k Keeper) HandleValidatorSignature(ctx sdk.Context, addr crypto.Address, power int64, signed bool) {
	fmt.Println("HandleValidatorSignature1")
	logger := k.Logger(ctx)
	height := ctx.BlockHeight()

	// fetch the validator public key
	consAddr := sdk.ConsAddress(addr)
	if _, err := k.GetPubkey(ctx, addr); err != nil {
		panic(fmt.Sprintf("Validator consensus-address %s not found: %s", consAddr, err.Error()))
	}

	// fetch signing info
	signInfo, found := k.GetValidatorSigningInfo(ctx, consAddr)
	if !found {
		panic(fmt.Sprintf("Expected signing info for validator %s but not found", consAddr))
	}

	// this is a relative index, so it counts blocks the validator *should* have signed
	// will use the 0-value default signing info if not present, except for start height
	index := signInfo.IndexOffset % k.SignedBlocksWindow(ctx)
	signInfo.IndexOffset++

	// Update signed block bit array & counter
	// This counter just tracks the sum of the bit array
	// That way we avoid needing to read/write the whole array each time
	previous := k.GetValidatorMissedBlockBitArray(ctx, consAddr, index)
	missed := !signed
	switch {
	case !previous && missed:
		// Array value has changed from not missed to missed
		k.SetValidatorMissedBlockBitArray(ctx, consAddr, index, true)
	case previous && !missed:
		// Array value has changed from missed to not missed
		k.SetValidatorMissedBlockBitArray(ctx, consAddr, index, false)
	default:
		// Array value at this index has not changed, no need to update counter
	}
	if missed { // increment counter
		signInfo.MissedBlocksCounter++
	} else { // set counter to 0
		signInfo.MissedBlocksCounter = 0
	}

	validator, err := k.sk.GetValidatorByConsAddr(ctx, consAddr)
	if err == nil && !validator.IsInactivated() {
		k.sk.HandleValidatorSignature(ctx, validator.ValKey, missed)
	}

	if missed {
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeLiveness,
				sdk.NewAttribute(types.AttributeKeyAddress, consAddr.String()),
				sdk.NewAttribute(types.AttributeKeyMissedBlocks, fmt.Sprintf("%d", signInfo.MissedBlocksCounter)),
				sdk.NewAttribute(types.AttributeKeyHeight, fmt.Sprintf("%d", height)),
			),
		)

		logger.Info(
			fmt.Sprintf("Absent validator %s at height %d, %d missed, threshold %d", consAddr, height, signInfo.MissedBlocksCounter, k.MinSignedPerWindow(ctx)))
	}

	minHeight := signInfo.StartHeight + k.SignedBlocksWindow(ctx)
	maxMissed := k.SignedBlocksWindow(ctx) - k.MinSignedPerWindow(ctx)

	// if we are past the minimum height and the validator has missed too many blocks, punish them
	if height > minHeight && signInfo.MissedBlocksCounter > maxMissed {
		validator, err := k.sk.GetValidatorByConsAddr(ctx, consAddr)
		if err == nil && !validator.IsInactivated() {

			// Downtime confirmed: slash and inactivate the validator
			logger.Info(fmt.Sprintf("Validator %s past min height of %d and below signed blocks threshold of %d",
				consAddr, minHeight, k.MinSignedPerWindow(ctx)))

			ctx.EventManager().EmitEvent(
				sdk.NewEvent(
					types.EventTypeInactivate,
					sdk.NewAttribute(types.AttributeKeyAddress, consAddr.String()),
					sdk.NewAttribute(types.AttributeKeyPower, fmt.Sprintf("%d", power)),
					sdk.NewAttribute(types.AttributeKeyReason, types.AttributeValueMissingSignature),
					sdk.NewAttribute(types.AttributeKeyInactivated, consAddr.String()),
				),
			)
			k.sk.Inactivate(ctx, validator.ValKey)

			signInfo.InactiveUntil = ctx.BlockHeader().Time.Add(k.DowntimeInactiveDuration(ctx))

			// We need to reset the counter & array so that the validator won't be immediately inactivated for downtime upon rebonding.
			signInfo.MissedBlocksCounter = 0
			signInfo.IndexOffset = 0
			k.clearValidatorMissedBlockBitArray(ctx, consAddr)
		} else {
			// Validator was (a) not found or (b) already inactivated, don't slash
			logger.Info(
				fmt.Sprintf("Validator %s would have been inactivated for downtime, but was either not found in store or already inactivated", consAddr),
			)
		}
	}

	fmt.Println("HandleValidatorSignature2", consAddr.String(), signInfo)
	// Set the updated signing info
	k.SetValidatorSigningInfo(ctx, consAddr, signInfo)
}
