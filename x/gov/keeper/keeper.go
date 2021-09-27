package keeper

import (
	"errors"
	"fmt"

	"github.com/KiraCore/sekai/x/gov/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Keeper struct {
	cdc            codec.BinaryCodec
	storeKey       sdk.StoreKey
	bk             types.BankKeeper
	proposalRouter types.ProposalRouter
}

func NewKeeper(storeKey sdk.StoreKey, cdc codec.BinaryCodec, bk types.BankKeeper) Keeper {
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
func (k Keeper) SetNetworkProperties(ctx sdk.Context, properties *types.NetworkProperties) error {
	err := k.ValidateNetworkProperties(ctx, properties)
	if err != nil {
		return err
	}
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixNetworkProperties)
	prefixStore.Set([]byte("property"), k.cdc.MustMarshal(properties))
	return nil
}

// GetNetworkProperties get network properties from KVStore
func (k Keeper) GetNetworkProperties(ctx sdk.Context) *types.NetworkProperties {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixNetworkProperties)
	bz := prefixStore.Get([]byte("property"))

	properties := new(types.NetworkProperties)
	k.cdc.MustUnmarshal(bz, properties)
	return properties
}

func (k Keeper) ValidateNetworkProperties(ctx sdk.Context, properties *types.NetworkProperties) error {

	if properties.MinTxFee == 0 {
		return fmt.Errorf("min_tx_fee should not be ZERO")
	}
	if properties.MaxTxFee == 0 {
		return fmt.Errorf("max_tx_fee should not be ZERO")
	}
	if properties.MaxTxFee < properties.MinTxFee {
		return fmt.Errorf("max_tx_fee should not be lower than min_tx_fee")
	}
	// TODO: for now skipping few of validations
	// if properties.VoteQuorum == 0 {
	// 	return fmt.Errorf("vote_quorum should not be zero")
	// }
	// if properties.ProposalEndTime == 0 {
	// 	return fmt.Errorf("proposal_end_time should not be zero")
	// }
	// if properties.ProposalEnactmentTime == 0 {
	// 	return fmt.Errorf("proposal_enactment_time should not be zero")
	// }
	// if properties.MinProposalEndBlocks == 0 {
	// 	return fmt.Errorf("min_proposal_end_blocks should not be zero")
	// }
	// if properties.MinProposalEnactmentBlocks == 0 {
	// 	return fmt.Errorf("min_proposal_enactment_blocks should not be zero")
	// }
	// if properties.MischanceRankDecreaseAmount == 0 {
	// 	return fmt.Errorf("mischance_rank_decrease_amount should not be zero")
	// }
	// if properties.MaxMischance == 0 {
	// 	return fmt.Errorf("max_mischance should not be zero")
	// }
	// if properties.InactiveRankDecreasePercent == 0 {
	// 	return fmt.Errorf("inactive_rank_decrease_percent should not be zero")
	// }
	// if properties.InactiveRankDecreasePercent == 0 {
	// 	return fmt.Errorf("inactive_rank_decrease_percent should not be zero")
	// }
	if properties.InactiveRankDecreasePercent > 100 {
		return fmt.Errorf("inactive_rank_decrease_percent should not be lower than 100%%")
	}
	// if properties.MinValidators == 0 {
	// 	return fmt.Errorf("min_validators should not be zero")
	// }
	// if properties.PoorNetworkMaxBankSend == 0 {
	// 	return fmt.Errorf("min_validators should not be zero")
	// }
	// if properties.JailMaxTime == 0 {
	// 	return fmt.Errorf("jail_max_time should not be zero")
	// }
	// fee := k.GetExecutionFee(ctx, (&types.MsgHandleIdentityRecordsVerifyRequest{}).Type())
	// maxFee := properties.MinTxFee
	// if fee != nil {
	// 	if maxFee < fee.ExecutionFee {
	// 		maxFee = fee.ExecutionFee
	// 	}
	// 	if maxFee < fee.FailureFee {
	// 		maxFee = fee.FailureFee
	// 	}
	// }
	// if properties.MinIdentityApprovalTip < maxFee*2 {
	// 	return fmt.Errorf("min_identity_approval_tip should not be bigger or equal than 2x approval fee")
	// }
	// if properties.UniqueIdentityKeys == "" {
	// 	return fmt.Errorf("unique_identity_keys should not be empty")
	// }
	// monikerExists := false
	// if properties.UniqueIdentityKeys != FormalizeIdentityRecordKey(properties.UniqueIdentityKeys) {
	// 	return fmt.Errorf("unique identity keys on network property should be formailzed with lowercase keys")
	// }
	// uniqueKeys := strings.Split(properties.UniqueIdentityKeys, ",")
	// for _, key := range uniqueKeys {
	// 	if !ValidateIdentityRecordKey(key) {
	// 		return fmt.Errorf("invalid identity record key exists, key=%s", key)
	// 	}
	// 	if key == "moniker" {
	// 		monikerExists = true
	// 	}
	// }
	// if !monikerExists {
	// 	return fmt.Errorf("moniker should be exist in unique keys list")
	// }
	return nil
}

// GetNetworkProperty get single network property by key
func (k Keeper) GetNetworkProperty(ctx sdk.Context, property types.NetworkProperty) (types.NetworkPropertyValue, error) {
	properties := k.GetNetworkProperties(ctx)
	switch property {
	case types.MinTxFee:
		return types.NetworkPropertyValue{Value: properties.MinTxFee}, nil
	case types.MaxTxFee:
		return types.NetworkPropertyValue{Value: properties.MaxTxFee}, nil
	case types.VoteQuorum:
		return types.NetworkPropertyValue{Value: properties.VoteQuorum}, nil
	case types.ProposalEndTime:
		return types.NetworkPropertyValue{Value: properties.ProposalEndTime}, nil
	case types.ProposalEnactmentTime:
		return types.NetworkPropertyValue{Value: properties.ProposalEnactmentTime}, nil
	case types.MinProposalEndBlocks:
		return types.NetworkPropertyValue{Value: properties.MinProposalEndBlocks}, nil
	case types.MinProposalEnactmentBlocks:
		return types.NetworkPropertyValue{Value: properties.MinProposalEnactmentBlocks}, nil
	case types.EnableForeignFeePayments:
		return types.NetworkPropertyValue{Value: BoolToInt(properties.EnableForeignFeePayments)}, nil
	case types.MischanceRankDecreaseAmount:
		return types.NetworkPropertyValue{Value: properties.MischanceRankDecreaseAmount}, nil
	case types.MaxMischance:
		return types.NetworkPropertyValue{Value: properties.MaxMischance}, nil
	case types.MischanceConfidence:
		return types.NetworkPropertyValue{Value: properties.MischanceConfidence}, nil
	case types.InactiveRankDecreasePercent:
		return types.NetworkPropertyValue{Value: properties.InactiveRankDecreasePercent}, nil
	case types.PoorNetworkMaxBankSend:
		return types.NetworkPropertyValue{Value: properties.PoorNetworkMaxBankSend}, nil
	case types.MinValidators:
		return types.NetworkPropertyValue{Value: properties.MinValidators}, nil
	case types.JailMaxTime:
		return types.NetworkPropertyValue{Value: properties.JailMaxTime}, nil
	case types.EnableTokenWhitelist:
		return types.NetworkPropertyValue{Value: BoolToInt(properties.EnableTokenWhitelist)}, nil
	case types.EnableTokenBlacklist:
		return types.NetworkPropertyValue{Value: BoolToInt(properties.EnableTokenBlacklist)}, nil
	case types.MinIdentityApprovalTip:
		return types.NetworkPropertyValue{Value: properties.MinIdentityApprovalTip}, nil
	case types.UniqueIdentityKeys:
		return types.NetworkPropertyValue{StrValue: properties.UniqueIdentityKeys}, nil
	default:
		return types.NetworkPropertyValue{}, errors.New("trying to fetch network property that does not exist")
	}
}

// SetNetworkProperty set single network property by key
func (k Keeper) SetNetworkProperty(ctx sdk.Context, property types.NetworkProperty, value types.NetworkPropertyValue) error {
	properties := k.GetNetworkProperties(ctx)
	switch property {
	case types.MinTxFee:
		properties.MinTxFee = value.Value
	case types.MaxTxFee:
		properties.MaxTxFee = value.Value
	case types.VoteQuorum:
		properties.VoteQuorum = value.Value
	case types.ProposalEndTime:
		properties.ProposalEndTime = value.Value
	case types.ProposalEnactmentTime:
		properties.ProposalEnactmentTime = value.Value
	case types.MinProposalEndBlocks:
		properties.MinProposalEndBlocks = value.Value
	case types.MinProposalEnactmentBlocks:
		properties.MinProposalEnactmentBlocks = value.Value
	case types.EnableForeignFeePayments:
		if value.Value > 0 {
			properties.EnableForeignFeePayments = true
		}
		properties.EnableForeignFeePayments = false
	case types.MischanceRankDecreaseAmount:
		properties.MischanceRankDecreaseAmount = value.Value
	case types.MaxMischance:
		properties.MaxMischance = value.Value
	case types.MischanceConfidence:
		properties.MischanceConfidence = value.Value
	case types.InactiveRankDecreasePercent:
		properties.InactiveRankDecreasePercent = value.Value
	case types.PoorNetworkMaxBankSend:
		properties.PoorNetworkMaxBankSend = value.Value
	case types.MinValidators:
		properties.MinValidators = value.Value
	case types.JailMaxTime:
		properties.JailMaxTime = value.Value
	case types.EnableTokenBlacklist:
		properties.EnableTokenBlacklist = IntToBool(value.Value)
	case types.EnableTokenWhitelist:
		properties.EnableTokenWhitelist = IntToBool(value.Value)
	case types.MinIdentityApprovalTip:
		properties.MinIdentityApprovalTip = value.Value
	case types.UniqueIdentityKeys:
		properties.UniqueIdentityKeys = value.StrValue
	default:
		return errors.New("trying to set network property that does not exist")
	}
	return k.SetNetworkProperties(ctx, properties)
}

// SetExecutionFee set fee by execution function name
func (k Keeper) SetExecutionFee(ctx sdk.Context, fee *types.ExecutionFee) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixExecutionFee)
	key := []byte(fee.TransactionType)
	prefixStore.Set(key, k.cdc.MustMarshal(fee))
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
	k.cdc.MustUnmarshal(bz, fee)
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
		k.cdc.MustUnmarshal(bz, fee)
	}
	return fees
}
