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

	tokenInfo := a.keeper.GetTokenInfo(ctx, p.Denom)
	if tokenInfo != nil {
		tokenInfo.Name = p.Name
		tokenInfo.Symbol = p.Symbol
		tokenInfo.Icon = p.Icon
		tokenInfo.Description = p.Description
		tokenInfo.Website = p.Website
		tokenInfo.Social = p.Social
		tokenInfo.Inactive = p.Inactive
		tokenInfo.FeeRate = p.FeeRate
		tokenInfo.FeeEnabled = p.FeeEnabled
		tokenInfo.StakeCap = p.StakeCap
		tokenInfo.StakeMin = p.StakeMin
		tokenInfo.StakeEnabled = p.StakeEnabled
		return a.keeper.UpsertTokenInfo(ctx, *tokenInfo)
	}

	return a.keeper.UpsertTokenInfo(ctx, tokenstypes.NewTokenInfo(
		p.Denom, p.TokenType, p.FeeRate, p.FeeEnabled, p.Supply, p.SupplyCap, p.StakeCap, p.StakeMin, p.StakeEnabled, p.Inactive,
		p.Symbol, p.Name, p.Icon, p.Decimals,
		p.Description, p.Website, p.Social, p.Holders, p.MintingFee, p.Owner, p.OwnerEditDisabled, p.NftMetadata, p.NftHash,
	))
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
