package types

func NewDataRegistryEntry(hash, reference, encoding string, size uint64) DataRegistryEntry {
	return DataRegistryEntry{
		Hash:      hash,
		Reference: reference,
		Encoding:  encoding,
		Size_:     size,
	}
}
