package keeper

import (
	"errors"

	"github.com/KiraCore/sekai/x/gov/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Keeper struct {
	cdc            codec.BinaryMarshaler
	storeKey       sdk.StoreKey
	bk             types.BankKeeper
	proposalRouter types.ProposalRouter
}

func NewKeeper(storeKey sdk.StoreKey, cdc codec.BinaryMarshaler, bk types.BankKeeper) Keeper {
	return Keeper{
		cdc:      cdc,
		storeKey: storeKey,
		bk:       bk,
	}
}

// BondDenom returns the denom that is basically used for fee payment
func (k Keeper) BondDenom(ctx sdk.Context) string {
	return "ukex"
}

func (k *Keeper) SetProposalRouter(proposalRouter types.ProposalRouter) {
	k.proposalRouter = proposalRouter
}

func (k Keeper) GetProposalRouter() types.ProposalRouter {
	return k.proposalRouter
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
	case types.MinProposalEndBlocks:
		return properties.MinProposalEndBlocks, nil
	case types.MinProposalEnactmentBlocks:
		return properties.MinProposalEnactmentBlocks, nil
	case types.EnableForeignFeePayments:
		return BoolToInt(properties.EnableForeignFeePayments), nil
	case types.MischanceRankDecreaseAmount:
		return properties.MischanceRankDecreaseAmount, nil
	case types.MaxMischance:
		return properties.MaxMischance, nil
	case types.MischanceConfidence:
		return properties.MischanceConfidence, nil
	case types.InactiveRankDecreasePercent:
		return properties.InactiveRankDecreasePercent, nil
	case types.PoorNetworkMaxBankSend:
		return properties.PoorNetworkMaxBankSend, nil
	case types.MinValidators:
		return properties.MinValidators, nil
	case types.JailMaxTime:
		return properties.JailMaxTime, nil
	case types.EnableTokenWhitelist:
		return BoolToInt(properties.EnableTokenWhitelist), nil
	case types.EnableTokenBlacklist:
		return BoolToInt(properties.EnableTokenBlacklist), nil
	case types.MinIdentityApprovalTip:
		return properties.MinIdentityApprovalTip, nil
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
	case types.MaxTxFee:
		properties.MaxTxFee = value
	case types.VoteQuorum:
		properties.VoteQuorum = value
	case types.ProposalEndTime:
		properties.ProposalEndTime = value
	case types.ProposalEnactmentTime:
		properties.ProposalEnactmentTime = value
	case types.MinProposalEndBlocks:
		properties.MinProposalEndBlocks = value
	case types.MinProposalEnactmentBlocks:
		properties.MinProposalEnactmentBlocks = value
	case types.EnableForeignFeePayments:
		if value > 0 {
			properties.EnableForeignFeePayments = true
		}
		properties.EnableForeignFeePayments = false
	case types.MischanceRankDecreaseAmount:
		properties.MischanceRankDecreaseAmount = value
	case types.MaxMischance:
		properties.MaxMischance = value
	case types.MischanceConfidence:
		properties.MischanceConfidence = value
	case types.InactiveRankDecreasePercent:
		properties.InactiveRankDecreasePercent = value
	case types.PoorNetworkMaxBankSend:
		properties.PoorNetworkMaxBankSend = value
	case types.MinValidators:
		properties.MinValidators = value
	case types.JailMaxTime:
		properties.JailMaxTime = value
	case types.EnableTokenBlacklist:
		properties.EnableTokenBlacklist = IntToBool(value)
	case types.EnableTokenWhitelist:
		properties.EnableTokenWhitelist = IntToBool(value)
	case types.MinIdentityApprovalTip:
		properties.MinIdentityApprovalTip = value
	default:
		return errors.New("trying to set network property that does not exist")
	}
	k.SetNetworkProperties(ctx, properties)
	return nil
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

// GetExecutionFees get fees from execution function name
func (k Keeper) GetExecutionFees(ctx sdk.Context) []*types.ExecutionFee {
	iterator := sdk.KVStorePrefixIterator(ctx.KVStore(k.storeKey), types.KeyPrefixExecutionFee)
	defer iterator.Close()
	fees := []*types.ExecutionFee{}
	for ; iterator.Valid(); iterator.Next() {
		bz := iterator.Value()
		fee := new(types.ExecutionFee)
		k.cdc.MustUnmarshalBinaryBare(bz, fee)
	}
	return fees
}
