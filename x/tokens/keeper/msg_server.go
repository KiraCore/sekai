package keeper

import (
	"context"
	"fmt"
	"strings"

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

func (k msgServer) UpsertTokenAlias(
	goCtx context.Context,
	msg *types.MsgUpsertTokenAlias,
) (*types.MsgUpsertTokenAliasResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	isAllowed := k.cgk.CheckIfAllowedPermission(ctx, msg.Proposer, govtypes.PermUpsertTokenAlias)
	if !isAllowed {
		return nil, errorsmod.Wrap(govtypes.ErrNotEnoughPermissions, govtypes.PermUpsertTokenAlias.String())
	}

	err := k.keeper.UpsertTokenAlias(ctx, *types.NewTokenAlias(
		msg.Symbol,
		msg.Name,
		msg.Icon,
		msg.Decimals,
		msg.Denoms,
		msg.Invalidated,
	))
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeUpsertTokenAlias,
			sdk.NewAttribute(types.AttributeKeyProposer, msg.Proposer.String()),
			sdk.NewAttribute(types.AttributeKeySymbol, msg.Symbol),
			sdk.NewAttribute(types.AttributeKeyName, msg.Name),
			sdk.NewAttribute(types.AttributeKeyIcon, msg.Icon),
			sdk.NewAttribute(types.AttributeKeyDecimals, fmt.Sprintf("%d", msg.Decimals)),
			sdk.NewAttribute(types.AttributeKeyDenoms, strings.Join(msg.Denoms, ",")),
		),
	)
	return &types.MsgUpsertTokenAliasResponse{}, err
}

func (k msgServer) UpsertTokenRate(goCtx context.Context, msg *types.MsgUpsertTokenRate) (*types.MsgUpsertTokenRateResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	err := msg.ValidateBasic()
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	isAllowed := k.cgk.CheckIfAllowedPermission(ctx, msg.Proposer, govtypes.PermUpsertTokenRate)
	if !isAllowed {
		return nil, errorsmod.Wrap(govtypes.ErrNotEnoughPermissions, govtypes.PermUpsertTokenRate.String())
	}

	err = k.keeper.UpsertTokenRate(ctx, *types.NewTokenRate(
		msg.Denom,
		msg.Rate,
		msg.FeePayments,
		msg.StakeCap,
		msg.StakeMin,
		msg.StakeToken,
		msg.Invalidated,
	))

	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeUpsertTokenRate,
			sdk.NewAttribute(types.AttributeKeyProposer, msg.Proposer.String()),
			sdk.NewAttribute(types.AttributeKeyDenom, msg.Denom),
			sdk.NewAttribute(types.AttributeKeyRate, msg.Rate.String()),
			sdk.NewAttribute(types.AttributeKeyFeePayments, fmt.Sprintf("%t", msg.FeePayments)),
		),
	)

	return &types.MsgUpsertTokenRateResponse{}, nil
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
