package keeper

import (
	"context"

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

	// TODO: get previous recovery - and check proof if already exists

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

	// TODO: set rotation history or something

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
