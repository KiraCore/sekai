package keeper_test

import (
	"time"

	"github.com/KiraCore/sekai/x/evidence/types"
	"github.com/KiraCore/sekai/x/staking"
	"github.com/KiraCore/sekai/x/staking/teststaking"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *KeeperTestSuite) TestHandleDoubleSign() {
	ctx := suite.ctx.WithIsCheckTx(false).WithBlockHeight(1)
	suite.populateValidators(ctx)

	power := int64(100)
	operatorAddr, val := valAddresses[0], pubkeys[0]
	tstaking := teststaking.NewHelper(suite.T(), ctx, suite.app.CustomStakingKeeper, suite.app.CustomGovKeeper)

	tstaking.CreateValidator(operatorAddr, val, true)

	// execute end-blocker and verify validator attributes
	staking.EndBlocker(ctx, suite.app.CustomStakingKeeper)

	// double sign less than max age
	evidence := &types.Equivocation{
		Height:           0,
		Time:             time.Unix(0, 0),
		Power:            power,
		ConsensusAddress: sdk.ConsAddress(val.Address()).String(),
	}
	suite.app.EvidenceKeeper.HandleEquivocationEvidence(ctx, evidence)

	// should be jailed and tombstoned
	validator, _ := suite.app.CustomStakingKeeper.GetValidator(ctx, operatorAddr)
	suite.True(validator.IsJailed())

	// submit duplicate evidence
	suite.app.EvidenceKeeper.HandleEquivocationEvidence(ctx, evidence)

	// require we cannot unjail
	suite.Error(suite.app.CustomSlashingKeeper.Activate(ctx, operatorAddr))

	// require we be able to unbond now
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
}

func (suite *KeeperTestSuite) TestHandleDoubleSign_TooOld() {
	ctx := suite.ctx.WithIsCheckTx(false).WithBlockHeight(1).WithBlockTime(time.Now())
	suite.populateValidators(ctx)

	power := int64(100)
	operatorAddr, val := valAddresses[0], pubkeys[0]
	tstaking := teststaking.NewHelper(suite.T(), ctx, suite.app.CustomStakingKeeper, suite.app.CustomGovKeeper)

	tstaking.CreateValidator(operatorAddr, val, true)

	// execute end-blocker and verify validator attributes
	staking.EndBlocker(ctx, suite.app.CustomStakingKeeper)

	evidence := &types.Equivocation{
		Height:           0,
		Time:             ctx.BlockTime(),
		Power:            power,
		ConsensusAddress: sdk.ConsAddress(val.Address()).String(),
	}

	cp := suite.app.BaseApp.GetConsensusParams(ctx)

	ctx = ctx.WithConsensusParams(cp)
	ctx = ctx.WithBlockTime(ctx.BlockTime().Add(cp.Evidence.MaxAgeDuration + 1))
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + cp.Evidence.MaxAgeNumBlocks + 1)
	suite.app.EvidenceKeeper.HandleEquivocationEvidence(ctx, evidence)

	validator, _ := suite.app.CustomStakingKeeper.GetValidator(ctx, operatorAddr)
	suite.False(validator.IsJailed())
}
