package config

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	AccountAddressPrefix   = "kira"
	AccountPubKeyPrefix    = "kirapub"
	ValidatorAddressPrefix = "kiravaloper"
	ValidatorPubKeyPrefix  = "kiravaloperpub"
	ConsNodeAddressPrefix  = "kiravalcons"
	ConsNodePubKeyPrefix   = "kiravalconspub"
)

func SetConfig() {
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(AccountAddressPrefix, AccountPubKeyPrefix)
	config.SetBech32PrefixForValidator(ValidatorAddressPrefix, ValidatorPubKeyPrefix)
	config.SetBech32PrefixForConsensusNode(ConsNodeAddressPrefix, ConsNodePubKeyPrefix)
	config.Seal()
}