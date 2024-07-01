package keeper

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"

	appparams "github.com/KiraCore/sekai/app/params"
	custodytypes "github.com/KiraCore/sekai/x/custody/types"
	govtypes "github.com/KiraCore/sekai/x/gov/types"
	"github.com/KiraCore/sekai/x/recovery/types"
	slashingtypes "github.com/KiraCore/sekai/x/slashing/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the recovery MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

var RecoveryFee = sdk.Coins{sdk.NewInt64Coin(appparams.DefaultDenom, 1000_000_000)}

// allow ANY user to register or modify existing recovery secret & verify if the nonce is correct
func (k msgServer) RegisterRecoverySecret(goCtx context.Context, msg *types.MsgRegisterRecoverySecret) (*types.MsgRegisterRecoverySecretResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// check if validator recovery token exists
	_, err := k.GetRecoveryToken(ctx, msg.Address)
	if err == nil {
		return nil, types.ErrAddressHasValidatorRecoveryToken
	}

	// check previous recovery and check proof if already exists
	oldRecord, err := k.Keeper.GetRecoveryRecord(ctx, msg.Address)
	if err == nil { // recovery record already exists
		bz, err := hex.DecodeString(msg.Proof)
		if err != nil {
			return nil, err
		}
		hash := sha256.Sum256(bz)
		if hex.EncodeToString(hash[:]) != oldRecord.Challenge {
			return nil, types.ErrInvalidProof
		}
	}

	k.Keeper.SetRecoveryRecord(ctx, types.RecoveryRecord{
		Address:   msg.Address,
		Challenge: msg.Challenge,
		Nonce:     msg.Nonce,
	})

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Address),
		),
	)

	return &types.MsgRegisterRecoverySecretResponse{}, nil
}

func (k msgServer) RotateValidatorByHalfRRTokenHolder(goCtx context.Context, msg *types.MsgRotateValidatorByHalfRRTokenHolder) (*types.MsgRotateValidatorByHalfRRTokenHolderResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// check if validator recovery token exists
	recoveryToken, err := k.GetRecoveryToken(ctx, msg.Address)
	if err != nil {
		return nil, types.ErrRecoveryTokenDoesNotExist
	}

	// check rr token amount
	rrHolder := sdk.MustAccAddressFromBech32(msg.RrHolder)
	balances := k.bk.GetAllBalances(ctx, rrHolder)
	rrAmount := balances.AmountOf(recoveryToken.Token)
	supply := k.bk.GetSupply(ctx, recoveryToken.Token)
	if rrAmount.Mul(sdk.NewInt(2)).LT(supply.Amount) {
		return nil, types.ErrNotEnoughRRTokenAmountForRotation
	}

	rotation := k.GetRotationHistory(ctx, msg.Recovery)
	if rotation.Rotated != "" {
		return nil, types.ErrTargetAddressAlreadyHasRotationHistory
	}

	// set rotation history
	k.SetRotationHistory(ctx, types.Rotation{
		Address: msg.Address,
		Rotated: msg.Recovery,
	})

	addr := sdk.MustAccAddressFromBech32(msg.Address)
	rotatedAddr := sdk.MustAccAddressFromBech32(msg.Recovery)

	// - recovery module
	k.DeleteRecoveryToken(ctx, recoveryToken)
	recoveryToken.Address = msg.Recovery
	k.SetRecoveryToken(ctx, recoveryToken)

	// - multistaking
	info := k.msk.GetCompoundInfoByAddress(ctx, msg.Address)
	k.msk.RemoveCompoundInfo(ctx, info)
	info.Delegator = msg.Recovery
	k.msk.SetCompoundInfo(ctx, info)

	pools := k.msk.GetAllStakingPools(ctx)
	for _, pool := range pools {
		isDelegator := k.msk.IsPoolDelegator(ctx, pool.Id, addr)
		if isDelegator {
			k.msk.RemovePoolDelegator(ctx, pool.Id, addr)
			k.msk.SetPoolDelegator(ctx, pool.Id, rotatedAddr)
		}
	}

	rewards := k.msk.GetDelegatorRewards(ctx, addr)
	if !rewards.IsZero() {
		k.msk.RemoveDelegatorRewards(ctx, addr)
		k.msk.SetDelegatorRewards(ctx, rotatedAddr, rewards)
	}

	stpool, found := k.msk.GetStakingPoolByValidator(ctx, sdk.ValAddress(addr).String())
	if found {
		k.msk.RemoveStakingPool(ctx, stpool)
		stpool.Validator = sdk.ValAddress(rotatedAddr).String()
		k.msk.SetStakingPool(ctx, stpool)
	}

	// - staking
	validator, err := k.sk.GetValidator(ctx, sdk.ValAddress(addr))
	if err == nil {
		k.sk.RemoveValidator(ctx, validator)
		validator.ValKey = sdk.ValAddress(rotatedAddr)
		k.sk.AddValidator(ctx, validator)
	}

	// - gov:councilor
	councilor, found := k.gk.GetCouncilor(ctx, addr)
	if found {
		k.gk.DeleteCouncilor(ctx, councilor)
		councilor.Address = rotatedAddr
		k.gk.SaveCouncilor(ctx, councilor)
	}

	// - gov:identity_records
	records := k.gk.GetIdRecordsByAddress(ctx, addr)
	for _, record := range records {
		k.gk.DeleteIdentityRecordById(ctx, record.Id)
		record.Address = msg.Recovery
		k.gk.SetIdentityRecord(ctx, record)
	}

	requests := k.gk.GetIdRecordsVerifyRequestsByRequester(ctx, addr)
	for _, req := range requests {
		k.gk.DeleteIdRecordsVerifyRequest(ctx, req.Id)
		req.Address = msg.Recovery
		k.gk.SetIdentityRecordsVerifyRequest(ctx, req)
	}

	requests = k.gk.GetIdRecordsVerifyRequestsByApprover(ctx, addr)
	for _, req := range requests {
		k.gk.DeleteIdRecordsVerifyRequest(ctx, req.Id)
		req.Verifier = msg.Recovery
		k.gk.SetIdentityRecordsVerifyRequest(ctx, req)
	}

	// - gov:network_actor
	actor, found := k.gk.GetNetworkActorByAddress(ctx, addr)
	if found {
		k.gk.DeleteNetworkActor(ctx, actor)
		for _, role := range actor.Roles {
			k.gk.UnassignRoleFromActor(ctx, actor, role)
		}
		for _, perm := range actor.Permissions.Whitelist {
			k.gk.DeleteWhitelistAddressPermKey(ctx, actor, govtypes.PermValue(perm))
		}

		actor.Address = rotatedAddr
		k.gk.SaveNetworkActor(ctx, actor)
		for _, role := range actor.Roles {
			k.gk.AssignRoleToActor(ctx, actor, role)
		}
		for _, perm := range actor.Permissions.Whitelist {
			k.gk.SetWhitelistAddressPermKey(ctx, actor, govtypes.PermValue(perm))
		}
	}

	// - gov:proposals
	proposals, err := k.gk.GetProposals(ctx)
	if err != nil {
		return nil, err
	}

	for _, p := range proposals {
		content, ok := p.GetContent().(*slashingtypes.ProposalSlashValidator)
		if ok {
			if content.Offender == sdk.ValAddress(addr).String() {
				content.Offender = sdk.ValAddress(rotatedAddr).String()
				any, err := codectypes.NewAnyWithValue(msg)
				if err != nil {
					return nil, err
				}

				p.Content = any
				k.gk.SaveProposal(ctx, p)
			}
		}
	}

	// - gov:votes
	for _, p := range proposals {
		vote, found := k.gk.GetVote(ctx, p.ProposalId, addr)
		if found {
			k.gk.DeleteVote(ctx, vote)
			vote.Voter = rotatedAddr
			k.gk.SaveVote(ctx, vote)
		}
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Address),
		),
	)

	return &types.MsgRotateValidatorByHalfRRTokenHolderResponse{}, nil
}

// allow ANY KIRA address that knows the recovery secret or has a sufficient number of RR tokens to rotate the address
func (k msgServer) RotateRecoveryAddress(goCtx context.Context, msg *types.MsgRotateRecoveryAddress) (*types.MsgRotateRecoveryAddressResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// check if validator recovery token exists
	_, err := k.GetRecoveryToken(ctx, msg.Address)
	if err == nil {
		return nil, types.ErrAddressHasValidatorRecoveryToken
	}

	// Pay 1000 KEX
	feePayer := sdk.MustAccAddressFromBech32(msg.FeePayer)
	err = k.bk.SendCoinsFromAccountToModule(ctx, feePayer, types.ModuleName, RecoveryFee)
	if err != nil {
		return nil, err
	}

	record, err := k.Keeper.GetRecoveryRecord(ctx, msg.Address)
	if err != nil {
		return nil, err
	}

	bz, err := hex.DecodeString(msg.Proof)
	if err != nil {
		return nil, err
	}
	hash := sha256.Sum256(bz)
	if hex.EncodeToString(hash[:]) != record.Challenge {
		return nil, types.ErrInvalidProof
	}

	addr := sdk.MustAccAddressFromBech32(msg.Address)
	rotatedAddr := sdk.MustAccAddressFromBech32(msg.Recovery)

	rotation := k.GetRotationHistory(ctx, msg.Recovery)
	if rotation.Rotated != "" {
		return nil, types.ErrTargetAddressAlreadyHasRotationHistory
	}

	// set rotation history
	k.SetRotationHistory(ctx, types.Rotation{
		Address: msg.Address,
		Rotated: msg.Recovery,
	})

	// - account module
	acc := k.ak.GetAccount(ctx, addr)
	if acc == nil {
		return nil, types.ErrAccountDoesNotExists
	}
	rotatedAcc := k.ak.GetAccount(ctx, rotatedAddr)
	if rotatedAcc != nil {
		return nil, types.ErrRotatedAccountAlreadyExists
	}

	// - bank module
	balances := k.bk.GetAllBalances(ctx, addr)
	if balances.IsAllPositive() {
		err = k.bk.SendCoins(ctx, addr, rotatedAddr, balances)
		if err != nil {
			return nil, err
		}
	}

	// - collectives module
	contributers := k.ck.GetAllCollectiveContributers(ctx)
	for _, cc := range contributers {
		if cc.Address == msg.Address {
			k.ck.DeleteCollectiveContributer(ctx, cc.Name, cc.Address)
			cc.Address = msg.Recovery
			k.ck.SetCollectiveContributer(ctx, cc)
		}
	}

	// - gov:councilor
	councilor, found := k.gk.GetCouncilor(ctx, addr)
	if found {
		k.gk.DeleteCouncilor(ctx, councilor)
		councilor.Address = rotatedAddr
		k.gk.SaveCouncilor(ctx, councilor)
	}

	// - gov:identity_records
	records := k.gk.GetIdRecordsByAddress(ctx, addr)
	for _, record := range records {
		k.gk.DeleteIdentityRecordById(ctx, record.Id)
		record.Address = msg.Recovery
		k.gk.SetIdentityRecord(ctx, record)
	}

	requests := k.gk.GetIdRecordsVerifyRequestsByRequester(ctx, addr)
	for _, req := range requests {
		k.gk.DeleteIdRecordsVerifyRequest(ctx, req.Id)
		req.Address = msg.Recovery
		k.gk.SetIdentityRecordsVerifyRequest(ctx, req)
	}

	requests = k.gk.GetIdRecordsVerifyRequestsByApprover(ctx, addr)
	for _, req := range requests {
		k.gk.DeleteIdRecordsVerifyRequest(ctx, req.Id)
		req.Verifier = msg.Recovery
		k.gk.SetIdentityRecordsVerifyRequest(ctx, req)
	}

	// - gov:network_actor
	actor, found := k.gk.GetNetworkActorByAddress(ctx, addr)
	if found {
		k.gk.DeleteNetworkActor(ctx, actor)
		for _, role := range actor.Roles {
			k.gk.UnassignRoleFromActor(ctx, actor, role)
		}
		for _, perm := range actor.Permissions.Whitelist {
			k.gk.DeleteWhitelistAddressPermKey(ctx, actor, govtypes.PermValue(perm))
		}

		actor.Address = rotatedAddr
		k.gk.SaveNetworkActor(ctx, actor)
		for _, role := range actor.Roles {
			k.gk.AssignRoleToActor(ctx, actor, role)
		}
		for _, perm := range actor.Permissions.Whitelist {
			k.gk.SetWhitelistAddressPermKey(ctx, actor, govtypes.PermValue(perm))
		}
	}

	// - gov:proposals
	proposals, err := k.gk.GetProposals(ctx)
	if err != nil {
		return nil, err
	}

	for _, p := range proposals {
		content, ok := p.GetContent().(*slashingtypes.ProposalSlashValidator)
		if ok {
			if content.Offender == sdk.ValAddress(addr).String() {
				content.Offender = sdk.ValAddress(rotatedAddr).String()
				any, err := codectypes.NewAnyWithValue(msg)
				if err != nil {
					return nil, err
				}

				p.Content = any
				k.gk.SaveProposal(ctx, p)
			}
		}
	}

	// - gov:votes
	for _, p := range proposals {
		vote, found := k.gk.GetVote(ctx, p.ProposalId, addr)
		if found {
			k.gk.DeleteVote(ctx, vote)
			vote.Voter = rotatedAddr
			k.gk.SaveVote(ctx, vote)
		}
	}

	// - multistaking
	info := k.msk.GetCompoundInfoByAddress(ctx, msg.Address)
	k.msk.RemoveCompoundInfo(ctx, info)
	info.Delegator = msg.Recovery
	k.msk.SetCompoundInfo(ctx, info)

	pools := k.msk.GetAllStakingPools(ctx)
	for _, pool := range pools {
		isDelegator := k.msk.IsPoolDelegator(ctx, pool.Id, addr)
		if isDelegator {
			k.msk.RemovePoolDelegator(ctx, pool.Id, addr)
			k.msk.SetPoolDelegator(ctx, pool.Id, rotatedAddr)
		}
	}

	rewards := k.msk.GetDelegatorRewards(ctx, addr)
	if !rewards.IsZero() {
		k.msk.RemoveDelegatorRewards(ctx, addr)
		k.msk.SetDelegatorRewards(ctx, rotatedAddr, rewards)
	}

	stpool, found := k.msk.GetStakingPoolByValidator(ctx, sdk.ValAddress(addr).String())
	if found {
		k.msk.RemoveStakingPool(ctx, stpool)
		stpool.Validator = sdk.ValAddress(rotatedAddr).String()
		k.msk.SetStakingPool(ctx, stpool)
	}

	// - spending
	sppools := k.spk.GetAllSpendingPools(ctx)
	for _, pool := range sppools {
		info := k.spk.GetClaimInfo(ctx, pool.Name, addr)
		if info != nil {
			k.spk.RemoveClaimInfo(ctx, *info)
			info.Account = msg.Recovery
			k.spk.SetClaimInfo(ctx, *info)
		}
	}

	// - staking
	validator, err := k.sk.GetValidator(ctx, sdk.ValAddress(addr))
	if err == nil {
		k.sk.RemoveValidator(ctx, validator)
		validator.ValKey = sdk.ValAddress(rotatedAddr)
		k.sk.AddValidator(ctx, validator)
	}

	// - custody
	settings := k.custodyk.GetCustodyInfoByAddress(ctx, addr)
	if settings != nil {
		k.custodyk.DropCustodyRecord(ctx, addr)
		k.custodyk.SetCustodyRecord(ctx, custodytypes.CustodyRecord{
			Address:         rotatedAddr,
			CustodySettings: settings,
		})
	}

	custodians := k.custodyk.GetCustodyCustodiansByAddress(ctx, addr)
	if custodians != nil {
		k.custodyk.DropCustodyCustodiansByAddress(ctx, addr)
		k.custodyk.AddToCustodyCustodians(ctx, custodytypes.CustodyCustodiansRecord{
			Address:           rotatedAddr,
			CustodyCustodians: custodians,
		})
	}

	whitelist := k.custodyk.GetCustodyWhiteListByAddress(ctx, addr)
	if whitelist != nil {
		k.custodyk.DropCustodyWhiteListByAddress(ctx, addr)
		k.custodyk.AddToCustodyWhiteList(ctx, custodytypes.CustodyWhiteListRecord{
			Address:          rotatedAddr,
			CustodyWhiteList: whitelist,
		})
	}

	limits := k.custodyk.GetCustodyLimitsByAddress(ctx, addr)
	if limits != nil {
		k.custodyk.DropCustodyLimitsByAddress(ctx, addr)
		k.custodyk.AddToCustodyLimits(ctx, custodytypes.CustodyLimitRecord{
			Address:       rotatedAddr,
			CustodyLimits: limits,
		})
	}

	limitsStatus := k.custodyk.GetCustodyLimitsStatusByAddress(ctx, addr)
	if limitsStatus != nil {
		k.custodyk.DropCustodyLimitsStatus(ctx, addr)
		k.custodyk.AddToCustodyLimitsStatus(ctx, custodytypes.CustodyLimitStatusRecord{
			Address:         rotatedAddr,
			CustodyStatuses: limitsStatus,
		})
	}

	txPool := k.custodyk.GetCustodyPoolByAddress(ctx, addr)
	if txPool != nil {
		k.custodyk.DropCustodyPool(ctx, addr)
		k.custodyk.AddToCustodyPool(ctx, custodytypes.CustodyPool{
			Address:      rotatedAddr,
			Transactions: txPool,
		})
	}

	// nothing to do with following modules
	// - basket
	// - distributor
	// - evidence
	// - slashing
	// - tokens
	// - ubi

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Address),
		),
	)

	return &types.MsgRotateRecoveryAddressResponse{}, nil
}

// mint `rr/<moniker>` tokens and deposit them to the validator account.
// This function will require putting up a bond in the amount of `validator_recovery_bond` otherwise should fail
func (k msgServer) IssueRecoveryTokens(goCtx context.Context, msg *types.MsgIssueRecoveryTokens) (*types.MsgIssueRecoveryTokensResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	addr := sdk.MustAccAddressFromBech32(msg.Address)

	// check if validator and previously not issued token
	_, err := k.Keeper.GetRecoveryToken(ctx, msg.Address)
	if err == nil {
		return nil, types.ErrRecoveryTokenAlreadyExists
	}

	// KEX token spend
	properties := k.gk.GetNetworkProperties(ctx)
	amount := sdk.NewInt(int64(properties.ValidatorRecoveryBond)).Mul(sdk.NewInt(1000_000))
	coins := sdk.NewCoins(sdk.NewCoin(appparams.DefaultDenom, amount))
	err = k.bk.SendCoinsFromAccountToModule(ctx, addr, types.ModuleName, coins)
	if err != nil {
		return nil, err
	}

	records, err := k.gk.GetIdRecordsByAddressAndKeys(ctx, addr, []string{"moniker"})
	if err != nil {
		return nil, err
	}
	if len(records) != 1 {
		return nil, types.ErrInvalidMoniker
	}

	denom := fmt.Sprintf("rr/%s", strings.ToLower(records[0].Value))

	// issue 10'000'000 tokens
	recoveryTokenAmount := sdk.NewInt(10_000_000).Mul(sdk.NewInt(1000_000))
	recoveryCoins := sdk.NewCoins(sdk.NewCoin(denom, recoveryTokenAmount))
	err = k.tk.MintCoins(ctx, types.ModuleName, recoveryCoins)
	if err != nil {
		return nil, err
	}

	err = k.bk.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr, recoveryCoins)
	if err != nil {
		return nil, err
	}

	k.Keeper.SetRecoveryToken(ctx, types.RecoveryToken{
		Address:          msg.Address,
		Token:            denom,
		RrSupply:         recoveryTokenAmount,
		UnderlyingTokens: coins,
	})

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Address),
		),
	)

	return &types.MsgIssueRecoveryTokensResponse{}, nil
}

// burn tokens and redeem underlying tokens
func (k msgServer) BurnRecoveryTokens(goCtx context.Context, msg *types.MsgBurnRecoveryTokens) (*types.MsgBurnRecoveryTokensResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	addr := sdk.MustAccAddressFromBech32(msg.Address)
	recoveryToken, err := k.GetRecoveryTokenByDenom(ctx, msg.RrCoin.Denom)
	if err != nil {
		return nil, err
	}

	redeemAmount := sdk.Coins{}
	for _, coin := range recoveryToken.UnderlyingTokens {
		amount := coin.Amount.Mul(msg.RrCoin.Amount).Quo(recoveryToken.RrSupply)
		if amount.IsPositive() {
			redeemAmount = redeemAmount.Add(sdk.NewCoin(coin.Denom, amount))
		}
	}

	if !redeemAmount.IsZero() {
		err = k.bk.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr, redeemAmount)
		if err != nil {
			return nil, err
		}
	}

	err = k.bk.SendCoinsFromAccountToModule(ctx, addr, types.ModuleName, sdk.NewCoins(msg.RrCoin))
	if err != nil {
		return nil, err
	}

	err = k.tk.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(msg.RrCoin))
	if err != nil {
		return nil, err
	}

	recoveryToken.RrSupply = recoveryToken.RrSupply.Sub(msg.RrCoin.Amount)
	recoveryToken.UnderlyingTokens = sdk.Coins(recoveryToken.UnderlyingTokens).Sub(redeemAmount...)

	if recoveryToken.RrSupply.IsZero() {
		k.Keeper.DeleteRecoveryToken(ctx, recoveryToken)
	} else {
		k.Keeper.SetRecoveryToken(ctx, recoveryToken)
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Address),
		),
	)
	return &types.MsgBurnRecoveryTokensResponse{}, nil
}

// claim RR token holder rewards
func (k msgServer) ClaimRRHolderRewards(goCtx context.Context, msg *types.MsgClaimRRHolderRewards) (*types.MsgClaimRRHolderRewardsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	addr := sdk.MustAccAddressFromBech32(msg.Sender)
	k.ClaimRewards(ctx, addr)
	return &types.MsgClaimRRHolderRewardsResponse{}, nil
}

// register RR token holder
func (k msgServer) RegisterRRTokenHolder(goCtx context.Context, msg *types.MsgRegisterRRTokenHolder) (*types.MsgRegisterRRTokenHolderResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	addr := sdk.MustAccAddressFromBech32(msg.Holder)
	k.Keeper.RegisterRRTokenHolder(ctx, addr)
	return &types.MsgRegisterRRTokenHolderResponse{}, nil
}
