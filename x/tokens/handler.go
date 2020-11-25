package tokens

import (
	"time"

	keeper2 "github.com/KiraCore/sekai/x/gov/keeper"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"

	customgovtypes "github.com/KiraCore/sekai/x/gov/types"
	"github.com/KiraCore/sekai/x/tokens/keeper"
	"github.com/KiraCore/sekai/x/tokens/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewHandler returns new instance of handler
func NewHandler(ck keeper.Keeper, cgk types.CustomGovKeeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case *types.MsgUpsertTokenAlias:
			return handleUpsertTokenAlias(ctx, ck, cgk, msg)
		case *types.MsgUpsertTokenRate:
			return handleUpsertTokenRate(ctx, ck, cgk, msg)

		// Proposals
		case *types.MsgProposalUpsertTokenAlias:
			return handleProposalUpsertTokenAlias(ctx, ck, cgk, msg)
		case *types.MsgProposalUpsertTokenRates:
			return handleProposalUpsertTokenRates(ctx, ck, cgk, msg)
		default:
			return nil, errors.Wrapf(errors.ErrUnknownRequest, "unrecognized %s message type: %T", types.ModuleName, msg)
		}
	}
}

func handleProposalUpsertTokenRates(ctx sdk.Context, ck keeper.Keeper, cgk types.CustomGovKeeper, msg *types.MsgProposalUpsertTokenRates) (*sdk.Result, error) {
	isAllowed := cgk.CheckIfAllowedPermission(ctx, msg.Proposer, customgovtypes.PermCreateUpsertTokenRateProposal)
	if !isAllowed {
		return nil, errors.Wrap(customgovtypes.ErrNotEnoughPermissions, customgovtypes.PermCreateUpsertTokenRateProposal.String())
	}

	return &sdk.Result{}, nil
}

func handleProposalUpsertTokenAlias(ctx sdk.Context, ck keeper.Keeper, cgk types.CustomGovKeeper, msg *types.MsgProposalUpsertTokenAlias) (*sdk.Result, error) {
	isAllowed := cgk.CheckIfAllowedPermission(ctx, msg.Proposer, customgovtypes.PermCreateUpsertTokenAliasProposal)
	if !isAllowed {
		return nil, errors.Wrap(customgovtypes.ErrNotEnoughPermissions, customgovtypes.PermCreateUpsertTokenAliasProposal.String())
	}

	return CreateAndSaveProposalWithContent(ctx, cgk, types.NewProposalUpsertTokenAlias(
		msg.Symbol,
		msg.Name,
		msg.Icon,
		msg.Decimals,
		msg.Denoms,
	))
}

func handleUpsertTokenAlias(
	ctx sdk.Context,
	ck keeper.Keeper,
	cgk types.CustomGovKeeper,
	msg *types.MsgUpsertTokenAlias,
) (*sdk.Result, error) {
	isAllowed := cgk.CheckIfAllowedPermission(ctx, msg.Proposer, customgovtypes.PermUpsertTokenAlias)
	if !isAllowed {
		return nil, errors.Wrap(customgovtypes.ErrNotEnoughPermissions, "PermUpsertTokenAlias")
	}

	err := ck.UpsertTokenAlias(ctx, *types.NewTokenAlias(
		msg.Symbol,
		msg.Name,
		msg.Icon,
		msg.Decimals,
		msg.Denoms,
	))
	return &sdk.Result{}, err
}

func handleUpsertTokenRate(ctx sdk.Context, ck keeper.Keeper, cgk types.CustomGovKeeper, msg *types.MsgUpsertTokenRate) (*sdk.Result, error) {
	err := msg.ValidateBasic()
	if err != nil {
		return nil, errors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	isAllowed := cgk.CheckIfAllowedPermission(ctx, msg.Proposer, customgovtypes.PermUpsertTokenRate)
	if !isAllowed {
		return nil, errors.Wrap(customgovtypes.ErrNotEnoughPermissions, "PermUpsertTokenRate")
	}

	err = ck.UpsertTokenRate(ctx, *types.NewTokenRate(
		msg.Denom,
		msg.Rate,
		msg.FeePayments,
	))

	if err != nil {
		return nil, errors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}
	return &sdk.Result{}, nil
}

func CreateAndSaveProposalWithContent(ctx sdk.Context, ck types.CustomGovKeeper, content customgovtypes.Content) (*sdk.Result, error) {
	blockTime := ctx.BlockTime()
	proposalID, err := ck.GetNextProposalID(ctx)
	if err != nil {
		return nil, err
	}

	properties := ck.GetNetworkProperties(ctx)

	proposal, err := customgovtypes.NewProposal(
		proposalID,
		content,
		blockTime,
		blockTime.Add(time.Minute*time.Duration(properties.ProposalEndTime)),
		blockTime.Add(time.Minute*time.Duration(properties.ProposalEnactmentTime)),
	)

	ck.SaveProposal(ctx, proposal)
	ck.AddToActiveProposals(ctx, proposal)

	return &sdk.Result{
		Data: keeper2.ProposalIDToBytes(proposalID),
	}, nil
}
