package rosetta

type SigningPayload struct {
	Address           string            `json:"address,omitempty"`
	AccountIdentifier AccountIdentifier `json:"account_identifier,omitempty"`
	HexBytes          string            `json:"hex_bytes"`
	SignatureType     SignatureType     `json:"signature_type,omitempty"`
}

type PublicKey struct {
	HexBytes  string    `json:"hex_bytes"`
	CurveType CurveType `json:"curve_type"`
}

type Signature struct {
	SigningPayload SigningPayload `json:"signing_payload"`
	PublicKey      PublicKey      `json:"public_key"`
	SignatureType  SignatureType  `json:"signature_type"`
	HexBytes       string         `json:"hex_bytes"`
}
