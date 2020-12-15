package common

// DataReferenceEntry is a struct to be used for a data reference
type DataReferenceEntry struct {
	Hash      string `json:"hash"`
	Reference string `json:"reference"`
	Encoding  string `json:"encoding"`
	Size      uint64 `json:"size"`
}

// DataRefs will save all data references
var DataRefs = make(map[string]DataReferenceEntry)

// UpdateKey will save a data reference
func UpdateKey(key string, reference DataReferenceEntry) {
	DataRefs[key] = reference
}
