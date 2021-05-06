package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/crypto"

	"github.com/KiraCore/sekai/x/slashing/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// HandleValidatorSignature handles a validator signature, must be called once per validator per block.
func (k Keeper) HandleValidatorSignature(ctx sdk.Context, addr crypto.Address, power int64, signed bool) {
	logger := k.Logger(ctx)
	height := ctx.BlockHeight()

	// fetch the validator public key
	consAddr := sdk.ConsAddress(addr)
	if _, err := k.GetPubkey(ctx, addr); err != nil {
		panic(fmt.Sprintf("Validator consensus-address %s not found: %s", consAddr.String(), err.Error()))
	}

	validator, err := k.sk.GetValidatorByConsAddr(ctx, consAddr)
	if err != nil {
		panic(fmt.Sprintf("Validator not found by consensus-address: %s", consAddr.String()))
	}

	if !validator.IsActive() {
		// if validator is not active, stop performance counter and status change by validator signature
		return
	}

	// fetch signing info
	signInfo, found := k.GetValidatorSigningInfo(ctx, consAddr)
	if !found {
		panic(fmt.Sprintf("Expected signing info for validator %s but not found", consAddr))
	}

	properties := k.gk.GetNetworkProperties(ctx)

	// Update uptime counter
	missed := !signed
	if missed { // increment counter
		signInfo.MissedBlocksCounter++
		// increment mischance only when missed blocks are bigger than mischance confidence
		if signInfo.MischanceConfidence >= int64(properties.MischanceConfidence) {
			signInfo.Mischance++
		} else {
			signInfo.MischanceConfidence++
		}
	} else { // set counter to 0
		signInfo.Mischance = 0
		signInfo.MischanceConfidence = 0
		signInfo.ProducedBlocksCounter++
		signInfo.LastPresentBlock = ctx.BlockHeight()
	}

	// handle staking module's validator object update actions
	k.sk.HandleValidatorSignature(ctx, validator.ValKey, missed, signInfo.Mischance)

	if missed {
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeLiveness,
				sdk.NewAttribute(types.AttributeKeyAddress, consAddr.String()),
				sdk.NewAttribute(types.AttributeKeyMischance, fmt.Sprintf("%d", signInfo.Mischance)),
				sdk.NewAttribute(types.AttributeKeyLastPresentBlock, fmt.Sprintf("%d", signInfo.LastPresentBlock)),
				sdk.NewAttribute(types.AttributeKeyMissedBlocks, fmt.Sprintf("%d", signInfo.MissedBlocksCounter)),
				sdk.NewAttribute(types.AttributeKeyProducedBlocks, fmt.Sprintf("%d", signInfo.ProducedBlocksCounter)),
				sdk.NewAttribute(types.AttributeKeyHeight, fmt.Sprintf("%d", height)),
			),
		)

		logger.Info(
			fmt.Sprintf("Absent validator %s at height %d, %d mischance, %d missed blocks, %d produced blocks, threshold %d", consAddr, height, signInfo.Mischance, signInfo.MissedBlocksCounter, signInfo.ProducedBlocksCounter, properties.MaxMischance))
	}

	// if mischance overflow max value, we punish them
	if signInfo.Mischance > int64(properties.MaxMischance) {
		// Downtime confirmed: slash and inactivate the validator
		logger.Info(fmt.Sprintf("Validator %s past max mischance threshold of %d",
			consAddr, properties.MaxMischance))

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
	}

	// Set the updated signing info
	k.SetValidatorSigningInfo(ctx, consAddr, signInfo)
}
