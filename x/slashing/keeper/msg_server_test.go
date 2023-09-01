package keeper_test

import (
	"testing"
	"time"

	simapp "github.com/KiraCore/sekai/app"
	"github.com/KiraCore/sekai/x/slashing/keeper"
	"github.com/KiraCore/sekai/x/slashing/types"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestMsgServerRefuteSlashingProposal(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	addr := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
	valAddr := sdk.ValAddress(addr)

	proposalID, err := app.CustomGovKeeper.CreateAndSaveProposalWithContent(ctx, "title", "description", &types.ProposalSlashValidator{
		Offender:         valAddr.String(),
		StakingPoolId:    1,
		MisbehaviourTime: time.Time{},
		MisbehaviourType: "DoubleSign",
		JailPercentage:   10,
		Colluders:        []string{},
		Refutation:       "",
	})
	require.NoError(t, err)

	msgServer := keeper.NewMsgServerImpl(app.CustomSlashingKeeper)

	_, err = msgServer.RefuteSlashingProposal(sdk.WrapSDKContext(ctx), &types.MsgRefuteSlashingProposal{
		Sender:     addr.String(),
		Validator:  valAddr.String(),
		Refutation: "https://refutation.text",
	})
	require.NoError(t, err)

	proposal, found := app.CustomGovKeeper.GetProposal(ctx, proposalID)
	require.True(t, found)

	contentInterface := proposal.GetContent()
	content := contentInterface.(*types.ProposalSlashValidator)

	require.Equal(t, *content, types.ProposalSlashValidator{
		Offender:         valAddr.String(),
		StakingPoolId:    1,
		MisbehaviourTime: time.Time{},
		MisbehaviourType: "DoubleSign",
		JailPercentage:   10,
		Colluders:        []string(nil),
		Refutation:       "https://refutation.text",
	})
}
