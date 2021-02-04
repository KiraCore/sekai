package gov

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/KiraCore/sekai/x/gov/keeper"
	customgovtypes "github.com/KiraCore/sekai/x/gov/types"
)

func NewHandler(ck keeper.Keeper) sdk.Handler {
	msgServer := keeper.NewMsgServerImpl(ck)

	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case *customgovtypes.MsgSetNetworkProperties:
			res, err := msgServer.SetNetworkProperties(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *customgovtypes.MsgSetExecutionFee:
			res, err := msgServer.SetExecutionFee(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		// Permission Related
		case *customgovtypes.MsgWhitelistPermissions:
			res, err := msgServer.WhitelistPermissions(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *customgovtypes.MsgBlacklistPermissions:
			res, err := msgServer.BlacklistPermissions(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		// Councilor Related
		case *customgovtypes.MsgClaimCouncilor:
			res, err := msgServer.ClaimCouncilor(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		// Role Related
		case *customgovtypes.MsgWhitelistRolePermission:
			res, err := msgServer.WhitelistRolePermission(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *customgovtypes.MsgBlacklistRolePermission:
			res, err := msgServer.BlacklistRolePermission(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *customgovtypes.MsgRemoveWhitelistRolePermission:
			res, err := msgServer.RemoveWhitelistRolePermission(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *customgovtypes.MsgRemoveBlacklistRolePermission:
			res, err := msgServer.RemoveBlacklistRolePermission(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *customgovtypes.MsgCreateRole:
			res, err := msgServer.CreateRole(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *customgovtypes.MsgAssignRole:
			res, err := msgServer.AssignRole(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *customgovtypes.MsgRemoveRole:
			res, err := msgServer.RemoveRole(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		// Proposal related
		case *customgovtypes.MsgProposalAssignPermission:
			res, err := msgServer.ProposalAssignPermission(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *customgovtypes.MsgProposalSetNetworkProperty:
			res, err := msgServer.ProposalSetNetworkProperty(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *customgovtypes.MsgProposalUpsertDataRegistry:
			res, err := msgServer.ProposalUpsertDataRegistry(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *customgovtypes.MsgProposalSetPoorNetworkMessages:
			res, err := msgServer.ProposalSetPoorNetworkMsgs(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *customgovtypes.MsgVoteProposal:
			res, err := msgServer.VoteProposal(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		default:
			return nil, errors.Wrapf(errors.ErrUnknownRequest, "unrecognized %s message type: %T", customgovtypes.ModuleName, msg)
		}
	}
}
