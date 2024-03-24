package keeper

import (
	"errors"
	"fmt"

	"github.com/KiraCore/sekai/x/gov/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Keeper struct {
	cdc            codec.BinaryCodec
	storeKey       storetypes.StoreKey
	bk             types.BankKeeper
	proposalRouter types.ProposalRouter
}

func NewKeeper(storeKey storetypes.StoreKey, cdc codec.BinaryCodec, bk types.BankKeeper) Keeper {
	return Keeper{
		cdc:      cdc,
		storeKey: storeKey,
		bk:       bk,
	}
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
	if !properties.InactiveRankDecreasePercent.IsNil() && properties.InactiveRankDecreasePercent.GT(sdk.OneDec()) {
		return fmt.Errorf("inactive_rank_decrease_percent should not be lower than 100%%")
	}
	// if properties.MinValidators == 0 {
	// 	return fmt.Errorf("min_validators should not be zero")
	// }
	// if properties.PoorNetworkMaxBankSend == 0 {
	// 	return fmt.Errorf("min_validators should not be zero")
	// }
	// if properties.UnjailMaxTime == 0 {
	// 	return fmt.Errorf("unjail_max_time should not be zero")
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
	case types.MinimumProposalEndTime:
		return types.NetworkPropertyValue{Value: properties.MinimumProposalEndTime}, nil
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
		return types.NetworkPropertyValue{StrValue: properties.InactiveRankDecreasePercent.String()}, nil
	case types.PoorNetworkMaxBankSend:
		return types.NetworkPropertyValue{Value: properties.PoorNetworkMaxBankSend}, nil
	case types.MinValidators:
		return types.NetworkPropertyValue{Value: properties.MinValidators}, nil
	case types.UnjailMaxTime:
		return types.NetworkPropertyValue{Value: properties.UnjailMaxTime}, nil
	case types.EnableTokenWhitelist:
		return types.NetworkPropertyValue{Value: BoolToInt(properties.EnableTokenWhitelist)}, nil
	case types.EnableTokenBlacklist:
		return types.NetworkPropertyValue{Value: BoolToInt(properties.EnableTokenBlacklist)}, nil
	case types.MinIdentityApprovalTip:
		return types.NetworkPropertyValue{Value: properties.MinIdentityApprovalTip}, nil
	case types.UniqueIdentityKeys:
		return types.NetworkPropertyValue{StrValue: properties.UniqueIdentityKeys}, nil
	case types.UbiHardcap:
		return types.NetworkPropertyValue{Value: properties.UbiHardcap}, nil
	case types.ValidatorsFeeShare:
		return types.NetworkPropertyValue{StrValue: properties.ValidatorsFeeShare.String()}, nil
	case types.InflationRate:
		return types.NetworkPropertyValue{StrValue: properties.InflationRate.String()}, nil
	case types.InflationPeriod:
		return types.NetworkPropertyValue{Value: properties.InflationPeriod}, nil
	case types.UnstakingPeriod:
		return types.NetworkPropertyValue{Value: properties.UnstakingPeriod}, nil
	case types.MaxDelegators:
		return types.NetworkPropertyValue{Value: properties.MaxDelegators}, nil
	case types.MinDelegationPushout:
		return types.NetworkPropertyValue{Value: properties.MinDelegationPushout}, nil
	case types.SlashingPeriod:
		return types.NetworkPropertyValue{Value: properties.SlashingPeriod}, nil
	case types.MaxJailedPercentage:
		return types.NetworkPropertyValue{StrValue: properties.MaxJailedPercentage.String()}, nil
	case types.MaxSlashingPercentage:
		return types.NetworkPropertyValue{StrValue: properties.MaxSlashingPercentage.String()}, nil
	case types.MinCustodyReward:
		return types.NetworkPropertyValue{Value: properties.MinCustodyReward}, nil
	case types.MaxCustodyBufferSize:
		return types.NetworkPropertyValue{Value: properties.MaxCustodyBufferSize}, nil
	case types.MaxCustodyTxSize:
		return types.NetworkPropertyValue{Value: properties.MaxCustodyTxSize}, nil
	case types.AbstentionRankDecreaseAmount:
		return types.NetworkPropertyValue{Value: properties.AbstentionRankDecreaseAmount}, nil
	case types.MaxAbstention:
		return types.NetworkPropertyValue{Value: properties.MaxAbstention}, nil
	case types.MinCollectiveBond:
		return types.NetworkPropertyValue{Value: properties.MinCollectiveBond}, nil
	case types.MinCollectiveBondingTime:
		return types.NetworkPropertyValue{Value: properties.MinCollectiveBondingTime}, nil
	case types.MaxCollectiveOutputs:
		return types.NetworkPropertyValue{Value: properties.MaxCollectiveOutputs}, nil
	case types.MinCollectiveClaimPeriod:
		return types.NetworkPropertyValue{Value: properties.MinCollectiveClaimPeriod}, nil
	case types.ValidatorRecoveryBond:
		return types.NetworkPropertyValue{Value: properties.ValidatorRecoveryBond}, nil
	case types.MaxAnnualInflation:
		return types.NetworkPropertyValue{StrValue: properties.MaxAnnualInflation.String()}, nil
	case types.MinDappBond:
		return types.NetworkPropertyValue{Value: properties.MinDappBond}, nil
	case types.MaxDappBond:
		return types.NetworkPropertyValue{Value: properties.MaxDappBond}, nil
	case types.DappBondDuration:
		return types.NetworkPropertyValue{Value: properties.DappBondDuration}, nil
	case types.DappVerifierBond:
		return types.NetworkPropertyValue{StrValue: properties.DappVerifierBond.String()}, nil
	case types.DappAutoDenounceTime:
		return types.NetworkPropertyValue{Value: properties.DappAutoDenounceTime}, nil
	case types.DappMischanceRankDecreaseAmount:
		return types.NetworkPropertyValue{Value: properties.DappMischanceRankDecreaseAmount}, nil
	case types.DappMaxMischance:
		return types.NetworkPropertyValue{Value: properties.DappMaxMischance}, nil
	case types.DappInactiveRankDecreasePercent:
		return types.NetworkPropertyValue{Value: properties.DappInactiveRankDecreasePercent}, nil
	case types.DappPoolSlippageDefault:
		return types.NetworkPropertyValue{StrValue: properties.DappPoolSlippageDefault.String()}, nil
	case types.MintingFtFee:
		return types.NetworkPropertyValue{Value: properties.MintingFtFee}, nil
	case types.MintingNftFee:
		return types.NetworkPropertyValue{Value: properties.MintingNftFee}, nil
	case types.VetoThreshold:
		return types.NetworkPropertyValue{StrValue: properties.VetoThreshold.String()}, nil
	case types.AutocompoundIntervalNumBlocks:
		return types.NetworkPropertyValue{Value: properties.AutocompoundIntervalNumBlocks}, nil
	case types.BridgeAddress:
		return types.NetworkPropertyValue{StrValue: properties.BridgeAddress}, nil
	case types.BridgeCosmosEthereumExchangeRate:
		return types.NetworkPropertyValue{StrValue: properties.BridgeCosmosEthereumExchangeRate.String()}, nil
	case types.BridgeEthereumCosmosExchangeRate:
		return types.NetworkPropertyValue{StrValue: properties.BridgeEthereumCosmosExchangeRate.String()}, nil
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
	case types.MinimumProposalEndTime:
		properties.MinimumProposalEndTime = value.Value
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
		decValue, err := sdk.NewDecFromStr(value.StrValue)
		if err != nil {
			return err
		}
		properties.InactiveRankDecreasePercent = decValue
	case types.PoorNetworkMaxBankSend:
		properties.PoorNetworkMaxBankSend = value.Value
	case types.MinValidators:
		properties.MinValidators = value.Value
	case types.UnjailMaxTime:
		properties.UnjailMaxTime = value.Value
	case types.EnableTokenBlacklist:
		properties.EnableTokenBlacklist = IntToBool(value.Value)
	case types.EnableTokenWhitelist:
		properties.EnableTokenWhitelist = IntToBool(value.Value)
	case types.MinIdentityApprovalTip:
		properties.MinIdentityApprovalTip = value.Value
	case types.UniqueIdentityKeys:
		properties.UniqueIdentityKeys = value.StrValue
	case types.UbiHardcap:
		properties.UbiHardcap = value.Value
	case types.ValidatorsFeeShare:
		decValue, err := sdk.NewDecFromStr(value.StrValue)
		if err != nil {
			return err
		}
		properties.ValidatorsFeeShare = decValue
	case types.InflationRate:
		decValue, err := sdk.NewDecFromStr(value.StrValue)
		if err != nil {
			return err
		}
		properties.InflationRate = decValue
	case types.InflationPeriod:
		properties.InflationPeriod = value.Value
	case types.UnstakingPeriod:
		properties.UnstakingPeriod = value.Value
	case types.MaxDelegators:
		properties.MaxDelegators = value.Value
	case types.MinDelegationPushout:
		properties.MinDelegationPushout = value.Value
	case types.SlashingPeriod:
		properties.SlashingPeriod = value.Value
	case types.MaxJailedPercentage:
		decValue, err := sdk.NewDecFromStr(value.StrValue)
		if err != nil {
			return err
		}
		properties.MaxJailedPercentage = decValue
	case types.MaxSlashingPercentage:
		decValue, err := sdk.NewDecFromStr(value.StrValue)
		if err != nil {
			return err
		}
		properties.MaxSlashingPercentage = decValue
	case types.MinCustodyReward:
		properties.MinCustodyReward = value.Value
	case types.MaxCustodyBufferSize:
		properties.MaxCustodyBufferSize = value.Value
	case types.MaxCustodyTxSize:
		properties.MaxCustodyTxSize = value.Value
	case types.AbstentionRankDecreaseAmount:
		properties.AbstentionRankDecreaseAmount = value.Value
	case types.MaxAbstention:
		properties.MaxAbstention = value.Value
	case types.MinCollectiveBond:
		properties.MinCollectiveBond = value.Value
	case types.MinCollectiveBondingTime:
		properties.MinCollectiveBondingTime = value.Value
	case types.MaxCollectiveOutputs:
		properties.MaxCollectiveOutputs = value.Value
	case types.MinCollectiveClaimPeriod:
		properties.MinCollectiveClaimPeriod = value.Value
	case types.ValidatorRecoveryBond:
		properties.ValidatorRecoveryBond = value.Value
	case types.MaxAnnualInflation:
		decValue, err := sdk.NewDecFromStr(value.StrValue)
		if err != nil {
			return err
		}
		properties.MaxAnnualInflation = decValue
	case types.MinDappBond:
		properties.MinDappBond = value.Value
	case types.MaxDappBond:
		properties.MaxDappBond = value.Value
	case types.DappBondDuration:
		properties.DappBondDuration = value.Value
	case types.DappVerifierBond:
		decValue, err := sdk.NewDecFromStr(value.StrValue)
		if err != nil {
			return err
		}
		properties.DappVerifierBond = decValue
	case types.DappAutoDenounceTime:
		properties.DappAutoDenounceTime = value.Value
	case types.DappMischanceRankDecreaseAmount:
		properties.DappMischanceRankDecreaseAmount = value.Value
	case types.DappMaxMischance:
		properties.DappMaxMischance = value.Value
	case types.DappInactiveRankDecreasePercent:
		properties.DappInactiveRankDecreasePercent = value.Value
	case types.DappPoolSlippageDefault:
		decValue, err := sdk.NewDecFromStr(value.StrValue)
		if err != nil {
			return err
		}
		properties.DappPoolSlippageDefault = decValue
	case types.MintingFtFee:
		properties.MintingFtFee = value.Value
	case types.MintingNftFee:
		properties.MintingNftFee = value.Value
	case types.VetoThreshold:
		decValue, err := sdk.NewDecFromStr(value.StrValue)
		if err != nil {
			return err
		}
		properties.VetoThreshold = decValue
	case types.AutocompoundIntervalNumBlocks:
		properties.AutocompoundIntervalNumBlocks = value.Value
	case types.BridgeAddress:
		properties.BridgeAddress = value.StrValue
	case types.BridgeCosmosEthereumExchangeRate:
		decValue, err := sdk.NewDecFromStr(value.StrValue)
		if err != nil {
			return err
		}
		properties.BridgeCosmosEthereumExchangeRate = decValue
	case types.BridgeEthereumCosmosExchangeRate:
		decValue, err := sdk.NewDecFromStr(value.StrValue)
		if err != nil {
			return err
		}
		properties.BridgeEthereumCosmosExchangeRate = decValue

	default:
		return errors.New("trying to set network property that does not exist")
	}
	return k.SetNetworkProperties(ctx, properties)
}

// SetExecutionFee set fee by execution function name
func (k Keeper) SetExecutionFee(ctx sdk.Context, fee types.ExecutionFee) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixExecutionFee)
	key := []byte(fee.TransactionType)
	prefixStore.Set(key, k.cdc.MustMarshal(&fee))
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

// GetExecutionFees get all execution fees
func (k Keeper) GetExecutionFees(ctx sdk.Context) []types.ExecutionFee {
	iterator := sdk.KVStorePrefixIterator(ctx.KVStore(k.storeKey), types.KeyPrefixExecutionFee)
	defer iterator.Close()
	fees := []types.ExecutionFee{}
	for ; iterator.Valid(); iterator.Next() {
		bz := iterator.Value()
		fee := types.ExecutionFee{}
		k.cdc.MustUnmarshal(bz, &fee)
		fees = append(fees, fee)
	}
	return fees
}
