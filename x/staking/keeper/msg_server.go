package keeper

import (
	"context"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"

	govkeeper "github.com/KiraCore/sekai/x/gov/keeper"
	customgovtypes "github.com/KiraCore/sekai/x/gov/types"
	"github.com/KiraCore/sekai/x/staking/types"
	"github.com/cosmos/cosmos-sdk/types/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type msgServer struct {
	keeper    Keeper
	govKeeper govkeeper.Keeper
}

// NewMsgServerImpl returns an implementation of the bank MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper, govKeeper govkeeper.Keeper) types.MsgServer {
	return &msgServer{
		keeper:    keeper,
		govKeeper: govKeeper,
	}
}

var _ types.MsgServer = msgServer{}

func (k msgServer) ClaimValidator(goCtx context.Context, msg *types.MsgClaimValidator) (*types.MsgClaimValidatorResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	isAllowed := govkeeper.CheckIfAllowedPermission(ctx, k.govKeeper, sdk.AccAddress(msg.ValKey), customgovtypes.PermClaimValidator)
	if !isAllowed {
		return nil, errors.Wrap(customgovtypes.ErrNotEnoughPermissions, "PermClaimValidator")
	}

	pk, ok := msg.PubKey.GetCachedValue().(cryptotypes.PubKey)
	if !ok {
		return nil, errors.Wrap(errors.ErrInvalidPubKey, "Invalid pub key")
	}

	validator, err := types.NewValidator(msg.Moniker, msg.Website, msg.Social, msg.Identity, msg.Commission, msg.ValKey, pk)
	if err != nil {
		return nil, err
	}

	k.keeper.AddPendingValidator(ctx, validator)

	return &types.MsgClaimValidatorResponse{}, nil
}

func (k msgServer) ProposalUnjailValidator(ctx context.Context, validator *types.MsgProposalUnjailValidator) (*types.MsgProposalUnjailValidatorResponse, error) {
	panic("implement me")
}
