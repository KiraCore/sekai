package types

import (
	"github.com/KiraCore/sekai/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewMsgChangeCosmosEthereum(from sdk.AccAddress, to string, inAmount sdk.Coins, outAmount int64) *MsgChangeCosmosEthereum {
	return &MsgChangeCosmosEthereum{from, to, inAmount, outAmount}
}

func (m *MsgChangeCosmosEthereum) Route() string {
	return ModuleName
}

func (m *MsgChangeCosmosEthereum) Type() string {
	return types.MsgTypeChangeCosmosEthereum
}

func (m *MsgChangeCosmosEthereum) ValidateBasic() error {
	return nil
}

func (m *MsgChangeCosmosEthereum) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m *MsgChangeCosmosEthereum) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{
		m.From,
	}
}

func NewMsgChangeEthereumCosmos(addr sdk.AccAddress, from string, to sdk.AccAddress, inAmount int64, outAmount sdk.Coins) *MsgChangeEthereumCosmos {
	return &MsgChangeEthereumCosmos{addr, from, to, inAmount, outAmount}
}

func (m *MsgChangeEthereumCosmos) Route() string {
	return ModuleName
}

func (m *MsgChangeEthereumCosmos) Type() string {
	return types.MsgTypeChangeCosmosEthereum
}

func (m *MsgChangeEthereumCosmos) ValidateBasic() error {
	return nil
}

func (m *MsgChangeEthereumCosmos) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m *MsgChangeEthereumCosmos) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{
		m.Addr,
	}
}
