package keeper

import (
	"context"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	appparams "github.com/KiraCore/sekai/app/params"
	govtypes "github.com/KiraCore/sekai/x/gov/types"
	"github.com/KiraCore/sekai/x/tokens/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

type msgServer struct {
	keeper Keeper
	cgk    types.CustomGovKeeper
}

// NewMsgServerImpl returns an implementation of the bank MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper, cgk types.CustomGovKeeper) types.MsgServer {
	return &msgServer{
		keeper: keeper,
		cgk:    cgk,
	}
}

var _ types.MsgServer = msgServer{}

func (k msgServer) UpsertTokenInfo(goCtx context.Context, msg *types.MsgUpsertTokenInfo) (*types.MsgUpsertTokenInfoResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	err := msg.ValidateBasic()
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	isAllowed := k.cgk.CheckIfAllowedPermission(ctx, msg.Proposer, govtypes.PermUpsertTokenInfo)
	if !isAllowed {
		return nil, errorsmod.Wrap(govtypes.ErrNotEnoughPermissions, govtypes.PermUpsertTokenInfo.String())
	}

	err = k.keeper.UpsertTokenInfo(ctx, *types.NewTokenInfo(
		msg.Denom,
		msg.Rate,
		msg.FeePayments,
		msg.StakeCap,
		msg.StakeMin,
		msg.StakeToken,
		msg.Invalidated,
		msg.Symbol,
		msg.Name,
		msg.Icon,
		msg.Decimals,
	))

	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeUpsertTokenInfo,
			sdk.NewAttribute(types.AttributeKeyProposer, msg.Proposer.String()),
			sdk.NewAttribute(types.AttributeKeyDenom, msg.Denom),
			sdk.NewAttribute(types.AttributeKeyRate, msg.Rate.String()),
			sdk.NewAttribute(types.AttributeKeyFeePayments, fmt.Sprintf("%t", msg.FeePayments)),
		),
	)

	return &types.MsgUpsertTokenInfoResponse{}, nil
}

func (k msgServer) EthereumTx(goCtx context.Context, msg *types.MsgEthereumTx) (*types.MsgEthereumTxResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	err := msg.ValidateBasic()
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	if msg.TxType == "NativeSend" {
		sender, err := sdk.AccAddressFromBech32(msg.Sender)
		if err != nil {
			return nil, err
		}
		recipient := msg.AsTransaction().To()
		value := sdk.NewIntFromBigInt(msg.AsTransaction().Value())
		cutUnit := sdk.NewInt(1000_000_000_000)
		balance := value.Quo(cutUnit)
		amount := sdk.NewCoin(appparams.DefaultDenom, balance)

		err = k.keeper.bankKeeper.SendCoins(ctx, sender, sdk.AccAddress(recipient.Bytes()), sdk.Coins{amount})
		if err != nil {
			return nil, err
		}
	} else {
		return nil, types.ErrUnimplementedTxType
	}

	return &types.MsgEthereumTxResponse{}, nil
}
