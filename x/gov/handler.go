package gov

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/KiraCore/sekai/x/gov/keeper"
	"github.com/KiraCore/sekai/x/gov/types"
)

func NewHandler(ck keeper.Keeper) sdk.Handler {
	msgServer := keeper.NewMsgServerImpl(ck)

	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case *types.MsgSetNetworkProperties:
			res, err := msgServer.SetNetworkProperties(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgSetExecutionFee:
			res, err := msgServer.SetExecutionFee(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		// Permission Related
		case *types.MsgWhitelistPermissions:
			res, err := msgServer.WhitelistPermissions(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgBlacklistPermissions:
			res, err := msgServer.BlacklistPermissions(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		// Councilor Related
		case *types.MsgClaimCouncilor:
			res, err := msgServer.ClaimCouncilor(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		// Role Related
		case *types.MsgWhitelistRolePermission:
			res, err := msgServer.WhitelistRolePermission(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgBlacklistRolePermission:
			res, err := msgServer.BlacklistRolePermission(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgRemoveWhitelistRolePermission:
			res, err := msgServer.RemoveWhitelistRolePermission(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgRemoveBlacklistRolePermission:
			res, err := msgServer.RemoveBlacklistRolePermission(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgCreateRole:
			res, err := msgServer.CreateRole(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgAssignRole:
			res, err := msgServer.AssignRole(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgRemoveRole:
			res, err := msgServer.RemoveRole(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		// Proposal related
		case *types.MsgProposalAssignPermission:
			res, err := msgServer.ProposalAssignPermission(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgProposalSetNetworkProperty:
			res, err := msgServer.ProposalSetNetworkProperty(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgProposalUpsertDataRegistry:
			res, err := msgServer.ProposalUpsertDataRegistry(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgProposalSetPoorNetworkMessages:
			res, err := msgServer.ProposalSetPoorNetworkMessages(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgVoteProposal:
			res, err := msgServer.VoteProposal(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgProposalCreateRole:
			res, err := msgServer.ProposalCreateRole(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		// identity registrar related
		case *types.MsgCreateIdentityRecord:
			res, err := msgServer.CreateIdentityRecord(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgEditIdentityRecord:
			res, err := msgServer.EditIdentityRecord(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgRequestIdentityRecordsVerify:
			res, err := msgServer.RequestIdentityRecordsVerify(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgApproveIdentityRecords:
			res, err := msgServer.ApproveIdentityRecords(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgCancelIdentityRecordsVerifyRequest:
			res, err := msgServer.CancelIdentityRecordsVerifyRequest(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		default:
			return nil, errors.Wrapf(errors.ErrUnknownRequest, "unrecognized %s message type: %T", types.ModuleName, msg)
		}
	}
}
