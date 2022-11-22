package keeper

import (
	"context"

	"github.com/KiraCore/sekai/x/collectives/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
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

// CreateCollective defines a method for creating collective.
// allow ANY user to create Staking Collective even if they have no roles or
// permissions enabling that
func (k msgServer) CreateCollective(
	goCtx context.Context,
	msg *types.MsgCreateCollective,
) (*types.MsgCreateCollectiveResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	collective := k.keeper.GetCollective(ctx, msg.Name)
	if collective.Name != "" {
		return nil, types.ErrCollectiveAlreadyExists
	}

	properties := k.keeper.gk.GetNetworkProperties(ctx)
	if len(msg.SpendingPools) > int(properties.MaxCollectiveOutputs) {
		return nil, types.ErrNumberOfSpendingPoolsBiggerThanMaxOutputs
	}

	bondsValue := k.keeper.GetBondsValue(ctx, msg.Bonds)

	// check if initial bond is lower than 10%
	minCollectiveBond := sdk.NewDec(int64(properties.MinCollectiveBond)).Mul(sdk.NewDec(1000_000))
	if bondsValue.LT(minCollectiveBond.Quo(sdk.NewDec(10))) { // MinCollectiveBond is in KEX
		return nil, types.InitialBondLowerThanTenPercentOfMinimumBond
	}

	collectiveStatus := types.CollectiveInactive
	if bondsValue.GTE(minCollectiveBond) {
		collectiveStatus = types.CollectiveActive
	}

	k.keeper.SetCollective(ctx, types.Collective{
		Name:             msg.Name,
		Description:      msg.Description,
		Status:           collectiveStatus,
		DepositWhitelist: msg.DepositWhitelist,
		OwnersWhitelist:  msg.OwnersWhitelist,
		SpendingPools:    msg.SpendingPools,
		ClaimStart:       msg.ClaimStart,
		ClaimPeriod:      msg.ClaimPeriod,
		ClaimEnd:         msg.ClaimEnd,
		VoteQuorum:       msg.VoteQuorum,
		VotePeriod:       msg.VotePeriod,
		VoteEnactment:    msg.VoteEnactment,
		Donations:        sdk.NewCoins(),
		Rewards:          sdk.NewCoins(),
		LastDistribution: uint64(ctx.BlockTime().Unix()),
	})

	// create contribute contributor here
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}
	err = k.keeper.bk.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, msg.Bonds)
	if err != nil {
		return nil, err
	}

	k.keeper.SetCollectiveContributer(ctx, types.CollectiveContributor{
		Address:      msg.Sender,
		Name:         msg.Name,
		Bonds:        msg.Bonds,
		Locking:      0,
		Donation:     sdk.ZeroDec(),
		DonationLock: false,
	})

	return &types.MsgCreateCollectiveResponse{}, nil
}

// ContributeCollective defines a method for putting bonds on collective.
// can be sent by any whitelisted “contributor” account that wants to add
// tokens to the Staking Collective during or after creation process
func (k msgServer) ContributeCollective(
	goCtx context.Context,
	msg *types.MsgBondCollective,
) (*types.MsgBondCollectiveResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}
	err = k.keeper.bk.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, msg.Bonds)
	if err != nil {
		return nil, err
	}

	collective := k.keeper.GetCollective(ctx, msg.Name)
	if collective.Name == "" {
		return nil, types.ErrCollectiveDoesNotExist
	}

	// check if the user is whitelisted user for the collective
	if !collective.DepositWhitelist.Any {
		isWhitelisted := false
		for _, addr := range collective.DepositWhitelist.Accounts {
			if addr == msg.Sender {
				isWhitelisted = true
			}
		}
		actor, found := k.keeper.gk.GetNetworkActorByAddress(ctx, sender)
		if found {
			for _, role := range collective.DepositWhitelist.Roles {
				for _, arole := range actor.Roles {
					if arole == role {
						isWhitelisted = true
					}
				}
			}
		}
		if !isWhitelisted {
			return nil, types.ErrNotWhitelistedForCollectiveDeposit
		}
	}

	cc := k.keeper.GetCollectiveContributer(ctx, msg.Name, msg.Sender)
	if cc.Name != "" {
		cc.Bonds = sdk.NewCoins(cc.Bonds...).Add(msg.Bonds...)
	} else {
		cc = types.CollectiveContributor{
			Address:      msg.Sender,
			Name:         msg.Name,
			Bonds:        msg.Bonds,
			Locking:      0,
			Donation:     sdk.ZeroDec(),
			DonationLock: false,
		}
	}
	k.keeper.SetCollectiveContributer(ctx, cc)
	return &types.MsgBondCollectiveResponse{}, nil
}

// DonateCollective defines a method to set lock and donation for bonds on the
// collection - allows to lock staking derivatives for a specific time period
// and donating a defined percentage of staking rewards to the collective.
func (k msgServer) DonateCollective(
	goCtx context.Context,
	msg *types.MsgDonateCollective,
) (*types.MsgDonateCollectiveResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	cc := k.keeper.GetCollectiveContributer(ctx, msg.Name, msg.Sender)
	if cc.Name != "" {
		return nil, types.ErrNotCollectiveContributer
	}

	// The maximum time period should be **NO greater than 1 year** from the latest block time **OR** the latest time the locking & donation transaction
	// was sent.
	oneYear := 86400 * 365
	if msg.Locking > uint64(ctx.BlockTime().Unix())+uint64(oneYear) {
		return nil, types.ErrLockPeriodCannotExceedOneYear
	}

	// Depositors can always extend the unlock date (indefinitely) but never decrease it.
	if cc.Locking > msg.Locking {
		return nil, types.ErrLockPeriodCanOnlyBeIncreased
	}

	cc.Locking = msg.Locking

	// Depositors should also have the ability to “lock” the donation amount using a dedicated `donation-lock` field until the “locking” period passes.
	// If the locking period is extended the “donation lock” should also persist and remain not changeable.
	if cc.DonationLock && cc.Donation != msg.Donation {
		return nil, types.ErrDonationLocked
	}

	// TODO: All donations should be subtracted from the amounts being sent to the spending pools.

	k.keeper.SetCollectiveContributer(ctx, cc)
	return &types.MsgDonateCollectiveResponse{}, nil
}

// WithdrawCollective can be sent by any whitelisted “contributor” to withdraw
// their tokens (unless locking is enabled)
func (k msgServer) WithdrawCollective(
	goCtx context.Context,
	msg *types.MsgWithdrawCollective,
) (*types.MsgWithdrawCollectiveResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	cc := k.keeper.GetCollectiveContributer(ctx, msg.Name, msg.Sender)
	if cc.Name != "" {
		return nil, types.ErrNotCollectiveContributer
	}

	// After the collective has been created, all whitelisted contributors should be allowed to voluntarily “lock”
	// their staking derivatives for a defined time period by providing the date (UNIX timestamp) at which ALL deposited tokens
	// can be withdrawn. Once the date is set it should NOT be possible to “unlock”
	// the tokens until that specific date passes.
	if cc.Locking > uint64(ctx.BlockTime().Unix()) {
		return nil, types.ErrBondsLockedOnTheCollective
	}

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}
	err = k.keeper.bk.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sender, cc.Bonds)
	if err != nil {
		return nil, err
	}
	k.keeper.DeleteCollectiveContributer(ctx, msg.Name, msg.Sender)
	return &types.MsgWithdrawCollectiveResponse{}, nil
}