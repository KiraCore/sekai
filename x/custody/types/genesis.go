package types

// DefaultGenesis returns the default CustomGo genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		MaxCustodyBufferSize: 10,
		MaxCustodyTxSize:     8192,
	}
}
