package types

import (
	"github.com/KiraCore/sekai/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.Msg = &MsgDisableBasketDeposits{}

// NewMsgUpsertTokenRate returns an instance of MsgUpserTokenRate
func NewMsgUpsertTokenRate(proposer sdk.AccAddress) *MsgDisableBasketDeposits {
	return &MsgDisableBasketDeposits{
		Sender: proposer.String(),
	}
}

// Route returns route
func (m *MsgDisableBasketDeposits) Route() string {
	return ModuleName
}

// Type returns return message type
func (m *MsgDisableBasketDeposits) Type() string {
	return types.MsgTypeUpsertTokenRate
}

// ValidateBasic returns basic validation result
func (m *MsgDisableBasketDeposits) ValidateBasic() error {
	return nil
}

// GetSignBytes returns to sign bytes
func (m *MsgDisableBasketDeposits) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

// GetSigners returns signers
func (m *MsgDisableBasketDeposits) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

// TODO:
//   // DisableBasketWithdraws - emergency function & permission to disable one or all withdrawals of one or all token in the basket
//   rpc DisableBasketWithdraws(MsgDisableBasketWithdraws) returns (MsgDisableBasketWithdrawsResponse);
//   // DisableBasketSwaps - emergency function & permission to disable one or all withdrawals of one or all token in the basket
//   rpc DisableBasketSwaps(MsgDisableBasketSwaps) returns (MsgDisableBasketSwapsResponse);
//   // BasketTokenMint - to mint basket tokens
//   rpc BasketTokenMint(MsgBasketTokenMint) returns (MsgBasketTokenMintResponse);
//   // BasketTokenBurn - to burn basket tokens and redeem underlying aggregate tokens
//   rpc BasketTokenBurn(MsgBasketTokenBurn) returns (MsgBasketTokenBurnResponse);
//   // BasketTokenSwap - to swap one or many of the basket tokens for one or many others
//   rpc BasketTokenSwap(MsgBasketTokenSwap) returns (MsgBasketTokenSwapResponse);
//   // BasketClaimRewards - to force staking derivative `SDB` basket to claim outstanding rewards of one all or many aggregate `V<ID>` tokens
//   rpc BasketClaimRewards(MsgBasketClaimRewards) returns (MsgBasketClaimRewardsResponse);
