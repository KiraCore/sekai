package rosetta

type ErrorDetails interface{}

type Error struct {
	Code        int          `json:"code"`
	Message     string       `json:"message"`
	Description string       `json:"description,omitempty"`
	Retriable   bool         `json:"retriable"`
	Details     ErrorDetails `json:"details,omitempty"`
}
