package keeper

import (
	"errors"

	"github.com/KiraCore/sekai/x/gov/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Keeper struct {
	cdc      codec.BinaryMarshaler
	storeKey sdk.StoreKey
}

func NewKeeper(storeKey sdk.StoreKey, cdc codec.BinaryMarshaler) Keeper {
	return Keeper{cdc: cdc, storeKey: storeKey}
}

// BondDenom returns the denom that is basically used for fee payment
func (k Keeper) BondDenom(ctx sdk.Context) string {
	return "ukex"
}

// SetNetworkProperties set network properties on KVStore
func (k Keeper) SetNetworkProperties(ctx sdk.Context, properties *types.NetworkProperties) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixNetworkProperties)
	prefixStore.Set([]byte("property"), k.cdc.MustMarshalBinaryBare(properties))
}

// GetNetworkProperties get network properties from KVStore
func (k Keeper) GetNetworkProperties(ctx sdk.Context) *types.NetworkProperties {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixNetworkProperties)
	bz := prefixStore.Get([]byte("property"))

	properties := new(types.NetworkProperties)
	k.cdc.MustUnmarshalBinaryBare(bz, properties)
	return properties
}

// GetNetworkProperty get single network property by key
func (k Keeper) GetNetworkProperty(ctx sdk.Context, property types.NetworkProperty) (uint64, error) {
	properties := k.GetNetworkProperties(ctx)
	switch property {
	case types.MinTxFee:
		return properties.MinTxFee, nil
	case types.MaxTxFee:
		return properties.MaxTxFee, nil
	case types.VoteQuorum:
		return properties.VoteQuorum, nil
	case types.ProposalEndTime:
		return properties.ProposalEndTime, nil
	case types.ProposalEnactmentTime:
		return properties.ProposalEnactmentTime, nil
	case types.EnableForeignFeePayments:
		if properties.EnableForeignFeePayments {
			return 1, nil
		}
		return 0, nil
	default:
		return 0, errors.New("trying to fetch network property that does not exist")
	}
}

// SetNetworkProperty set single network property by key
func (k Keeper) SetNetworkProperty(ctx sdk.Context, property types.NetworkProperty, value uint64) error {
	properties := k.GetNetworkProperties(ctx)
	switch property {
	case types.MinTxFee:
		properties.MinTxFee = value
		return nil
	case types.MaxTxFee:
		properties.MaxTxFee = value
		return nil
	case types.VoteQuorum:
		properties.VoteQuorum = value
		return nil
	case types.ProposalEndTime:
		properties.ProposalEndTime = value
		return nil
	case types.ProposalEnactmentTime:
		properties.ProposalEnactmentTime = value
		return nil
	case types.EnableForeignFeePayments:
		if value > 0 {
			properties.EnableForeignFeePayments = true
		}
		properties.EnableForeignFeePayments = false
		return nil
	default:
		return errors.New("trying to set network property that does not exist")
	}
}

// SetExecutionFee set fee by execution function name
func (k Keeper) SetExecutionFee(ctx sdk.Context, fee *types.ExecutionFee) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixExecutionFee)
	key := []byte(fee.TransactionType)
	prefixStore.Set(key, k.cdc.MustMarshalBinaryBare(fee))
}

// GetExecutionFee get fee from execution function name
func (k Keeper) GetExecutionFee(ctx sdk.Context, txType string) *types.ExecutionFee {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixExecutionFee)
	key := []byte(txType)
	if !prefixStore.Has(key) {
		return nil
	}
	bz := prefixStore.Get([]byte(txType))

	fee := new(types.ExecutionFee)
	k.cdc.MustUnmarshalBinaryBare(bz, fee)
	return fee
}
