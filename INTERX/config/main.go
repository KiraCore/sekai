package config

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// AccountAddressPrefix is the prefix for account address
	AccountAddressPrefix = "kira"
	// AccountPubKeyPrefix is the prefix for account public key
	AccountPubKeyPrefix = "kirapub"
	// ValidatorAddressPrefix is the prefix for validator address
	ValidatorAddressPrefix = "kiravaloper"
	// ValidatorPubKeyPrefix is the prefix for validator public key
	ValidatorPubKeyPrefix = "kiravaloperpub"
	// ConsNodeAddressPrefix is the prefix for cons node address
	ConsNodeAddressPrefix = "kiravalcons"
	// ConsNodePubKeyPrefix is the prefix for cons node public key
	ConsNodePubKeyPrefix = "kiravalconspub"
)

// SetConfig is a function to set configuration for cosmos sdk
func SetConfig() {
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(AccountAddressPrefix, AccountPubKeyPrefix)
	config.SetBech32PrefixForValidator(ValidatorAddressPrefix, ValidatorPubKeyPrefix)
	config.SetBech32PrefixForConsensusNode(ConsNodeAddressPrefix, ConsNodePubKeyPrefix)
	config.Seal()
}
