package keeper

import (
	"context"

	govtypes "github.com/KiraCore/sekai/x/gov/types"
	"github.com/KiraCore/sekai/x/layer2/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type msgServer struct {
	keeper Keeper
}

// NewMsgServerImpl returns an implementation of the bank MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{
		keeper: keeper,
	}
}

var _ types.MsgServer = msgServer{}

func (k msgServer) CreateDappProposal(goCtx context.Context, msg *types.MsgCreateDappProposal) (*types.MsgCreateDappProposalResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	properties := k.keeper.gk.GetNetworkProperties(ctx)
	addr := sdk.MustAccAddressFromBech32(msg.Sender)

	// permission check PermCreateDappProposalWithoutBond
	isAllowed := k.keeper.CheckIfAllowedPermission(ctx, addr, govtypes.PermCreateDappProposalWithoutBond)
	if !isAllowed {
		minDappBond := properties.MinDappBond
		if msg.Bond.Denom != k.keeper.BondDenom(ctx) {
			return nil, types.ErrInvalidDappBondDenom
		}
		// check 1% of properties.MinDappBond
		if msg.Bond.Amount.Mul(sdk.NewInt(100)).LT(sdk.NewInt(int64(minDappBond)).Mul(sdk.NewInt(1000_000))) {
			return nil, types.ErrLowAmountToCreateDappProposal
		}
	}

	// send initial bond to module account
	if msg.Bond.IsPositive() {
		err := k.keeper.bk.SendCoinsFromAccountToModule(ctx, addr, types.ModuleName, sdk.Coins{msg.Bond})
		if err != nil {
			return nil, err
		}
	}

	dapp := k.keeper.GetDapp(ctx, msg.Dapp.Name)
	if dapp.Name != "" {
		return nil, types.ErrDappAlreadyExists
	}

	// create dapp object
	msg.Dapp.TotalBond = msg.Bond
	msg.Dapp.CreationTime = uint64(ctx.BlockTime().Unix())
	msg.Dapp.Status = types.Bootstrap
	k.keeper.SetDapp(ctx, msg.Dapp)
	k.keeper.SetUserDappBond(ctx, types.UserDappBond{
		DappName: msg.Dapp.Name,
		User:     msg.Sender,
		Bond:     msg.Bond,
	})

	return &types.MsgCreateDappProposalResponse{}, nil
}

func (k msgServer) BondDappProposal(goCtx context.Context, msg *types.MsgBondDappProposal) (*types.MsgBondDappProposalResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	dapp := k.keeper.GetDapp(ctx, msg.DappName)
	if dapp.Name == "" {
		return nil, types.ErrDappDoesNotExist
	}

	if k.keeper.BondDenom(ctx) != msg.Bond.Denom {
		return nil, types.ErrInvalidDappBondDenom
	}

	// send initial bond to module account
	addr := sdk.MustAccAddressFromBech32(msg.Sender)
	err := k.keeper.bk.SendCoinsFromAccountToModule(ctx, addr, types.ModuleName, sdk.Coins{msg.Bond})
	if err != nil {
		return nil, err
	}

	properties := k.keeper.gk.GetNetworkProperties(ctx)
	if dapp.TotalBond.Amount.GTE(sdk.NewInt(int64(properties.MaxDappBond)).Mul(sdk.NewInt(1000_000))) {
		return nil, types.ErrMaxDappBondReached
	}

	dapp.TotalBond = dapp.TotalBond.Add(msg.Bond)
	k.keeper.SetDapp(ctx, dapp)

	userDappBond := k.keeper.GetUserDappBond(ctx, msg.DappName, msg.Sender)
	if userDappBond.User != "" {
		userDappBond.Bond = userDappBond.Bond.Add(msg.Bond)
	} else {
		userDappBond = types.UserDappBond{
			User:     msg.Sender,
			DappName: msg.DappName,
			Bond:     msg.Bond,
		}
	}
	k.keeper.SetUserDappBond(ctx, userDappBond)

	return &types.MsgBondDappProposalResponse{}, nil
}

func (k msgServer) ReclaimDappBondProposal(goCtx context.Context, msg *types.MsgReclaimDappBondProposal) (*types.MsgReclaimDappBondProposalResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	userDappBond := k.keeper.GetUserDappBond(ctx, msg.DappName, msg.Sender)
	if userDappBond.DappName == "" {
		return nil, types.ErrUserDappBondDoesNotExist
	}
	if userDappBond.Bond.Denom != msg.Bond.Denom {
		return nil, types.ErrInvalidDappBondDenom
	}
	if userDappBond.Bond.Amount.LT(msg.Bond.Amount) {
		return nil, types.ErrNotEnoughUserDappBond
	}

	userDappBond.Bond.Amount = userDappBond.Bond.Amount.Sub(msg.Bond.Amount)
	k.keeper.SetUserDappBond(ctx, userDappBond)

	// send tokens back to user
	addr := sdk.MustAccAddressFromBech32(msg.Sender)
	err := k.keeper.bk.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr, sdk.Coins{msg.Bond})
	if err != nil {
		return nil, err
	}

	return &types.MsgReclaimDappBondProposalResponse{}, nil
}

func (k msgServer) ExitDapp(goCtx context.Context, msg *types.MsgExitDapp) (*types.MsgExitDappResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	operator := k.keeper.GetDappOperator(ctx, msg.DappName, msg.Sender)
	if operator.DappName == "" {
		return nil, types.ErrNotDappOperator
	}
	k.keeper.DeleteDappOperator(ctx, msg.DappName, msg.Sender)

	return &types.MsgExitDappResponse{}, nil
}

func (k msgServer) PauseDappTx(goCtx context.Context, msg *types.MsgPauseDappTx) (*types.MsgPauseDappTxResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	operator := k.keeper.GetDappOperator(ctx, msg.DappName, msg.Sender)
	if operator.DappName == "" {
		return nil, types.ErrNotDappOperator
	}
	if operator.Status != types.OperatorActive {
		return nil, types.ErrDappOperatorNotActive
	}
	operator.Status = types.OperatorPaused
	k.keeper.SetDappOperator(ctx, operator)

	return &types.MsgPauseDappTxResponse{}, nil
}

func (k msgServer) UnPauseDappTx(goCtx context.Context, msg *types.MsgUnPauseDappTx) (*types.MsgUnPauseDappTxResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	operator := k.keeper.GetDappOperator(ctx, msg.DappName, msg.Sender)
	if operator.DappName == "" {
		return nil, types.ErrNotDappOperator
	}
	if operator.Status != types.OperatorPaused {
		return nil, types.ErrDappOperatorNotPaused
	}
	operator.Status = types.OperatorActive
	k.keeper.SetDappOperator(ctx, operator)

	return &types.MsgUnPauseDappTxResponse{}, nil
}

func (k msgServer) ReactivateDappTx(goCtx context.Context, msg *types.MsgReactivateDappTx) (*types.MsgReactivateDappTxResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	operator := k.keeper.GetDappOperator(ctx, msg.DappName, msg.Sender)
	if operator.DappName == "" {
		return nil, types.ErrNotDappOperator
	}
	if operator.Status != types.OperatorDeactivatived {
		return nil, types.ErrDappOperatorNotDeactivated
	}
	operator.Status = types.OperatorActive
	k.keeper.SetDappOperator(ctx, operator)

	return &types.MsgReactivateDappTxResponse{}, nil
}

func (k msgServer) ExecuteDappTx(goCtx context.Context, msg *types.MsgExecuteDappTx) (*types.MsgExecuteDappTxResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	session := k.keeper.GetDappSession(ctx, msg.DappName)
	if session.DappName == "" {
		return nil, types.ErrNoDappSessionExists
	}

	if session.Leader != msg.Sender {
		return nil, types.ErrNotDappSessionLeader
	}
	session.Status = types.Ongoing

	return &types.MsgExecuteDappTxResponse{}, nil
}

func (k msgServer) DenounceLeaderTx(goCtx context.Context, msg *types.MsgDenounceLeaderTx) (*types.MsgDenounceLeaderTxResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	k.keeper.SetDappLeaderDenouncement(ctx, types.DappLeaderDenouncement{
		DappName:     msg.DappName,
		Leader:       msg.Leader,
		Sender:       msg.Sender,
		Denouncement: msg.DenounceText,
	})

	return &types.MsgDenounceLeaderTxResponse{}, nil
}

func (k msgServer) TransferDappTx(goCtx context.Context, msg *types.MsgTransferDappTx) (*types.MsgTransferDappTxResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	_ = ctx

	return &types.MsgTransferDappTxResponse{}, nil
}

// TODO:
// ### c**) Team & Investors Incentives**
// Here are a few examples of ways in which “issuance” and “pool” configuration parameters can be used:
// - Fair Launch - no extra tokens issued and all LP coins are immediately unlocked (`pool.drip` set to 0).
// - User Assisted Launch - LP Spending Pool is configured to slowly distribute LP tokens, the `issuance.premint` is set to
// a small reasonable amount while the `issuance.postmint` is not used.
// This enables small teams that need to hire a few developers to establish a token treasury
// and sell their stake to users that are locked in the LP.
// - Investor Assisted Launch - LP Spending Pool is configured to slowly distribute LP tokens
// while premint and postmint enable the creation of treasury and sale of SAFT agreements for large-scale projects.
// The `issuance.time` parameter can be used to clearly define the time when investor tokens will be issued during the “postmint” event
// while the `issuance.deposit` address can be set up by the team as a Spending Pool to easily distribute tokens to their
// rightful owners as well as configure an **optional** “drip” if needed to not scare the LP token holders with an immediate increase of
// the token supply.

// TODO:
// ### d**) Optional Execution & Operators Incentives**

// While raising the launch proposal the deployers must provide `executors_min` and `executors_max` parameters that
// define how many validators are needed to run the dApp. For example in the case of the DEX or a game - one or two validators might be sufficient
// as executors while in the case of the bridge and MPC it would be expected to see at least 21+ validators collaborating on securing BTC or ETH bridge
// address with ECDSA TSS or a multisig. While there is a limitation in terms of how many nodes might want to run the code
// there is no limitation in terms of how many nodes might want to participate in verification (fisherman).
// To make it worth a while for validators to execute the dApp code we will utilize the dApp LP pool to create those incentives.
// Besides the impermanent loss, the LP token holders will incur a default `1% fee` on all swaps,
// deposits, and redemptions configurable in the [Network Properties](https://www.notion.so/de74fe4b731a47df86683f2e9eefa793)
// as `dapp_pool_fee` parameter.
// The fixed fee will be applied after the swap from where `50%` of the corresponding tokens must be **burned** (deminted),
// `25%` given as a reward to liquidity providers and the remaining `25%` will be split between **ACTIVE** dApp executors, and verifiers (fisherman).
// Additionally, the premint and postmint tokens can be used to incentivize operators before dApp starts to generate revenue.

func (k msgServer) RedeemDappPoolTx(goCtx context.Context, msg *types.MsgRedeemDappPoolTx) (*types.MsgRedeemDappPoolTxResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	_ = ctx

	return &types.MsgRedeemDappPoolTxResponse{}, nil
}

func (k msgServer) SwapDappPoolTx(goCtx context.Context, msg *types.MsgSwapDappPoolTx) (*types.MsgSwapDappPoolTxResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	_ = ctx

	dapp := k.keeper.GetDapp(ctx, msg.DappName)
	if dapp.Name != "" {
		return nil, types.ErrDappDoesNotExist
	}

	// TODO: Uniswap v2 like pool implementation
	// dapp.TotalBond
	// dapp.Pool.Ratio
	// dapp.Issurance.Deposit
	// dapp.Issurance.Premint
	// dapp.Issurance.Postmint
	// dapp.Issurance.Time

	// TODO: If the KEX collateral in the pool falls below dapp_liquidation_threshold (by default set to 100’000 KEX)
	// then the dApp will enter a depreciation phase lasting dapp_liquidation_period (by default set to 2419200, that is ~28d)
	// after which the execution will be stopped.
	// On uniswap pool, people will be able to buy some LP tokens?

	return &types.MsgSwapDappPoolTxResponse{}, nil
}

func (k msgServer) ConvertDappPoolTx(goCtx context.Context, msg *types.MsgConvertDappPoolTx) (*types.MsgConvertDappPoolTxResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	_ = ctx

	return &types.MsgConvertDappPoolTxResponse{}, nil
}

func (k msgServer) MintCreateFtTx(goCtx context.Context, msg *types.MsgMintCreateFtTx) (*types.MsgMintCreateFtTxResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	_ = ctx

	return &types.MsgMintCreateFtTxResponse{}, nil
}

func (k msgServer) MintCreateNftTx(goCtx context.Context, msg *types.MsgMintCreateNftTx) (*types.MsgMintCreateNftTxResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	_ = ctx

	return &types.MsgMintCreateNftTxResponse{}, nil
}

func (k msgServer) MintIssueTx(goCtx context.Context, msg *types.MsgMintIssueTx) (*types.MsgMintIssueTxResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	_ = ctx

	return &types.MsgMintIssueTxResponse{}, nil
}

func (k msgServer) MintBurnTx(goCtx context.Context, msg *types.MsgMintBurnTx) (*types.MsgMintBurnTxResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	_ = ctx

	return &types.MsgMintBurnTxResponse{}, nil
}

// TODO: implement - step1
// Until when the Dapp start the session?
// Should it be by governance?
// Or until when, should wait for operators to join?

// TODO: implement - step2
//   rpc RedeemDappPoolTx(MsgRedeemDappPoolTx) returns (MsgRedeemDappPoolTxResponse);
//   rpc SwapDappPoolTx(MsgSwapDappPoolTx) returns (MsgSwapDappPoolTxResponse);
//   rpc ConvertDappPoolTx(MsgConvertDappPoolTx) returns (MsgConvertDappPoolTxResponse);
//   rpc UpsertDappProposalTx(MsgUpsertDappProposalTx) returns (MsgUpsertDappProposalTxResponse);
//   rpc VoteUpsertDappProposalTx(MsgVoteUpsertDappProposalTx) returns (MsgVoteUpsertDappProposalTxResponse);
//   rpc TransferDappTx(MsgTransferDappTx) returns (MsgTransferDappTxResponse);
//   rpc MintCreateFtTx(MsgMintCreateFtTx) returns (MsgMintCreateFtTxResponse);
//   rpc MintCreateNftTx(MsgMintCreateNftTx) returns (MsgMintCreateNftTxResponse);
//   rpc MintIssueTx(MsgMintIssueTx) returns (MsgMintIssueTxResponse);
//   rpc MintBurnTx(MsgMintBurnTx) returns (MsgMintBurnTxResponse);
