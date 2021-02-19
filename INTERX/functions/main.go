package functions

import (
	functionmeta "github.com/KiraCore/sekai/function_meta"
	sekaitypes "github.com/KiraCore/sekai/types"
)

// GetKiraFunctions is a function to get all kira functions registered
func GetKiraFunctions() sekaitypes.FunctionList {
	return functionmeta.GetFunctionList()
}

// GetInterxFunctions is a function to get all interx functions registered
func GetInterxFunctions() InterxFunctionList {
	return interxFunctions
}
