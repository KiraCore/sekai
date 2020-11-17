package functions

import (
	functionmeta "github.com/KiraCore/sekai/function_meta"
	sekaitypes "github.com/KiraCore/sekai/types"
)

// GetAllFunctions is a function to get all functions registered
func GetAllFunctions() sekaitypes.FunctionList {
	return functionmeta.GetFunctionList()
}
