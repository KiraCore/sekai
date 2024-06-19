package tokens

import (
	kiratypes "github.com/KiraCore/sekai/types"
	"github.com/KiraCore/sekai/x/gov/types"
	"github.com/KiraCore/sekai/x/tokens/keeper"
	tokenstypes "github.com/KiraCore/sekai/x/tokens/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type ApplyUpsertTokenInfosProposalHandler struct {
	keeper keeper.Keeper
}

func NewApplyUpsertTokenInfosProposalHandler(keeper keeper.Keeper) *ApplyUpsertTokenInfosProposalHandler {
	return &ApplyUpsertTokenInfosProposalHandler{keeper: keeper}
}

func (a ApplyUpsertTokenInfosProposalHandler) ProposalType() string {
	return kiratypes.ProposalTypeUpsertTokenInfos
}

func (a ApplyUpsertTokenInfosProposalHandler) Apply(ctx sdk.Context, proposalID uint64, proposal types.Content, slash sdk.Dec) error {
	p := proposal.(*tokenstypes.ProposalUpsertTokenInfo)

	rate := tokenstypes.NewTokenInfo(
		p.Denom, p.Rate, p.FeeEnabled, p.StakeCap, p.StakeMin, p.StakeEnabled, p.Inactive,
		p.Symbol, p.Name, p.Icon, p.Decimals,
	)
	return a.keeper.UpsertTokenInfo(ctx, rate)
}

type ApplyWhiteBlackChangeProposalHandler struct {
	keeper keeper.Keeper
}

func NewApplyWhiteBlackChangeProposalHandler(keeper keeper.Keeper) *ApplyWhiteBlackChangeProposalHandler {
	return &ApplyWhiteBlackChangeProposalHandler{keeper: keeper}
}

func (a ApplyWhiteBlackChangeProposalHandler) ProposalType() string {
	return kiratypes.ProposalTypeTokensWhiteBlackChange
}

func (a ApplyWhiteBlackChangeProposalHandler) Apply(ctx sdk.Context, proposalID uint64, proposal types.Content, slash sdk.Dec) error {
	p := proposal.(*tokenstypes.ProposalTokensWhiteBlackChange)

	if p.IsBlacklist {
		if p.IsAdd {
			a.keeper.AddTokensToBlacklist(ctx, p.Tokens)
		} else {
			a.keeper.RemoveTokensFromBlacklist(ctx, p.Tokens)
		}
	} else {
		if p.IsAdd {
			a.keeper.AddTokensToWhitelist(ctx, p.Tokens)
		} else {
			a.keeper.RemoveTokensFromWhitelist(ctx, p.Tokens)
		}
	}
	return nil
}
