package types

import "github.com/cosmos/cosmos-sdk/codec/types"

var _ types.UnpackInterfacesMessage = GenesisState{}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (data GenesisState) UnpackInterfaces(unpacker types.AnyUnpacker) error {
	for _, p := range data.Validators {
		err := p.UnpackInterfaces(unpacker)
		if err != nil {
			return err
		}
	}
	return nil
}
