package keeper

import (
	"context"
	"crypto/sha256"
	"encoding/hex"

	govtypes "github.com/KiraCore/sekai/x/gov/types"
	"github.com/KiraCore/sekai/x/recovery/types"
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

// allow ANY user to register or modify existing recovery secret & verify if the nonce is correct
func (k msgServer) RegisterRecoverySecret(goCtx context.Context, msg *types.MsgRegisterRecoverySecret) (*types.MsgRegisterRecoverySecretResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

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

// allow ANY KIRA address that knows the recovery secret or has a sufficient number of RR tokens to rotate the address
func (k msgServer) RotateRecoveryAddress(goCtx context.Context, msg *types.MsgRotateRecoveryAddress) (*types.MsgRotateRecoveryAddressResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

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

	// TODO: set rotation history or something

	// - account module
	acc := k.ak.GetAccount(ctx, addr)
	rotatedAcc := k.ak.GetAccount(ctx, rotatedAddr)
	if rotatedAcc != nil {
		return nil, types.ErrRotatedAccountAlreadyExists
	}
	err = acc.SetPubKey(nil)
	if err != nil {
		return nil, err
	}
	err = acc.SetAddress(rotatedAddr)
	if err != nil {
		return nil, err
	}
	k.ak.SetAccount(ctx, acc)

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
			k.gk.RemoveRoleFromActor(ctx, actor, role)
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

	// - gov:votes
	proposals, err := k.gk.GetProposals(ctx)
	if err != nil {
		return nil, err
	}
	for _, p := range proposals {
		vote, found := k.gk.GetVote(ctx, p.ProposalId, addr)
		if found {
			k.gk.DeleteVote(ctx, vote)
			vote.Voter = rotatedAddr
			k.gk.SaveVote(ctx, vote)
		}
	}

	// TODO:
	// - multistaking
	// - slashing
	// - spending
	// - staking
	// - tokens
	// - ubi
	// - custody

	// nothing to do with following modules
	// - basket
	// - distributor
	// - evidence

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Address),
		),
	)

	return &types.MsgRotateRecoveryAddressResponse{}, nil
}

// mint `rr_<moniker>` tokens and deposit them to the validator account.
// This function will require putting up a bond in the amount of `validator_recovery_bond` otherwise should fail
func (k msgServer) IssueRecoveryTokens(goCtx context.Context, msg *types.MsgIssueRecoveryTokens) (*types.MsgIssueRecoveryTokensResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: kex token spend
	// TODO: check if validator and previously not issued token
	k.Keeper.SetRecoveryToken(ctx, types.RecoveryToken{
		Address: msg.Address,
		Token:   msg.Address,
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

// burn tokens and redeem KEX
func (k msgServer) BurnRecoveryTokens(goCtx context.Context, msg *types.MsgBurnRecoveryTokens) (*types.MsgBurnRecoveryTokensResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: kex token recovvery
	k.Keeper.DeleteRecoveryToken(ctx, types.RecoveryToken{
		Address: msg.Address,
		Token:   msg.Address,
	})

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Address),
		),
	)
	return &types.MsgBurnRecoveryTokensResponse{}, nil
}
