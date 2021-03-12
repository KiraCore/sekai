package rosetta

type ErrorDetails interface{}

type Error struct {
	Code        int64        `json:"code"`
	Message     string       `json:"message"`
	Description string       `json:"description,omitempty"`
	Retriable   bool         `json:"retriable"`
	details     ErrorDetails `json:"details,omitempty"`
}
