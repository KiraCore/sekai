package teststaking

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	customgovkeeper "github.com/KiraCore/sekai/x/gov/keeper"
	govtypes "github.com/KiraCore/sekai/x/gov/types"
	"github.com/KiraCore/sekai/x/staking"
	"github.com/KiraCore/sekai/x/staking/keeper"
	stakingtypes "github.com/KiraCore/sekai/x/staking/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Helper is a structure which wraps the staking handler
// and provides methods useful in tests
type Helper struct {
	t         *testing.T
	h         sdk.Handler
	k         keeper.Keeper
	govKeeper customgovkeeper.Keeper

	Ctx sdk.Context
}

// NewHelper creates staking Handler wrapper for tests
func NewHelper(t *testing.T, ctx sdk.Context, k keeper.Keeper, govKeeper customgovkeeper.Keeper) *Helper {
	return &Helper{t, staking.NewHandler(k, govKeeper), k, govKeeper, ctx}
}

// CreateValidator calls handler to create a new staking validator
func (sh *Helper) CreateValidator(addr sdk.ValAddress, pk cryptotypes.PubKey, ok bool) {
	// create permission whitelisted actor
	actor := govtypes.NewNetworkActor(
		sdk.AccAddress(addr),
		nil,
		1,
		nil,
		govtypes.NewPermissions([]govtypes.PermValue{
			govtypes.PermClaimValidator,
		}, nil),
		1,
	)
	sh.govKeeper.SaveNetworkActor(sh.Ctx, actor)

	// claim validator
	sh.createValidator(addr, pk, ok)
}

// ClaimValidatorMsg returns a message used to create validator in this service.
func (sh *Helper) ClaimValidatorMsg(addr sdk.ValAddress, pk cryptotypes.PubKey) *stakingtypes.MsgClaimValidator {
	msg, err := stakingtypes.NewMsgClaimValidator("moniker", addr, pk)
	require.NoError(sh.t, err)
	return msg
}

func (sh *Helper) createValidator(addr sdk.ValAddress, pk cryptotypes.PubKey, ok bool) {
	msg, err := stakingtypes.NewMsgClaimValidator(
		fmt.Sprintf("%s-%d", "moniker", rand.Intn(100)),
		addr,
		pk,
	)
	require.NoError(sh.t, err)
	sh.Handle(msg, ok)
}

// Handle calls staking handler on a given message
func (sh *Helper) Handle(msg sdk.Msg, ok bool) *sdk.Result {
	res, err := sh.h(sh.Ctx, msg)
	if ok {
		require.NoError(sh.t, err)
		require.NotNil(sh.t, res)
	} else {
		require.Error(sh.t, err)
		require.Nil(sh.t, res)
	}
	return res
}

// CheckValidator asserts that a validor exists and has a given status (if status!="")
// and if has a right inactive flag.
func (sh *Helper) CheckValidator(addr sdk.ValAddress, status stakingtypes.ValidatorStatus) stakingtypes.Validator {
	v, err := sh.k.GetValidator(sh.Ctx, addr)
	require.NoError(sh.t, err)
	require.Equal(sh.t, status, v.Status)
	return v
}

// TurnBlock calls EndBlocker and updates the block time
func (sh *Helper) TurnBlock(newTime time.Time) sdk.Context {
	sh.Ctx = sh.Ctx.WithBlockTime(newTime)
	staking.EndBlocker(sh.Ctx, sh.k)
	return sh.Ctx
}

// TurnBlockTimeDiff calls EndBlocker and updates the block time by adding the
// duration to the current block time
func (sh *Helper) TurnBlockTimeDiff(diff time.Duration) sdk.Context {
	sh.Ctx = sh.Ctx.WithBlockTime(sh.Ctx.BlockHeader().Time.Add(diff))
	staking.EndBlocker(sh.Ctx, sh.k)
	return sh.Ctx
}
