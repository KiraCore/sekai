package types

type FunctionParameter struct {
	Type        string              `json:"type"`
	Optional    bool                `json:"optional"`
	Description string              `json:"description"`
	Fields      *FunctionParameters `json:"fields,omitempty"`
}

type FunctionParameters = map[string]FunctionParameter

type FunctionMeta struct {
	FunctionID  int64              `json:"function_id"`
	Description string             `json:"description"`
	Parameters  FunctionParameters `json:"parameters"`
}

type FunctionList = map[string]FunctionMeta
