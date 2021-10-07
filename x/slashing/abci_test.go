package slashing_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	simapp "github.com/KiraCore/sekai/app"
	"github.com/KiraCore/sekai/x/slashing"
	"github.com/KiraCore/sekai/x/staking"
	"github.com/KiraCore/sekai/x/staking/teststaking"
	stakingtypes "github.com/KiraCore/sekai/x/staking/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestBeginBlocker(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	pks := simapp.CreateTestPubKeys(1)
	simapp.AddTestAddrsFromPubKeys(app, ctx, pks, sdk.TokensFromConsensusPower(200, sdk.DefaultPowerReduction))
	addr, pk := sdk.ValAddress(pks[0].Address()), pks[0]
	tstaking := teststaking.NewHelper(t, ctx, app.CustomStakingKeeper, app.CustomGovKeeper)

	// bond the validator
	tstaking.CreateValidator(addr, pk, true)
	staking.EndBlocker(ctx, app.CustomStakingKeeper)

	val := abci.Validator{
		Address: pk.Address(),
	}

	// mark the validator as having signed
	req := abci.RequestBeginBlock{
		LastCommitInfo: abci.LastCommitInfo{
			Votes: []abci.VoteInfo{{
				Validator:       val,
				SignedLastBlock: true,
			}},
		},
	}

	slashing.BeginBlocker(ctx, req, app.CustomSlashingKeeper)

	info, found := app.CustomSlashingKeeper.GetValidatorSigningInfo(ctx, sdk.ConsAddress(pk.Address()))
	require.True(t, found)
	require.Equal(t, sdk.ConsAddress(pk.Address()).String(), info.Address)
	require.Equal(t, int64(0), info.StartHeight)
	require.Equal(t, time.Unix(0, 0).UTC(), info.InactiveUntil.UTC())
	require.Equal(t, int64(0), info.Mischance)
	require.Equal(t, int64(0), info.LastPresentBlock)
	require.Equal(t, int64(0), info.MissedBlocksCounter)
	require.Equal(t, int64(1), info.ProducedBlocksCounter)

	height := int64(0)

	// for 1000 blocks, mark the validator as having signed
	for ; height < 1000; height++ {
		ctx = ctx.WithBlockHeight(height)
		req = abci.RequestBeginBlock{
			LastCommitInfo: abci.LastCommitInfo{
				Votes: []abci.VoteInfo{{
					Validator:       val,
					SignedLastBlock: true,
				}},
			},
		}

		slashing.BeginBlocker(ctx, req, app.CustomSlashingKeeper)
	}

	// for 500 blocks, mark the validator as having not signed
	for ; height < 1500; height++ {
		ctx = ctx.WithBlockHeight(height)
		req = abci.RequestBeginBlock{
			LastCommitInfo: abci.LastCommitInfo{
				Votes: []abci.VoteInfo{{
					Validator:       val,
					SignedLastBlock: false,
				}},
			},
		}

		slashing.BeginBlocker(ctx, req, app.CustomSlashingKeeper)
	}

	// end block
	staking.EndBlocker(ctx, app.CustomStakingKeeper)

	// validator should be jailed
	validator, err := app.CustomStakingKeeper.GetValidatorByConsAddr(ctx, sdk.GetConsAddress(pk))
	require.NoError(t, err)
	require.Equal(t, stakingtypes.Inactive, validator.GetStatus())
}
