package types

type FunctionParameter struct {
	Type        string              `json:"type"`
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

func NewFunctionMeta(funcID int64, desc string, params FunctionParameters) FunctionMeta {
	return FunctionMeta{
		FunctionID:  funcID,
		Description: desc,
		Parameters:  params,
	}
}

func NewFunctionParameter(paramType string, desc string) FunctionParameter {
	return FunctionParameter{
		Type:        paramType,
		Description: desc,
	}
}
